package main

import (
    "errors"
    "fmt"
    //"encoding/json"
)

var weatheCache map[string]InWeatherRange
var cityCache map[string]InCity

func getCityCoord(city string) (float32, float32) {
    //61.499, 23.787
    lat, lon := getCity(city)
    return lat, lon
}

func getWeather(city string, mode int) (InWeatherRange, error) { // JSON
    // MODE: 0=summary, 1=detail
    // Convert name to coord
    lat, lon := getCityCoord(city)
    if lat==0 && lon==0 {
        var emptyWeather InWeatherRange
        err := errors.New("City not found")
        fmt.Println(err)
        return emptyWeather, err
    }
    // Check cache
    // Make request
    var r_obj InWeatherRange = getWeather(lat, lon)
    // Format response

    // Cache response
    // Return responce
    return r_obj, nil
}
