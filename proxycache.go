package main

import (
    "errors"
    "fmt"
    "log"
    "time"
    "github.com/JValtteri/weather/owm"
    "github.com/JValtteri/weather/om"
)

type Coords struct {
    lat float32
    lon float32
}

type OmWrapper struct {
    dt uint
    data om.WeatherRange
}

func omWrap(data om.WeatherRange) OmWrapper {
    var wrapper OmWrapper
    wrapper.data = data
    wrapper.dt = uint(time.Now().Unix())
    return wrapper
}

func omUnwrap(wrapper OmWrapper) (om.WeatherRange, uint) {
    return wrapper.data, wrapper.dt
}

var owmWeatherCache map[string]owm.WeatherRange  = make(map[string]owm.WeatherRange)
var omWeatherCache map[string]OmWrapper   = make(map[string]OmWrapper)
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

func GetOwmProxyWeather(city string) (owm.WeatherRange, error) {
    var lat, lon float32 = GetCityCoord(city)   // Convert name to coord
    if lat==0 && lon==0 {
        var emptyWeather owm.WeatherRange
        err := errors.New(fmt.Sprintf("r:%4vu:%4v: City [%s] not found", rqNum, uniqRqNum, city))
        log.Println(err)
        return emptyWeather, err
    }
    var r_obj owm.WeatherRange
    var ok bool
    r_obj, ok = searchOwmCacheWeather(city)        // Check cache
    if ok {
        log.Printf("r:%6vu:%4v: Get weather: %s at %.3f %.3f\n", rqNum, uniqRqNum, city, lat, lon)
        return r_obj, nil
    }
    uniqRqNum++
    log.Printf("r:%6vu:%4v: Get weather: %s at %.3f %.3f (New request)\n", rqNum, uniqRqNum, city, lat, lon)
    r_obj = owm.Forecast(lat, lon)              // Make request
    r_obj.Timestamp = uint(time.Now().Unix())
    addOwmCacheWeather(city, r_obj)                // Cache response
    return r_obj, nil
}

func GetOmProxyWeather(city string) (om.WeatherRange, error) {
    var lat, lon float32 = GetCityCoord(city)   // Convert name to coord
    if lat==0 && lon==0 {
        var emptyWeather om.WeatherRange
        err := errors.New(fmt.Sprintf("r:%4vu:%4v: City [%s] not found", rqNum, uniqRqNum, city))
        log.Println(err)
        return emptyWeather, err
    }
    var r_obj om.WeatherRange
    var ok bool
    r_obj, ok = searchOmCacheWeather(city)        // Check cache
    if ok {
        log.Printf("r:%6vu:%4v: Get weather: %s at %.3f %.3f\n", rqNum, uniqRqNum, city, lat, lon)
        return r_obj, nil
    }
    uniqRqNum++
    log.Printf("r:%6vu:%4v: Get weather: %s at %.3f %.3f (New request)\n", rqNum, uniqRqNum, city, lat, lon)
    r_obj = om.Forecast(lat, lon)                 // Make request
    addOmCacheWeather(city, r_obj)                // Cache response
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
    if len(cityCache) < CONFIG.CACHE_SIZE {
        cityCache[key] = coords
    } else {
        cullMap(&cityCache)
        cityCache[key] = coords
    }
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

func addOwmCacheWeather(key string, r_obj owm.WeatherRange) {
    if len(owmWeatherCache) < CONFIG.CACHE_SIZE {
        owmWeatherCache[key] = r_obj
    } else {
        cullMap(&owmWeatherCache)
        owmWeatherCache[key] = r_obj
    }
}

func addOmCacheWeather(key string, r_obj om.WeatherRange) {
    var wrapper OmWrapper = omWrap(r_obj)
    if len(omWeatherCache) < CONFIG.CACHE_SIZE {
        omWeatherCache[key] = wrapper
    } else {
        cullMap(&omWeatherCache)
        omWeatherCache[key] = wrapper
    }
}

// Throws away half (every 2nd) item of a given map
func cullMap[V Coords | owm.WeatherRange | om.WeatherRange | OmWrapper | []byte ](m *map[string]V) {
    var i int = 0
    for key, _ := range *m {
        if i == 0 {
            delete(*m, key)
        }
        i = (i + 1) % 2
    }
}

func searchOwmCacheWeather(key string) (owm.WeatherRange, bool) {
    var r_obj owm.WeatherRange
    var ok bool
    r_obj, ok = owmWeatherCache[key]
    if !ok {
        return owm.WeatherRange{}, false
    }
    var tagAge uint = (uint(time.Now().Unix()) - r_obj.Timestamp)
    if tagAge > (SECONDS_IN_HOUR*CONFIG.CACHE_AGE) {
        delete(owmWeatherCache, key)
        return r_obj, false
    }
    return r_obj, ok
}

func searchOmCacheWeather(key string) (om.WeatherRange, bool) {
    var r_obj om.WeatherRange
    var dt uint
    var wrapper OmWrapper
    var ok bool
    wrapper, ok = omWeatherCache[key]
    if !ok {
        return om.WeatherRange{}, false
    }
    r_obj, dt = omUnwrap(wrapper)
    if len(r_obj.Hourly.Time) == 0 {
        log.Println("OM Data range is empty!")
        return om.WeatherRange{}, false
    }
    var tagAge uint = (uint(time.Now().Unix()) - dt)
    //fmt.Printf("Now: \t%v\nDT:\t%v\nAge: \t%v\nMAX: \t%v\n", uint(time.Now().Unix()), dt, tagAge, SECONDS_IN_HOUR*CONFIG.CACHE_AGE)
    if tagAge > (SECONDS_IN_HOUR*CONFIG.CACHE_AGE) {
        delete(omWeatherCache, key)
        return r_obj, false
    }
    return r_obj, ok
}

func addCacheIcon(key string, data []byte) {
    if len(iconCache) < CONFIG.CACHE_SIZE {
        iconCache[key] = data
    } else {
        cullMap(&iconCache)
        iconCache[key] = data
    }
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
