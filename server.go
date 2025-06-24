package main

import (
    "fmt"
    "net/http"
    "log"
    "strings"
    "strconv"
    "unicode"
    "unicode/utf8"
    "encoding/base64"
)

var rqNum uint = 0
var uniqRqNum uint = 0

func server() {
    log.Println("Server UP")
    LoadConfig()
    http.HandleFunc("/", defaultRequest)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    http.Handle("/js/",  http.StripPrefix("/js/",  http.FileServer(http.Dir("./static/js"))))
    http.HandleFunc("/city",        cityOverviewRequest)
    http.HandleFunc("/city/detail", cityDetailRequest)
    if CONFIG.ENABLE_TLS {
        log.Fatal(http.ListenAndServeTLS(
            fmt.Sprintf( ":%s", CONFIG.SERVER_PORT),
            CONFIG.CERT_FILE,
            CONFIG.PRIVATE_KEY_FILE, nil))
    } else {
        log.Fatal(http.ListenAndServe( fmt.Sprintf(":%s", CONFIG.SERVER_PORT), nil))
    }
}

func defaultRequest(w http.ResponseWriter, request *http.Request) {
    var path []string = strings.Split(request.URL.Path, "/")
    if path[1] == "img" {
        var id string   = sanitize(path[2])
        var icon []byte = GetProxyIcon(id)
        setCorrs(w)
        w.Header().Set("Content-Type", "image/png")
        _, err := w.Write(icon)
        if err != nil {
            http.Error(w, "Failed to send image", http.StatusInternalServerError)
        }
    } else {
        http.ServeFile(w, request, "./static/index.html")
    }
}

func cityOverviewRequest(w http.ResponseWriter, request *http.Request) {
    rqNum++
    var f_obj WeekWeather
    var city string  = sanitize(getCookie(request))
    if city == "" {
        return
    }
    r_obj, err := GetProxyWeather(city)
    f_obj = mapDays(r_obj)
    setCorrs(w)
    if err != nil {
        http.NotFound(w, request)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, unloadJSON(f_obj))
}

func cityDetailRequest(w http.ResponseWriter, request *http.Request) {
    rqNum++
    var f_obj DayHours
    var city string  = sanitize(getCookie(request))
    if city == "" {
        return
    }
    var dayNo string = sanitize(request.URL.Query().Get("day"))
    var dayNumber int
    dayNumber, err := strconv.Atoi(dayNo)
    if err != nil {
        http.Error(w, "400 Bad Request", http.StatusBadRequest)
    }
    r_obj, err := GetProxyWeather(city)
    f_obj = mapHours(r_obj, dayNumber)
    setCorrs(w)
    if err != nil {
        http.NotFound(w, request)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, unloadJSON(f_obj))
}

func getCookie(request *http.Request) string {
    cookie, err := request.Cookie("city")
    if err != nil {
        return ""
    }
    value, err := base64.StdEncoding.DecodeString(cookie.Value)
    if err != nil {
        return ""
    }
    city := strings.TrimSpace(string(value))
    return city
}

func setCorrs(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", CONFIG.ORIGIN_URL)
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func sanitize(input string) string {
    var result strings.Builder
    for i := 0; i < len(input); {
        r, size := utf8.DecodeRuneInString(input[i:])
        if unicode.IsSpace(r) || unicode.IsLetter(r) || unicode.IsDigit(r) {
            result.WriteRune(r)
            i += size
        } else {
            i++
        }
    }
    return strings.ToLower(result.String())
}
