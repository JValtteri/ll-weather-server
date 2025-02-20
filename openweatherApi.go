package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "strings"
    "io"
    "os"
    "log"
)

const FORECAST_URL string = "https://api.openweathermap.org/data/2.5/forecast?lat={LAT}&lon={LON}&units={UNITS}&appid={API_KEY}"

const KEYFILE string = "key.txt"
var API_KEY string = "no-key"
var UNITS string = "metric"

func getWeather(lat float32, lon float32) InWeatherRange {
    var weather_obj InWeatherRange

    if API_KEY == "no-key" {
        updateKey()
    }
    fmt.Println("Get weather at", lat, lon)
    requestURL := makeURL(lat, lon, UNITS)
    fmt.Println(requestURL)
    raw_weather := makeRequest(requestURL)
    err := json.Unmarshal(raw_weather, &weather_obj)
    if err != nil {
        fmt.Println(err)
    }
    return weather_obj
}

func updateKey() {
    content, err := os.ReadFile(KEYFILE)
    if err != nil {
        log.Fatal(err)
    }
    API_KEY = string(content)
}

func makeURL(lat, lon float32, units string) string {
    url := ""
    url = strings.Replace(FORECAST_URL, "{LAT}", str_f(lat), 1)
    url = strings.Replace(url, "{LON}", str_f(lon), 1)
    url = strings.Replace(url, "{UNITS}", units, 1)
    url = strings.Replace(url, "{API_KEY}", API_KEY, 1)
    return url
}

// Make a request to chosen address
func makeRequest(address string) []byte {
    resp, err := http.Get(address)
    if err != nil {
        fmt.Sprintf("Welp! GET from %s failed", address)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        // handle error
    }
    return body
}


// UTILS

func str_f(f float32) string {
    s := fmt.Sprintf("%f", f)
    return s
}

func unloadJSON(object InWeatherRange) string {
    body, err := json.Marshal(object)
    if err != nil {
        fmt.Println(err)
    }
    return string(body)
}
