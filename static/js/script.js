import * as table from "./table.js?id=MTc1NDY0OTQ4Ng=="
import * as cookie from "./cookie.js?id=MTc1NDY0OTQ4Ng=="

const body          = document.getElementById('body');
// Inputs
const cityInput     = document.getElementById('city-name');
const cookieConsent = document.getElementById('accept');
const modelInput    = document.getElementById('model');
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
let dayIndex = 0;

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

/* Converts Base64 to str, via uint8
 */
function decode64(str) {
    return atob(str);
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str:  Name of the city to search
 * returns:  bool: ok
 */
async function fetchWeatherData() {
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

function activateUI() {
    dayButton.removeAttribute("disabled");      // Enables Day/Night button
    reloadBtn.removeAttribute("disabled");      // Enables Reload button
    dayTitle.setAttribute("hidden", "");        // Hide hour forecast title
    hoursForecast.setAttribute("hidden", "");   // Hide hour forecast table
    hoursForecast.innerHTML = "";               // Clear hour forecast table
    let lat = weekData.Coord.Lat.toFixed(2);
    let lon = weekData.Coord.Lon.toFixed(2);
    cityTitle.textContent = `${weekData.City}  –  ${lat}°, ${lon}°`;
    table.populateTable(weekData.Days, daysForecast, str_tf());
}

/* Sends sets the city name for weather request
 */
async function submitCity() {
    if (cityInput.value && cityInput.value != "City") {
        prepCookies();
        let ok = await fetchWeatherData();
        if (ok === true) {
            activateUI();
        }
    }
}

function prepCookies() {
    if (cookieConsent.checked) {
        cookie.setCookie("city", base64(cityInput.value), ttl*DAY);
        cookie.setCookie("model", base64(modelInput.value), ttl*DAY);
    } else {
        cookie.setCookie("city", base64(cityInput.value));   // Cookie expires in one second
        cookie.setCookie("model", base64(modelInput.value));
    }
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str: Name of the city to search
 * dayIndex: int: Index of the day requested
 */
async function fetchWeatherDetail() {
    try {
        const response = await fetch(`city/detail?day=${encodeURIComponent(dayIndex)}`);
        if (!response.ok) throw new Error('Network response was not ok');
        hourData = await response.json();
        table.populateTable(hourData.Hours, hoursForecast, '');
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

function toggleFullscreen() {
    if (!document.fullscreenElement) {
        makeFullscreen();
    } else {
        document.exitFullscreen();
        if (cookieConsent.checked === true) {
            localStorage.setItem("fullscreen", "")
        }
    }
}

function makeFullscreen() {
    document.querySelector("body").requestFullscreen()
                                  .catch((TypeError) => {});
    if (cookieConsent.checked === true) {
        localStorage.setItem("fullscreen", "true")
    }
}

async function reloadForecast(){
    prepCookies()
    await fetchWeatherData();
    table.populateTable(weekData.Days, daysForecast, str_tf());
    if (!hoursForecast.hasAttribute("hidden")) {
        await fetchWeatherDetail();
    }
}

/*
 * Buttons and Events
 */

// Buttons
const submitButton  = document.getElementById("submit-button");
const dayButton     = document.getElementById("day-night");
const fullscreenBtn = document.getElementById('fullscreen');
const reloadBtn     = document.getElementById('reload');

/* Submit button
 */
submitButton.addEventListener("click", () => {
    submitCity();
});

/* Reload button
 */
reloadBtn.addEventListener("click", () => {
    reloadForecast();
});

/* Day/Night button
 */
dayButton.addEventListener("click", () => {
    timeframe = (timeframe + 1) % 2     // Flip timeframe between 1 and 0 (day/night)
    table.populateTable(weekData.Days, daysForecast, str_tf());
    if (timeframe == 1) {
        dayButton.textContent = "DAY";
        dayButton.classList.remove('night');
    } else {
        dayButton.textContent = "NIGHT";
        dayButton.classList.add('night');
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

/* Weather model changed
 */
modelInput.addEventListener("change", () => {
    reloadForecast();
})

/* "Remember Me" clicked
 */
cookieConsent.addEventListener('click', () => {
    if (cookieConsent.checked) {
        cookie.setCookie("city", base64(cityInput.value), ttl*DAY);
        localStorage.setItem("consent", "true");
    } else {
        cookie.setCookie("city", "");       // Cookie is set as session cookie, so the browser should remove it after the session
        localStorage.clear();
    }
});

/* "Fullscreen" clicked
 */
fullscreenBtn.addEventListener("click", () => {
    toggleFullscreen();
})

/* Clicked anywhere on document
 */
body.addEventListener("click", () => {
     if (cookieConsent.checked) {
        if (localStorage.getItem("fullscreen") === "true" ) {
            makeFullscreen();
        }
    }
});

/* Click on a day
 * Opens detail view
 */
daysForecast.addEventListener("click", function(event) {
    if (event.target.tagName === "TH" || event.target.tagName === "TD") {
        var target = event.target;
    } else if (event.target.tagName === "IMG" || event.target.tagName === "TD") {
        var target = event.target.parentElement;
    } else if (event.target.tagName === "B") {
        return;
    }
    const columnIndex = Array.from(target.parentElement.cells).indexOf(target);
    if (columnIndex != 0) {
        dayIndex = columnIndex-1
        fetchWeatherDetail();
        hoursForecast.removeAttribute("hidden");
        dayTitle.removeAttribute("hidden");
        dayTitle.textContent = daysForecast.rows[0].cells[columnIndex].textContent
    }
});

// Check concent from cookie
if (localStorage.getItem("consent") === "true" ) {
    cookieConsent.checked = true
    modelInput.value = decode64(cookie.getCookie("model"))
}

// Auto search if concented
if (cookieConsent.checked) {
    if (localStorage.getItem("fullscreen") === "true" ) {
        makeFullscreen();
    }
    cookieConsent.checked = true;
    let ok = await fetchWeatherData();
    if (ok === true) {
        activateUI();
    }
}
