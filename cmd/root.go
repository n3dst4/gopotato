/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

var longDesc = `
Gopotato is a journal manager for Go. It is a rewrite of Potato.
`

var configFilePath string = ""

type CarryOverTodos struct {
	Enabled                   bool
	OnlyIncomplete            bool
	HeadingRegex              string
	HeadingRegexCaseSensitive bool
}

type Config struct {
	RootPath       string
	KeepDays       int
	CarryOverTodos CarryOverTodos
}

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopotato",
	Short: "Go rewrite of potato jornal manager",
	Long:  longDesc,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		txt, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		_, err = toml.Decode(string(txt), &config)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println("Config loaded:", config)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command called")
		arrange()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	rootCmd.PersistentFlags().StringVar(
		&configFilePath,
		"config",
		usr.HomeDir+"/.local/config/potato.toml",
		"config file (default is $HOME/.gopotato.yaml)",
	)
}
