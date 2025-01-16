//go:build e2e

package integration_test

import (
	"context"
	"fmt"
	. "github.com/MakeNowJust/heredoc/dot"
	"github.com/squeedee/geo/cmd"
	internaltesting "github.com/squeedee/geo/internal/testing"
	"os"
	"os/exec"
	"testing"
	"time"
)

var ApiKey = os.Getenv(cmd.ApiKeyName)

// Note: I'm not a fan of checking for minutiae in output created by a reliable framework like 'cobra',
// but I've included the full usage message for completeness.
var usageMessage = D(`Usage:
  geo [flags]

Examples:
  geo "Henrico, VA" 10001 "Seattle, WA"

Flags:
  -h, --help   help for geo`)

func TestIntegrationWithWorkingKey(t *testing.T) {
	internaltesting.MustCompileOnce(t)

	if ApiKey == "" {
		t.Fatalf("No %s set", cmd.ApiKeyName)
	}

	tests := map[string]struct {
		Args         []string
		ExpectOutput []string
		ExpectError  bool
	}{
		"'geo', no arguments provided -> exit code 1, help user and display usage": {
			ExpectOutput: []string{
				"No location arguments provided, please provide at least one location name, ZIP or Postal Code.",
				usageMessage,
			},
			ExpectError: true,
		},
		"'geo -h', help flag -> displays usage": {
			Args: []string{
				"-h",
			},
			ExpectOutput: []string{
				usageMessage,
			},
		},
		"'geo 23228', valid zip code -> Henrico, VA": {
			Args: []string{"23228"},
			ExpectOutput: []string{
				"'23228' results:",
				"  Name: Henrico County, US, 23228",
				"  Lat,Lon: 37.463800, -77.398000",
			},
		},
		"'geo 99999', invalid zip code -> exit code 1, display not found message": {
			Args: []string{"99999"},
			ExpectOutput: []string{
				"'99999' results:",
				"  No matches found.",
			},
			ExpectError: true,
		},
		"'geo \"Henrico, VA\"', valid place name -> loc(Henrico, VA)": {
			Args: []string{"Henrico, VA"},
			ExpectOutput: []string{
				"'Henrico, VA' results:",
				"  Name: Henrico, Virginia, US",
				"  Lat,Lon: 37.495702, -77.335257",
			},
		},
		"'geo \"San José, CA\"', special characters -> loc(San José, CA)": {
			Args: []string{"San José, CA"},
			ExpectOutput: []string{
				"'San José, CA' results:",
				"Name: San Jose, California, US",
				"Lat,Lon: 37.336166, -121.890591",
			},
		},
		"'geo \"not-a-place, NY\"', invalid place -> exit code 1, display not found message": {
			Args: []string{"not-a-place, NY"},
			ExpectOutput: []string{
				"'not-a-place, NY' results:",
				"  No matches found.",
			},
			ExpectError: true,
		},
		"'geo \"壚靁縅-lalala\"', invalid place and special characters -> exit code 1, display not found message": {
			Args: []string{"壚靁縅-lalala"},
			ExpectOutput: []string{
				"'壚靁縅-lalala' results:",
				"  No matches found.",
			},
			ExpectError: true,
		},
		"Test text layout": {
			Args: []string{"Henrico, VA", "10001", "Seattle, WA"},
			ExpectOutput: []string{
				"'Henrico, VA' results:\n",
				"  Name: Henrico, Virginia, US\n",
				"  Lat,Lon: 37.495702, -77.335257\n",
				"\n",
				"'10001' results:\n",
				"  Name: New York, US, 10001\n",
				"  Lat,Lon: 40.748400, -73.996700\n",
				"\n",
				"'Seattle, WA' results:\n",
				"  Name: Seattle, Washington, US\n",
				"  Lat,Lon: 47.603832, -122.330062\n",
				"\n",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
			defer cancel()

			geoCmd := exec.CommandContext(ctx, "../../build/geo", tc.Args...)
			geoCmd.Env = []string{
				fmt.Sprintf("%s=%s", cmd.ApiKeyName, ApiKey),
			}

			outStr, err := geoCmd.CombinedOutput()
			if err != nil && !tc.ExpectError {
				t.Logf("unexpected error: %s", err)
				t.Fail()
			} else if err == nil && tc.ExpectError {
				t.Logf("expected error did not occur")
				t.Fail()
			}

			outputMatcher := internaltesting.NewOutputMatcher(string(outStr))

			for _, expectedOutput := range tc.ExpectOutput {
				outputMatcher.MatchText(t, expectedOutput)
			}
		})
	}
}

