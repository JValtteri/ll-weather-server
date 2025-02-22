package main

import (
    "errors"
    "fmt"
    //"encoding/json"
)

type Coords struct {
    lat float32
    lon float32
}

var weatheCache map[string]InWeatherRange = make(map[string]InWeatherRange)
var cityCache map[string]Coords = make(map[string]Coords)

func getCityCoord(city string) (float32, float32) {
    //61.499, 23.787
    // Check cache
    lat, lon, ok := searchCacheCity(city)
    if ok {
        fmt.Println("Found cached:", city)
        return lat, lon
    }
    // Make request
    lat, lon = getCity(city)
    // Cache coords
    addCacheCity(city, lat, lon)
    return lat, lon
}

func getProxyWeather(city string, mode int) (InWeatherRange, error) {
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

func addCacheCity(key string, lat, lon float32) {
    var coords Coords
    coords.lat = lat
    coords.lon = lon
    cityCache[key] = coords
}

func searchCacheCity(key string) (float32, float32, bool) {
    coords, ok := cityCache[key]
    if !ok {
        return 0, 0, false
    }
    return coords.lat, coords.lon, ok
}

/*
func searchCacheCity(key string) (float32, float32, bool) {
    city_obj, ok := cityCache[key]
    if !ok {
        return 0, 0, false
    }
    var lat, lon float32 = getLatLon(city_obj)
    return lat, lon, ok
}
*/
