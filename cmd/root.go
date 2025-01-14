package cmd

import (
	"github.com/spf13/cobra"

	"os"
)

var RootCmd = &cobra.Command{
	Use:   "geo",
	Short: "Geo-locate place names and zip codes within the USA",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
