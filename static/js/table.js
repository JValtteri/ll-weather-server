
import * as color from "./color.js";

/* Tables
 */

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
        let element = func(columnElm, dataTarget, unit);
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
export function populateTable(days, table, prefix) {
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
    color.applyColors(table);
}
