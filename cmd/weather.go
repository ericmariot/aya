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
	"unicode"

	"github.com/guptarohit/asciigraph"
	"github.com/spf13/cobra"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather [city]",
	Short: "Get the weather for a specified city",
	Long: `Get the weather for a specified city. For example:

	aya weather criciuma
	aya weather san-diego
	aya weather london --graph
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			city, err := getCity()
			if err != nil {
				fmt.Println("Error: ", err)
			}

			graph, err := cmd.Flags().GetBool("graph")
			if err != nil {
				fmt.Println("Error: ", err)
			}

			normalized, err := normalizeCity(city)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			currentWeather(normalized, graph)
		} else {
			city := args[0]
			graph, err := cmd.Flags().GetBool("graph")
			if err != nil {
				fmt.Println("Error: ", err)
			}

			normalized, err := normalizeCity(city)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			currentWeather(normalized, graph)
		}

	},
}

func normalizeCity(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	return result, err
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
	lat, lon, _, err := cityToGeoLoc(city)
	if err != nil {
		return Weather{}, fmt.Errorf("error getting coordinates: %v", err)
	}

	fmt.Println("🌤️  Getting weather")
	res, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=" + lat + "&longitude=" + lon + "&current=temperature_2m,precipitation,is_day&hourly=temperature_2m,precipitation_probability,cloud_cover&timezone=auto&forecast_days=2")
	if err != nil {
		return Weather{}, fmt.Errorf("error getting coordinates: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("error getting coordinates: %v", err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return Weather{}, fmt.Errorf("error getting coordinates: %v", err)
	}

	return weather, nil
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
		fmt.Println("Error parsing time:", err)
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

func getIPAddress() (string, error) {
	resp, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["origin"], nil
}

func getCityFromIP(ip string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s/json", ip))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ipInfo IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return "", err
	}

	return ipInfo.City, nil
}

func getCity() (string, error) {
	ip, err := getIPAddress()
	if err != nil {
		fmt.Println("Error getting IP address:", err)
		return "", err
	}

	city, err := getCityFromIP(ip)
	if err != nil {
		fmt.Println("Error getting city from IP:", err)
		return "", err
	}

	return city, nil
}
