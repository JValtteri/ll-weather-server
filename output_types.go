package main

/*
 * Structures for returning JSON data to front end
 */

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/JValtteri/weather/om"
    "github.com/JValtteri/weather/owm"
)

type WeekWeather struct {
    City      string
    Timestamp uint
    Days      []DayWeather
    Coord     struct {
        Lat   float32
        Lon   float32
    }
}

type DayWeather struct {
    DayName    string
    Day        WeatherData
    Night      WeatherData
    Coord     struct {
        Lat   float32
        Lon   float32
    }
}

type DayHours struct {
    City      string
    Timestamp uint
    Hours     []WeatherData
    Coord     struct {
        Lat   float32
        Lon   float32
    }
}

type WeatherData struct {
    Title       string   // Day or Time
    Temp        struct {
        Temp    float32
        Feels   float32
        Bulb    float32
    }
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
    Radiation struct {
        Direct float32
        Diffuse float32
    }
    Visibility uint
    Uv         float32
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
    //dayname := time.Unix(int64( timestamp - ( SECONDS_IN_HOUR * 3 ) - 1 ), 0).Weekday().String()
    // offset by one second to make midnight part of the previous day
    dayname := time.Unix(int64( timestamp - 1 ), 0).Weekday().String()
    return dayname, uint(hour)
}

/* Maps raw_weather data from OWM in to a WeekWeather struct
 */
func mapOwmDays(raw_weather owm.WeatherRange) WeekWeather {
    // Expected timeslots
    // 3 6 9 12 15 18 21 (24/0)
    var week WeekWeather
    var rainSum float32
    var maxChance float32
    var day_no int = 0
    var newDay bool = true
    var sunrise uint = raw_weather.City.Sunrise
    var sunset uint = raw_weather.City.Sunset

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
            populateOwmData(
                &days[day_no].Day, &raw_weather.List[i], dayName, maxChance, rainSum, hour, sunrise, sunset,
            )
            newDay = false
        }
        switch hour {
            case 12:
                // Populate 12:00 data
                days[day_no].DayName           = dayName
                populateOwmData(
                    &days[day_no].Day, &raw_weather.List[i], dayName, maxChance, rainSum, hour, sunrise, sunset,
                )
            case 15:
                // Set final day rain
                days[day_no].Day.Rain.Chance   = toInt(maxChance*100)
                days[day_no].Day.Rain.Amount   = rainSum
                // Start counting night rain
                maxChance = 0
                rainSum = 0
            case 21:
                // Populate 21:00 data
                days[day_no].DayName           = dayName
                populateOwmData(
                    &days[day_no].Night, &raw_weather.List[i], dayName, maxChance, rainSum, hour, sunrise, sunset,
                )
            case 0:
                days[day_no].DayName           = dayName
                populateOwmData(
                    &days[day_no].Night, &raw_weather.List[i], dayName, maxChance, rainSum, hour,
                    raw_weather.City.Sunrise, raw_weather.City.Sunset,
                )
                day_no += 1
                newDay = true
            default:
                // TODO: Make sure there is data for last night
        }
    }
    // Populate metadata
    week.City = raw_weather.City.Name
    week.Timestamp = uint(time.Now().Unix())
    week.Days = days
    return week
}

/* Maps raw_weather data from OWM in to a DayHours struct
 */
