/*
Copyright © 2024 ericmariot <ericmariots@gmail.com>
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
)

type CityCoordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
	Name      string `json:"name"`
}

type Weather struct {
	Timezone string `json:"timezone"`
	Current  struct {
		Time        string  `json:"time"`
		Temperature float32 `json:"temperature_2m"`
		Rain        float32 `json:"precipitation"`
		IsDay       int8    `json:"is_day"`
	} `json:"current"`
	Hourly struct {
		Time        []string  `json:"time"`
		Temperature []float32 `json:"temperature_2m"`
		RainChance  []int8    `json:"precipitation_probability"`
		CloudCover  []int8    `json:"cloud_cover"`
	} `json:"hourly"`
}

func getWeather(city string) {
	fmt.Println("🌎 Getting coordinates for", strings.ToUpper(string(city[0]))+city[1:])
	lat, lon, _, err := cityToGeoLoc(city)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("🌤️  Getting weather")
	res, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=" + lat + "&longitude=" + lon + "&current=temperature_2m,precipitation,is_day&hourly=temperature_2m,precipitation_probability,cloud_cover&timezone=auto&forecast_days=2")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	fmt.Println(weather.Timezone, "TZ")
	fmt.Println("Last update:", parseToTime(weather.Current.Time).Format("15:04"))
	fmt.Printf("%.1f°C", weather.Current.Temperature)
	fmt.Println("")
	fmt.Println("")

	hourNow := parseToTime(weather.Current.Time).Hour()
	data := []float64{}
	for i := hourNow; i < hourNow+24; i++ {
		timeStr := parseToTime(weather.Hourly.Time[i]).Format("15:04")
		temperature := weather.Hourly.Temperature[i]
		data = append(data, float64(temperature))
		fmt.Printf("%s %.1f°C \n", timeStr, temperature)
	}

	graph := asciigraph.Plot(data)
	fmt.Println(graph)
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
func main() {
	getWeather("criciuma")
}

func parseToTime(timeString string) time.Time {
	layout := "2006-01-02T15:04"

	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}

	return t
}
