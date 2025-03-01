import * as table from "./table.js"
import * as cookie from "./cookie.js"

// Inputs
const cityInput     = document.getElementById('city-name');
const cookieConsent = document.getElementById('accept');
// Tables
const daysForecast  = document.getElementById("days-forecast");
const hoursForecast = document.getElementById("hours-forecast");
// Titles
const cityTitle     = document.getElementById("city-name-title");
const dayTitle      = document.getElementById("day-title");
// Json Data
let hourData = null;
let weekData = null;
//
let timeframe = 1;  // 1 = Day, 0 = Night

const DAY = 24*60*60*1000;
const SECOND = 1000;

const ttl = 30;     // cookie max life

/* Maps the timeframe variable to string value
 */
function str_tf() {
    if (timeframe=="1") {
        return "Day";
    } else {
        return "Night";
    }
}

/* Converts str to Base64, via uint8
*/
function base64(str) {
    const encoder = new TextEncoder();
    const utf8Bytes = encoder.encode(str);
    return btoa(String.fromCharCode(...utf8Bytes));
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str: Name of the city to search
 */
async function fetchWeatherData() {
    try {
        const request = `city`;

        const response = await fetch(request);
        if (!response.ok) throw new Error('Network response was not ok');
        weekData = await response.json();
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
    dayButton.removeAttribute("disabled");      // Enables Day/Night button
    dayTitle.setAttribute("hidden", "");        // Hide hour forecast title
    hoursForecast.setAttribute("hidden", "");   // Hide hour forecast table
    hoursForecast.innerHTML = "";               // Clear hour forecast table
    cityTitle.textContent = weekData.City;
    table.populateTable(weekData.Days, daysForecast, str_tf());
}

/* Sends sets the city name for weather request
 */
function submitCity() {
    if (cityInput.value && cityInput.value != "City") {
        if (cookieConsent.checked) {
            cookie.setCookie("city", base64(cityInput.value), ttl*DAY);
        } else {
            cookie.setCookie("city", base64(cityInput.value));   // Cookie expires in one second
        }
        fetchWeatherData();
    }
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str: Name of the city to search
 * dayIndex: int: Index of the day requested
 */
async function fetchWeatherDetail(cityName, dayIndex) {
    try {
        const response = await fetch(`city/detail?day=${encodeURIComponent(dayIndex)}`);
        if (!response.ok) throw new Error('Network response was not ok');
        hourData = await response.json();
        table.populateTable(hourData.Hours, hoursForecast, '');
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

function makeFullscreen() {
    document.querySelector("body").requestFullscreen();
}

/*
 * Buttons and Events
 */

// Buttons
const submitButton  = document.getElementById("submit-button");
const dayButton     = document.getElementById("day-night");
const fullscreenBtn = document.getElementById('fullscreen');

/* Submit button
 */
submitButton.addEventListener("click", () => {
    submitCity();
});

/* Day/Night button
 */
dayButton.addEventListener("click", () => {
    timeframe = (timeframe + 1) % 2     // Flip timeframe between 1 and 0 (day/night)
    table.populateTable(weekData.Days, daysForecast, str_tf());
    if (timeframe == 1) {
        dayButton.textContent = "DAY";
    } else {
        dayButton.textContent = "NIGHT";
    }
});

/* Clear input on click
 */
cityInput.addEventListener("click", () => {
    cityInput.value = "";
});

/* Search on Enter key
 */
cityInput.addEventListener('keydown', (event) => {
    if (event.key === "Enter") {
        submitCity();
    }
});

/* "Remember Me" clicked
 */
cookieConsent.addEventListener('click', () => {
    if (cookieConsent.checked) {
        cookie.setCookie("city", base64(cityInput.value), ttl*DAY);
        cookie.setCookie("consent", "true", ttl*DAY);
    } else {
        cookie.setCookie("city", "");       // Cookie is set as session cookie, so the browser should remove it after the session
        cookie.setCookie("consent", "");
    }
});

/* Set to fullscreen
 */
fullscreenBtn.addEventListener("click", () => {
    makeFullscreen();
})

/* Click on a day
 * Opens detail view
 */
daysForecast.addEventListener("click", function(event) {
    if (event.target.tagName === "TH" || event.target.tagName === "TD") {
        const columnIndex = Array.from(event.target.parentElement.cells).indexOf(event.target);
        fetchWeatherDetail(cityTitle.textContent, columnIndex-1);
        hoursForecast.removeAttribute("hidden");
        dayTitle.removeAttribute("hidden");
        dayTitle.textContent = daysForecast.rows[0].cells[columnIndex].textContent
    }
});

// Check concent from cookie
if (cookie.getCookie("consent") === "true" ) {
    cookieConsent.checked = true
}

// Auto search if concented
if (cookieConsent.checked) {
    cookieConsent.checked = true;
    fetchWeatherData();
}
