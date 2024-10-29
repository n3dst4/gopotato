/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/n3dst4/gopotato/utils"
	"github.com/spf13/cobra"
)

var baselineConfig = utils.BaselineConfig{
	RootPath: "~/Sync/notes",
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if config file already exists, exit
		if _, err := os.Stat(configFilePath); err == nil {
			color.Red("Config file already exists: %s", configFilePath)
			os.Exit(1)
		}
		color.Green("Creating config file: %s", configFilePath)
		file, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		err = toml.NewEncoder(file).Encode(baselineConfig)
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(fmt.Sprintf("# %s\n\n", configFilePath))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
