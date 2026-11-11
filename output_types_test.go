package main

import (
    "testing"

    "github.com/JValtteri/weather/om"
)

func testRange(hours int) om.WeatherRange {
    var r om.WeatherRange
    r.Hourly.Time = make([]uint, hours)
    return r
}

func testRangeWithForecast(timeHours int, forecastHours int) om.WeatherRange {
    r := testRange(timeHours)
    r.Hourly.Surface_pressure = make([]float32, timeHours)
    for i := 0; i < forecastHours && i < timeHours; i++ {
        r.Hourly.Surface_pressure[i] = 1013
    }
    return r
}

func testFullRange(timeHours int, forecastHours int, temp float32) om.WeatherRange {
    r := testRangeWithForecast(timeHours, forecastHours)
    h := &r.Hourly
    h.Temperature_2m = make([]float32, timeHours)
    h.Apparent_temperature = make([]float32, timeHours)
    h.Wet_bulb_temperature_2m = make([]float32, timeHours)
    h.Uv_index = make([]float32, timeHours)
    h.Direct_radiation = make([]float32, timeHours)
    h.Diffuse_radiation = make([]float32, timeHours)
    h.Cloud_cover = make([]int, timeHours)
    h.Cloud_cover_low = make([]int, timeHours)
    h.Cloud_cover_mid = make([]int, timeHours)
    h.Cloud_cover_high = make([]int, timeHours)
    h.Pressure_msl = make([]float32, timeHours)
    h.Relative_humidity_2m = make([]int, timeHours)
    h.Visibility = make([]float32, timeHours)
    h.Precipitation_probability = make([]int, timeHours)
    h.Precipitation = make([]float32, timeHours)
    h.Wind_speed_10m = make([]float32, timeHours)
    h.Wind_gusts_10m = make([]float32, timeHours)
    h.Wind_direction_10m = make([]int, timeHours)
    h.Weather_code = make([]int, timeHours)
    h.Is_day = make([]int, timeHours)
    for i := 0; i < timeHours; i++ {
        h.Relative_humidity_2m[i] = 50
        h.Is_day[i] = 1
        if i < forecastHours {
            h.Pressure_msl[i] = 1013
            h.Temperature_2m[i] = temp
        }
    }
    return r
}

func TestForecastHours(t *testing.T) {
    tests := []struct {
        name string
        r    om.WeatherRange
        want int
    }{
        {"time array only", testRange(72), 72},
        {"zero padded forecast", testRangeWithForecast(72, 48), 48},
        {"fully missing forecast", testRangeWithForecast(72, 0), 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := forecastHours(tt.r); got != tt.want {
                t.Errorf("Expected %d hours, got %d", tt.want, got)
            }
        })
    }

    // Surface pressure can be missing while pressure_msl is still available.
    r := testRange(72)
    r.Hourly.Surface_pressure = make([]float32, 72)
    r.Hourly.Pressure_msl = make([]float32, 72)
    for i := 0; i < 60; i++ {
        r.Hourly.Pressure_msl[i] = 1013
    }
    if got := forecastHours(r); got != 60 {
        t.Errorf("Expected 60 hours from pressure_msl, got %d", got)
    }
}

func TestOverviewSourceLength(t *testing.T) {
    tests := []struct {
        name       string
        ranges     []om.WeatherRange
        wantModels int
        wantLength int
    }{
        {
            "missing forecast is ignored",
            []om.WeatherRange{testRange(72), testRange(0)},
            1, 72,
        },
        {
            "forecast shorter than two days is ignored",
            []om.WeatherRange{testRange(47), testRange(72)},
            1, 72,
        },
        {
            "forecast of exactly two days is included",
            []om.WeatherRange{testRange(48), testRange(72)},
            2, 48,
        },
        {
            "zero padded short forecast is ignored",
            []om.WeatherRange{testRangeWithForecast(72, 72), testRangeWithForecast(72, 24)},
            1, 72,
        },
        {
            "zero padded long forecast is included",
            []om.WeatherRange{testRangeWithForecast(72, 72), testRangeWithForecast(72, 48)},
            2, 48,
        },
        {
            "fully missing forecast is ignored",
            []om.WeatherRange{testRangeWithForecast(72, 72), testRangeWithForecast(72, 0)},
            1, 72,
        },
        {
            "output matches the shortest forecast",
            []om.WeatherRange{testRange(72), testRange(96), testRange(120)},
            3, 72,
        },
        {
            "output matches the shortest valid data",
            []om.WeatherRange{
                testRangeWithForecast(120, 120),
                testRangeWithForecast(96, 96),
                testRangeWithForecast(72, 72),
            },
            3, 72,
        },
        {
            "all forecasts insufficient",
            []om.WeatherRange{testRange(0), testRange(24), testRange(47)},
            0, 0,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            source := overviewSource(tt.ranges)
            if len(source.ranges) != tt.wantModels {
                t.Errorf("Expected %d models, got %d", tt.wantModels, len(source.ranges))
            }
            if source.length() != tt.wantLength {
                t.Errorf("Expected length %d, got %d", tt.wantLength, source.length())
            }
        })
    }
}

func TestSingleModelSourceLength(t *testing.T) {
    source := singleModelSource(testRange(1))
    if source.empty() {
        t.Fatal("Expected a single model source to be non-empty")
    }
    if source.length() != 1 {
        t.Errorf("Expected length 1, got %d", source.length())
    }
}

func TestAggregateIgnoresInsufficientModels(t *testing.T) {
    long := testFullRange(72, 72, 10)
    short := testFullRange(72, 24, -50)
    missing := testFullRange(72, 0, -100)

    wd := aggregateOmAtIndex([]om.WeatherRange{long, short, missing}, 0, "00:00")
    if wd.Overview == nil {
        t.Fatal("Expected overview values")
    }
    if wd.Overview.Temp.Min != 10 || wd.Overview.Temp.Max != 10 {
        t.Errorf("Expected temp range 10..10, got %v..%v",
            wd.Overview.Temp.Min, wd.Overview.Temp.Max)
    }
    if wd.Temp.Temp != 10 {
        t.Errorf("Expected temp 10, got %v", wd.Temp.Temp)
    }
}
