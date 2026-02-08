
/*
 * Color highlights
 */

function colorTemp(element) {
    let value = parseInt(element.textContent.split('°')[0]);
    if (value > 23) {
        element.classList.add('hot')
    } else if (value <= -30) {
        element.classList.add('vheavy')
    } else if (value <= -20) {
        element.classList.add('arctic')
    } else if (value < 0) {
        element.classList.add('cold')
    }
}

function colorCloud(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 66) {
        element.classList.add('overcast')
    } else if (value > 33) {
        element.classList.add('broken-clouds')
    } else if (value > 0) {
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
    let value = parseFloat(element.textContent.split(' ')[0]);
    if (value >= 20) {
        element.classList.add('warn')
    } else if (value >= 7) {
        element.classList.add('vheavy')
    } else if (value > 1) {
        element.classList.add('heavy')
    } else if (value > 1/3) {
        element.classList.add('medium')
    } else if (value > 0.0) {
        element.classList.add('light')
    }
}

function colorWind(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value >= 20) {
        element.classList.add('warn')
    } else if (value >= 15) {
        element.classList.add('vheavy')
    } else if (value >= 10) {
        element.classList.add('heavy')
    } else if (value > 7) {
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

function idByTitle(table, name) {

    for ( let i=0 ; i < table.rows.length ; i++ ) {
        const text = table.rows[i].cells[0].innerText.split(" ")[0].trim();
        if (text === name) {
            return i;
        }
    }
    console.error(`Name: "${name}" not found in table`);
}

/*
 * Applies highlight colors to an entire table
 */
export function applyColors(table) {
    applyColorRow(table, idByTitle(table, "Temp."), colorTemp);
    applyColorRow(table, idByTitle(table, "Feels"), colorTemp);
    applyColorRow(table, idByTitle(table, "Wet"), colorTemp);

    applyColorRow(table, idByTitle(table, "Clouds"), colorCloud);
    applyColorRow(table, idByTitle(table, "Low"), colorCloud);
    applyColorRow(table, idByTitle(table, "Med"), colorCloud);
    applyColorRow(table, idByTitle(table, "High"), colorCloud);

    applyColorRow(table, idByTitle(table, "Rain%"), colorRainChance);
    applyColorRow(table, idByTitle(table, "Rain"), colorRain);
    applyColorRow(table, idByTitle(table, "Wind"), colorWind);
    applyColorRow(table, idByTitle(table, "Gust"), colorWind);
}

export function colorSun(element, sunUp) {
    if (!sunUp) {
        element.classList.add('night');
    } else {
        element.classList.remove('night');   // This should be unnecessary
    }
}


