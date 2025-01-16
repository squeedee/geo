# geo

A simple CLI for fetching a place's geolocation data.

# Building

Install dependencies:

* [Golang 1.23.1 or higher.](https://go.dev/dl/) 
* [make](https://www.gnu.org/software/make/manual/make.html)  
* [curl](https://curl.se/)

Run:

```shell
make build
```

# Usage

Set the API Key

```shell
export OPEN_WEATHER_API_KEY=<api-key>
```

Get help from 'geo' (ensure you built the binary):

```
build/geo -h 
```

Run the example:
```shell
build/geo "Henrico, VA" 10001 "Seattle, WA"
```

# Testing

## Unit tests

Unit tests are in any file named `*_test.go` that also does not have
the built tag: `//go:build e2e` at the top of the file.

Run the unit tests with:

```shell
make test
```

## Integration tests
Integration tests are in the [./test/integration/](./test/integration/) directory.


Ensure you have set the API Key:

```shell
export OPEN_WEATHER_API_KEY=<api-key>
```

Run the integration tests with:

```shell
make e2e
```

## All tests
Run all the tests with:

```shell
make test-all
```

# Notes:

## Unit test that hits the API
There's a "unit" test in [./internal/cmd/geo_test.go](./internal/cmd/geo_test.go)
that hit's the OpenWeather API. I took a shortcut in these unit tests
and did not mock the service, as I assume you really want to see
the integration tests.


## Stubbing with VCR
I have a trick (which I did not use) with unit tests in Go (and a few other languages like Ruby) where I use
a [vcr](https://github.com/dnaeon/go-vcr) library to avoid building stubs. Further, 'vcr' 
libraries enable us to switch between fast local iteration and slow, complex e2e.

There are definitely caveats with this style of testing, and it doesn't fit every component,
though it would work exceptionally well with this example.

This is something I would normally introduce, but it really only makes
sense with developer driven acceptance testing, and I didn't want to
push my agenda on that in this example :D
