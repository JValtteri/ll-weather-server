/*
 * Color highlights
 */

function clearColorClasses(element) {
    const classes = ['light', 'medium', 'heavy', 'warm', 'mild', 'hot', 'cold', 'arctic', 'clear-sky', 'broken-clouds', 'overcast', 'night', 'vheavy', 'warn', 'accent-green'];
    classes.forEach(c => element.classList.remove(c));
    element.style.backgroundColor = '';
    element.style.backgroundImage = '';
}

function hexToRgb(hex) {
    const h = hex.replace('#', '');
    const bigint = parseInt(h.length === 3 ? h.split('').map(c => c + c).join('') : h, 16);
    return [(bigint >> 16) & 255, (bigint >> 8) & 255, bigint & 255];
}
function rgbToHex(r,g,b) {
    return '#' + [Math.max(0,Math.min(255,Math.round(r))), Math.max(0,Math.min(255,Math.round(g))), Math.max(0,Math.min(255,Math.round(b)))].map(x => x.toString(16).padStart(2,'0')).join('');
}
function lerpColor(a,b,t) {
    const c1 = hexToRgb(a), c2 = hexToRgb(b);
    return rgbToHex(c1[0] + (c2[0]-c1[0])*t, c1[1] + (c2[1]-c1[1])*t, c1[2] + (c2[2]-c1[2])*t);
}

// Single source of truth for the color thresholds. The result is a CSS
// class name (or null when no rule matches).
function getColorClassForValue(value, path) {
    if ((path && path[1] === 'Temp') || (path && path[0] === 'Temp')) {
        if (value > 23) return 'hot';
        if (value >= 20) return 'warm';
        if (value <= -30) return 'vheavy';
        if (value <= -20) return 'arctic';
        if (value < 0) return 'cold';
        return null;
    }
    if (path && path[0] === 'Uv') {
        if (value >= 8) return 'vheavy';
        if (value >= 6) return 'warn';
        if (value >= 5) return 'hot';
        if (value >= 4) return 'warm';
        if (value >= 3) return 'mild';
        return null;
    }
    if (path && path[0] === 'Radiation') {
        if (value > 799) return 'accent-green';
        if (value > 399) return 'warm';
        if (value > 199) return 'mild';
        return null;
    }
    if (path && path[0] === 'Clouds') {
        if (value > 66) return 'overcast';
        if (value > 33) return 'broken-clouds';
        if (value > 0) return 'clear-sky';
        return null;
    }
    if (path && path[0] === 'Pressure') {
        if (value > 1040) return 'warm';
        if (value > 1022) return 'mild';
        if (value < 986) return 'vheavy';
        if (value < 1004) return 'light';
        return null;
    }
    if (path && path[0] === 'Humidity') {
        if (value > 90) return 'vheavy';
        if (value > 70) return 'medium';
        if (value < 30) return 'warm';
        return null;
    }
    if (path && path[0] === 'Visibility') {
        if (value === 0) return null;
        if (value < 1) return 'warn';
        if (value < 10) return 'mild';
        if (value < 30) return 'clear-sky';
        return null;
    }
    if (path && path[0] === 'Rain' && path[1] === 'Chance') {
        if (value > 70) return 'medium';
        if (value > 19) return 'light';
        return null;
    }
    if (path && path[0] === 'Rain') {
        if (value >= 20) return 'warn';
        if (value >= 7) return 'vheavy';
        if (value > 1) return 'heavy';
        if (value > 1/3) return 'medium';
        if (value > 0.0) return 'light';
        return null;
    }
    if (path && path[0] === 'Wind') {
        if (value >= 20) return 'warn';
        if (value >= 15) return 'vheavy';
        if (value >= 10) return 'heavy';
        if (value > 7) return 'medium';
        if (value >= 4) return 'light';
        return null;
    }
    return null;
}

function classToColor(cls) {
    switch(cls) {
        case 'light': return '#5272ad';
        case 'medium': return '#34548f';
        case 'heavy': return '#20407b';
        case 'warm': return '#976727';
        case 'mild': return '#81785b';
        case 'hot': return '#782818';
        case 'cold': return '#707898';
        case 'arctic': return '#74a1cf';
        case 'clear-sky': return '#408d83';
        case 'broken-clouds': return '#157162';
        case 'overcast': return '#1b4d46';
        case 'night': return '#2f2f2f';
        case 'vheavy': return '#61207b';
        case 'warn': return '#a72d2d';
        case 'accent-green': return getComputedStyle ? getComputedStyle(document.documentElement).getPropertyValue('--accent-color').trim() || '#408d83' : '#408d83';
        default: return '#34548f';
    }
}

