package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Config struct {
	City            string                     `json:"city"`
	LastIP          string                     `json:"last_ip"`
	CityCoordinates map[string]CityCoordinates `json:"city_coordinates"`
}

var configFilePath string

func initConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error finding home directory:", err)
		os.Exit(1)
	}
	configFilePath = filepath.Join(homeDir, ".aya.json")

	currentIP, err := getIPAddress()
	if err != nil {
		fmt.Println("error getting IP address:", err)
		return
	}

	var config Config
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Config file does not exist. Creating new one...")
		city, err := getCityFromIP(currentIP)
		if err != nil {
			fmt.Println("error getting city from IP:", err)
			return
		}

		config = Config{
			City:            normalizeCityName(city),
			LastIP:          currentIP,
			CityCoordinates: make(map[string]CityCoordinates),
		}
		saveConfig(config)
	} else {
		config, err = loadConfig()
		if err != nil {
			fmt.Println("error loading config:", err)
			return
		}

		if config.LastIP != currentIP {
			fmt.Println("IP address has changed. Updating city...")
			city, err := getCityFromIP(currentIP)
			if err != nil {
				fmt.Println("error getting city from IP:", err)
				return
			}
			config.City = normalizeCityName(city)
			config.LastIP = currentIP
			saveConfig(config)
		}
	}
}

func loadConfig() (Config, error) {
	var config Config
	file, err := os.Open(configFilePath)
	if err != nil {
		return config, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&config)

	if config.CityCoordinates == nil {
		config.CityCoordinates = make(map[string]CityCoordinates)
	}

	return config, err
}

func saveConfig(config Config) {
	file, err := os.Create(configFilePath)
	if err != nil {
		fmt.Println("error creating config file:", err)
		return
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		fmt.Println("error encoding config:", err)
	}
}

func getIPAddress() (string, error) {
	resp, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
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

	var ipInfo struct {
		City string `json:"city"`
	}

	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}

	return ipInfo.City, nil
}

func normalizeCityName(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func formatCityName(city string) string {
	city = strings.ReplaceAll(city, "-", " ")
	return cases.Title(language.English, cases.NoLower).String(city)
}
