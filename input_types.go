package main

import (
)

/*
 * Structures for parsing raw input JSON data
 */

// Range Data //

type InWeatherRange struct {
    timestamp uint          // timestamp for when data was received
    Cod       string        // TCP code (200 = ok)
    City struct {
        Name string
        Id   int
        Coord struct {
            Lat float32
            Lon float32
        }
        Country    string
        Population int
        Timezone   int
        Sunrise    uint
        Sunset     uint
    }
    List []InWeather
}

 // Weather Data //

type InWeather struct {
    Dt         uint         // POSIX time of forecast
    Visibility uint
    Pop        float32      // Probability for precipitation

    Main struct {
        Temp       float32
        Feels_like float32
        Sea_level  int      // Air pressure
        Grnd_level int      // Air pressure
        Humidity   float32
    }
    Weather []struct {
        Id          int
        Main        string  // One word description
        Description string  // Description
        Icon        string  // Icon ID for "openweathermap.org/img/wn/{ICON}@2x.png"
    }
    Clouds struct {
        All uint            // Cloud %
    }
    Wind struct {
        Speed float32
        Deg   int
        Gust  float32
    }
    Rain struct {
        Mm float32 `json:"3h"`
    }
    Snow struct {
        Mm float32 `json:"3h"`
    }
}

// City Object //

type InCity []struct {
    Name    string
    Lat     float32
    Lon     float32
    Country string
}