function getColorForValue(value, path) {
    const cls = getColorClassForValue(value, path);
    // No matching rule means no color, i.e. the default table background.
    return cls ? classToColor(cls) : null;
}

// Color gradient for overview cells - colored according to cell data min and max
export function applyOverviewGradient(element, range, path, unit='') {
    clearColorClasses(element);
    element.classList.add('overview-cell');
    element.setAttribute('data-overview', 'true');

    let minv = range.min;
    let maxv = range.max;
    // Ensure the lower value is always on the left
    if (minv > maxv) {
        const tmp = minv;
        minv = maxv;
        maxv = tmp;
    }

    // Each end of the gradient uses the same color rules as a regular cell.
    // A missing color means the default table background should show through.
    const cMin = getColorForValue(minv, path);
    const cMax = getColorForValue(maxv, path);

    if (!cMin && !cMax) {
        element.style.setProperty('background-color', 'transparent', 'important');
        element.style.setProperty('background-image', 'none', 'important');
        return;
    }

    // A zero-width range (including min=max=0) is a solid color.
    if (Math.abs(maxv - minv) < 1e-12) {
        element.style.setProperty('background-color', cMin, 'important');
        element.style.setProperty('background-image', 'none', 'important');
        return;
    }

    // Use transparent stops for ends without a color so the table background
    // remains visible there instead of being replaced by a fallback color.
    const startColor = cMin || 'transparent';
    const endColor = cMax || 'transparent';
    element.style.setProperty('background-image', `linear-gradient(90deg, ${startColor}, ${endColor})`, 'important');
    element.style.setProperty('background-color', cMin && cMax ? cMin : 'transparent', 'important');
    element.style.setProperty('background-size', '100% 100%', 'important');
}

/* Apply color to a row.
 * table: target table
 * row: index of the target row
 * func: function to run on each cell
 */
function applyColorRow(table, row, func) {
    if (row < 0 || !table.rows[row]) return;
    let target = table.rows[row];
    target.childNodes.forEach(cell => {
        if (cell.nodeType === 1 && cell.getAttribute && cell.getAttribute('data-overview') === 'true') {
            return;
        }
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

// Describes how each table row should be colored. The thresholds themselves
// live only in getColorClassForValue above.
const ROW_COLOR_RULES = [
    { titles: ['Temp.', 'Feels', 'Wet'], path: ['Temp'], parse: el => parseInt(el.textContent.split('°')[0]) },
    { titles: ['UV'], path: ['Uv'], parse: el => parseInt(el.textContent.split(' ')[0]) },
    { titles: ['Radiation', 'Diffuse'], path: ['Radiation'], parse: el => parseInt(el.textContent.split(' ')[0]) },
    { titles: ['Clouds', 'Low', 'Med', 'High'], path: ['Clouds'], parse: el => parseInt(el.textContent.split(' ')[0]) },
    { titles: ['Pressure'], path: ['Pressure'], parse: el => parseInt(el.textContent.split(' ')[0]) },
    { titles: ['Humidity'], path: ['Humidity'], parse: el => parseInt(el.textContent.split(' ')[0]) },
    { titles: ['Visibility'], path: ['Visibility'], parse: el => parseFloat(el.textContent.split(' ')[0]) },
    { titles: ['Rain%'], path: ['Rain', 'Chance'], parse: el => parseInt(el.textContent.split(' ')[0]) },
    { titles: ['Rain'], path: ['Rain'], parse: el => parseFloat(el.textContent.split(' ')[0]) },
    { titles: ['Wind', 'Gust'], path: ['Wind'], parse: el => parseInt(el.textContent.split(' ')[0]) },
];

/*
 * Applies highlight colors to an entire table
 */
export function applyColors(table) {
    for (const rule of ROW_COLOR_RULES) {
        const colorer = (element) => {
            clearColorClasses(element);
            const value = rule.parse(element);
            if (isNaN(value)) return;
            const cls = getColorClassForValue(value, rule.path);
            if (cls) {
                element.classList.add(cls);
                element.style.setProperty('background-color', classToColor(cls), 'important');
            }
        };
        for (const title of rule.titles) {
            try {
                applyColorRow(table, idByTitle(table, title), colorer);
            } catch (e) {}
        }
    }
}

export function colorSun(element, sunUp) {
    if (!sunUp) {
        element.classList.add('night');
        element.style.setProperty('background-color', classToColor('night'), 'important');
    } else {
        element.classList.remove('night');
        element.style.removeProperty('background-color');
    }
}
