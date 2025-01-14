package cmd_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/squeedee/geo/cmd"
	internalcmd "github.com/squeedee/geo/internal/cmd"
	"os"
	"testing"
)

var ApiKey = os.Getenv(cmd.ApiKeyName)

func TestDirectGeocoding_CoordsByLocationName(t *testing.T) {
	if ApiKey == "" {
		t.Fatalf("No %s set", cmd.ApiKeyName)
	}

	defaultFields := internalcmd.DirectGeocoding{
		Key: ApiKey,
	}

	tests := map[string]struct {
		name             string
		fields           internalcmd.DirectGeocoding
		expectedError    string
		expectedLocation *internalcmd.Locations
	}{
		"Richmond, well formed without country code": {
			fields: defaultFields,
			name:   "Richmond, VA",
			expectedLocation: &internalcmd.Locations{
				internalcmd.Location{
					Name:    "Richmond",
					Lat:     37.5385087,
					Lon:     -77.43428,
					Country: "US",
					State:   "Virginia",
				},
			},
		},
		"Richmond, well formed with country code": {
			fields: defaultFields,
			name:   "Richmond, VA, USA",
			expectedLocation: &internalcmd.Locations{
				internalcmd.Location{
					Name:    "Richmond",
					Lat:     37.5385087,
					Lon:     -77.43428,
					Country: "US",
					State:   "Virginia",
				},
			},
		},
		"Melbourne, well formed with country code": {
			fields: defaultFields,
			name:   "Melbourne, VIC, AU",
			expectedLocation: &internalcmd.Locations{
				internalcmd.Location{
					Name:    "Melbourne",
					Lat:     -37.8142176,
					Lon:     144.9631608,
					Country: "AU",
					State:   "Victoria",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			location, err := tc.fields.CoordsByLocationName(tc.name)
			if tc.expectedError == "" && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			} else if tc.expectedError != "" && err == nil {
				t.Fatalf("Expected error, but didn't get one")
			}

			if diff := cmp.Diff(tc.expectedLocation, location); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
