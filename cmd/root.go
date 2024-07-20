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
	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopotato",
	Short: "Go rewrite of potato jornal manager",
	Long:  longDesc,

	// PersistentPreRun
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		color.Cyan("Reading config file from %s", configFilePath)

		// read config file (confiFilePath has been set by cobra)
		txt, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// deserialize config file into &config, overwriting existing values
		_, err = toml.Decode(string(txt), &config)
		if err != nil {
			log.Fatal(err)
		}

		// pretty print config
		color.Cyan(pretty.Sprint("Config loaded:", config))

		// validate config
		validate = validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(config)
		if err != nil {
			log.Fatal(err)
		}
	},

	// Run
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
