
const submitButton  = document.getElementById("submit-button");
const dayButton     = document.getElementById("day-night");
const daysForecast  = document.getElementById("days-forecast");
let data = null;
let timeframe = 1;  // 1 = Day, 0 = Night

function populateTable(days) {
    let day_index = 1;
    let target = null;
    days.forEach(day => {
        if (timeframe=="1") {
            target = day.Day
        } else {
            target = day.Night
        }
        daysForecast.rows[1].cells[day_index].textContent = target.Description;
        daysForecast.rows[2].cells[day_index].textContent = target.Temp;
        daysForecast.rows[3].cells[day_index].textContent = target.Humidity;
        daysForecast.rows[4].cells[day_index].textContent = target.Clouds;
        day_index += 1;
    });
}

async function fetchWeatherData(cityName) {
    try {
        const response = await fetch(`http://localhost:3000/city?name=${encodeURIComponent(cityName)}`);
        console.log("Requested:", cityName)
        if (!response.ok) throw new Error('Network response was not ok');
        data = await response.json();
        // Enable Night button
        dayButton.removeAttribute("disabled");
        populateTable(data.Days, timeframe);
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error);
    }
}

submitButton.addEventListener("click", () => {
    console.log("Clikked");
    const cityInput = document.getElementById('city-name');
    if (cityInput.value) {
        fetchWeatherData(cityInput.value);
    }
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


