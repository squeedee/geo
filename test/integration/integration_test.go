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
		"geo 23228 -> Henrico, VA": {
			Args: []string{"23228"},
			ExpectOutput: []string{
				"Henrico, VA",
			},
		},
	}

	for name, tc := range tests {
		t.Helper()
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

//func TestApiKey(t *testing.T) {
//
//	tests := map[string]struct {
//		Key            *string // API key for OpenWeather. Use Nil for not set, and "" for set but empty.
//		Args         string  // single string representing all params to pass to the command
//		ExpectOutput   func(t *testing.T, s bufio.Scanner)
//		ExpectExitCode int
//	}{
//		"no params provided -> exit code 1, help user and display usage": {
//			Args:         "",
//			ExpectExitCode: 1,
//			ExpectOutput: func(t *testing.T, s bufio.Scanner) {
//				ScanFor(t, s, "No parameters provided, please provide at least one location name, ZIP or Postal Code.")
//			},
//		},
//	}
//
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//
//		}
//	}
//
//}
