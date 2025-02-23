package main

import (
    "fmt"
    "net/http"
    "log"
    "strings"
    "unicode"
    "unicode/utf8"
    "os"
)

const CONFIG_FILE string = "config.txt"
var ORIGIN_URL string = loadConfig()

func server() {
    fmt.Println("Server UP")
    http.HandleFunc("/", apiUp)
    http.HandleFunc("/city", cityOverviewRequest)
    http.HandleFunc("/city/detail", cityDetailRequest)
    log.Fatal(http.ListenAndServe(":3000", nil))
}

func apiUp(w http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(w, "API UP")
}

func cityOverviewRequest(w http.ResponseWriter, request *http.Request) {
    var city string
    var f_obj WeekWeather
    city = sanitize(request.URL.Query().Get("name"))
    r_obj, err := getProxyWeather(city, 0)
    f_obj = mapDays(r_obj)
    setCorrs(w)
    if err != nil {
        http.NotFound(w, request)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, unloadJSON(f_obj))
}

func cityDetailRequest(w http.ResponseWriter, request *http.Request) {
    var city string
    city = request.URL.Query().Get("name")
    fmt.Fprintf(w, "Detail Requested: [%s]", city)
}

func setCorrs(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", ORIGIN_URL)
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
    return result.String()
}

func loadConfig() string {
    content, err := os.ReadFile(CONFIG_FILE)
    if err != nil {
        log.Fatal(err)
    }
    return string(content)
}
