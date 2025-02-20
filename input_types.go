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
    List []InWeatherDay
}

 // Weather Data //

type InWeatherDay struct {
    Dt float64
    Visibility int
    Main struct {
        Temp float32
        Humidity float32
        Icon string
    }
    Weather []struct {
        Main string
        Description string
        Icon string
    }
    Clouds struct {
        All int
    }
}

