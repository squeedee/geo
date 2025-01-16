//go:build e2e

package cmd_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/squeedee/geo/cmd"
	internalcmd "github.com/squeedee/geo/internal/cmd"
	"net/http"
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
		expectStatusCode int
	}{
		"Richmond, well formed without country code is successfully found": {
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
			expectStatusCode: http.StatusOK,
		},
		"Richmond, well formed with country code is successfully found": {
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
			expectStatusCode: http.StatusOK,
		},
		"Unknown place is an empty location list": { // Note: The OpenWeather API is really inconsistent
			fields:           defaultFields,
			name:             "Fallafelville",
			expectedLocation: []internalcmd.NameResult{},
			expectStatusCode: http.StatusOK,
		},
		"Known place with bad api key is unauthorized": {
			fields: internalcmd.DirectGeocoding{
				Key: "invalid-key",
			},
			name:             "Richmond, VA, USA",
			expectStatusCode: http.StatusUnauthorized,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			location, code, err := tc.fields.LocationByName(tc.name)
			if tc.expectedError == "" && err != nil {
				t.Fatalf("LocationByName() Unexpected error: %v", err)
			} else if tc.expectedError != "" && err == nil {
				t.Fatalf("LocationByName() Expected error, but didn't get one")
			} else if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("LocationByName() Unexpected error: %s, expected %s", err, tc.expectedError)
			}

			if tc.expectStatusCode != code {
				t.Fatalf("LocationByName() Unexpected status code: %d, expected %d", code, tc.expectStatusCode)
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
		expectErrorKind  string
		expectStatusCode int
	}{
		"23228 is valid": {
			zip:    "23228",
			fields: defaultFields,
			expectedLocation: &internalcmd.ZipResult{
				Zip:     "23228",
				Name:    "Henrico County",
				Lat:     37.4638,
				Lon:     -77.398,
				Country: "US",
			},
			expectStatusCode: http.StatusOK,
		},
		"99999 is a 404 error": {
			fields:           defaultFields,
			zip:              "99999",
			expectedError:    "zip '99999' not found",
			expectStatusCode: http.StatusNotFound,
		},
		"99998 is a 404 error": { // variation check
			fields:           defaultFields,
			zip:              "99998",
			expectedError:    "zip '99998' not found",
			expectStatusCode: http.StatusNotFound,
		},
		"invalid api key provides unauthorized code": { // variation check
			fields: internalcmd.DirectGeocoding{
				Key: "invalid-key",
			},
			zip:              "23228",
			expectStatusCode: http.StatusUnauthorized,
			expectedError:    "zip '23228' not found",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			location, code, err := tc.fields.LocationByZip(tc.zip)
			if tc.expectedError == "" && err != nil {
				t.Fatalf("LocationByZip() Unexpected error: %v", err)
			} else if tc.expectedError != "" && err == nil {
				t.Fatalf("LocationByZip() Expected error, but didn't get one")
			} else if err != nil && tc.expectedError != err.Error() {
				t.Fatalf("LocationByZip() Unexpected error: %s, expected %s", err, tc.expectedError)
			}

			if tc.expectStatusCode != code {
				t.Fatalf("LocationByZip() Unexpected status code: %d, expected %d", code, tc.expectStatusCode)
			}

			if diff := cmp.Diff(tc.expectedLocation, location); diff != "" {
				t.Errorf("LocationByZip() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
