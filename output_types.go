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
    RainChance float32
    RainTotal  float32
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
        Chance float32 // %
        Amount float32 // mm
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
    var hour uint = 0
    var rainSum float32
    var maxChance float32
    var i int = 0
    var day_no int = 0
    var day DayWeather
    var emptyDay DayWeather
    days := make([]DayWeather, 5, 8)
    for i < len(raw_weather.List) {
        if day_no == len(days) {
            days = append(days, emptyDay)
        }
        day = days[day_no]
        _, hour = convertTime(raw_weather.List[i].Dt)
        if raw_weather.List[i].Pop > maxChance {
            maxChance = raw_weather.List[i].Pop
        }
        rainSum += raw_weather.List[i].Rain.Mm + raw_weather.List[i].Snow.Mm
        if hour == 12 {
            // Populate 12:00 data
            day.DayName, _        = convertTime(raw_weather.List[i].Dt)
            day.Day.Temp          = raw_weather.List[i].Main.Temp
            day.Day.Pressure      = raw_weather.List[i].Main.Sea_level
            day.Day.Humidity      = raw_weather.List[i].Main.Humidity
            day.Day.Description   = raw_weather.List[i].Weather[0].Description
            day.Day.IconID        = raw_weather.List[i].Weather[0].Icon
            day.Day.Clouds.Clouds = raw_weather.List[i].Clouds.All
            day.Day.Visibility    = raw_weather.List[i].Visibility
            day.Day.Wind.Speed    = raw_weather.List[i].Wind.Speed
            day.Day.Wind.Deg      = raw_weather.List[i].Wind.Deg
            day.Day.Wind.Dir      = windToStr(raw_weather.List[i].Wind.Deg)
            days[day_no] = day  // This should be done differently.
        } else if hour == 21 {
            // Move index to night 22:00
            day.Night.Temp          = raw_weather.List[i].Main.Temp
            day.Night.Humidity      = raw_weather.List[i].Main.Humidity
            day.Night.Description   = raw_weather.List[i].Weather[0].Description
            day.Night.IconID        = raw_weather.List[i].Weather[0].Icon
            day.Night.Clouds.Clouds = raw_weather.List[i].Clouds.All
            day.Night.Visibility    = raw_weather.List[i].Visibility
            day.Night.Wind.Speed    = raw_weather.List[i].Wind.Speed
            day.Night.Wind.Deg      = raw_weather.List[i].Wind.Deg
            day.Night.Wind.Dir      = windToStr(raw_weather.List[i].Wind.Deg)

            day.RainChance        = maxChance
            day.RainTotal         = rainSum
            days[day_no] = day  // This should be done differently.
            day_no += 1
        }
        i += 1
    }
    // Populate metadata
    week.City = raw_weather.City.Name
    week.Timestamp = uint(time.Now().Unix())
    week.Days = days
    return week
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
