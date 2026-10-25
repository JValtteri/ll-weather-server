import * as color from "./color.js?id=MTc3MDUxMTU4Nwk";


/* Function to construct a path to a specific element in JSON
 * dataTarget: obj:    'day' object
 * path:       []str:  path to element expressed as a str list
 */
function constructJsonPath(dataTarget, path) {
    for (let i = 0; i < path.length; ++i) {
        if (dataTarget === undefined || dataTarget === null) {
            return undefined;
        }
        dataTarget = dataTarget[path[i]];
    }
    return dataTarget;
}

/* Returns the data object that holds the Overview values, and the path
 * relative to that object.
 * Day forecasts store the data separately for 'Day' and 'Night', while hour
 * forecasts store it directly on the 'day' object.
 */
function getOverviewData(day, path) {
    if (!day) return null;
    let data = day;
    let subPath = Array.isArray(path) ? path.slice() : [];
    if (subPath.length > 0 && (subPath[0] === 'Day' || subPath[0] === 'Night')) {
        const timeframe = day[subPath[0]];
        if (timeframe && typeof timeframe === 'object') {
            data = timeframe;
            subPath = subPath.slice(1);
        }
    }
    if (!data || !data.Overview) return null;
    return { overview: data.Overview, path: subPath };
}

function isOverview(day, path) {
    return getOverviewData(day, path) !== null;
}

function setTooltip(element, text) {
    if (text !== undefined && text !== null && text !== '') {
        element.title = text;
    }
}

/* Function to add text content to a cell
 * To be used as a parameter for populateRow()
 */
function addText(element, day, path, unit='') {
    if (isOverview(day, path)) {
        let text = constructJsonPath(day, path);
        element.textContent = `${text ?? ''}${unit}`;
        return element;
    }
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
    if (sun) {
        color.colorSun(element, sun.SunUp);
    }
    return element;
}


/* Function to add int content to a cell
 * Rounds numbers to integers. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addInt(element, day, path, unit='') {
    if (isOverview(day, path)) {
        return addIntOverview(element, day, path, unit);
    }
    let num = constructJsonPath(day, path);
    if (num === undefined || num === null) {
        num = 0;
    }
    element.textContent = `${num.toFixed(0)}${unit}`;
    return element;
}

function getOverviewForPath(day, path) {
    const context = getOverviewData(day, path);
    if (!context) return null;
    const ov = context.overview;
    const p = context.path;
    if (p[0] === 'Temp') {
        if (p[1] === 'Temp') return ov.Temp;
        if (p[1] === 'Feels') return ov.Feels;
        if (p[1] === 'Bulb') return ov.Bulb;
    }
    if (p[0] === 'Uv') return ov.Uv;
    if (p[0] === 'Radiation') {
        if (p[1] === 'Direct') return ov.Direct;
        if (p[1] === 'Diffuse') return ov.Diffuse;
    }
    if (p[0] === 'Clouds') {
        if (p[1] === 'Clouds') return ov.Clouds;
        if (p[1] === 'Low') return ov.CloudsLo;
        if (p[1] === 'Mid') return ov.CloudsMed;
        if (p[1] === 'High') return ov.CloudsHigh;
    }
    if (p[0] === 'Pressure') return ov.Pressure;
    if (p[0] === 'Humidity') return ov.Humidity;
    if (p[0] === 'Visibility') return ov.Visibility;
    if (p[0] === 'Rain') {
        if (p[1] === 'Chance') return ov.RainChance;
        if (p[1] === 'Amount') return ov.RainAmount;
    }
    if (p[0] === 'Wind') {
        if (p[1] === 'Speed') return ov.WindSpeed;
        if (p[1] === 'Gust') return ov.WindGust;
    }
    return null;
}

/* Returns the path relative to the object that holds the Overview values.
 * Used for coloring the gradient, which expects paths without the Day/Night
 * prefix.
 */
function getOverviewPath(day, path) {
    const context = getOverviewData(day, path);
    return context ? context.path : path;
}

/* Formats a min-max range. If both ends would display as the same number,
 * only that single value is shown.
 */
function formatOverviewRange(minDisp, maxDisp, unit='') {
    const suffix = unit ? ` ${unit.trim()}` : '';
    if (minDisp === maxDisp) {
        return `${minDisp}${suffix}`;
    }
    return `${minDisp} - ${maxDisp}${suffix}`;
}

function addIntOverview(element, day, path, unit='') {
    const ov = getOverviewForPath(day, path);
    if (!ov) {
        let num = constructJsonPath(day, path);
        if (num === undefined || num === null) num = 0;
        element.textContent = `${num.toFixed(0)}${unit}`;
        return element;
    }
    const minv = ov.Min, maxv = ov.Max, avg = ov.Avg;
    let minDisp = Math.round(minv);
    let maxDisp = Math.round(maxv);
    element.textContent = formatOverviewRange(minDisp, maxDisp, unit);
    setTooltip(element, `Average: ${avg.toFixed(1)}${unit ? unit.trim() : ''}`);
    color.applyOverviewGradient(element, {min: minv, max: maxv, avg: avg}, getOverviewPath(day, path), unit);
    return element;
}

