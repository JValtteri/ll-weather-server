
const submitButton  = document.getElementById("submit-button");
const dayButton     = document.getElementById("day-night");
const daysForecast  = document.getElementById("days-forecast");
const cityInput     = document.getElementById('city-name');
const cityTitle     = document.getElementById("city-name-title");
const fullscreenBtn = document.getElementById('fullscreen');

let data = null;
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

function addNum(element, num, unit='') {
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
function populateTable(days, table) {
    // Create rows and titles
    table.innerHTML = "";
    populateRow(days, table, 'th', ['DayName'],                    "",          addText, ''); //Day name
    populateRow(days, table, 'td', [str_tf()],                     "",          addImage, '');// Icons
    populateRow(days, table, 'td', [str_tf(), 'Description'],      "Desc.",     addText, '');
    populateRow(days, table, 'td', [str_tf(), 'Temp'],             "Temp.",     addNum, '°C');
    populateRow(days, table, 'td', ['RainChance'],                 "Rain%",     addNum, ' %'); // Chance
    populateRow(days, table, 'td', ['RainTotal'],                  "Rain [mm]", addNum, ' mm'); // Total
    populateRow(days, table, 'td', [str_tf(), 'Clouds','Clouds'],  "Clouds %",  addNum, ' %'); // Total
                                                                                      // Layers
    populateRow(days, table, 'td', [str_tf(), 'Wind','Speed'],     "Wind",      addNum, ' m/s');
    populateRow(days, table, 'td', [str_tf(), 'Wind','Dir'],       "Direction", addText, '');
    populateRow(days, table, 'td', [str_tf(), 'Pressure'],         "Pressure",  addNum, ' hPa');
    populateRow(days, table, 'td', [str_tf(), 'Humidity'],         "Humidity",  addNum, ' %');
    //populateRow(days, table, 'td', [str_tf(), 'Visibility'],      "Visibility", addText);
}

/* Makes a request for weather data and populates the table with the data
 * cityName: str: Name of the city to search
 */
async function fetchWeatherData(cityName) {
    try {
        const response = await fetch(`city?name=${encodeURIComponent(cityName)}`);
        console.log("Requested:", cityName)
        if (!response.ok) throw new Error('Network response was not ok');
        data = await response.json();
        // Enable Night button
        dayButton.removeAttribute("disabled");
        cityTitle.textContent = data.City;
        populateTable(data.Days, daysForecast);
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

function submitCity() {
    if (cityInput.value && cityInput.value != "City") {
        fetchWeatherData(cityInput.value);
    }
}

function makeFullscreen() {
    document.querySelector("body").requestFullscreen();
}

submitButton.addEventListener("click", () => {
    console.log("Clikked");
    submitCity();
});

dayButton.addEventListener("click", () => {
    console.log("Clikked");
    timeframe = (timeframe + 1) % 2     // Flip timeframe between 1 and 0 (day/night)
    populateTable(data.Days);
    if (timeframe == 1) {
        dayButton.textContent = "DAY";
    } else {
        dayButton.textContent = "NIGHT";
    }
});

cityInput.addEventListener("click", () => {
    cityInput.value = "";
});

cityInput.addEventListener('keydown', (event) => {
    if (event.key === "Enter") {
        submitCity();
    }
});

fullscreenBtn.addEventListener("click", () => {
    makeFullscreen();
})
