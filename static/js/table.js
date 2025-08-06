
import * as color from "./color.js";


/* Function to construct a path to a specific element in JSON
 * dataTarget: obj:    'day' object
 * path:       []str:  path to element expressed as a str list
 */
function constructJsonPath(dataTarget, path) {
    for (let i = 0; i < path.length; ++i) {
        dataTarget = dataTarget[path[i]];
    }
    return dataTarget;
}

/* Function to add text content to a cell
 * To be used as a parameter for populateRow()
 */
function addText(element, day, path, _='', _0=0) {
    let text = constructJsonPath(day, path);
    element.textContent = text;
    return element;
}

/* Function to day/night indicator to title
 * To be used as a parameter for populateRow()
 */
function addSun(element, day, path, _='', _0=0) {
    let text = constructJsonPath(day, path);
    let sunPath = [];
    let sun = null;
    if (path.length > 1) {
        sunPath = [path[0]];
        sun = constructJsonPath(day, sunPath);
    } else {
        sun = day;
    }
    element.textContent = text;
    color.colorSun(element, sun.SunUp);
    return element;
}


/* Function to add number content to a cell
 * Rounds numbers to integers. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addNum(element, day, path, unit='', decimals=0) {
    let num = constructJsonPath(day, path);
    if (!num) {
        num = 0;
    }
    element.textContent = `${num.toFixed(decimals)}${unit}`;
    return element;
}

/* Function to add an image element from data.day
 * To be used as a parameter for populateRow()
 */
function addImage(element, day, path, _='', _0=0) {
    let image = constructJsonPath(day, path);
    let imgElm = document.createElement('img');
    imgElm.src = "img/"+image.IconID
    imgElm.alt = image.Description;
    element.appendChild(imgElm);
    return element;
}

/* Toggles element 'hidden' status on or off
 */
function toggleHide(element) {
    if (element.hasAttribute('hidden')) {
        element.removeAttribute('hidden');
    } else {
        element.setAttribute('hidden', '');
    }
}

/* Creates a new row <tr>
 * Inputs:
 * days:        []obj:    'day' objects
 * table:       element:  target table element
 * elementType: str:      the type of html element to add i.e 'tr', 'td'
 * parameters:  []str:    containing the path from day object to chosen field
 * rowTitle:    str:      A title for the row (placed in first column)
 * func:        function: a function used to add content to the created cell
 * unit:        string:   unit sign of the value
 * masterOf     string:   adds a show/hide button that controls the visibility of the group of that name
 * slaveTo      string:   is set as hidden, and belongs to group by this name
 */
function populateRow(days, table, elementType, path, rowTitle, func, unit='', masterOf='', slaveTo='') {
    if (path[0] === '') {
        path.splice(0, 1);
    }
    // Create new row
    let rowElm = document.createElement('tr');
    if (slaveTo) {
        rowElm.setAttribute("hidden", '');
        rowElm.classList.add(slaveTo);
    }
    table.appendChild(rowElm);
    // Title cells
    let titleElm = document.createElement(elementType);
    if (masterOf) {
        titleElm.textContent = rowTitle + " +";
        titleElm.setAttribute("id", masterOf);
        titleElm.addEventListener("click", function() {
            const subjectLines = Array.from( document.getElementsByClassName(masterOf) );
            subjectLines.forEach(line => {
                toggleHide(line);
            });
        });

    } else {
        titleElm.textContent = rowTitle;
    }
    rowElm.appendChild(titleElm);
    // Data cells
    let columnElm;
    days.forEach(day => {
        columnElm = document.createElement(elementType);
        let element = func(columnElm, day, path, unit, decimals);
        rowElm.appendChild(element);
    });
}

/* Creates all new rows <tr>, by calling populateRow()
 * Inputs:
 * days:        []obj:    'day' objects
 * table:       element:  target table element
 * parefix:     []str:    containing the path from day object to chosen field
 */
export function populateTable(days, table, path) {
    // Paths
    const title = [path, 'Title'];
    const image = [path];
    const desc = [path, 'Description'];
    const temp = [path, 'Temp', 'Temp'];
    const feels = [path, 'Temp', 'Feels'];
    const bulb = [path, 'Temp', 'Bulb'];
    const uv = [path, 'Uv'];
    const radiation = [path, 'Radiation', 'Direct'];
    const diffuse = [path, 'Radiation', 'Diffuse'];

    const clouds = [path, 'Clouds','Clouds'];
    const cloudsLo = [path, 'Clouds','Low'];
    const cloudsMed = [path, 'Clouds','Mid'];
    const cloudsHigh = [path, 'Clouds','High'];

    const chance = [path, 'Rain', 'Chance'];
    const amount = [path, 'Rain', 'Amount'];

    const speed = [path, 'Wind','Speed'];
    const gust = [path, 'Wind','Gust'];
    const direction = [path, 'Wind','Dir'];

    const pressure = [path, 'Pressure'];
    const humidity = [path, 'Humidity'];
    const visibility = [path, 'Visibility'];


    // Create rows and titles
    table.innerHTML = "";
    populateRow(days, table, 'th', title,      "",           addSun, '');                  //Day name with sun up indication
    populateRow(days, table, 'td', image,      "",           addImage, '');                // Icons
    populateRow(days, table, 'td', desc,       "Desc.",      addText, '');
    populateRow(days, table, 'td', temp,       "Temp.",      addNum, '°C', "temp");

    populateRow(days, table, 'td', feels,      "Feels",      addNum, '°C', '', "temp");
    populateRow(days, table, 'td', bulb,       "Wet Bulb",   addNum, '°C', '', "temp");
    populateRow(days, table, 'td', uv,         "UV",         addNum, '', '', "temp");

    populateRow(days, table, 'td', radiation,  "Radiation",  addNum, ' W/m²', 'radiation', "temp");
    populateRow(days, table, 'td', diffuse,    "Diffuse",    addNum, ' W/m²', '', "radiation");

    populateRow(days, table, 'td', clouds,     "Clouds %",   addNum, ' %', "clouds");      // Total
    populateRow(days, table, 'td', cloudsLo,   "Low %",      addNum, ' %', '', "clouds");
    populateRow(days, table, 'td', cloudsMed,  "Med %",      addNum, ' %', '', "clouds");
    populateRow(days, table, 'td', cloudsHigh, "High %",     addNum, ' %', '', "clouds");

    populateRow(days, table, 'td', visibility, "Visibility", addText, " m", '', "clouds");

    populateRow(days, table, 'td', chance,     "Rain%",      addNum, ' %');                // Chance
    populateRow(days, table, 'td', amount,     "Rain [mm]",  addNum, ' mm', "atmo");       // Total
    populateRow(days, table, 'td', pressure,   "Pressure",    addNum, ' hPa', '', "atmo");
    populateRow(days, table, 'td', humidity,   "Humidity",    addNum, ' %', '', "atmo");

    populateRow(days, table, 'td', speed,      "Wind",       addNum, ' m/s', "wind");
    populateRow(days, table, 'td', gust,       "Gust",       addText, ' m/s', "", "wind");
    populateRow(days, table, 'td', direction,  "Direction",  addText, '');

    color.applyColors(table);
}
