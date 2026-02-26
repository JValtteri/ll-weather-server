import * as cookie from "./cookie.js?id=MTc1NDk1NDIwMA"
import * as api from "./api.js?id=MTc1ODMxODUzNA"
import * as util from "./utils.js?id=MTc1ODMxODUzNA"
import * as table from "./table.js?id=MTc3MTc1NzgxMgr"

/*
 * UI elements amd high level UI functions
 */

// Inputs
export const cityInput     = document.getElementById('city-name');
export const cookieConsent = document.getElementById('accept');
export const modelInput    = document.getElementById('model');
export const citiesLst     = document.getElementById('cities');

// Tables
export const daysForecast  = document.getElementById("days-forecast");
export const hoursForecast = document.getElementById("hours-forecast");

// Titles
export const cityTitle     = document.getElementById("city-name-title");
export const dayTitle      = document.getElementById("day-title");

// Buttons
export const submitButton  = document.getElementById("submit-button");
export const dayButton     = document.getElementById("day-night");
export const fullscreenBtn = document.getElementById('fullscreen');
export const reloadBtn     = document.getElementById('reload');

const logo                 = document.getElementById('logo');

let cities = [];
let timeframe = 1;  // 1 = Day, 0 = Night


/* Sends sets the city name for weather request
 */
export async function submitCity() {
    if (cityInput.value && cityInput.value != "City") {
        let city = util.base64(cityInput.value);
        let model = util.base64(modelInput.value);
        cookie.prepCookies(cookieConsent.checked, city, model, cookie.ttl*util.DAY);
        const ok = await loadWeek();
        if ( ok ) {
            util.saveSearchEntry(city);
        }
    }
}

export async function loadWeek() {
    cityInput.value = util.decode64(cookie.getCookie("city"));
    logo.removeAttribute("hidden");
    const ok = await api.fetchWeatherData();
    if (ok) {
        activateUI();
    }
    logo.setAttribute("hidden", "");
    return ok;
}

export function toggleFullscreen() {
    if (!document.fullscreenElement) {
        makeFullscreen();
    } else {
        document.exitFullscreen();
        if (cookieConsent.checked === true) {
            localStorage.setItem("fullscreen", "")
        }
    }
}

export async function reloadForecast(){
    let city = cookie.getCookie("city");
    let model = util.base64(modelInput.value);
    cookie.prepCookies(cookieConsent.checked, city, model, cookie.ttl*util.DAY);
    loadWeek();
    table.populateTable(api.weekData.Days, daysForecast, util.str_tf(timeframe));
    if (!hoursForecast.hasAttribute("hidden")) {
        await api.fetchWeatherDetail();
        table.populateTable(api.hourData.Hours, hoursForecast, '');
    }
}

export function flipTime() {
    timeframe = (timeframe + 1) % 2     // Flip timeframe between 1 and 0 (day/night)
    table.populateTable(api.weekData.Days, daysForecast, util.str_tf(timeframe));
    if (timeframe == 1) {
        dayButton.textContent = "DAY";
        dayButton.classList.remove('night');
    } else {
        dayButton.textContent = "NIGHT";
        dayButton.classList.add('night');
    }
}

export function selectColumn(event) {
    if (event.target.tagName === "TH" || event.target.tagName === "TD") {
        var target = event.target;
    } else if (event.target.tagName === "IMG" || event.target.tagName === "TD") {
        var target = event.target.parentElement;
    } else if (event.target.tagName === "B") {
        return;
    }
    const columnIndex = Array.from(target.parentElement.cells).indexOf(target);
    if (columnIndex != 0) {
        api.setDayIndex(columnIndex-1);
        openWeatherDetail();
        hoursForecast.removeAttribute("hidden");
        dayTitle.removeAttribute("hidden");
        dayTitle.textContent = daysForecast.rows[0].cells[columnIndex].textContent
    }
}

export function loadPrefrences() {
    cookieConsent.checked = true;
    modelInput.value = util.decode64(cookie.getCookie("model"));
}

export function rememberMe() {
    if (cookieConsent.checked) {
        cookie.setCookie("city", util.base64(cityInput.value), cookie.ttl*util.DAY);
        localStorage.setItem("consent", "true");
    } else {
        cookie.setCookie("city", "");       // Cookie is set as session cookie, so the browser should remove it after the session
        localStorage.clear();
    }
}

export function fullscreenPrefrence() {
    if (cookieConsent.checked) {
        if (localStorage.getItem("fullscreen") === "true" ) {
            makeFullscreen();
        }
    }
}

export function populateSearchHistory() {
    citiesLst.innerHTML = '';
    util.populateSearchHistory(citiesLst);
}

export function updateSuggestions() {
    cities = util.loadSearchHistory();
}

export function suggest(value) {
    updateSuggestions();
    if (value != "") {
        citiesLst.innerHTML = '';
    }
    cities.forEach( city => {
        // Comparing if value is a fragment of a city
        if ((( util.decode64(city).toLowerCase()).indexOf(value.toLowerCase())) > -1) {
            let node = document.createElement("option");
            let val = document.createTextNode(util.decode64(city));
            node.appendChild(val);
            citiesLst.appendChild(node);
        }
    });
}

function activateUI() {
    dayButton.removeAttribute("disabled");      // Enables Day/Night button
    reloadBtn.removeAttribute("disabled");      // Enables Reload button
    dayTitle.setAttribute("hidden", "");        // Hide hour forecast title
    hoursForecast.setAttribute("hidden", "");   // Hide hour forecast table
    hoursForecast.innerHTML = "";               // Clear hour forecast table
    let lat = api.weekData.Coord.Lat.toFixed(2);
    let lon = api.weekData.Coord.Lon.toFixed(2);
    cityTitle.textContent = `${api.weekData.City}  –  ${lat}°, ${lon}°`;
    table.populateTable(api.weekData.Days, daysForecast, util.str_tf(timeframe));
}

async function openWeatherDetail() {
    let ok = await api.fetchWeatherDetail();
    if (ok) {
        table.populateTable(api.hourData.Hours, hoursForecast, '');
    }
}

function makeFullscreen() {
    document.querySelector("body").requestFullscreen()
                                  .catch((TypeError) => {});
    if (cookieConsent.checked === true) {
        localStorage.setItem("fullscreen", "true")
    }
}
