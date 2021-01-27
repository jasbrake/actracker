package model

import (
	"log"
	"net"

	"github.com/jasbrake/acpinger"
	"github.com/jasbrake/actracker/db"
	"github.com/jasbrake/actracker/geoip"
	"github.com/jasbrake/actracker/lib/pgnet"
)

type Player struct {
	Name         string     `json:"name"`
	IP           pgnet.CIDR `json:"ip"`
	Country      string     `json:"country"`
	CountryISO   string     `json:"country_iso" db:"country_iso"`
	ClientNumber int        `json:"clientnumber"`
	Ping         int        `json:"ping"`
	Team         string     `json:"team"`
	Kills        int        `json:"kills"`
	Deaths       int        `json:"deaths"`
	Flags        int        `json:"flags"`
	Teamkills    int        `json:"teamkills"`
	Accuracy     int        `json:"accuracy"`
	Health       int        `json:"health"`
	Armour       int        `json:"armour"`
	GunSelected  int        `json:"gun_selected" db:"gun_selected"`
	Role         int        `json:"role"`
	State        int        `json:"state"`
}

func NewPlayer(p acpinger.Player) Player {
	_, cidr, err := net.ParseCIDR(p.IP)
	if err != nil {
		log.Printf("Failed to parse IP for player %+v: %s\n", p, err)
	}

	return Player{
		Name:         p.Name,
		IP:           pgnet.CIDR{IPNet: cidr},
		ClientNumber: p.ClientNumber,
		Ping:         p.Ping,
		Team:         p.Team,
		Kills:        p.Frags,
		Deaths:       p.Deaths,
		Flags:        p.Flagscore,
		Teamkills:    p.Teamkills,
		Accuracy:     p.Accuracy,
		Health:       p.Health,
		Armour:       p.Armour,
		GunSelected:  p.GunSelected,
		Role:         p.Role,
		State:        p.State,
	}
}

func (p Player) SaveForGame(game int) error {
	statement := `INSERT INTO game_player (game_id, name, ip, country, country_iso, team, kills, deaths, flags, teamkills, accuracy, gun_selected)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := db.Conn.Exec(statement, game, p.Name, p.IP, p.Country, p.CountryISO, p.Team, p.Kills, p.Deaths, p.Flags, p.Teamkills, p.Accuracy, p.GunSelected)
	return err
}

func (p *Player) UpdateLocation() {
	geo, err := geoip.DB.Country(p.IP.IP)
	if err != nil {
		log.Printf("failed to find GeoIP for %s: %s", p.IP, err)
	} else {
		p.Country = geo.Country.Names["en"]
		p.CountryISO = geo.Country.IsoCode
	}
}

// GetPlayerNames gets all players names that start with the provided pattern
func GetPlayerNames(startsWith string) ([]string, error) {
	var names []string

	query := "SELECT DISTINCT(name) FROM game_player WHERE name LIKE $1 || '%' LIMIT 10;"

	stmt, err := db.Conn.Preparex(query)
	if err != nil {
		return names, err
	}

	err = stmt.Select(&names, startsWith)
	return names, err
}
