package main

import (
    "fmt"
    "net/http"
    "log"
)

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
    city = request.URL.Query().Get("name")
    fmt.Fprintf(w, "Requested: [%s]", city)
}

func cityDetailRequest(w http.ResponseWriter, request *http.Request) {
    var city string
    city = request.URL.Query().Get("name")
    fmt.Fprintf(w, "Detail Requested: [%s]", city)
}
