# AC Tracker

## Features
- Live game details
- Past games and player scores saved to DB

### Future Features
- Search (by name, server, IP range, date range)
- Show past games by server
- Player ladder
- Highlight and separate competitive games ("inters") and tournaments from public games

## API

All endpoints return JSON.

### `/servers`

Returns a list of servers with active games along with the live game details.

## Development

In order to run for development, I use `go run main.go` with the following environment vars set. You might benefit from putting this into a bash script to run it easily.

```
GIN_MODE=debug
GEOIP_DB_PATH=/path/to/geoip-dbs/GeoLite2-Country.mmdb
PINGER_COUNT=100
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=actracker
DB_SSLMODE=disable
APP_PORT=3000
SLEEP_SECONDS=10
MAX_SLEEP_SECONDS=600
```
