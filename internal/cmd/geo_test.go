//go:build e2e

package cmd_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/squeedee/geo/cmd"
	internalcmd "github.com/squeedee/geo/internal/cmd"
	"os"
	"testing"
)

var ApiKey = os.Getenv(cmd.ApiKeyName)

func TestLocationByName(t *testing.T) {
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
		expectedLocation []internalcmd.NameResult
	}{
		"Richmond, well formed without country code": {
			fields: defaultFields,
			name:   "Richmond, VA",
			expectedLocation: []internalcmd.NameResult{
				{
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
			expectedLocation: []internalcmd.NameResult{
				{
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
			expectedLocation: []internalcmd.NameResult{
				{
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
			location, err := tc.fields.LocationByName(tc.name)
			if tc.expectedError == "" && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			} else if tc.expectedError != "" && err == nil {
				t.Fatalf("Expected error, but didn't get one")
			}

			if diff := cmp.Diff(tc.expectedLocation, location); diff != "" {
				t.Errorf("LocationByName() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLocationByZip(t *testing.T) {
	if ApiKey == "" {
		t.Fatalf("No %s set", cmd.ApiKeyName)
	}

	defaultFields := internalcmd.DirectGeocoding{
		Key: ApiKey,
	}

	tests := map[string]struct {
		zip              string
		fields           internalcmd.DirectGeocoding
		expectedError    string
		expectedLocation *internalcmd.ZipResult
	}{
		"23228 well formed without country": {
			zip:    "23228",
			fields: defaultFields,
			expectedLocation: &internalcmd.ZipResult{
				Zip:     "23228",
				Name:    "Henrico County",
				Lat:     37.4638,
				Lon:     -77.398,
				Country: "US",
			},
		},
		"10001 well formed without country": { // variation test
			zip:    "10001",
			fields: defaultFields,
			expectedLocation: &internalcmd.ZipResult{
				Zip:     "10001",
				Name:    "New York",
				Lat:     40.7484,
				Lon:     -73.9967,
				Country: "US",
			},
		},
		"99999 is an error": {
			fields:        defaultFields,
			zip:           "99999",
			expectedError: "zip '99999' not found",
		},
		"99998 is an error": { // variation check
			fields:        defaultFields,
			zip:           "99998",
			expectedError: "zip '99998' not found",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			location, err := tc.fields.LocationByZip(tc.zip)
			if tc.expectedError == "" && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			} else if tc.expectedError != "" && err == nil {
				t.Fatalf("Expected error, but didn't get one")
			}

			if diff := cmp.Diff(tc.expectedLocation, location); diff != "" {
				t.Errorf("LocationByZip() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
