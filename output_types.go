package main

/*
 * Structures for returning JSON data to front end
 */

import (
    "encoding/json"
    "fmt"
    "log"
    "sort"
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

type OverviewVals struct {
    Min    float32
    Max    float32
    Median float32
    Count  int
    values []float32
}

type Overview struct {
    Temp         OverviewVals
    Feels        OverviewVals
    Bulb         OverviewVals
    Uv           OverviewVals
    Direct       OverviewVals
    Diffuse      OverviewVals
    Clouds       OverviewVals
    CloudsLo     OverviewVals
    CloudsMed    OverviewVals
    CloudsHigh   OverviewVals
    Pressure     OverviewVals
    Humidity     OverviewVals
    Visibility   OverviewVals
    RainChance   OverviewVals
    RainAmount   OverviewVals
    WindSpeed    OverviewVals
    WindGust     OverviewVals
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
    Overview   *Overview
}

const REFERENCE_TIME uint = 1740096000  // UTC (+0:00)
const SECONDS_IN_DAY uint = 86400;
const SECONDS_IN_HOUR uint = 3600;
/* Models with less than two days of data are too short to contribute to an
 * aggregated range forecast.
 */
const minOverviewForecastHours int = int(2 * (SECONDS_IN_DAY / SECONDS_IN_HOUR));
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

/* A weatherSource wraps one or more model runs for the same location.
 * Regular forecasts use a single model; overview forecasts aggregate
 * several models and expose their min/max/median values.
 */
type weatherSource struct {
    ranges       []om.WeatherRange
    withOverview bool
}

func singleModelSource(r om.WeatherRange) weatherSource {
    return weatherSource{ranges: []om.WeatherRange{r}}
}

func overviewSource(ranges []om.WeatherRange) weatherSource {
    valid := make([]om.WeatherRange, 0, len(ranges))
    for _, r := range ranges {
        // Ignore models with missing or too short forecasts instead of
        // treating their missing values as zero.
        if usableForecast(r) {
            valid = append(valid, r)
        }
    }
    return weatherSource{ranges: valid, withOverview: true}
}

func usableForecast(r om.WeatherRange) bool {
    return forecastHours(r) >= minOverviewForecastHours
}

/* Returns the number of hours that contain actual forecast data. Open-Meteo
 * pads missing data with nulls, which unmarshal to zero. Surface pressure is
 * never zero, so the first zero marks the end of the forecast. Pressure_msl
 * is used as a fallback if surface pressure is unavailable.
 */
func forecastHours(r om.WeatherRange) int {
    length := len(r.Hourly.Time)
    hasPressure := false
    for _, pressure := range [][]float32{r.Hourly.Surface_pressure, r.Hourly.Pressure_msl} {
        if len(pressure) == 0 {
            continue
        }
        hasPressure = true
        limit := len(pressure)
        if limit > length {
            limit = length
        }
        end := 0
        for end < limit && pressure[end] != 0 {
            end++
        }
        if end > 0 {
            return end
        }
    }
    if hasPressure {
        return 0
    }
    return length
}

func (s weatherSource) empty() bool {
    return len(s.ranges) == 0
}

func (s weatherSource) reference() *om.WeatherRange {
    return &s.ranges[0]
}

/* Returns the number of hours to include in the output. An aggregated range
 * is limited to the shortest forecast of the models that have enough data.
 */
func (s weatherSource) length() int {
    if s.empty() {
        return 0
    }
    length := len(s.ranges[0].Hourly.Time)
    if s.withOverview {
        length = forecastHours(s.ranges[0])
    }
    for _, r := range s.ranges[1:] {
        modelLength := len(r.Hourly.Time)
        if s.withOverview {
            modelLength = forecastHours(r)
        }
        if modelLength < length {
            length = modelLength
        }
    }
    return length
}

/* Returns the weather data for a single time slot. Regular forecasts use the
 * running rain totals for the day; overview forecasts aggregate the models.
 */
func (s weatherSource) weatherAt(index int, title string, rainChance int, rainAmount float32) WeatherData {
    if s.withOverview {
        return aggregateOmAtIndex(s.ranges, index, title)
    }
    var wd WeatherData
    populateOmData(&wd, &s.ranges[0], index, title, rainChance, rainAmount)
    return wd
}

/* Maps raw_weather data from OM in to a WeekWeather struct
 */
func mapOmDays(source weatherSource, city string, lat float32, lon float32) WeekWeather {
    // Expected timeslots hourly
    // 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 (24/0)
    if source.empty() {
        return WeekWeather{}
    }
    ref := source.reference()
    var week WeekWeather
    DT_OFFSET = ref.UTC_OFFSET_SECONDS       // Update timezone offset
    days := make([]DayWeather, 1, 24)
    dayNo := 0
    newDay := true
    rain := newRainAccumulator(len(source.ranges))
    length := source.length()
    for i := 0; i < length; i++ {
        dayName, hour := convertTime(ref.Hourly.Time[i])
        // Trim away past days
        if timeInPast(ref.Hourly.Time[i]) {
            continue
        }
        if hour == 6 {
            // Start counting rain for the day from 6:00
            rain.reset()
        }
        rain.add(source.ranges, i)
        rainChance, rainAmount := rain.current()
        if newDay {
            // Populate data for a day fragment
            days[dayNo].DayName = dayName
            days[dayNo].Day = source.weatherAt(i, dayName, rainChance, rainAmount)
            newDay = false
        }
        if hour == 12 {
            // Populate 12:00 data
            days[dayNo].DayName = dayName
            days[dayNo].Day = source.weatherAt(i, dayName, rainChance, rainAmount)
        } else if hour == 15 {
            // Set final day rain
            if source.withOverview {
                chance, amount := rain.medians()
                days[dayNo].Day.Rain.Chance = toInt(chance.Median)
                days[dayNo].Day.Rain.Amount = amount.Median
                if days[dayNo].Day.Overview != nil {
                    days[dayNo].Day.Overview.RainChance = chance
                    days[dayNo].Day.Overview.RainAmount = amount
                }
            } else {
                days[dayNo].Day.Rain.Chance = rainChance
                days[dayNo].Day.Rain.Amount = rainAmount
            }
            // Start counting night rain
            rain.reset()
        } else if hour == 23 {
            // Populate 23:00 data
            days[dayNo].DayName = dayName
            days[dayNo].Night = source.weatherAt(i, dayName, rainChance, rainAmount)
        } else if hour == 0 {
            days[dayNo].DayName = dayName
            days[dayNo].Night = source.weatherAt(i, dayName, rainChance, rainAmount)
            if len(ref.Hourly.Surface_pressure) > i+3 && ref.Hourly.Surface_pressure[i+1] > 0 {
                // If pressure shows zero, rest of the data is null
                days = append(days, DayWeather{})
            } else {
                break
            }
            dayNo++
            newDay = true
        }
    }
    // Populate metadata
    week.Timestamp = uint(time.Now().Unix())
    week.Days = days
    week.City = city
    week.Coord.Lat = lat
    week.Coord.Lon = lon
    return week
}

/* Accumulates rain over a day for each model. */
type rainAccumulator struct {
    maxChance []int
    sum       []float32
}

func newRainAccumulator(models int) *rainAccumulator {
    return &rainAccumulator{
        maxChance: make([]int, models),
        sum:       make([]float32, models),
    }
}

func (r *rainAccumulator) reset() {
    for i := range r.maxChance {
        r.maxChance[i] = 0
        r.sum[i] = 0
    }
}

func (r *rainAccumulator) add(ranges []om.WeatherRange, index int) {
    for i := range r.maxChance {
        if index >= forecastHours(ranges[i]) {
            continue
        }
        if chance := ranges[i].Hourly.Precipitation_probability[index]; chance > r.maxChance[i] {
            r.maxChance[i] = chance
        }
        r.sum[i] += ranges[i].Hourly.Precipitation[index]
    }
}

/* Returns the rain totals for a single-model forecast. */
func (r *rainAccumulator) current() (int, float32) {
    return r.maxChance[0], r.sum[0]
}

/* Returns the median rain totals across all models. */
func (r *rainAccumulator) medians() (OverviewVals, OverviewVals) {
    chance := initOverviewVals()
    amount := initOverviewVals()
    for i := range r.maxChance {
        addVal(&chance, float32(r.maxChance[i]))
        addVal(&amount, r.sum[i])
    }
    finalizeOverviewVals(&chance)
    finalizeOverviewVals(&amount)
    return chance, amount
}

/* Lists every value set in the overview so they can be processed uniformly. */
func (o *Overview) values() []*OverviewVals {
    return []*OverviewVals{
        &o.Temp, &o.Feels, &o.Bulb, &o.Uv,
        &o.Direct, &o.Diffuse, &o.Clouds, &o.CloudsLo,
        &o.CloudsMed, &o.CloudsHigh, &o.Pressure, &o.Humidity,
        &o.Visibility, &o.RainChance, &o.RainAmount,
        &o.WindSpeed, &o.WindGust,
    }
}

func (o *Overview) finalize() {
    for _, vals := range o.values() {
        finalizeOverviewVals(vals)
    }
}

func collectOverview(r_objs []om.WeatherRange, index int) Overview {
    var ov Overview
    for _, vals := range ov.values() {
        *vals = initOverviewVals()
    }
    for _, r := range r_objs {
        length := forecastHours(r)
        if length < minOverviewForecastHours || index >= length {
            continue
        }
        h := &r.Hourly
        addVal(&ov.Temp, h.Temperature_2m[index])
        addVal(&ov.Feels, h.Apparent_temperature[index])
        addVal(&ov.Bulb, h.Wet_bulb_temperature_2m[index])
        addVal(&ov.Uv, h.Uv_index[index])
        addVal(&ov.Direct, h.Direct_radiation[index])
        addVal(&ov.Diffuse, h.Diffuse_radiation[index])
        addVal(&ov.Clouds, float32(h.Cloud_cover[index]))
        addVal(&ov.CloudsLo, float32(h.Cloud_cover_low[index]))
        addVal(&ov.CloudsMed, float32(h.Cloud_cover_mid[index]))
        addVal(&ov.CloudsHigh, float32(h.Cloud_cover_high[index]))
        addVal(&ov.Pressure, h.Pressure_msl[index])
        addVal(&ov.Humidity, float32(h.Relative_humidity_2m[index]))
        addVal(&ov.Visibility, h.Visibility[index])
        addVal(&ov.RainChance, float32(h.Precipitation_probability[index]))
        addVal(&ov.RainAmount, h.Precipitation[index])
        addVal(&ov.WindSpeed, h.Wind_speed_10m[index])
        addVal(&ov.WindGust, h.Wind_gusts_10m[index])
    }
    ov.finalize()
    return ov
}

func aggregateOmAtIndex(r_objs []om.WeatherRange, index int, title string) WeatherData {
    ov := collectOverview(r_objs, index)
    var wd WeatherData
    // Pick scalar values (icon, wind direction, sun) from the first usable model.
    populated := false
    for i := range r_objs {
        length := forecastHours(r_objs[i])
        if length < minOverviewForecastHours || index >= length {
            continue
        }
        populateOmData(&wd, &r_objs[i], index, title,
            toInt(ov.RainChance.Median), ov.RainAmount.Median)
        populated = true
        break
    }
    if !populated {
        wd.Title = title
        wd.Rain.Chance = toInt(ov.RainChance.Median)
        wd.Rain.Amount = ov.RainAmount.Median
    }
    // Display values are the model medians; min/max live in Overview.
    wd.Temp.Temp = ov.Temp.Median
    wd.Temp.Feels = ov.Feels.Median
    wd.Temp.Bulb = ov.Bulb.Median
    wd.Uv = ov.Uv.Median
    wd.Radiation.Direct = ov.Direct.Median
    wd.Radiation.Diffuse = ov.Diffuse.Median
    wd.Clouds.Clouds = uint(ov.Clouds.Median + 0.5)
    wd.Clouds.Low = uint(ov.CloudsLo.Median + 0.5)
    wd.Clouds.Mid = uint(ov.CloudsMed.Median + 0.5)
    wd.Clouds.High = uint(ov.CloudsHigh.Median + 0.5)
    wd.Pressure = toInt(ov.Pressure.Median)
    wd.Humidity = ov.Humidity.Median
    wd.Visibility = uint(toInt(ov.Visibility.Median))
    wd.Rain.Chance = toInt(ov.RainChance.Median)
    wd.Rain.Amount = ov.RainAmount.Median
    wd.Wind.Speed = ov.WindSpeed.Median
    wd.Wind.Gust = ov.WindGust.Median
    wd.Overview = &ov
    return wd
}

/* Maps raw_weather data from OM in to a DayHours struct
 */
func mapOmHours(source weatherSource, dayIndex int) DayHours {
    if source.empty() {
        return DayHours{}
    }
    ref := source.reference()
    dayNo := 0
    hours := make([]WeatherData, 0, 24)
    length := source.length()
    for i := 0; i < length; i++ {
        _, hour := convertTime(ref.Hourly.Time[i])
        // Trim away past hours
        if timeInPast(ref.Hourly.Time[i]) {
            continue
        }
        if len(ref.Hourly.Surface_pressure) > i+3 && ref.Hourly.Surface_pressure[i+1] == 0 {
            // If pressure shows zero, rest of the data is null
            break
        }
        if dayIndex == dayNo {
            title := fmt.Sprintf("%v:00", hour)
            rainChance := ref.Hourly.Precipitation_probability[i]
            rainAmount := ref.Hourly.Precipitation[i]
            hours = append(hours, source.weatherAt(i, title, rainChance, rainAmount))
            if hour == 0 {
                break
            }
        }
        if hour == 0 {
            dayNo++
        }
    }
    var day DayHours
    day.Hours = hours
    day.Coord.Lat = ref.Latitude
    day.Coord.Lon = ref.Longtitude
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
func populateOmData(target *WeatherData, source *om.WeatherRange, index int, title string, rainChance int, rainAmount float32) {
    target.Title         = title
    target.Rain.Chance   = rainChance
    target.Rain.Amount   = rainAmount
    target.SunUp         = toBool(source.Hourly.Is_day[index])
    target.Temp.Temp     = source.Hourly.Temperature_2m[index]
    target.Temp.Feels    = source.Hourly.Apparent_temperature[index]
    target.Temp.Bulb     = source.Hourly.Wet_bulb_temperature_2m[index]
    target.Uv            = source.Hourly.Uv_index[index]
    target.Pressure      = toInt(source.Hourly.Pressure_msl[index])
    target.Humidity      = float32(source.Hourly.Relative_humidity_2m[index])
    target.Clouds.Clouds = uint(source.Hourly.Cloud_cover[index])
    target.Clouds.Low    = uint(source.Hourly.Cloud_cover_low[index])
    target.Clouds.Mid    = uint(source.Hourly.Cloud_cover_mid[index])
    target.Clouds.High   = uint(source.Hourly.Cloud_cover_high[index])
    target.Visibility    = uint(toInt(source.Hourly.Visibility[index]))
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
        return "NW"
    case 8:
        return "N"
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

func initOverviewVals() OverviewVals {
    return OverviewVals{Min: 1e12, Max: -1e12, Median: 0, Count: 0}
}

func addVal(o *OverviewVals, v float32) {
    if o.Count == 0 {
        o.Min = v
        o.Max = v
    } else {
        if v < o.Min {
            o.Min = v
        }
        if v > o.Max {
            o.Max = v
        }
    }
    o.values = append(o.values, v)
    o.Count++
}

func median(values []float32) float32 {
    if len(values) == 0 {
        return 0
    }
    sorted := append([]float32(nil), values...)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
    mid := len(sorted) / 2
    if len(sorted)%2 == 0 {
        return (sorted[mid-1] + sorted[mid]) / 2
    }
    return sorted[mid]
}

func finalizeOverviewVals(o *OverviewVals) {
    if o.Count > 0 {
        o.Median = median(o.values)
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
