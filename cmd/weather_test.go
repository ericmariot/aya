/*
Copyright © 2024 github.com/ericmariot <ericmariots@gmail.com>
*/
package cmd

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	Response *http.Response
	Err      error
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

type MockConfigLoader struct {
	Config  Config
	LoadErr error
	SaveErr error
}

func (m *MockConfigLoader) LoadConfig() (Config, error) {
	return m.Config, m.LoadErr
}

func (m *MockConfigLoader) SaveConfig(config Config) error {
	m.Config = config
	return m.SaveErr
}

type MockGeoLocator struct {
	Lat    string
	Lon    string
	Name   string
	GeoErr error
}

func (m *MockGeoLocator) CityToGeoLoc(city string) (string, string, string, error) {
	return m.Lat, m.Lon, m.Name, m.GeoErr
}

func TestGetCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		config      Config
		loadErr     error
		geoLat      string
		geoLon      string
		geoName     string
		geoErr      error
		expectedLat string
		expectedLon string
		expectedErr error
	}{
		{
			name: "City in config",
			city: "testCity",
			config: Config{
				CityCoordinates: map[string]CityCoordinates{
					"testCity": {Latitude: "12.34", Longitude: "56.78", Name: "testCity"},
				},
			},
			expectedLat: "12.34",
			expectedLon: "56.78",
		},
		{
			name:        "City not in config, geo lookup success",
			city:        "newCity",
			geoLat:      "98.76",
			geoLon:      "54.32",
			geoName:     "newCity",
			config:      Config{CityCoordinates: make(map[string]CityCoordinates)},
			expectedLat: "98.76",
			expectedLon: "54.32",
		},
		{
			name:        "City not in config, geo lookup error",
			city:        "errorCity",
			geoErr:      errors.New("geo lookup failed"),
			config:      Config{CityCoordinates: make(map[string]CityCoordinates)},
			expectedErr: errors.New("geo lookup failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfigLoader := &MockConfigLoader{
				Config:  tt.config,
				LoadErr: tt.loadErr,
			}

			mockGeoLocator := &MockGeoLocator{
				Lat:    tt.geoLat,
				Lon:    tt.geoLon,
				Name:   tt.geoName,
				GeoErr: tt.geoErr,
			}

			lat, lon, err := getCoordinates(tt.city, mockConfigLoader, mockGeoLocator)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLat, lat)
				assert.Equal(t, tt.expectedLon, lon)
			}
		})
	}
}

func TestParseToTime(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{"2006-01-02T15:04", time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC), false},
		{"2024-05-16T12:30", time.Date(2024, 5, 16, 12, 30, 0, 0, time.UTC), false},
		{"invalid-time", time.Time{}, true},
		{"2024-05-16", time.Time{}, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := parseToTime(test.input)
			t.Log("Result:", result)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
