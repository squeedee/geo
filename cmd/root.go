package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"os"
)

const ApiKeyName = "OPEN_WEATHER_API_KEY"

var RootCmd = &cobra.Command{
	Use:   "geo",
	Short: "Geo-locate place names and zip codes within the USA",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("No location arguments provided, please provide at least one location name, ZIP or Postal Code.\n\n")
			_ = cmd.Usage()
			os.Exit(1)
		}

	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
}
