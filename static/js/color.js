
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

function colorUV(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value >= 8) {
        element.classList.add('vheavy')
    } else if (value >= 6) {
        element.classList.add('warn')
    } else if (value >= 5) {
        element.classList.add('hot')
    } else if (value >= 4) {
        element.classList.add('warm')
    } else if (value >= 3) {
        element.classList.add('mild')
    }
}

function colorRadiation(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 799) {
        element.classList.add('accent-green')
    } else if (value > 399) {
        element.classList.add('warm')
    } else if (value > 199) {
        element.classList.add('mild')
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

function colorPressure(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 1040) {
        element.classList.add('warm')
    } else if (value > 1022) {
        element.classList.add('mild')
    } else if (value < 1004) {
        element.classList.add('light')
    } else if (value < 986) {
        element.classList.add('vheavy')
    }
}

function colorHumidity(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 90) {
        element.classList.add('vheavy')
    } else if (value > 70) {
        element.classList.add('medium')
    } else if (value < 30) {
        element.classList.add('warm')
    }
}

function colorVisibility(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value < 1) {
        element.classList.add('warn')
    } else if (value < 10) {
        element.classList.add('mild')
    } else if (value < 30) {
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

    applyColorRow(table, idByTitle(table, "UV"), colorUV);
    applyColorRow(table, idByTitle(table, "Radiation"), colorRadiation);
    applyColorRow(table, idByTitle(table, "Diffuse"), colorRadiation);

    applyColorRow(table, idByTitle(table, "Clouds"), colorCloud);
    applyColorRow(table, idByTitle(table, "Low"), colorCloud);
    applyColorRow(table, idByTitle(table, "Med"), colorCloud);
    applyColorRow(table, idByTitle(table, "High"), colorCloud);

    applyColorRow(table, idByTitle(table, "Pressure"), colorPressure);
    applyColorRow(table, idByTitle(table, "Humidity"), colorHumidity);
    applyColorRow(table, idByTitle(table, "Visibility"), colorVisibility);

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
