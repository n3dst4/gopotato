/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var arrangeCmd = &cobra.Command{
	Use:   "arrange",
	Short: "Idempotently creates today's journal and organises older journals",
	Long: `Creates today's journal in the format YYYY-MM-DD.md, and moves older
	journals to YYYY-MM/ if they are older than the configured threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		arrange()
	},
}

var journalFilenameRegex *regexp.Regexp
var monthArchiveFoldernameRegex *regexp.Regexp

func init() {
	rootCmd.AddCommand(arrangeCmd)
	journalFilenameRegex, _ = regexp.Compile(
		`^(\d{4})[-_](\d{2})[-_](\d{2}).(?:md|markdown)$`,
	)
	monthArchiveFoldernameRegex, _ = regexp.Compile(
		`^(\d{4})[-_](\d{2})$`,
	)
}

// does the daily chores
func arrange() {
	validateConfig()
	err := os.Chdir(config.JournalsPath)
	defer os.Chdir(config.RootPath)
	if err != nil {
		log.Fatal(err)
	}

	archiveOldJournals()
	archiveOldMonths()
	createTodaysJournal()
}

// archiveOldJournals moves journals older than the configured threshold to a
// YYYY-MM/ folder
func archiveOldJournals() {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	days := config.KeepDays * -1
	cutoffDate := now.AddDate(0, 0, days)

	for _, file := range entries {
		matches := journalFilenameRegex.FindStringSubmatch(file.Name())

		if file.IsDir() {
			continue
		}
		if len(matches) == 0 {
			color.Red("Skipping unrecognized file %s", file.Name())
			continue
		}

		fileYear, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatalf("Invalid file year: %v\n", err)
		}
		fileMonth, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatalf("Invalid file month: %v\n", err)
		}
		fileDay, err := strconv.Atoi(matches[3])
		if err != nil {
			log.Fatalf("Invalid file day: %v\n", err)
		}

		fileDate := time.Date(fileYear, time.Month(fileMonth), fileDay, 0, 0, 0, 0, time.UTC)

		if fileDate.Before(cutoffDate) {
			folderPath := config.JournalsPath + "/" + fileDate.Format("2006-01")
			fileName := fileDate.Format("2006-01-02") + ".md"
			color.Yellow("Archiving old file %s to %s/%s", file.Name(), folderPath, fileName)
			os.MkdirAll(folderPath, os.ModePerm)
			os.Rename(file.Name(), fmt.Sprintf("%s/%s", folderPath, fileName))
		}
	}
}

func archiveOldMonths() {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	months := config.KeepMonths * -1
	// Set to the first of the month to avoid partial months
	now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	cutoffDate := now.AddDate(0, months, 0)

	for _, folder := range entries {
		matches := monthArchiveFoldernameRegex.FindStringSubmatch(folder.Name())

		if len(matches) == 0 || !folder.IsDir() {
			continue
		}

		folderYear, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatalf("Invalid file year: %v\n", err)
		}
		folderMonth, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatalf("Invalid file month: %v\n", err)
		}

		folderDate := time.Date(folderYear, time.Month(folderMonth), 1, 0, 0, 0, 0, time.UTC)

		if folderDate.Before(cutoffDate) {
			archivePath := config.JournalsPath + "/" + folderDate.Format("2006") + "/" + folderDate.Format("2006-01")
			color.Yellow("Archiving old folder %s to %s", folder.Name(), archivePath)
			os.MkdirAll(archivePath, os.ModePerm)
			// merge everything in folder.Name to archivePath/folderName
			entries, err := os.ReadDir(folder.Name())
			if err != nil {
				log.Fatalf("Failed to read directory %s: %v\n", folder.Name(), err)
			}

			for _, entry := range entries {
				// if entry.IsDir() {
				// 	continue
				// }
				// move the file to the new location
				newPath := filepath.Join(archivePath, entry.Name())
				// os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
				if err := os.Rename(filepath.Join(folder.Name(), entry.Name()), newPath); err != nil {
					log.Fatalf("Failed to move file %s to %s: %v\n", entry.Name(), newPath, err)
				}
			}
			// remove the old folder
			if err := os.RemoveAll(folder.Name()); err != nil {
				log.Fatalf("Failed to remove old folder %s: %v\n", folder.Name(), err)
			}
		}
	}

}

// createTodaysJournal creates a new journal file if one does not already exist
func createTodaysJournal() {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	todayDateString := time.Now().Format("2006-01-02")
	todayFileName := todayDateString + ".md"

	var todayFile fs.DirEntry
	for _, file := range entries {
		if file.Name() == todayFileName && file.Type().IsRegular() {
			todayFile = file
		}
	}

	if todayFile == nil {
		labelColor := color.New(color.FgGreen)
		pathColor := color.New(color.FgCyan)
		labelColor.Printf("Creating new journal file: ")
		pathColor.Printf("%s\n", config.JournalsPath+"/"+todayFileName)
		file, err := os.Create(todayFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("# %s\n\n", todayDateString))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		labelColor := color.New(color.FgBlue)
		pathColor := color.New(color.FgCyan)
		labelColor.Printf("Existing journal file: ")
		pathColor.Printf("%s\n", config.JournalsPath+"/"+todayFileName)
	}
	fmt.Println()

}
