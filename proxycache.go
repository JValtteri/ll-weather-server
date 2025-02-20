package main

import (
    //"fmt"
    //"encoding/json"
)

func getCityCoord(city string) (float32, float32) {
    return 61.499, 23.787
}

func getSummaryWeather(city string) InWeatherRange { // JSON
    // Convert name to coord
    lat, lon := getCityCoord(city)
    // Check cache
    // Make request
    r_obj := getWeather(lat, lon)
    // Format response
    // Cache response
    // Return responce
    return r_obj
}

func getDetailWeather(city string) InWeatherRange { // JSON
    // Convert name to coord
    lat, lon := getCityCoord(city)
    // Check cache
    // Make request
    r_obj := getWeather(lat, lon)
    // Format response
    // Cache response
    // Return responce
    return r_obj
}
