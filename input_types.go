package main

import (
)

/*
 * Structures for parsing raw input JSON data
 */

// Range Data //

type InWeatherRange struct {
    Cod string
    City struct {
        Name string
        Id int
        Coord struct {
            Lat float32
            Lon float32
        }
    }
    List []InWeather
}

 // Weather Data //

type InWeather struct {
    Dt uint
    Visibility uint
    Main struct {
        Temp float32
        Humidity float32
    }
    Weather []struct {
        Main string
        Description string
        Icon string
    }
    Clouds struct {
        All uint
    }
    Wind struct {
        Speed float32
        Deg int
    }
}

