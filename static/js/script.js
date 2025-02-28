
const submitButton  = document.getElementById("submit-button");
const dayButton     = document.getElementById("day-night");
const daysForecast  = document.getElementById("days-forecast");
const hoursForecast = document.getElementById("hours-forecast");
const cityInput     = document.getElementById('city-name');
const cityTitle     = document.getElementById("city-name-title");
const dayTitle      = document.getElementById("day-title");
const fullscreenBtn = document.getElementById('fullscreen');

let hourData = null;
let weekData = null;
let timeframe = 1;  // 1 = Day, 0 = Night

/* Maps the timeframe variable to string value
 */
function str_tf() {
    if (timeframe=="1") {
        return "Day";
    } else {
        return "Night";
    }
}

/* Function to add text content to a cell
 * To be used as a parameter for populateRow()
 */
function addText(element, text, unit='') {
    element.textContent = text;
    return element;
}

/* Function to add number content to a cell
 * Rounds numbers to integers. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addNum(element, num=0, unit='') {
    element.textContent = `${num.toFixed(0)}${unit}`;
    return element;
}

/* Function to add an image element from data.day
 * To be used as a parameter for populateRow()
 */
function addImage(element, dataTarget, _='') {
    let imgElm = document.createElement('img');
    imgElm.src = "img/"+dataTarget.IconID
    imgElm.alt = dataTarget.Description;
    element.appendChild(imgElm);
    return element;
}

/* Creates a new row <tr>
 * Inputs:
 * days:        []obj:    'day' objects
 * table:       element:  target table element
 * elementType: str:      the type of html element to add i.e 'tr', 'td'
 * parameters:  []str:    containing the path from day object to chosen field
 * rowTitle:    str:      A title for the row (placed in first column)
 * func:        function: a function used to add content to the created cell
 */
function populateRow(days, table, elementType, parameters, rowTitle, func, unit='') {
    if (parameters[0] === '') {
        parameters.splice(0, 1);
    }
    // Create new row
    let rowElm = document.createElement('tr');
    table.appendChild(rowElm);
    // Title cells
    let titleElm = document.createElement(elementType);
    titleElm.textContent = rowTitle;
    rowElm.appendChild(titleElm);
    // Data cells
    let columnElm;
    days.forEach(day => {
        let dataTarget = day;
        columnElm = document.createElement(elementType);
        for (let i = 0; i < parameters.length; ++i) {
            dataTarget = dataTarget[parameters[i]];
        }
        element = func(columnElm, dataTarget, unit);
        rowElm.appendChild(element);
    });
}

/* Creates all new rows <tr>, by calling populateRow()
 * Inputs:
 * days:        []obj:    'day' objects
 * table:       element:  target table element
 * elementType: str:      the type of html element to add i.e 'tr', 'td'
 * parameters:  []str:    containing the path from day object to chosen field
 * rowTitle:    str:      A title for the row (placed in first column)
 * func:        function: a function used to add content to the created cell
 */
