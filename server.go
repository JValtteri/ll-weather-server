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
    "github.com/JValtteri/weather/owm"
    "github.com/JValtteri/weather/om"
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
    http.HandleFunc("/model",       modelChangeRequest)
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
    var city string  = getCookie(request, "city")
    var legacy string = sanitize(getCookie(request, "alternate"))
    city = sanitize(decode64(city))
    if city == "" {
        return
    }
    var err error
    if legacy == "true" {
        var r_obj owm.WeatherRange
        r_obj, err = GetOwmProxyWeather(city)
        f_obj = mapOwmDays(r_obj)
    } else {
        var r_obj om.WeatherRange
        r_obj, err = GetOmProxyWeather(city)
        f_obj = mapOmDays(r_obj)
    }
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
    var city string  = getCookie(request, "city")
    var legacy string = sanitize(getCookie(request, "alternate"))
    city = sanitize(decode64(city))
    if city == "" {
        return
    }
    var dayNo string = sanitize(request.URL.Query().Get("day"))
    var dayNumber int
    dayNumber, err := strconv.Atoi(dayNo)
    if err != nil {
        http.Error(w, "400 Bad Request", http.StatusBadRequest)
    }
    if legacy == "true" {
        var r_obj owm.WeatherRange
        r_obj, err = GetOwmProxyWeather(city)
        f_obj = mapOwmHours(r_obj, dayNumber)
    } else {
        var r_obj om.WeatherRange
        r_obj, err = GetOmProxyWeather(city)
        f_obj = mapOmHours(r_obj, dayNumber)
    }
    setCorrs(w)
    if err != nil {
        http.NotFound(w, request)
    }
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, unloadJSON(f_obj))
}

func modelChangeRequest(w http.ResponseWriter, request *http.Request) {
    rqNum++
    var model string = decode64(request.URL.Query().Get("model"))
    if !whitelist(model, modelList) {
        log.Printf("r:%6vu:%4v: Bad Model: %s\n", rqNum, uniqRqNum, model)
        return
    }
    // TODO:  ADD whitelist to catch any funnybusiness
    CONFIG.MODEL = model
    om.Config("", CONFIG.UNITS, CONFIG.MODEL)  // Update to use the requested model
    log.Printf("r:%6vu:%4v: Get Model: %s\n", rqNum, uniqRqNum, model)
}

func getCookie(request *http.Request, cookieName string) string {
    cookie, err := request.Cookie(cookieName)
    if err != nil {
        return ""
    }
    return cookie.Value
}

func decode64(str64 string) string {
    value, err := base64.StdEncoding.DecodeString(str64)
    if err != nil {
        return ""
    }
    decoded := strings.TrimSpace(string(value))
    return decoded
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

func whitelist(input string, list []string) bool {
    for _, v := range list {
        if v == input {
            return true
        }
    }
    return false
}

var modelList = []string{
    "",
    "best_match",
    "icon_eu",
    "gem_global",
    "gfs_graphcast025",
    "knmi_harmonie_arome_europe",
    "dmi_harmonie_arome_europe",
    "meteofrance_arpege_europe",
    "dmi_seamless",
    "ukmo_global_deterministic_10km",
    "jma_gsm",
    "metno_seamless",
    "ecmwf_ifs025",
    "ecmwf_aifs025_single",
    "cma_grapes_global",
}
