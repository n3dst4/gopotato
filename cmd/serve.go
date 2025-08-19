/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Not implemented yet",
	Long: `One day this will launch a webserver.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		viper.Get("config")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
