
/*
 * Color highlights
 */

function colorTemp(element) {
    let value = parseInt(element.textContent.split('°')[0]);
    if (value > 23) {
        element.classList.add('hot')
    } else if (value < -20) {
        element.classList.add('arctic')
    } else if (value < 0) {
        element.classList.add('light')
    }
}

function colorCloud(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 50) {
        element.classList.add('broken-clouds')
    } else if (value > 10) {
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
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value > 3) {
        element.classList.add('heavy')
    } else if (value > 1) {
        element.classList.add('medium')
    } else if (value > 0) {
        element.classList.add('light')
    }
}

function colorWind(element) {
    let value = parseInt(element.textContent.split(' ')[0]);
    if (value >= 15) {
        element.classList.add('heavy')
    } else if (value > 9) {
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

/*
 * Applies highlight colors to an entire table
 */
export function applyColors(table) {
    applyColorRow(table, 3, colorTemp);
    applyColorRow(table, 4, colorCloud);
    applyColorRow(table, 5, colorRainChance);
    applyColorRow(table, 6, colorRain);
    applyColorRow(table, 7, colorWind);
}
