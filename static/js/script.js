import * as ui from "./ui.js?id=MTc1ODMxODUzNA"

const body          = document.getElementById('body');

/*
 * Events
 */


/* Writing on city input
 */
ui.cityInput.addEventListener("keyup", (event) => {
    // Ignore special keys
    if (event.key.length === 1) {
        ui.suggest(ui.cityInput.value);
    }
});

/* Submit button
 */
ui.submitButton.addEventListener("click", () => {
    ui.submitCity();
});

/* Reload button
 */
ui.reloadBtn.addEventListener("click", () => {
    ui.reloadForecast();
});

/* Day/Night button
 */
ui.dayButton.addEventListener("click", () => {
    ui.flipTime();
});

/* Clear input on click
 */
ui.cityInput.addEventListener("click", () => {
    ui.cityInput.value = "";
    ui.populateSearchHistory();
});

/* Search on Enter key
 */
ui.cityInput.addEventListener('keydown', (event) => {
    if (event.key === "Enter") {
        ui.submitCity();
    }
});

/* Weather model changed
 */
ui.modelInput.addEventListener("change", () => {
    ui.reloadForecast();
})

/* "Remember Me" clicked
 */
ui.cookieConsent.addEventListener('click', () => {
    ui.rememberMe();
});

/* "Fullscreen" clicked
 */
ui.fullscreenBtn.addEventListener("click", () => {
    ui.toggleFullscreen();
})

/* Clicked anywhere on document
 */
body.addEventListener("click", () => {
    ui.fullscreenPrefrence();
});

/* Click on a day
 * Opens detail view
 */
ui.daysForecast.addEventListener("click", function(event) {
    ui.selectColumn(event);
});

// Check concent from cookie
if (localStorage.getItem("consent") === "true" ) {
    ui.loadPrefrences();
}

// Auto search if concented
if (ui.cookieConsent.checked) {
    if (localStorage.getItem("fullscreen") === "true" ) {
        ui.makeFullscreen();
    }
    ui.cookieConsent.checked = true;
    ui.loadWeek();
    ui.updateSuggestions();
}
