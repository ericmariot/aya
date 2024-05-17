/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
)

var weatherCmd = &cobra.Command{
	Use:   "weather [city]",
	Short: "Get the weather for a specified city",
	Long: `Get the weather for a specified city. For example:

	aya weather
	aya weather criciuma
	aya weather san-diego
	aya weather london --graph
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var city string
		if len(args) == 0 {
			config, err := loadConfig()
			if err != nil {
				fmt.Println("error loading config:", err)
				return
			}
			city = config.City
		} else {
			city = normalizeCityName(args[0])
		}
		graph, err := cmd.Flags().GetBool("graph")
		if err != nil {
			fmt.Println("error: ", err)
		}

		currentWeather(city, graph)
	},
}

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.Flags().BoolP("graph", "g", false, "plot a graph of the weather forecast of the next 24hours")
}

type CityCoordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	Name      string `json:"name"`
}

type Weather struct {
	Timezone string `json:"timezone"`
	Current  struct {
		Time        string  `json:"time"`
		Temperature float64 `json:"temperature_2m"`
		Rain        float64 `json:"precipitation"`
		IsDay       int64   `json:"is_day"`
	} `json:"current"`
	Hourly struct {
		Time        []string  `json:"time"`
		Temperature []float64 `json:"temperature_2m"`
		RainChance  []int64   `json:"precipitation_probability"`
		CloudCover  []int64   `json:"cloud_cover"`
	} `json:"hourly"`
}

type IPInfo struct {
	City string `json:"city"`
}

func getWeather(city string) (Weather, error) {
	fmt.Println("🌎 Getting coordinates for", strings.ToUpper(string(city[0]))+city[1:])
	lat, lon, err := getCoordinates(city)
	if err != nil {
		return Weather{}, fmt.Errorf("error getting coordinates: %v", err)
	}

	fmt.Println("🌤️  Getting weather")
	res, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=" + lat + "&longitude=" + lon + "&current=temperature_2m,precipitation,is_day&hourly=temperature_2m,precipitation_probability,cloud_cover&timezone=auto&forecast_days=2")
	if err != nil {
		return Weather{}, fmt.Errorf("error getting weather: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("error getting coordinates: %v", err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return Weather{}, fmt.Errorf("error converting json: %v", err)
	}

	return weather, nil
}

func getCoordinates(city string) (string, string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", "", err
	}

	coords, ok := config.CityCoordinates[city]
	if ok {
		return coords.Latitude, coords.Longitude, nil
	}

	lat, lon, name, err := cityToGeoLoc(city)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	config.CityCoordinates[city] = CityCoordinates{Latitude: lat, Longitude: lon, Name: name}
	saveConfig(config)

	return lat, lon, nil
}

func cityToGeoLoc(city string) (string, string, string, error) {
	res, err := http.Get("https://nominatim.openstreetmap.org/search?city=" + city + "&format=json&limit=1")
	if err != nil {
		return "", "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", "", err
	}

	var cityCoord []CityCoordinates
	err = json.Unmarshal(body, &cityCoord)
	if err != nil {
		return "", "", "", err
	}
	lat, lon, name := cityCoord[0].Latitude, cityCoord[0].Longitude, cityCoord[0].Name

	return lat, lon, name, nil

}

func parseToTime(timeString string) time.Time {
	layout := "2006-01-02T15:04"

	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println("error parsing time:", err)
	}

	return t
}

func currentWeather(city string, graph bool) {
	weather, err := getWeather(city)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(weather.Timezone, "TZ")
	fmt.Println("Last update:", parseToTime(weather.Current.Time).Format("15:04"))
	fmt.Printf("Current: %.1f°C\n\n", weather.Current.Temperature)
	if graph {
		fmt.Printf("\n")
		plotGraph(weather)
	}
}

func plotGraph(weather Weather) {
	hourNow := parseToTime(weather.Current.Time).Hour()
	data := []float64{}
	labels := []string{}

	for i := hourNow; i < hourNow+25; i++ {
		idx := i % len(weather.Hourly.Time)
		timeStr := parseToTime(weather.Hourly.Time[idx]).Format("15:04")
		temperature := weather.Hourly.Temperature[idx]
		data = append(data, float64(temperature))
		labels = append(labels, timeStr)
	}

	graph := asciigraph.Plot(data, asciigraph.Width(110), asciigraph.Height(10), asciigraph.Caption("Temperature Over Time"))

	var result strings.Builder
	result.WriteString(graph)
	result.WriteString("\n")
	labelWidth := 100 / len(labels)
	initialOffset := 5
	result.WriteString(strings.Repeat(" ", initialOffset))

	for i, label := range labels {
		if i%2 == 0 {
			result.WriteString(fmt.Sprintf("%-*s", labelWidth, label))
		} else {
			result.WriteString(strings.Repeat(" ", labelWidth))
		}
	}
	result.WriteString("\n")

	fmt.Println(result.String())
}
