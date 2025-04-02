/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
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

func init() {
	rootCmd.AddCommand(arrangeCmd)
	journalFilenameRegex, _ = regexp.Compile(
		`^(\d{4})[-_](\d{2})[-_](\d{2}).(?:md|markdown)$`,
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
	oneMonthAgo := now.AddDate(0, -1, 0)

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

		if fileDate.Before(oneMonthAgo) {
			folderPath := config.JournalsPath + "/" + fileDate.Format("2006-01")
			fileName := fileDate.Format("2006-01-02") + ".md"
			color.Yellow("Archiving old file %s to %s/%s", file.Name(), folderPath, fileName)
			os.MkdirAll(folderPath, os.ModePerm)
			os.Rename(file.Name(), fmt.Sprintf("%s/%s", folderPath, fileName))
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
		color.Green("Creating new journal file: %s", (config.JournalsPath + "/" + todayFileName))
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
		color.Blue("Existing journal file: %s", config.JournalsPath+"/"+todayFileName)
	}

}