func TestIntegrationApiKey(t *testing.T) {
	internaltesting.MustCompileOnce(t)

	if ApiKey == "" {
		t.Fatalf("No %s set", cmd.ApiKeyName)
	}

	tests := map[string]struct {
		ApiKey       *string
		Args         []string
		ExpectOutput []string
		ExpectError  bool
	}{
		"API Key not provided -> provide user guidance": {
			ApiKey: nil,
			Args:   []string{"23228"},
			ExpectOutput: []string{
				fmt.Sprintf("'%s' not set. Please visit 'https://openweathermap.org/api' and obtain an API key.", cmd.ApiKeyName),
				fmt.Sprintf("Set the key before runing 'geo' with:\n\texport %s=<your openweather api key>", cmd.ApiKeyName),
			},
			ExpectError: true,
		},
		"API Key Empty -> provide user guidance": {
			ApiKey: internaltesting.StrPtr(""),
			Args:   []string{"23228"},
			ExpectOutput: []string{
				fmt.Sprintf("'%s' not set. Please visit 'https://openweathermap.org/api' and obtain an API key.", cmd.ApiKeyName),
				fmt.Sprintf("Set the key before runing 'geo' with:\n\texport %s=<your openweather api key>", cmd.ApiKeyName),
			},
			ExpectError: true,
		},
		"API Key Invalid, zip code -> provide user guidance": { // Testing zip code path NOTE1
			ApiKey: internaltesting.StrPtr("invalid-key"),
			Args:   []string{"23228"},
			ExpectOutput: []string{
				fmt.Sprintf("'%s' is invalid. Please ensure you have the correct key from 'https://openweathermap.org/api'.", cmd.ApiKeyName),
			},
			ExpectError: true,
		},
		"API Key Invalid, location name -> provide user guidance": { // Testing place code path NOTE1
			ApiKey: internaltesting.StrPtr("invalid-key"),
			Args:   []string{"Henrico, Va"},
			ExpectOutput: []string{
				fmt.Sprintf("'%s' is invalid. Please ensure you have the correct key from 'https://openweathermap.org/api'.", cmd.ApiKeyName),
			},
			ExpectError: true,
		},
		"API Key Invalid, multiple arguments -> provide user guidance": {
			ApiKey: internaltesting.StrPtr("invalid-key"),
			Args:   []string{"23228", "Henrico, Va"},
			ExpectOutput: []string{
				fmt.Sprintf("'%s' is invalid. Please ensure you have the correct key from 'https://openweathermap.org/api'.", cmd.ApiKeyName),
			},
			ExpectError: true,
		},
		"API Key valid -> no error": {
			ApiKey:      internaltesting.StrPtr(ApiKey),
			Args:        []string{"23228"},
			ExpectError: false,
		},
	}

	for name, tc := range tests {
		t.Helper()
		t.Run(name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
			defer cancel()

			// Ensure this subshell does not have the apikey, which would poison the test command's subshell.
			err := os.Unsetenv(cmd.ApiKeyName)
			if err != nil {
				t.Fatal("failed to clear the API key from our environment")
			}

			geoCmd := exec.CommandContext(ctx, "../../build/geo", tc.Args...)

			if tc.ApiKey != nil {
				geoCmd.Env = []string{
					fmt.Sprintf("%s=%s", cmd.ApiKeyName, *tc.ApiKey),
				}
			}

			outStr, err := geoCmd.CombinedOutput()
			if err != nil && !tc.ExpectError {
				t.Logf("unexpected error: %s", err)
				t.Fail()
			} else if err == nil && tc.ExpectError {
				t.Logf("expected error did not occur")
				t.Fail()
			}

			outputMatcher := internaltesting.NewOutputMatcher(string(outStr))

			for _, expectedOutput := range tc.ExpectOutput {
				outputMatcher.MatchText(t, expectedOutput)
			}
		})
	}
}
