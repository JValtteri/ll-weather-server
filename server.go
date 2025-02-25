package main

import (
    "fmt"
    "net/http"
    "log"
    "strings"
    "encoding/json"
    "unicode"
    "unicode/utf8"
    "os"
)

type Config struct {
    ORIGIN_URL       string
    SERVER_PORT      string
    ENABLE_TLS       bool
    CERT_FILE        string
    PRIVATE_KEY_FILE string
    UNITS            string
    COUNTRY_CODE     string
}

const CONFIG_FILE string = "config.json"
var CONFIG Config

func server() {
    log.Println("Server UP")
    loadConfig(&CONFIG)
    LoadAPIConfig(CONFIG.UNITS, CONFIG.COUNTRY_CODE)
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
    var city string
    var f_obj WeekWeather
    city = sanitize(request.URL.Query().Get("name"))
    r_obj, err := GetProxyWeather(city, 0)
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
    return result.String()
}

func loadConfig(CONFIG *Config) {
    raw_config, err := os.ReadFile(CONFIG_FILE)
    if err != nil {
        log.Fatal(err)
    }
    err = json.Unmarshal(raw_config, &CONFIG)
    if err != nil {
        log.Println(err)
    }
    log.Printf("Server url/port: %s:%s\n", CONFIG.ORIGIN_URL, CONFIG.SERVER_PORT)
    if CONFIG.ENABLE_TLS {
        log.Println("TLS is Enabled")
    } else {
        log.Println("HTTP-Only mode")
    }
}
