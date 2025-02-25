
const submitButton  = document.getElementById("submit-button");
const dayButton     = document.getElementById("day-night");
const daysForecast  = document.getElementById("days-forecast");
const cityInput     = document.getElementById('city-name');
const cityTitle     = document.getElementById("city-name-title");
const fullscreenBtn = document.getElementById('fullscreen');

let data = null;
let timeframe = 1;  // 1 = Day, 0 = Night

function str_tf() {
    if (timeframe=="1") {
        return "Day";
    } else {
        return "Night";
    }
}

function populateRow(days, table, elementType, parameters, rowTitle) {
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
        columnElm.textContent = dataTarget;
        rowElm.appendChild(columnElm);
    });
}

// Create Image row
function populateImageRow(days, table, parameter, rowTitle) {
    // Create new row
    let rowElm = document.createElement('tr');
    table.appendChild(rowElm);
    // Title cells
    let titleElm = document.createElement('td');
    titleElm.textContent = rowTitle;
    rowElm.appendChild(titleElm);
    // Data cells
    let columnElm;
    days.forEach(day => {
        let dataTarget = day;
        columnElm = document.createElement('td');
        let imgElm = document.createElement('img');
        imgElm.src = "img/"+dataTarget[parameter].IconID
        imgElm.alt = dataTarget[parameter].Description;
        columnElm.appendChild(imgElm);
        rowElm.appendChild(columnElm);
    });
}

function populateTable(days) {
    // Create rows and titles
    daysForecast.innerHTML = "";
    populateRow(days, daysForecast, 'th', ['DayName'],                  "");          //Day name
    populateImageRow(days, daysForecast, str_tf(), "");                                         // Icons
    populateRow(days, daysForecast, 'td', [str_tf(), 'Description'],      "Desc.");
    populateRow(days, daysForecast, 'td', [str_tf(), 'Temp'],             "Temp.");
    populateRow(days, daysForecast, 'td', ['RainChance'],                 "Rain%");     // Chance
    populateRow(days, daysForecast, 'td', ['RainTotal'],                  "Rain [mm]"); // Total
    populateRow(days, daysForecast, 'td', [str_tf(), 'Clouds','Clouds'],  "Clouds %");  // Total
                                                                                      // Layers
    populateRow(days, daysForecast, 'td', [str_tf(), 'Wind','Speed'],     "Wind");
    populateRow(days, daysForecast, 'td', [str_tf(), 'Wind','Dir'],       "Direction");
    populateRow(days, daysForecast, 'td', [str_tf(), 'Pressure'],         "Pressure");
    populateRow(days, daysForecast, 'td', [str_tf(), 'Humidity'],         "Humidity");
    //populateRow(days, daysForecast, 'td', [str_tf(), 'Visibility'],       "Visibility");
}

async function fetchWeatherData(cityName) {
    try {
        const response = await fetch(`city?name=${encodeURIComponent(cityName)}`);
        console.log("Requested:", cityName)
        if (!response.ok) throw new Error('Network response was not ok');
        data = await response.json();
        // Enable Night button
        dayButton.removeAttribute("disabled");
        cityTitle.textContent = data.City;
        populateTable(data.Days, timeframe);
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
