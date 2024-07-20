/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func arrange() {
	// read current directory
	err := os.Chdir(config.JournalsPath)
	defer os.Chdir(config.RootPath)
	if err != nil {
		log.Fatal(err)
	}

	archiveOldJournals()
	createTodaysJournal()
}

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

func createTodaysJournal() {
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// const todayDate = new Date().toISOString().slice(0, 10);
	// const todayFileName = `${todayDate}.md`;
	todayDateString := time.Now().Format("2006-01-02")
	todayFileName := todayDateString + ".md"
	// const todayFile = Array.from(journalFiles).find(
	//   (file) => file.name === todayFileName
	// );
	var todayFile fs.DirEntry
	for _, file := range entries {
		if file.Name() == todayFileName && file.Type().IsRegular() {
			todayFile = file
		}
	}
	if todayFile == nil {
		color.Green("Creating new journal file: %s", todayFileName)
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
		color.Blue("Existing journal file: %s", todayFileName)
	}

	// const fullTodayFileName = join(journalsPath, todayFileName);
	// // if not, create a new file with today's date
	// if (!todayFile) {
	//   console.log(chalk.green(`Creating new journal file: ${fullTodayFileName}`));
	//   const todayFile = Deno.createSync(fullTodayFileName);
	//   writeAllSync(todayFile, new TextEncoder().encode(`# ${todayDate}\n\n`));
	// } else {
	//   console.log(chalk.blue(`Existing journal file: ${fullTodayFileName}`));
	// }
}

// arrangeCmd represents the arrange command
var arrangeCmd = &cobra.Command{
	Use:   "arrange",
	Short: "Idempotently creates today's journal and organises older journals",
	Long: `Creates today's journal in the format YYYY-MM-DD.md, and moves older
	journals to YYYY-MM/ if they are older than the configured threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		arrange()
	},
}

func init() {
	rootCmd.AddCommand(arrangeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// arrangeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// arrangeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