/* Function to add int content to a cell
 * Rounds numbers to integer/1000. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addNumPerK(element, day, path, unit='') {
    if (isOverview(day, path)) {
        const ov = getOverviewForPath(day, path);
        if (ov) {
            const minv = ov.Min/1000, maxv = ov.Max/1000, avg = ov.Avg/1000;
            element.textContent = formatOverviewRange(minv.toFixed(1), maxv.toFixed(1), unit);
            setTooltip(element, `Average: ${avg.toFixed(1)}${unit ? unit.trim() : ''}`);
            color.applyOverviewGradient(element, {min: ov.Min, max: ov.Max, avg: ov.Avg}, getOverviewPath(day, path), unit);
            return element;
        }
    }
    let num = constructJsonPath(day, path);
    if (num === undefined || num === null) {
        num = 0;
    }
    const scaled = num/1000;
    element.textContent = `${scaled.toFixed(1)}${unit}`;
    return element;
}


/* Function to add number content to a cell
 * Rounds numbers to one decimal. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addFloat(element, day, path, unit='') {
    if (isOverview(day, path)) {
        const ov = getOverviewForPath(day, path);
        if (ov) {
            element.textContent = formatOverviewRange(ov.Min.toFixed(1), ov.Max.toFixed(1), unit);
            setTooltip(element, `Average: ${ov.Avg.toFixed(1)}${unit ? unit.trim() : ''}`);
            color.applyOverviewGradient(element, {min: ov.Min, max: ov.Max, avg: ov.Avg}, getOverviewPath(day, path), unit);
            return element;
        }
    }
    let num = constructJsonPath(day, path);
    if (num === undefined || num === null) {
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

/* Creates a new row <tr> and populates it
 * Inputs:
 * days:        []obj:    'day' objects
 * table:       element:  target table element
 * elementType: str:      the type of html element to add inside row i.e 'th', 'td'
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
    let rowElm = createNewRow(table, slaveTo);
    addTitleCell(rowElm, elementType, rowTitle, masterOf);

    // Data cells
    days.forEach(day => {
        addDataCell(rowElm, elementType, day, path, func, unit);
    });
}

/* Creates a new row element <tr>
 * Inputs:
 * table:       element:  target table element
 * slaveTo      string:   is set as hidden, and belongs to group by this name
 */
function createNewRow(table, slaveTo) {
    let rowElm = document.createElement('tr');
    if (slaveTo) {
        rowElm.setAttribute("hidden", '');
        rowElm.classList.add(slaveTo);
    }
    table.appendChild(rowElm);
    return rowElm;
}

/* Creates a new tile column element <td> or <th> and populates it
 * Inputs:
 * rowElem:     element:  row element <tr>
 * elementType: str:      the type of html element to add i.e 'tr', 'td'
 * rowTitle:    str:      A title for the row (placed in first column)
 * masterOf     string:   adds a show/hide button that controls the visibility of the group of that name
 */
function addTitleCell(rowElm, elementType, rowTitle, masterOf) {
    let titleElm = document.createElement(elementType);
    titleElm.classList.add('static');
    if (masterOf) {
        titleElm.innerHTML = rowTitle + "  <b>+</b>"; // could use &ensp; for longer space;
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
}

/* Creates a new cell <td> or <th>
 * Inputs:
 * rowElem:     element:  row element <tr>
 * elementType: str:      the type of html element to add i.e 'td', 'th'
 * day          obj:      'day' object
 * path:        []str:    path to element expressed as a str list
 * func:        function: a function used to add content to the created cell
 * unit:        string:   unit sign of the value
 */
function addDataCell(rowElm, elementType, day, path, func, unit) {
    let columnElm = document.createElement(elementType);
    let element = func(columnElm, day, path, unit);
    rowElm.appendChild(element);
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

    populateRow(days, table, 'td', clouds,     "Clouds",   addInt, ' %', "clouds");      // Total
    populateRow(days, table, 'td', cloudsHigh, "High %",     addInt, ' %', '', "clouds");
    populateRow(days, table, 'td', cloudsMed,  "Med %",      addInt, ' %', '', "clouds");
    populateRow(days, table, 'td', cloudsLo,   "Low %",      addInt, ' %', '', "clouds");

    populateRow(days, table, 'td', pressure,   "Pressure",    addInt, ' hPa', '', "clouds");
    populateRow(days, table, 'td', humidity,   "Humidity",    addInt, ' %', '', "clouds");

    populateRow(days, table, 'td', visibility, "Visibility", addNumPerK, " km", '', "clouds");

    populateRow(days, table, 'td', chance,     "Rain%",      addFloat, ' %');                // Chance
    populateRow(days, table, 'td', amount,     "Rain [mm]",  addFloat, ' mm');       // Total

    populateRow(days, table, 'td', speed,      "Wind",       addInt, ' m/s', "wind");
    populateRow(days, table, 'td', gust,       "Gust",       addInt, ' m/s', "", "wind");
    populateRow(days, table, 'td', direction,  "Direction",  addText, '');

    color.applyColors(table);
}
