package main

/*
 * Structures for returning JSON data to front end
 */

import (
    "time"
)

type WeekWeather struct {
    City      string
    Timestamp uint
    Days      []DayWeather
}

type DayWeather struct {
    DayName    string
    RainChance int
    RainTotal  int
    Day        WeatherData
    Night      WeatherData
}

type WeatherData struct {
    Temp        float32
    Description string

    Wind struct {
        Speed float32
        Gust  float32
        Deg   int
        Dir   string
    }
    Humidity float32
    Pressure int       // Sea level
    Clouds struct {
        Clouds uint
        Low    uint    // Cloud layers
        Mid    uint
        High   uint
    }
    Rain struct {
        Chance int // %
        Amount int // mm
    }
    Visibility uint
    IconID     string
}

const REFERENCE_TIME uint = 1740096000
const SECONDS_IN_DAY uint = 86400;
const SECONDS_IN_HOUR uint = 3600;

func convertTime(timestamp uint) (string, uint) {
    calibratedTime := timestamp - REFERENCE_TIME
    hour := (calibratedTime % SECONDS_IN_DAY) / SECONDS_IN_HOUR
    dayname := time.Unix(int64(timestamp), 0).Weekday().String()
    return dayname, hour
}

func mapDays(raw_weather InWeatherRange) WeekWeather {
    // Expected timeslots
    // 0 3 6 9 12 15 18 21 (24)
    var week WeekWeather
    var rainSum float32
    var maxChance float32
    var i int = 0
    var day_no int = 0
    var newDay bool = true
    days := make([]DayWeather, 5, 8)
    for i < len(raw_weather.List) {
        dayName, hour := convertTime(raw_weather.List[i].Dt)
        if raw_weather.List[i].Pop > maxChance {
            maxChance = raw_weather.List[i].Pop
        }
        rainSum += raw_weather.List[i].Rain.Mm + raw_weather.List[i].Snow.Mm
        if newDay {
            // TODO: Make sure there is some data for first day
            newDay = false
        }
        if hour == 12 {
            // Populate 12:00 data
            days[day_no].DayName           = dayName
            populateData(&days[day_no].Day, &raw_weather.List[i])
        } else if hour == 21 {
            // Populate 21:00 data
            populateData(&days[day_no].Night, &raw_weather.List[i])
            days[day_no].RainChance        = toInt(maxChance*100)
            days[day_no].RainTotal         = toInt(rainSum)
            day_no += 1
            newDay = true
        } else {
            // TODO: Make sure there is data for last night
        }
        i += 1
    }
    // Populate metadata
    week.City = raw_weather.City.Name
    week.Timestamp = uint(time.Now().Unix())
    week.Days = days
    return week
}

func populateData(target *WeatherData, source *InWeather) {
    target.Temp          = source.Main.Temp
    target.Pressure      = source.Main.Sea_level
    target.Humidity      = source.Main.Humidity
    target.Description   = source.Weather[0].Description
    target.IconID        = source.Weather[0].Icon
    target.Clouds.Clouds = source.Clouds.All
    target.Visibility    = source.Visibility
    target.Wind.Speed    = source.Wind.Speed
    target.Wind.Deg      = source.Wind.Deg
    target.Wind.Dir      = windToStr(source.Wind.Deg)
}

func windToStr(wind int) string {
    if wind == 0 {
        return "N"
    }
    switch (wind+(45/2))/45 {
    case 1:
        return "NE"
    case 2:
        return "E"
    case 3:
        return "SE"
    case 4:
        return "S"
    case 5:
        return "SW"
    case 6:
        return "W"
    case 7:
        return "W"
    case 8:
        return "NW"
    default:
        return "N"
    }
}

func toInt(n float32) int {
    n += 0.5
    return int(n)
}
