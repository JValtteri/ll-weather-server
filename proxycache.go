package main

import (
    "errors"
    "fmt"
    "log"
    "time"
    "github.com/JValtteri/weather/owm"
)

type Coords struct {
    lat float32
    lon float32
}

var weatherCache map[string]owm.WeatherRange = make(map[string]owm.WeatherRange)
var cityCache map[string]Coords           = make(map[string]Coords)
var iconCache map[string][]byte           = make(map[string][]byte)

func GetCityCoord(city string) (float32, float32) {
    var lat, lon float32
    var ok bool
    lat, lon, ok = searchCacheCity(city)        // Check cache
    if ok {
        return lat, lon
    }
    lat, lon = owm.Coord(city)                  // Make request
    addCacheCity(city, lat, lon)                // Cache coords
    return lat, lon
}

func GetProxyWeather(city string) (owm.WeatherRange, error) {
    var lat, lon float32 = GetCityCoord(city)   // Convert name to coord
    if lat==0 && lon==0 {
        var emptyWeather owm.WeatherRange
        err := errors.New(fmt.Sprintf("r:%4vu:%4v: City [%s] not found", rqNum, uniqRqNum, city))
        log.Println(err)
        return emptyWeather, err
    }
    var r_obj owm.WeatherRange
    var ok bool
    r_obj, ok = searchCacheWeather(city)        // Check cache
    if ok {
        log.Printf("r:%6vu:%4v: Get weather: %s at %.3f %.3f\n", rqNum, uniqRqNum, city, lat, lon)
        return r_obj, nil
    }
    uniqRqNum++
    log.Printf("r:%6vu:%4v: Get weather: %s at %.3f %.3f (New request)\n", rqNum, uniqRqNum, city, lat, lon)
    r_obj = owm.Forecast(lat, lon)              // Make request
    r_obj.Timestamp = uint(time.Now().Unix())
    addCacheWeather(city, r_obj)                // Cache response
    return r_obj, nil
}

func GetProxyIcon(id string) []byte {
    var icon []byte
    var ok bool
    icon, ok = searchCacheIcon(id)
    if ok {
        return icon
    }
    rqNum++
    uniqRqNum++
    log.Printf("r:%6vu:%4v: Get new icon: %v (New request)\n", rqNum, uniqRqNum, id)
    icon = owm.Icon(id)                         // Make request
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
    var coords Coords
    var ok bool
    coords, ok = cityCache[key]
    if !ok {
        return 0, 0, false
    }
    return coords.lat, coords.lon, ok
}

func addCacheWeather(key string, r_obj owm.WeatherRange) {
    weatherCache[key] = r_obj
}

func searchCacheWeather(key string) (owm.WeatherRange, bool) {
    var r_obj owm.WeatherRange
    var ok bool
    r_obj, ok = weatherCache[key]
    if !ok {
        var emptyResponse owm.WeatherRange
        return emptyResponse, false
    }
    var tagAge uint = (uint(time.Now().Unix()) - r_obj.Timestamp)
    if tagAge > (SECONDS_IN_HOUR*CONFIG.CACHE_AGE) {
        delete(weatherCache, key)
        var emptyResponse owm.WeatherRange
        return emptyResponse, false
    }
    return r_obj, ok
}

func addCacheIcon(key string, data []byte) {
    iconCache[key] = data
}

func searchCacheIcon(key string) ([]byte, bool) {
    var data []byte
    var ok bool
    data, ok = iconCache[key]
    if !ok {
        var no_data []byte
        return no_data, false
    }
    return data, ok
}
