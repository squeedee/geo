package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	internalcmd "github.com/squeedee/geo/internal/cmd"
	"strconv"

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

		var apiKey string
		var found bool
		if apiKey, found = os.LookupEnv(ApiKeyName); !found {
			fmt.Printf("'%s' not set. Please visit 'https://openweathermap.org/api' and obtain an API key.", ApiKeyName)
			fmt.Printf("Set the key before runing 'geo' with:\n\texport %s=<your openweather api key>", ApiKeyName)
			os.Exit(1)
		}

		g := internalcmd.DirectGeocoding{
			Key: apiKey,
		}

		for _, arg := range args {

			fmt.Printf("'%s' results:\n", arg)

			if _, conversionErr := strconv.Atoi(arg); conversionErr == nil {
				loc, err := g.LocationByZip(arg)
				if err != nil {
					fmt.Printf("unexpected error when getting the location '%s': %s\n", arg, err)
					os.Exit(1)
				}
				fmt.Printf("  Name: %s, %s, %s\n", loc.Name, loc.Country, loc.Zip)
				fmt.Printf("  Lat,Lon: %f, %f\n", loc.Lat, loc.Lon)
			} else { // non-numeric, use name
				locations, err := g.LocationByName(arg)
				if err != nil {
					fmt.Printf("unexpected error when getting the location '%s': %s\n", arg, err)
					os.Exit(1)
				}
				for _, loc := range locations {
					fmt.Printf("  Name: %s, %s, %s\n", loc.Name, loc.State, loc.Country)
					fmt.Printf("  Lat,Lon: %f, %f\n", loc.Lat, loc.Lon)
				}
			}

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
