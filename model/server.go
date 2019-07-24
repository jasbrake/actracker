package model

import (
	"fmt"
	"log"
	"net"

	"github.com/jasbrake/actracker/db"
	"github.com/jasbrake/actracker/lib/pgnet"
)

// Server represents a server and it's current game status
type Server struct {
	Key            string     `json:"key"`
	IP             pgnet.INET `json:"ip"`
	Port           int        `json:"port"`
	Description    string     `json:"description"`
	MaxClients     int        `json:"max_clients" db:"max_clients"`
	CountryISO     string     `json:"country_iso" db:"country_iso"`
	Country        string     `json:"country"`
	Match          bool       `json:"match"`
	CustomConnect  string     `json:"custom_connect" db:"custom_connect"`
	OnMasterServer bool       `json:"on_ms" db:"on_ms"`
	Disabled       bool       `json:"disabled"`
	Dead           bool       `json:"dead"`
	TimeoutCount   int        `json:"timeout_count"`
	Game           Game       `json:"game"`
}

// NewServer creates a server with the key in the ip:port format
func NewServer(ip string, port int) *Server {
	return &Server{
		Key:  fmt.Sprintf("%s:%d", ip, port),
		IP:   pgnet.INET{IP: net.ParseIP(ip)},
		Port: port,
	}
}

// Save saves the server into the DB
func (s *Server) Save() error {
	statement := `
	INSERT INTO servers (key, ip, port, description, max_clients, country, country_iso, match, custom_connect, on_ms, disabled, dead)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	ON CONFLICT (key) DO UPDATE SET
	description=EXCLUDED.description, max_clients=EXCLUDED.max_clients`
	_, err := db.Conn.Exec(statement, s.Key, s.IP, s.Port, s.Description, s.MaxClients, s.Country, s.CountryISO, s.Match, s.CustomConnect, s.OnMasterServer, s.Disabled, s.Dead)
	if err != nil {
		return err
	}
	return nil
}

// GetServers loads the servers from the DB
func GetServers() ([]*Server, error) {
	query := `SELECT * FROM servers`
	servers := make([]*Server, 0)
	rows, err := db.Conn.Queryx(query)
	if err != nil {
		return servers, err
	}

	for rows.Next() {
		s := &Server{}
		err := rows.StructScan(s)
		if err != nil {
			log.Println(err)
		}
		servers = append(servers, s)
	}
	return servers, nil
}
