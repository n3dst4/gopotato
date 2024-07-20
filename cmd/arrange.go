/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var tomlString = `rootPath = "~/Sync/notes"
keepDays = 6

[carryOverTodos]
enabled = true
onlyIncomplete = true
headingRegex = "^#to[- ]?dos? *$"
headingRegexCaseSensitive = false`

func arrange() {

	var config struct {
		RootPath       string
		KeepDays       int
		CarryOverTodos struct {
			Enabled                   bool
			OnlyIncomplete            bool
			HeadingRegex              string
			HeadingRegexCaseSensitive bool
		}
	}

	txt, err := os.ReadFile("/home/ndc/.config/potato.toml")

	if err != nil {
		log.Fatal(err)
	}

	// convert txt to string
	var txtString = string(txt)

	md, err := toml.Decode(txtString, &config)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(md)

}

// arrangeCmd represents the arrange command
var arrangeCmd = &cobra.Command{
	Use:   "arrange",
	Short: "Idempotently creates today's journal and organises older journals",
	Long: `Creates today's journal in the format YYYY-MM-DD.md, and moves older
	journals to YYYY-MM/ if they are older than the configured threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		// arrange()
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
