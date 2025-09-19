/*
 * Miscellanious functions
 */

export const DAY = 24*60*60*1000;
export const SECOND = 1000;
const maxHistory = 5;

/* Maps the timeframe variable to string value
 */
export function str_tf(tf) {
    if (tf=="1") {
        return "Day";
    } else {
        return "Night";
    }
}

/* Converts str to Base64, via uint8
 */
export function base64(str) {
    const encoder = new TextEncoder();
    const utf8Bytes = encoder.encode(str);
    return btoa(String.fromCharCode(...utf8Bytes));
}

/* Converts Base64 to str, via uint8
 */
export function decode64(str) {
    return atob(str);
}


export function saveSearchEntry(city) {
    let json = loadHistory();
    if (json.cities.includes(city)) {
        return;
    }
    json.cities.push(city);
    if (json.cities.length > maxHistory) {
        json.cities = json.cities.slice(1)
    }
    saveHistory(json);
}

export function loadSearchHistory() {
    let json = loadHistory();
    return json.cities;
}

export function populateSearchHistory(element) {
    let json = loadHistory();
    populateSelect(element, json.cities, suggestionInsert);
}

/* Returns citys search history JSON
 */
function loadHistory() {
    let json = JSON.parse(localStorage.getItem("cities"));
    if (json === null) {
        return {cities: []};
    }
    return json;
}

/* Saves a boat to local JSON
 */
function saveHistory(json) {
    localStorage.setItem("cities", JSON.stringify(json));
}

/*
 * Populates a select (drop-down) element with options
 * element: Target HTML element <select>
 * json:    Data to use
 * func:    Function to process the given data
 */
function populateSelect(element, json, func) {
    for (const item in json) {
        func(element, json, item);
    }
}

/* Function used to customize uif.populateSelect() function
 */
function suggestionInsert(element, json, item) {
    const name = json[item];
    const option = document.createElement('option');
    option.value = decode64(name);
    element.appendChild(option);
}
