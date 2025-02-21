package main

import (
    //"fmt"
    //"encoding/json"
)

func getCityCoord(city string) (float32, float32) {
    return 61.499, 23.787
}

func getSummaryWeather(city string) WeekWeather { // JSON
    // Convert name to coord
    lat, lon := getCityCoord(city)
    // Check cache
    // Make request
    var r_obj InWeatherRange = getWeather(lat, lon)
    // Format response
    var f_obj WeekWeather = mapDays(r_obj)
    // Cache response
    // Return responce
    return f_obj
}

func getDetailWeather(city string) InWeatherRange { // JSON
    // Convert name to coord
    lat, lon := getCityCoord(city)
    // Check cache
    // Make request
    var r_obj InWeatherRange = getWeather(lat, lon)
    // Format response
    // Cache response
    // Return responce
    return r_obj
}
