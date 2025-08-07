
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
function addText(element, day, path, _='') {
    let text = constructJsonPath(day, path);
    element.textContent = text;
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


/* Function to add number content to a cell
 * Rounds numbers to integers. Adds the optional unit parameter
 * To be used as a parameter for populateRow()
 */
function addNum(element, day, path, unit='') {
    let num = constructJsonPath(day, path);
    if (!num) {
        num = 0;
    }
    element.textContent = `${num.toFixed(0)}${unit}`;
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

/* Creates a new row <tr>
 * Inputs:
 * days:        []obj:    'day' objects
 * table:       element:  target table element
 * elementType: str:      the type of html element to add i.e 'tr', 'td'
 * parameters:  []str:    containing the path from day object to chosen field
 * rowTitle:    str:      A title for the row (placed in first column)
 * func:        function: a function used to add content to the created cell
 */
function populateRow(days, table, elementType, path, rowTitle, func, unit='') {
    if (path[0] === '') {
        path.splice(0, 1);
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
    // Create rows and titles
    table.innerHTML = "";
    populateRow(days, table, 'th', [path, 'Title'],            "",          addSun, ''); //Day name with sun up indication
    populateRow(days, table, 'td', [path],                     "",          addImage, '');// Icons
    populateRow(days, table, 'td', [path, 'Description'],      "Desc.",     addText, '');
    populateRow(days, table, 'td', [path, 'Temp'],             "Temp.",     addNum, '°C');
    populateRow(days, table, 'td', [path, 'Clouds','Clouds'],  "Clouds %",  addNum, ' %'); // Total
    populateRow(days, table, 'td', [path, 'Rain', 'Chance'],   "Rain%",     addNum, ' %'); // Chance
    populateRow(days, table, 'td', [path, 'Rain', 'Amount'],   "Rain [mm]", addNum, ' mm'); // Total
                                                                                      // Layers
    populateRow(days, table, 'td', [path, 'Wind','Speed'],     "Wind",      addNum, ' m/s');
    populateRow(days, table, 'td', [path, 'Wind','Dir'],       "Direction", addText, '');
    //populateRow(days, table, 'td', [path, 'Pressure'],         "Pressure",  addNum, ' hPa');
    //populateRow(days, table, 'td', [path, 'Humidity'],         "Humidity",  addNum, ' %');
    //populateRow(days, table, 'td', [path, 'Visibility'],      "Visibility", addText);
    color.applyColors(table);
}