function populateTable(days, table, prefix) {
    // Create rows and titles
    table.innerHTML = "";
    populateRow(days, table, 'th', [prefix, 'Title'],            "",          addText, ''); //Day name
    populateRow(days, table, 'td', [prefix],                     "",          addImage, '');// Icons
    populateRow(days, table, 'td', [prefix, 'Description'],      "Desc.",     addText, '');
    populateRow(days, table, 'td', [prefix, 'Temp'],             "Temp.",     addNum, '°C');
    populateRow(days, table, 'td', [prefix, 'Clouds','Clouds'],  "Clouds %",  addNum, ' %'); // Total
    populateRow(days, table, 'td', [prefix, 'Rain', 'Chance'],   "Rain%",     addNum, ' %'); // Chance
    populateRow(days, table, 'td', [prefix, 'Rain', 'Amount'],   "Rain [mm]", addNum, ' mm'); // Total
                                                                                      // Layers
    populateRow(days, table, 'td', [prefix, 'Wind','Speed'],     "Wind",      addNum, ' m/s');
    populateRow(days, table, 'td', [prefix, 'Wind','Dir'],       "Direction", addText, '');
    //populateRow(days, table, 'td', [prefix, 'Pressure'],         "Pressure",  addNum, ' hPa');
    //populateRow(days, table, 'td', [prefix, 'Humidity'],         "Humidity",  addNum, ' %');
    //populateRow(days, table, 'td', [prefix, 'Visibility'],      "Visibility", addText);
    applyColors(table);
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str: Name of the city to search
 */
async function fetchWeatherData(cityName) {
    try {
        const response = await fetch(`city?name=${encodeURIComponent(cityName)}`);
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
    populateTable(weekData.Days, daysForecast, str_tf());
}

function submitCity() {
    if (cityInput.value && cityInput.value != "City") {
        fetchWeatherData(cityInput.value);
    }
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str: Name of the city to search
 */
async function fetchWeatherDetail(cityName, dayIndex) {
    try {
        const response = await fetch(`city/detail?name=${encodeURIComponent(cityName)}&day=${encodeURIComponent(dayIndex)}`);
        console.log("Requested:", cityName)
        if (!response.ok) throw new Error('Network response was not ok');
        hourData = await response.json();
        populateTable(hourData.Hours, hoursForecast, '');
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

function makeFullscreen() {
    document.querySelector("body").requestFullscreen();
}

/*
 * Color highlights
 */


function colorTemp(element) {
    let value = parseInt(element.textContent.split('°')[0]);
    if (value > 20) {
        element.classList.add('hot')
    } else if (value < -20) {
        element.classList.add('arctic')
    }
}

function colorCloud(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 50) {
        element.classList.add('broken-clouds')
    } else if (value > 10) {
        element.classList.add('clear-sky')
    }
}

function colorRainChance(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 70) {
        element.classList.add('medium')
    } else if (value > 19) {
        element.classList.add('light')
    }
}

function colorRain(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 3) {
        element.classList.add('heavy')
    } else if (value > 1) {
        element.classList.add('medium')
    } else if (value > 0) {
        element.classList.add('light')
    }
}

function colorWind(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value >= 15) {
        element.classList.add('heavy')
    } else if (value > 9) {
        element.classList.add('medium')
    } else if (value >= 4) {
        element.classList.add('light')
    }
}

/* Apply color to a row.
 * table: target table
 * row: index of the target row
 * func: function to run on each cell
 */
function applyColorRow(table, row, func) {
    let target = table.rows[row];
    target.childNodes.forEach(cell => {
        func(cell);
    });
}

/*
 * Applies highlight colors to an entire table
 */
function applyColors(table) {
    applyColorRow(table, 3, colorTemp);
    applyColorRow(table, 4, colorCloud);
    applyColorRow(table, 5, colorRainChance);
    applyColorRow(table, 6, colorRain);
    applyColorRow(table, 7, colorWind);
}

/*
 * Buttons
 */

/* Submit button
 */
submitButton.addEventListener("click", () => {
    submitCity();
});

/* Day/Night button
 */
dayButton.addEventListener("click", () => {
    timeframe = (timeframe + 1) % 2     // Flip timeframe between 1 and 0 (day/night)
    populateTable(weekData.Days, daysForecast, str_tf());
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

/* Set to fullscreen
 */
fullscreenBtn.addEventListener("click", () => {
    makeFullscreen();
})

/* Handle a click on a day
 * Opens detail view
 */
daysForecast.addEventListener("click", function(event) {
    if (event.target.tagName === "TH" || event.target.tagName === "TD") {
        const columnIndex = Array.from(event.target.parentElement.cells).indexOf(event.target);
        console.log(`Column clicked: ${columnIndex}`);
        fetchWeatherDetail(cityTitle.textContent, columnIndex-1);
        hoursForecast.removeAttribute("hidden");
        dayTitle.removeAttribute("hidden");
        dayTitle.textContent = daysForecast.rows[0].cells[columnIndex].textContent
    }
  });
