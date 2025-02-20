package main

/*
 * Structures for returning JSON data to front end
 */

import (
)

type WeekWeather struct {
    City string
    Days []DayWeather
}

type DayWeather struct {
    Day int
    Time int
    Temp_day float32
    Temp_night float32
    Description string
    Humidity float32
    Clouds float32
    Visibility int
    Windspeed float32
    IconID string
}
