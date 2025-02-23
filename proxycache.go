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
var cityCache map[string]Coords           = make(map[string]Coords)
var iconCache map[string][]byte           = make(map[string][]byte)

func getCityCoord(city string) (float32, float32) {
    //61.499, 23.787
    lat, lon, ok := searchCacheCity(city)       // Check cache
    if ok {
        fmt.Println("Found cached name:", city)
        return lat, lon
    }
    lat, lon = getCity(city)                    // Make request
    addCacheCity(city, lat, lon)                // Cache coords
    return lat, lon
}

func getProxyWeather(city string, mode int) (InWeatherRange, error) {
    // MODE: 0=summary, 1=detail
    lat, lon := getCityCoord(city)              // Convert name to coord
    if lat==0 && lon==0 {
        var emptyWeather InWeatherRange
        err := errors.New("City not found")
        fmt.Println(err)
        return emptyWeather, err
    }
    var r_obj InWeatherRange
    var ok bool
    r_obj, ok = searchCacheWeather(city)        // Check cache
    if ok {
        fmt.Println("Found cached data:", city)
        return r_obj, nil
    }
    r_obj = getWeather(lat, lon)                // Make request
    r_obj.timestamp = uint(time.Now().Unix())
    addCacheWeather(city, r_obj)                // Cache response
    return r_obj, nil
}

func getProxyIcon(id string) []byte {
    var icon []byte
    var ok bool
    icon, ok = searchCacheIcon(id)
    if ok {
        return icon
    }
    icon = getIcon(id)                          // Make request
    addCacheIcon(id, icon)                      // Cache response
    return icon
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

func addCacheIcon(key string, data []byte) {
    iconCache[key] = data
}

func searchCacheIcon(key string) ([]byte, bool) {
    data, ok := iconCache[key]
    if !ok {
        var no_data []byte
        return no_data, false
    }
    return data, ok
}
