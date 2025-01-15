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
		"no arguments provided -> exit code 1, help user and display usage": {
			ExpectOutput: []string{
				"No location arguments provided, please provide at least one location name, ZIP or Postal Code.",
				usageMessage,
			},
			ExpectError: true,
		},
		"geo -h -> display usage": {
			Args: []string{
				"-h",
			},
			ExpectOutput: []string{
				usageMessage,
			},
		},
		"geo 23228 -> Henrico, VA": {
			Args: []string{"23228"},
			ExpectOutput: []string{
				"'23228' results:",
				"  Name: Henrico County, US, 23228",
				"  Lat,Lon: 37.463800, -77.398000",
			},
		},
		"geo 'Henrico, VA' -> loc(Henrico, VA)": {
			Args: []string{"Henrico, VA"},
			ExpectOutput: []string{
				"'Henrico, VA' results:",
				"  Name: Henrico, Virginia, US",
				"  Lat,Lon: 37.495702, -77.335257",
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
				t.Fatalf("unexpected error: %s", err)
			} else if err == nil && tc.ExpectError {
				t.Fatalf("expected error did not occur: %s", err)
			}

			outputMatcher := internaltesting.NewOutputMatcher(string(outStr))

			for _, expectedOutput := range tc.ExpectOutput {
				outputMatcher.MatchText(t, expectedOutput)
			}
		})
	}
}

func TestIntegrationWithMissingKey(t *testing.T) {
	internaltesting.MustCompileOnce(t)

	tests := map[string]struct {
		Args         []string
		ExpectOutput []string
		ExpectError  bool
	}{
		"geo 23228 -> API Key not provided": {
			Args: []string{"23228"},
			ExpectOutput: []string{
				fmt.Sprintf("'%s' not set. Please visit 'https://openweathermap.org/api' and obtain an API key.", cmd.ApiKeyName),
				fmt.Sprintf("Set the key before runing 'geo' with:\n\texport %s=<your openweather api key>", cmd.ApiKeyName),
			},
			ExpectError: true,
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

			outStr, err := geoCmd.CombinedOutput()
			if err != nil && !tc.ExpectError {
				t.Fatalf("unexpected error: %s", err)
			} else if err == nil && tc.ExpectError {
				t.Fatalf("expected error did not occur: %s", err)
			}

			outputMatcher := internaltesting.NewOutputMatcher(string(outStr))

			for _, expectedOutput := range tc.ExpectOutput {
				outputMatcher.MatchText(t, expectedOutput)
			}
		})
	}

}

//func TestIntegrationWithInvalidKey(t *testing.T) {
//	internaltesting.MustCompileOnce(t)
//
//	if ApiKey == "" {
//		t.Fatalf("No %s set", cmd.ApiKeyName)
//	}
//
//	tests := map[string]struct {
//		Args         []string
//		ExpectOutput []string
//		ExpectError  bool
//	}{
//		"no arguments provided -> exit code 1, help user and display usage": {
//			ExpectOutput: []string{
//				"No location arguments provided, please provide at least one location name, ZIP or Postal Code.",
//				usageMessage,
//			},
//			ExpectError: true,
//		},
//		"geo 23228 -> Henrico, VA": {
//			Args: []string{"23228"},
//			ExpectOutput: []string{
//				"Henrico, VA",
//			},
//		},
//	}
//
//	for name, tc := range tests {
//		t.Helper()
//		t.Run(name, func(t *testing.T) {
//
//			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
//			defer cancel()
//
//			geoCmd := exec.CommandContext(ctx, "../../build/geo", tc.Args...)
//			geoCmd.Env = []string{
//				fmt.Sprintf("%s=%s", cmd.ApiKeyName, ApiKey),
//			}
//
//			outStr, err := geoCmd.CombinedOutput()
//			if err != nil && !tc.ExpectError {
//				t.Fatalf("unexpected error: %s", err)
//			} else if err == nil && tc.ExpectError {
//				t.Fatalf("expected error did not occur: %s", err)
//			}
//
//			outputMatcher := internaltesting.NewOutputMatcher(string(outStr))
//
//			for _, expectedOutput := range tc.ExpectOutput {
//				outputMatcher.MatchText(t, expectedOutput)
//			}
//		})
//	}
//
//}
//
