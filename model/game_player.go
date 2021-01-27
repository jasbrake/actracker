package model

import (
	"time"

	"github.com/jasbrake/actracker/db"
	"github.com/jasbrake/actracker/lib/pgnet"
)

// GamePlayer is a combination of the Game and Player models
type GamePlayer struct {
	GameID      int        `json:"game_id" db:"game_id"`
	Key         string     `json:"key"`
	Mode        int        `json:"mode"`
	PlayerCount int        `json:"player_count" db:"player_count"`
	Map         string     `json:"map" db:"map"`
	Mastermode  int        `json:"mastermode"`
	Password    bool       `json:"password"`
	TimeEnded   time.Time  `json:"time_ended" db:"time_ended"`
	Name        string     `json:"name"`
	IP          pgnet.CIDR `json:"ip"`
	Country     string     `json:"country"`
	CountryISO  string     `json:"country_iso" db:"country_iso"`
	Team        string     `json:"team"`
	Kills       int        `json:"kills"`
	Deaths      int        `json:"deaths"`
	Flags       int        `json:"flags"`
	Teamkills   int        `json:"teamkills"`
	Accuracy    int        `json:"accuracy"`
	GunSelected int        `json:"gun_selected" db:"gun_selected"`
}

// GetPlayerGames gets games for a player
func GetPlayerGames(name string) ([]GamePlayer, error) {
	query := `SELECT g.id as game_id, g."key", g."map", g.mastermode, g."mode", g.player_count, g.time_ended, p.name, p.ip, p.country, p.country_iso, p.team, p.kills, p.deaths, p.teamkills, p.accuracy, p.gun_selected FROM game_player AS p
JOIN game AS g ON p.game_id = g.id
WHERE name=$1
ORDER BY g.time_ended DESC
LIMIT 10;
`
	players := make([]GamePlayer, 0)
	rows, err := db.Conn.Queryx(query, name)
	if err != nil {
		return players, err
	}

	for rows.Next() {
		p := GamePlayer{}
		err := rows.StructScan(&p)
		if err != nil {
			return players, err
		}
		players = append(players, p)
	}
	return players, nil
}
