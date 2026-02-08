import * as color from "./color.js?id=MTc3MDUxMTU4Nw";


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
function addText(element, day, path, unit='') {
    let text = constructJsonPath(day, path);
    element.textContent = `${text}${unit}`;
    return element;
}

/* Function to day/night indicator to title
 * To be used as a parameter for populateRow()
 */
function addSun(element, day, path, _='') {
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


/* Function to add int content to a cell
 * Rounds numbers to integers. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addInt(element, day, path, unit='') {
    let num = constructJsonPath(day, path);
    if (!num) {
        num = 0;
    }
    element.textContent = `${num.toFixed(0)}${unit}`;
    return element;
}

/* Function to add number content to a cell
 * Rounds numbers to one decimal. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addFloat(element, day, path, unit='') {
    let num = constructJsonPath(day, path);
    if (!num) {
        num = 0;
    }
    element.textContent = `${num.toFixed(1)}${unit}`;
    return element;
}

/* Function to add an image element from data.day
 * To be used as a parameter for populateRow()
 */
function addImage(element, day, path, _='') {
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

/* Toggles between '×' and '+' for expanding rows
 */
function toggleX(element) {
    if (element.innerHTML.slice(-5, -4) === "+") {
        element.innerHTML = element.innerHTML.replace("<b>+</b>", "<b>&times</b>");
    } else {
        element.innerHTML = element.innerHTML.replace("<b>×</b>", "<b>+</b>"); // &times;
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
        titleElm.innerHTML = rowTitle + "&ensp; <b>+</b>";
        titleElm.setAttribute("id", masterOf);
        titleElm.addEventListener("click", function() {
            const subjectLines = Array.from( document.getElementsByClassName(masterOf) );
            subjectLines.forEach(line => {
                toggleHide(line);
            });
            toggleX(titleElm);
        });

    } else {
        titleElm.textContent = rowTitle;
    }
    rowElm.appendChild(titleElm);
    // Data cells
    let columnElm;
    days.forEach(day => {
        columnElm = document.createElement(elementType);
        let element = func(columnElm, day, path, unit);
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
    populateRow(days, table, 'td', desc,       "Desc.",      addText, '', '', 'void');
    populateRow(days, table, 'td', temp,       "Temp.",      addInt, '°C', "temp");

    populateRow(days, table, 'td', feels,      "Feels",      addInt, '°C', '', "temp");
    populateRow(days, table, 'td', bulb,       "Wet Bulb",   addInt, '°C', '', "temp");
    populateRow(days, table, 'td', uv,         "UV",         addInt, '', '', "temp");

    populateRow(days, table, 'td', radiation,  "Radiation",  addInt, ' W/m²', 'radiation', "temp");
    populateRow(days, table, 'td', diffuse,    "Diffuse",    addInt, ' W/m²', '', "radiation");

    populateRow(days, table, 'td', clouds,     "Clouds %",   addInt, ' %', "clouds");      // Total
    populateRow(days, table, 'td', cloudsHigh, "High %",     addInt, ' %', '', "clouds");
    populateRow(days, table, 'td', cloudsMed,  "Med %",      addInt, ' %', '', "clouds");
    populateRow(days, table, 'td', cloudsLo,   "Low %",      addInt, ' %', '', "clouds");

    populateRow(days, table, 'td', pressure,   "Pressure",    addInt, ' hPa', '', "clouds");
    populateRow(days, table, 'td', humidity,   "Humidity",    addInt, ' %', '', "clouds");

    populateRow(days, table, 'td', visibility, "Visibility", addInt, " m", '', "clouds");

    populateRow(days, table, 'td', chance,     "Rain%",      addFloat, ' %');                // Chance
    populateRow(days, table, 'td', amount,     "Rain [mm]",  addFloat, ' mm');       // Total

    populateRow(days, table, 'td', speed,      "Wind",       addInt, ' m/s', "wind");
    populateRow(days, table, 'td', gust,       "Gust",       addInt, ' m/s', "", "wind");
    populateRow(days, table, 'td', direction,  "Direction",  addText, '');

    color.applyColors(table);
}
