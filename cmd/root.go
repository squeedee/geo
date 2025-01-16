package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	internalcmd "github.com/squeedee/geo/internal/cmd"
	"net/http"
	"strconv"

	"os"
)

const ApiKeyName = "OPEN_WEATHER_API_KEY"

var RootCmd = &cobra.Command{
	Use:     "geo",
	Short:   "Geo-locate place names and zip codes within the USA",
	Example: "  geo \"Henrico, VA\" 10001 \"Seattle, WA\"",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("No location arguments provided, please provide at least one location name, ZIP or Postal Code.\n\n")
			_ = cmd.Usage()
			os.Exit(1)
		}

		apiKey, apiKeyFound := os.LookupEnv(ApiKeyName)
		if !apiKeyFound || apiKey == "" {
			fmt.Printf("'%s' not set. Please visit 'https://openweathermap.org/api' and obtain an API key.", ApiKeyName)
			fmt.Printf("Set the key before runing 'geo' with:\n\texport %s=<your openweather api key>", ApiKeyName)
			os.Exit(1)
		}

		g := internalcmd.DirectGeocoding{
			Key: apiKey,
		}

		for _, arg := range args {

			fmt.Printf("'%s' results:\n", arg)

			if _, conversionErr := strconv.Atoi(arg); conversionErr == nil { // numeric, use zip
				loc, code, err := g.LocationByZip(arg)
				if code == http.StatusUnauthorized {
					fmt.Printf("'%s' is invalid. Please ensure you have the correct key from 'https://openweathermap.org/api'.\n", ApiKeyName)
					os.Exit(1)
				}
				if loc == nil {
					fmt.Println("  No matches found.")
					os.Exit(1)
				}
				if err != nil {
					fmt.Printf("unexpected error when getting the location '%s': %s\n", arg, err)
					os.Exit(1)
				}
				fmt.Printf("  Name: %s, %s, %s\n", loc.Name, loc.Country, loc.Zip)
				fmt.Printf("  Lat,Lon: %f, %f\n\n", loc.Lat, loc.Lon)
			} else { // non-numeric, use name
				locations, code, err := g.LocationByName(arg)
				if code == http.StatusUnauthorized {
					fmt.Printf("'%s' is invalid. Please ensure you have the correct key from 'https://openweathermap.org/api'.\n", ApiKeyName)
					os.Exit(1)
				}
				if err != nil {
					fmt.Printf("unexpected error when getting the location '%s': %s\n", arg, err)
					os.Exit(1)
				}
				if len(locations) == 0 {
					fmt.Printf("  No matches found.\n")
					os.Exit(1)
				}
				for _, loc := range locations {
					fmt.Printf("  Name: %s, %s, %s\n", loc.Name, loc.State, loc.Country)
					fmt.Printf("  Lat,Lon: %f, %f\n\n", loc.Lat, loc.Lon)
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
