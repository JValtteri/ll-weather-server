
/* Cookie(s)
 */

/* Creates a new cookie
 * name: str: cookie name
 * value: str: cookie value
 * ttl: int: cookies time-to-live in days. leave empty for session only
 */
export function setCookie(name, value, ttl) {
    const d = new Date();
    d.setTime(d.getTime() + (ttl));
    let expires = "expires="+ d.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";SameSite=Lax" + ";path=/";
}

/* Returns the value of a named cookie
 */
export function getCookie(name) {
    name = name + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let cookies = decodedCookie.split(';');
    for(let i = 0; i < cookies.length; i++) {
        let c = cookies[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

export function clearSiteCookies() {
    // Get all cookies from the document
    const cookies = document.cookie.split(";");

    // Loop through each cookie string
    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i];

        // Split to get individual name and value pairs
        const eqPos = cookie.indexOf("=");
        const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;

        // Set the cookie with an expiration date in the past
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/";
    }
}