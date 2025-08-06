package main

import (
    "log"
    "os"
    "encoding/json"
    "github.com/JValtteri/weather/owm"
    "github.com/JValtteri/weather/om"
)

type Config struct {
    ORIGIN_URL       string
    SERVER_PORT      string
    ENABLE_TLS       bool
    CERT_FILE        string
    PRIVATE_KEY_FILE string
    UNITS            string
    COUNTRY_CODE     string
    CACHE_AGE        uint  // Hours
    CACHE_SIZE       int
    MODEL            string
}

type ApiKeys struct {
    OPENWEATHERMAP string
}

const CONFIG_FILE string = "config.json"
const API_KEYFILE string = "api_keys.json"
var CONFIG Config
var API_KEYS ApiKeys

func LoadConfig() {
    raw_config := readConfig(CONFIG_FILE)
    unmarshal(raw_config, &CONFIG)
    log.Printf("Server url/port: %s:%s\n", CONFIG.ORIGIN_URL, CONFIG.SERVER_PORT)
    if CONFIG.ENABLE_TLS {
        log.Println("TLS is Enabled")
    } else {
        log.Println("HTTP-Only mode")
    }
    raw_config = readConfig(API_KEYFILE)
    unmarshal(raw_config, &API_KEYS)
    log.Println("Loaded API keys")
    owm.Config(API_KEYS.OPENWEATHERMAP, CONFIG.UNITS, CONFIG.COUNTRY_CODE)
    om.Config(""                     , CONFIG.UNITS, CONFIG.MODEL)
}

func readConfig(fileName string) []byte {
    raw_config, err := os.ReadFile(fileName)
    if err != nil {
        log.Fatal(err)
    }
    return raw_config
}

func unmarshal(data []byte, config any) {
    err := json.Unmarshal(data, &config)
    if err != nil {
        log.Fatal(err)
    }
}
