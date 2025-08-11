/*
 * API communication and data
 */

// Json Data
export let hourData = null;
export let weekData = null;

// variable
export let dayIndex = 0;

/* Makes a request for weather data and populates the table with the data
 * returns:  bool: ok
 */
export async function fetchWeatherData() {
    try {
        const request = `city`;
        const response = await fetch(request);
        if (!response.ok) throw new Error('Network response was not ok');
        weekData = await response.json();
    } catch (error) {
        return false;
    }
    return true;
}

/* Makes a request for weather data and populates the table with the data
 */
export async function fetchWeatherDetail() {
    try {
        const response = await fetch(`city/detail?day=${encodeURIComponent(dayIndex)}`);
        if (!response.ok) throw new Error('Network response was not ok');
        hourData = await response.json();
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
        return false
    }
    return true;
}

/* Setter for dayIndex
 */
export function setDayIndex(value) {
    dayIndex = value;
}
