/*
 * Miscellanious functions
 */

export const DAY = 24*60*60*1000;
export const SECOND = 1000;

/* Maps the timeframe variable to string value
 */
export function str_tf(tf) {
    if (tf=="1") {
        return "Day";
    } else {
        return "Night";
    }
}

/* Converts str to Base64, via uint8
 */
export function base64(str) {
    const encoder = new TextEncoder();
    const utf8Bytes = encoder.encode(str);
    return btoa(String.fromCharCode(...utf8Bytes));
}

/* Converts Base64 to str, via uint8
 */
export function decode64(str) {
    return atob(str);
}

