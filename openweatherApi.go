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
const CITY_URL string = "https://api.openweathermap.org/geo/1.0/direct?q={CITY_NAME},{COUNTRY_CODE}&limit={LIMIT}&appid={API_KEY}"

const KEYFILE string = "key.txt"
var API_KEY string = "no-key"
var UNITS string = "metric"
var COUNTRY string = "FI"
var CITY_LIMIT string = "1"

func getWeather(lat float32, lon float32) InWeatherRange {
    var weather_obj InWeatherRange

    if API_KEY == "no-key" {
        updateKey()
    }
    fmt.Println("Get weather at", lat, lon)
    requestURL := makeWeatherURL(lat, lon, UNITS)
    fmt.Println(requestURL)
    var raw_weather []byte = makeRequest(requestURL)
    err := json.Unmarshal(raw_weather, &weather_obj)
    if err != nil {
        fmt.Println(err)
    }
    return weather_obj
}

func getCity(name string) (float32, float32) {
    var city_obj InCity
    if API_KEY == "no-key" {
        updateKey()
    }
    fmt.Println("Get city:", name)
    requestURL := makeCityURL(name, COUNTRY, CITY_LIMIT)
    fmt.Println(requestURL)

    raw_city := makeRequest(requestURL)
    err := json.Unmarshal(raw_city, &city_obj)
    if err != nil {
        fmt.Println("JSON Marshal error:", err)
    }
    if len(city_obj) == 0 {
        return 0.0, 0.0
    }
    var lat float32 = city_obj[0].Lat
    var lon float32 = city_obj[0].Lon
    fmt.Println("City:", city_obj[0].Name, lat, lon)
    return lat, lon
}

func updateKey() {
    content, err := os.ReadFile(KEYFILE)
    if err != nil {
        log.Fatal(err)
    }
    API_KEY = string(content)
}

func makeWeatherURL(lat, lon float32, units string) string {
    url := ""
    url = strings.Replace(FORECAST_URL, "{LAT}", str_f(lat), 1)
    url = strings.Replace(url, "{LON}", str_f(lon), 1)
    url = strings.Replace(url, "{UNITS}", units, 1)
    url = strings.Replace(url, "{API_KEY}", API_KEY, 1)
    return url
}

func makeCityURL(name, country, limit string) string {
    url := ""
    url = strings.Replace(CITY_URL, "{CITY_NAME}", name, 1)
    url = strings.Replace(url, "{COUNTRY_CODE}", country, 1)
    url = strings.Replace(url, "{LIMIT}", limit, 1)
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
        fmt.Sprintf("Error reading response body:", err)
    }
    return body
}


// UTILS

func str_f(f float32) string {
    s := fmt.Sprintf("%f", f)
    return s
}

func unloadJSON(object any) string {
    body, err := json.Marshal(object)
    if err != nil {
        fmt.Println(err)
    }
    return string(body)
}
