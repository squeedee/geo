package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ZipResult struct {
	Zip     string  `json:"zip"`
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

type NameResult struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type DirectGeocoding struct {
	Key string // OpenWeather API Key
}

// LocationByName returns the coordinates of a named location.
// see https://openweathermap.org/api/geocoding-api#direct_name
func (g *DirectGeocoding) LocationByName(name string) ([]NameResult, int, error) {
	uri := g.buildNameLookupUri(name)

	result, err := http.Get(uri)

	var statusCode = 0
	if result != nil {
		statusCode = result.StatusCode
	}

	if err != nil || result == nil || statusCode < 200 || statusCode >= 300 {
		return nil, statusCode, err
	}

	resultBody, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, statusCode, err
	}

	var locations []NameResult
	err = json.Unmarshal(resultBody, &locations)
	if err != nil {
		return nil, statusCode, err
	}

	return locations, statusCode, nil
}

func (g *DirectGeocoding) LocationByZip(zip string) (*ZipResult, int, error) {
	uri := g.buildZipLookupUri(zip)

	result, err := http.Get(uri)
	var statusCode = 0
	if result != nil {
		statusCode = result.StatusCode
	}

	if err != nil || result == nil {
		return nil, statusCode, err
	}

	resultBody, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, statusCode, err
	}

	location := &ZipResult{}
	err = json.Unmarshal(resultBody, location)
	if err != nil {
		return nil, statusCode, err
	}

	if location.Zip == "" {
		return nil, statusCode, fmt.Errorf("zip '%s' not found", zip)
	}

	return location, statusCode, nil
}

func (g *DirectGeocoding) buildNameLookupUri(name string) string {
	uri := url.URL{
		Scheme: "http",
		Host:   "api.openweathermap.org",
		Path:   "geo/1.0/direct",
	}

	// Default USA
	if len(strings.Split(name, ",")) < 3 {
		name = fmt.Sprintf("%s, USA", name)
	}

	q := uri.Query()
	q.Add("q", name)
	q.Add("appid", g.Key)

	uri.RawQuery = q.Encode()
	return uri.String()
}

func (g *DirectGeocoding) buildZipLookupUri(zip string) string {
	uri := url.URL{
		Scheme: "http",
		Host:   "api.openweathermap.org",
		Path:   "geo/1.0/zip",
	}

	q := uri.Query()
	q.Add("zip", zip)
	q.Add("appid", g.Key)

	uri.RawQuery = q.Encode()
	return uri.String()
}
