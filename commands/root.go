// Commands for the potato system
package commands

import (
	"log"
	"os"

	cfg "github.com/n3dst4/gopotato/config"
	"github.com/n3dst4/gopotato/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// constants
const LONGDESC = `
Gopotato is a journal manager for Go. It is a rewrite of Potato.
`

// values set at runtime
var config *cfg.Config
var configFilePath string = ""

// var usr *user.User

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopotato",
	Short: "Go rewrite of potato jornal manager",
	Long:  LONGDESC,

	// PersistentPreRun - always runs before any subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Parent().Name() == "completion" {
			return
		}

		configFilePath = utils.TildeToHomeDir(configFilePath)
		color.Cyan("Reading config file from %s", configFilePath)
		// read config file (configFilePath has been set by cobra)
		txt, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		config, err = cfg.ParseConfig(string(txt))
		if err != nil {
			log.Fatal(err)
		}
	},

	// Run
	Run: func(cmd *cobra.Command, args []string) {
		arrange()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configFilePath,
		"config",
		"c",
		"~/.local/config/potato.toml",
		"config file (default is $HOME/.gopotato.yaml)",
	)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
