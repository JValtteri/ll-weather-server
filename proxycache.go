package main

import (
    "errors"
    "fmt"
    //"encoding/json"
)

//var weatheCache map[string]InWeatherRange
//var cityCache map[string]InCity

func getCityCoord(city string) (float32, float32) {
    //61.499, 23.787
    lat, lon := getCity(city)
    return lat, lon
}

func getSummaryWeather(city string, mode int) (WeekWeather, error) { // JSON
    // MODE: 0=summary, 1=detail
    // Convert name to coord
    lat, lon := getCityCoord(city)
    if lat==0 && lon==0 {
        var emptyWeather WeekWeather
        err := errors.New("City not found")
        fmt.Println(err)
        return emptyWeather, err
    }
    // Check cache
    // Make request
    var r_obj InWeatherRange = getWeather(lat, lon)
    // Format response
    var f_obj WeekWeather = mapDays(r_obj)
    // Cache response
    // Return responce
    return f_obj, nil
}
