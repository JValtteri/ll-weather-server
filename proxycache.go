package main

import (
    "errors"
    "fmt"
    "time"
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
        fmt.Println("Found cached name:", city)
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
    var r_obj InWeatherRange
    var ok bool
    // Check cache
    r_obj, ok = searchCacheWeather(city)
    if ok {
        fmt.Println("Found cached data:", city)
        return r_obj, nil
    }
    // Make request
    r_obj = getWeather(lat, lon)
    // Cache response
    r_obj.timestamp = uint(time.Now().Unix())
    addCacheWeather(city, r_obj)
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

func addCacheWeather(key string, r_obj InWeatherRange) {
    weatheCache[key] = r_obj
}

func searchCacheWeather(key string) (InWeatherRange, bool) {
    r_obj, ok := weatheCache[key]
    if !ok {
        var emptyResponse InWeatherRange
        return emptyResponse, false
    }
    // check age
    if (uint(time.Now().Unix()) - r_obj.timestamp) > (SECONDS_IN_HOUR*12) {
        fmt.Println("Purge old weather data")
        delete(weatheCache, key)
        var emptyResponse InWeatherRange
        return emptyResponse, false
    }
    return r_obj, ok
}
