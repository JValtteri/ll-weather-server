package main

import (
    "fmt"
)

func main() {
    fmt.Println("### Low Latency Weather Server ###")
    printAttribution()
    server()
    fmt.Println("Shutting down")
}

func printAttribution() {
    fmt.Println("Data Sources:")
    fmt.Println("  https://open-meteo.org/")
    fmt.Println("  https://openweathermap.org/api")
}
