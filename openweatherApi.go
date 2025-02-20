package main

import (
    //"net/http"
    //"encoding/json"
    "fmt"
    //"io"
    "os"
    "log"
)

const forecastUrl string = "https://api.openweathermap.org/data/2.5/forecast?lat={LAT}&lon={LON}&units={UNITS}&appid={API_KEY}"

const KEYFILE string = "key.txt"
var API_KEY string = "no-key"
var UNITS string = "metric"

func getWeather(lat float32, lon float32) {
    fmt.Println("Get weather at %s %s", lat, lon)
}

func updateKey() {
    content, err := os.ReadFile(KEYFILE)
    if err != nil {
        log.Fatal(err)
    }
    API_KEY = string(content)
}


