# AC Tracker

A service that periodically pings [AssaultCube](https://assault.cubers.net/) servers, provides an API for live games and tracks past games.

The frontend is located in a separate repository [here](https://github.com/jasbrake/actracker-web).

## Features
- Live game details
- Past games and player scores saved to DB

### Future Features
- Search (by name, server, IP range, date range)
- Show past games by server
- Player ladder
- Highlight and separate competitive games ("inters") and tournaments from public games
- Parse clan tags out of names to show online clan members

## API

All endpoints return JSON.

### `/servers`

Returns a list of servers with active games along with the live game details.

### `/player/:name`

Returns a list of recent games for a player name.

### `/player_autocomplete`

Expects a query parameter `q`.

Returns a list of player names that start with the query.


## Development

In order to run for development, create a `.env` file with your environment variables then run `./dev.sh`.

Example `.env` file:

```
# APP
APP_PORT=3000
PINGER_COUNT=50
SLEEP_SECONDS=10
MAX_SLEEP_SECONDS=300
#GIN_MODE="release"

# DB
DB_HOST="localhost"
DB_PORT=5432
DB_USER="user"
DB_PASSWORD="password"
DB_NAME="actracker"
DB_SSLMODE="require"

# GEOIP
GEOIP_DB_PATH="/usr/share/GeoIP/GeoLite2-Country.mmdb"
```
