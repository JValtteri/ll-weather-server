package main

import (
    "testing"
    "time"
    "github.com/JValtteri/weather/owm"
)

func setup() {
    CONFIG.CACHE_SIZE = 5
    CONFIG.CACHE_AGE = 1
    CONFIG.NETWORK = false
    owm.Config("", "", "", false)  // Prevent requests
}

func TestAddCityCache(t *testing.T) {
    setup()
    // cityCache
    addCacheCity("a", -1.25, 2.55)
    addCacheCity("b", 3.25, -2.55)
    addCacheCity("c", 0, 0)

    exa, eba := []float32{3.25, -2.55}, true
    a, b, c := searchCacheCity("b")
    if a != exa[0] || c !=  eba {
        t.Errorf("City Cache 'b' %v %v %v was %v %v %v", exa[0], exa[1], eba, a, b, c)
    }
    exb, ebb := []float32{-1.25, 2.55}, true
    a, b, c = searchCacheCity("a")
    if b != exb[1] || c != ebb {
        t.Errorf("City Cache 'a' %v %v %v was %v %v %v", exb[0], exb[1], ebb, a, b, c)
    }
    exc, ebc := []float32{0.0, 0.0}, true
    a, b, c = searchCacheCity("c")
    if a != exc[0] || c !=  ebc {
        t.Errorf("City Cache 'c' %v %v %v was %v %v %v", exc[0], exc[1], ebc, a, b, c)
    }
    exd, ebd := []float32{0.0, 0.0}, false
    a, b, c = searchCacheCity("d")
    if b != exd[1] || c !=  ebd {
        t.Errorf("City Cache 'd' %v %v %v was %v %v %v", exd[0], exd[1], ebd, a, b, c)
    }
}

func TestCityCacheLimit(t *testing.T) {
    setup()
    CONFIG.CACHE_SIZE = 2
    // cityCache
    addCacheCity("a", -1.25, 2.55)
    addCacheCity("c", 0, 0)
    addCacheCity("b", 3.25, -2.55)
    // Check the last item was found
    exa, eba := []float32{3.25, -2.55}, true
    a, b, c := searchCacheCity("b")
    if a != exa[0] || c !=  eba {
        t.Errorf("City Cache 'b' %v %v %v was %v %v %v", exa[0], exa[1], eba, a, b, c)
    }
    // Check the limit works
    if len(cityCache) > int(CONFIG.CACHE_SIZE) {
        t.Errorf("Cache size was %v > limit %v", len(cityCache), CONFIG.CACHE_SIZE)
    }
}

func TestAddWeatherCache(t *testing.T) {
    setup()
    //weatherCache
    var obj1 owm.WeatherRange
    obj1.Timestamp = uint(time.Now().Unix())
    obj1.City.Id = 123
    addCacheWeather("1st", obj1)
    var obj2 owm.WeatherRange
    obj2.Timestamp = uint(time.Now().Unix())
    obj2.City.Id = 231
    addCacheWeather("2nd", obj2)
    var obj3 owm.WeatherRange
    obj3.Timestamp = 0
    obj3.City.Id = 231
    addCacheWeather("2nd", obj3)
    fob, ok := searchCacheWeather("1st")
    if fob.City.Id != obj1.City.Id || !ok {
        t.Errorf("WeatherCache 1st not found, %v != %v, ok:%v", fob.City.Id, obj1.City.Id, ok)
    }
    fob, ok = searchCacheWeather("2nd")
    if fob.City.Id != obj2.City.Id || !ok {
        t.Errorf("WeatherCache 2nd not found %v != %v, ok:%v", fob.City.Id, obj1.City.Id, ok)
    }
    if ok {
        t.Errorf("WeatherCache violated CACHE_AGE")
    }
}

func TestWeatherCacheLimit(t *testing.T) {
    setup()
    CONFIG.CACHE_SIZE = 1
    // cityCache
    var obj1 owm.WeatherRange
    obj1.City.Id = 123
    addCacheWeather("1st", obj1)
    var obj2 owm.WeatherRange
    obj2.Timestamp = uint(time.Now().Unix())
    obj2.City.Id = 231
    addCacheWeather("2nd", obj2)
    // Check the last item was found
    fob, ok := searchCacheWeather("2nd")
    if fob.City.Id != obj2.City.Id || !ok {
        t.Errorf("WeatherCache 1st not found")
    }
    // Check the limit works
    if len(weatherCache) > int(CONFIG.CACHE_SIZE) {
        t.Errorf("Cache size was %v > limit %v", len(weatherCache), CONFIG.CACHE_SIZE)
    }
}

func TestCityCoord(t *testing.T) {
    setup()
    addCacheCity("Atlantis", 89.9, -90.0)
    lat, lon := GetCityCoord("Atlantis")
    if lat != 89.9 || lon != -90.0 {
        t.Errorf("Atlantis not found")
    }
}

func TestGetWeather(t *testing.T) {
    setup()
    addCacheCity("Atlantis", 89.9, -90)
    var atlantis owm.WeatherRange
    atlantis.City.Id = 123
    atlantis.Timestamp = uint(time.Now().Unix())
    addCacheWeather("Atlantis", atlantis)
    city, err := GetProxyWeather("Atlantis")
    if err != nil {
        t.Errorf("Get Atlantis got an Error: %v", err)
    }
    if city.City.Id != atlantis.City.Id {
        t.Errorf("Atlantis not found")
    }
}
