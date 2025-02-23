# Low Latency Weather Server

A weather server, with *golang* backend and *JS* front.

Intended to be dockerized

## Requirements

- [go](https://go.dev/) 1.19 or newer

## Planned and Completed Features

- [x] Caching weather data
- [x] Light weight
- [x] Show weather preview for noon and night (22:00), for astronomy purposes
- [ ] Seamless update of data, but cached data for lighting fast initial response
- [ ] Expand day previews to show detailed forecast for the day
- [x] Search for cities within Finland
- [ ] Show at a glance
  - [x] Temperature
  - [x] Total cloud %
  - [x] Humidity data
  - [x] Weather Icons
  - [ ] Wind direction and speed
- [ ] Location in URL for easy bookmarking
- [x] Configfile
- [ ] Dockerized
- [ ] (uses multiple data sources?)

## Configuration

Config files `key.txt` and `config.txt` should be placed in the server root directory.

| Desc. | filename | Notes |
| -- | -- | -- |
| Keyfile | `key.txt` | Contains the API key for Openweathermap |
| Configfile | `config.txt` | Contains the base URL:Port for the server |

###### Example config
```
http://localhost:3000
```

## Running server

As long as [*go*](https://go.dev/) is installed, you can run the server with this command
```
go run .
```

For long term deployments you should build the server
```
go build
```

and then run the generated executable `ll-weather-server`
