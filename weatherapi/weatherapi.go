package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		FeelsLikeC float64 `json:"feelslike_c"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				FeelsLikeC   float64 `json:"feelslike_c"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func getWeather() {
	city := "Criciuma"

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + key + "&q=" + city + "&days=1&aqi=no&alerts=no")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Weather API not available.")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%v - %v, %v \n%v \nActual: %v°C Sensation: %v°C\n\n", location.Country, location.Name, location.Region, current.Condition.Text, current.TempC, current.FeelsLikeC)

	for _, hour := range hours {
		date := time.Unix((hour.TimeEpoch), 0)

		if date.Before(time.Now().Add(-1 * time.Hour)) {
			continue
		}

		fmt.Printf("%v - %v°C %v %v\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
	}
}
