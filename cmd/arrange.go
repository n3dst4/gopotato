/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func arrange() {

}

// arrangeCmd represents the arrange command
var arrangeCmd = &cobra.Command{
	Use:   "arrange",
	Short: "Idempotently creates today's journal and organises older journals",
	Long: `Creates today's journal in the format YYYY-MM-DD.md, and moves older
	journals to YYYY-MM/ if they are older than the configured threshold.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("arrange command called")
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
