# Low Latency Weather Server

A weather server, with **Go** backend and **JS** front.

## Requirements

- [**Go**](https://go.dev/) 1.19 or newer

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
| Keyfile | `key.txt` | Should contain your API key for Openweathermap |
| Configfile | `config.json` | Server configuration |

###### Example config

```json
{
  "ORIGIN_URL": "http://localhost",
  "SERVER_PORT": "3000",
  "ENABLE_TLS": false,
  "CERT_FILE": "",
  "PRIVATE_KEY_FILE": ""
}
```

`"ORIGIN_URL"` should be the address the server is accessed at

#### HTTPS/TLS

To enable HTTPS, you need a valid certificate file and private key file.
Set `"ENABLE_TLS": true` and give the cert and private key file names to `"CERT_FILE":` and `"PRIVATE_KEY_FILE":` respectively.

## Running server

As long as [**Go**](https://go.dev/) is installed, you can run the server with this command
```
go run .
```

For long term deployments you should build the server
```
go build
```

and then run the generated executable `ll-weather-server`