func mapOwmHours(raw_weather owm.WeatherRange, dayIndex int) DayHours {
    var day_no int = 0
    var sunrise uint        = raw_weather.City.Sunrise
    var sunset uint         = raw_weather.City.Sunset
    hours := make([]WeatherData, 0, 8)
    for i := 0 ; i < len(raw_weather.List) ; i++ {
        _, hour := convertTime(raw_weather.List[i].Dt)
        if dayIndex == day_no {
            var hourData WeatherData
            var title string        = fmt.Sprintf("%v:00", hour)
            var rainChance float32  = raw_weather.List[i].Pop*100
            var rainSum float32     = raw_weather.List[i].Rain.Mm + raw_weather.List[i].Snow.Mm
            populateOwmData(
                &hourData, &raw_weather.List[i], title, rainChance, rainSum,
                hour, sunrise, sunset,
            )
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

/* Maps raw_weather data from OM in to a WeekWeather struct
 */
func mapOmDays(raw_weather om.WeatherRange, city string, lat float32, lon float32) WeekWeather {
    // Expected timeslots hourly
    // 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 (24/0)
    var week WeekWeather
    var rainSum float32
    var maxChance int
    var day_no int = 0
    var newDay bool = true
    DT_OFFSET = raw_weather.UTC_OFFSET_SECONDS       // Update tiemzone offset
    days := make([]DayWeather, 1, 24)
    for i := 0 ; i < len(raw_weather.Hourly.Time) ; i++ {
        dayName, hour := convertTime(raw_weather.Hourly.Time[i])
        // Trim away past days
        if timeInPast(raw_weather.Hourly.Time[i]) {
            continue
        }
        if hour == 6 {
            // Start counting rain for the day from 6:00
            maxChance = 0
            rainSum = 0
        }
        if raw_weather.Hourly.Precipitation_probability[i] > maxChance {
            maxChance = raw_weather.Hourly.Precipitation_probability[i]
        }
        rainSum += raw_weather.Hourly.Precipitation[i]
        if newDay {
            // Populate data for a day fragment
            days[day_no].DayName = dayName      // // // Something occationally causes an out of index error here // // //
            populateOmData(&days[day_no].Day, &raw_weather, i, dayName, maxChance, rainSum)
            newDay = false

        }
        switch hour {
            case 12:
                // Populate 12:00 data
                days[day_no].DayName = dayName
                populateOmData(&days[day_no].Day, &raw_weather, i, dayName, maxChance, rainSum)
            case 15:
                // Set final day rain
                days[day_no].Day.Rain.Chance   = maxChance
                days[day_no].Day.Rain.Amount   = rainSum
                // Start counting night rain
                maxChance = 0
                rainSum = 0
            case 23:
                // Populate 21:00 data
                days[day_no].DayName = dayName
                populateOmData(&days[day_no].Night, &raw_weather, i, dayName, maxChance, rainSum)
            case 0:
                days[day_no].DayName = dayName
                populateOmData(&days[day_no].Night, &raw_weather, i, dayName, maxChance, rainSum)
                day_no += 1
                newDay = true
                if len(raw_weather.Hourly.Surface_pressure) > i+3 {
                    // If pressure shows zero, rest of the data is null
                    if raw_weather.Hourly.Surface_pressure[i+1] > 0 {
                        days = append(days, DayWeather{})
                    } else {
                        break
                    }
                }
        }
    }
    // Populate metadata
    //week.City = raw_weather.City.Name
    week.Timestamp = uint(time.Now().Unix())
    week.Days = days
    week.City = city
    week.Coord.Lat = lat
    week.Coord.Lon = lon
    return week
}

/* Maps raw_weather data from OM in to a DayHours struct
 */
func mapOmHours(raw_weather om.WeatherRange, dayIndex int) DayHours {
    var day_no int = 0
    hours := make([]WeatherData, 0, 24)
    for i:=0 ; i< len(raw_weather.Hourly.Time) ; i++ {
        _, hour := convertTime(raw_weather.Hourly.Time[i])
        // Trim away past hours
        if timeInPast(raw_weather.Hourly.Time[i]) {
            continue
        }

        if len(raw_weather.Hourly.Surface_pressure) > i+3 {
            // If pressure shows zero, rest of the data is null
            if raw_weather.Hourly.Surface_pressure[i+1] == 0 {
                break
            }
        }

        if dayIndex == day_no {
            var hourData WeatherData
            var title string    = fmt.Sprintf("%v:00", hour)
            var maxChance int   = raw_weather.Hourly.Precipitation_probability[i]
            var rainSum float32 = raw_weather.Hourly.Precipitation[i]
            populateOmData(&hourData, &raw_weather, i, title, maxChance, rainSum)
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
    day.Hours     = hours
    day.Coord.Lat = raw_weather.Latitude
    day.Coord.Lon = raw_weather.Longtitude
    //day.City = raw_weather.City.Name
    day.Timestamp = uint(time.Now().Unix())
    return day
}

func timeInPast(dt uint) bool {
    now := uint(time.Now().Unix())
    // Consider data points over 1 hour old to be "too old"
    if dt + SECONDS_IN_HOUR < now {
        return true
    }
    return false
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

/* Copies data from from OWM InWeather to WeatherData
 */
func populateOwmData(
    target *WeatherData,
    source *owm.Weather,
    title string,
    maxChance float32,
    rainSum float32,
    hour uint,
    sunrise uint,
    sunset uint,
) {
    target.Title         = title
    target.Rain.Chance   = toInt(maxChance*100)
    target.Rain.Amount   = rainSum
    target.SunUp         = sunUp(hour, sunrise, sunset)

    target.Temp.Temp     = source.Main.Temp
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

/* Copies data from from OM InWeather to WeatherData
 */
func populateOmData(target *WeatherData, source *om.WeatherRange, index int, dayName string, maxChance int, rainSum float32) {
    target.Title         = dayName
    target.Rain.Chance   = maxChance
    target.Rain.Amount   = rainSum
    target.SunUp         = toBool(source.Hourly.Is_day[index])
    target.Temp.Temp     = source.Hourly.Temperature_2m[index]
    target.Temp.Feels    = source.Hourly.Apparent_temperature[index]
    target.Temp.Bulb     = source.Hourly.Wet_bulb_temperature_2m[index]
    target.Uv            = source.Hourly.Uv_index[index]
    target.Pressure      = toInt(source.Hourly.Pressure_msl[index])
    target.Humidity      = float32(source.Hourly.Relative_humidity_2m[index])
    //target.Description   = source.Weather[0].Description[index]
    target.Clouds.Clouds = uint(source.Hourly.Cloud_cover[index])
    target.Clouds.Low    = uint(source.Hourly.Cloud_cover_low[index])
    target.Clouds.Mid    = uint(source.Hourly.Cloud_cover_mid[index])
    target.Clouds.High   = uint(source.Hourly.Cloud_cover_high[index])
    target.Visibility    = uint(toInt(source.Hourly.Visibility[index]))
    target.SunUp         = toBool(source.Hourly.Is_day[index])
    target.Radiation.Direct  = source.Hourly.Direct_radiation[index]
    target.Radiation.Diffuse = source.Hourly.Diffuse_radiation[index]
    target.Wind.Speed    = source.Hourly.Wind_speed_10m[index]
    target.Wind.Gust     = source.Hourly.Wind_gusts_10m[index]
    target.Wind.Deg      = source.Hourly.Wind_direction_10m[index]
    target.Wind.Dir      = windToStr(target.Wind.Deg)
    target.IconID        = wmoToOwm(source.Hourly.Weather_code[index], target.SunUp)
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

func toBool(n int) bool {
    if n == 1 {
        return true
    } else {
        return false
    }
}

func unloadJSON(object any) string {
    body, err := json.Marshal(object)
    if err != nil {
        log.Println("JSON response marshalling error:", err)
    }
    return string(body)
}

/* Translates WNO WW codes to OWM weather image codes
 */
func wmoToOwm(ww int, day bool) string {
    var iconId string
    if day {
        iconId = OwmCodes[ww] + "d"
    } else {
        iconId = OwmCodes[ww] + "n"
    }
    return iconId
}

// Conversion table for WNO WW codes to OWM weather image codes
var OwmCodes = map[int]string{
    // Clear
    0: "01",
    // Clouds
    1: "02",
    2: "03",
    3: "04",
    // Fog
    45: "50",
    48: "50",
    // Drizzle
    51: "10",
    53: "10",
    55: "10",
    // Rain
    61: "10",
    63: "09",
    65: "09",
    // Freezing Rain
    66: "09",
    67: "09",
    // Snow
    71: "13",
    73: "13",
    75: "13",
    // Hail
    77: "13",
    // Showers
    80: "10",
    81: "10",
    82: "09",
    // Snow
    85: "13",
    86: "13",
    // Thunder
    95: "11",
    96: "11",
    99: "11",
}
