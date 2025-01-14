package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Location struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type Locations []Location

type DirectGeocoding struct {
	Key string // OpenWeather API Key
}

// CoordsByLocationName returns the coordinates of a named location.
// see https://openweathermap.org/api/geocoding-api#direct_name
func (g *DirectGeocoding) CoordsByLocationName(name string) (*Locations, error) {
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

	fmt.Println(uri.String())
	result, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}

	resultBody, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	locations := &Locations{}
	err = json.Unmarshal(resultBody, locations)
	if err != nil {
		return locations, err
	}

	return locations, nil

}
