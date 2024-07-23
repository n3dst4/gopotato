// Commands for the potato system
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
	"github.com/kr/pretty"

	"github.com/n3dst4/gopotato/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// constants
const LONGDESC = `
Gopotato is a journal manager for Go. It is a rewrite of Potato.
`

// values set at runtime
// var config *cfg.Config

var config = &utils.Config{
	KeepDays:     7,
	JournalsPath: "journals",
	PagesPath:    "pages",
	CarryOverTodos: utils.CarryOverTodos{
		OnlyIncomplete: true,
	},
}

var configFilePath string = ""

// var usr *user.User

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopotato",
	Short: "Go rewrite of potato jornal manager",
	Long:  LONGDESC,

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

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("initConfig", configFilePath)
	if configFilePath != "" {
		configFilePath = utils.TildeToHomeDir(configFilePath)
		viper.SetConfigFile(configFilePath)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(filepath.Join(home, ".local", "config"))
		viper.SetConfigType("toml")
		viper.SetConfigName("potato")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&config)
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(config)
	if err != nil {
		log.Fatal(err)
	}
	config.RootPath = utils.TildeToHomeDir(config.RootPath)
	config.JournalsPath = fmt.Sprintf("%s/%s", config.RootPath, config.JournalsPath)
	config.PagesPath = fmt.Sprintf("%s/%s", config.RootPath, config.PagesPath)

	color.Cyan(pretty.Sprint("Config loaded:", config))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
