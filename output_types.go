package main

/*
 * Structures for returning JSON data to front end
 */

import (
    "time"
    "fmt"
    "log"
    "encoding/json"
    "github.com/JValtteri/weather/owm"
)

type WeekWeather struct {
    City      string
    Timestamp uint
    Days      []DayWeather
}

type DayWeather struct {
    DayName    string
    Day        WeatherData
    Night      WeatherData
}

type DayHours struct {
    City      string
    Timestamp uint
    Hours     []WeatherData
}

type WeatherData struct {
    Title       string   // Day or Time
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
        Chance int     // %
        Amount float32 // mm
    }
    Visibility uint
    IconID     string
    SunUp      bool
}

const REFERENCE_TIME uint = 1740096000  // UTC (+0:00)
const SECONDS_IN_DAY uint = 86400;
const SECONDS_IN_HOUR uint = 3600;
var DT_OFFSET int = 0;

/* Converts Posix timestamp and time offset (timezone) to
 * dayname string
 * hour uint
 */
func convertTime(timestamp uint) (string, uint) {
    calibratedTime := int(timestamp) - int(REFERENCE_TIME) + DT_OFFSET
    hour := (calibratedTime % int(SECONDS_IN_DAY)) / int(SECONDS_IN_HOUR)
    // dayname calculation is offset by three hours and one second 3:00:01, to make
    // 00:00 day is counted as part of previous day, not the next
    dayname := time.Unix(int64( timestamp - ( SECONDS_IN_HOUR * 3 ) - 1 ), 0).Weekday().String()
    return dayname, uint(hour)
}

/* Maps raw_weather data in to a WeekWeather struct
 */
func mapDays(raw_weather owm.WeatherRange) WeekWeather {
    // Expected timeslots
    // 3 6 9 12 15 18 21 (24/0)
    var week WeekWeather
    var rainSum float32
    var maxChance float32
    var day_no int = 0
    var newDay bool = true
    DT_OFFSET = raw_weather.City.Timezone       // Update tiemzone offset
    days := make([]DayWeather, 5, 8)
    for i:=0 ; i<len(raw_weather.List) ; i++ {
        if day_no == 5 {
            break
        }
        dayName, hour := convertTime(raw_weather.List[i].Dt)
        if hour == 6 {
            // Start counting rain for the day from 6:00
            maxChance = 0
            rainSum = 0
        }
        if raw_weather.List[i].Pop > maxChance {
            maxChance = raw_weather.List[i].Pop
        }
        rainSum += raw_weather.List[i].Rain.Mm + raw_weather.List[i].Snow.Mm
        if newDay {
            // Populate data for a day fragment
            days[day_no].DayName           = dayName
            days[day_no].Day.Title         = dayName
            days[day_no].Day.Rain.Chance   = toInt(maxChance*100)
            days[day_no].Day.Rain.Amount   = rainSum
            days[day_no].Day.SunUp         = sunUp(hour, raw_weather.City.Sunrise, raw_weather.City.Sunset)
            populateData(&days[day_no].Day, &raw_weather.List[i])
            newDay = false

        }
        if hour == 12 {
            // Populate 12:00 data
            days[day_no].DayName           = dayName
            days[day_no].Day.Title         = dayName
            days[day_no].Day.Rain.Chance   = toInt(maxChance*100)
            days[day_no].Day.Rain.Amount   = rainSum
            days[day_no].Day.SunUp         = sunUp(hour, raw_weather.City.Sunrise, raw_weather.City.Sunset)
            populateData(&days[day_no].Day, &raw_weather.List[i])
        } else if hour == 15 {
            // Set final day rain
            days[day_no].Day.Rain.Chance   = toInt(maxChance*100)
            days[day_no].Day.Rain.Amount   = rainSum
            // Start counting night rain
            maxChance = 0
            rainSum = 0
        } else if hour == 21 {
            // Populate 21:00 data
            days[day_no].DayName           = dayName
            days[day_no].Night.Title       = dayName
            days[day_no].Night.Rain.Chance = toInt(maxChance*100)
            days[day_no].Night.Rain.Amount = rainSum
            days[day_no].Night.SunUp         = sunUp(hour, raw_weather.City.Sunrise, raw_weather.City.Sunset)
            populateData(&days[day_no].Night, &raw_weather.List[i])
            days[day_no].Night.Rain.Chance       = toInt(maxChance*100)
            days[day_no].Night.Rain.Amount       = rainSum
        } else if hour == 0 {
            days[day_no].DayName           = dayName
            days[day_no].Night.Title       = dayName
            days[day_no].Night.Rain.Chance = toInt(maxChance*100)
            days[day_no].Night.Rain.Amount = rainSum
            days[day_no].Night.SunUp       = sunUp(hour, raw_weather.City.Sunrise, raw_weather.City.Sunset)
            populateData(&days[day_no].Night, &raw_weather.List[i])
            days[day_no].Night.Rain.Chance       = toInt(maxChance*100)
            days[day_no].Night.Rain.Amount       = rainSum
            day_no += 1
            newDay = true
        } else {
            // TODO: Make sure there is data for last night
        }
    }
    // Populate metadata
    week.City = raw_weather.City.Name
    week.Timestamp = uint(time.Now().Unix())
    week.Days = days
    return week
}

/* Maps raw_weather data in to a DayHours struct
 */
func mapHours(raw_weather owm.WeatherRange, dayIndex int) DayHours {
    var day_no int = 0
    hours := make([]WeatherData, 0, 8)
    for i:=0 ; i< len(raw_weather.List) ; i++ {
        _, hour := convertTime(raw_weather.List[i].Dt)
        if dayIndex == day_no {
            var hourData WeatherData
            hourData.Title = fmt.Sprintf("%v:00", hour)
            hourData.Rain.Chance = toInt(raw_weather.List[i].Pop*100)
            hourData.Rain.Amount = raw_weather.List[i].Rain.Mm + raw_weather.List[i].Snow.Mm
            hourData.SunUp = sunUp(hour, raw_weather.City.Sunrise, raw_weather.City.Sunset)
            populateData(&hourData, &raw_weather.List[i])
            hours = append(hours, hourData)
            if hour == 0 {
                break
            }
        }
        if hour == 0 {
            day_no += 1
        }
    }
    var day DayHours
    day.Hours = hours
    day.City = raw_weather.City.Name
    day.Timestamp = uint(time.Now().Unix())
    return day
}

/* Returns bool true, if The Sun is up at given hour.
 */
func sunUp(hour uint, sunrise_dt uint, sunset_dt uint) bool {
    _, sunrise := convertTime(sunrise_dt)
    _, sunset := convertTime(sunset_dt)
    if sunset == 0 {
        sunset = 24
    }
    return sunrise < hour && hour < sunset
}

/* Copies data from InWeather to WeatherData
 */
func populateData(target *WeatherData, source *owm.Weather) {
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

/* Maps wind angle in to a string representation of it's direction
 */
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

/* Simple round float32 to int
 */
func toInt(n float32) int {
    n += 0.5
    return int(n)
}

func unloadJSON(object any) string {
    body, err := json.Marshal(object)
    if err != nil {
        log.Println("JSON response marshalling error:", err)
    }
    return string(body)
}
