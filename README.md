# Low Latency Weather Server

A weather server, with *golang* backend and JS front.

Intended to be dockerized


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

| Desc. | filename | Notes |
| -- | -- | -- |
| Keyfile | `key.txt` | Contains the API key for Openweathermap |
| Configfile | `config.txt` | Contains the base URL:Port for the server |

###### Example config
```
http://localhost:3000
```
