package main

import (
    "fmt"
    "net/http"
    "log"
)

const ORIGIN_URL string = "http://localhost:8000"

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
    city = request.URL.Query().Get("name")
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
