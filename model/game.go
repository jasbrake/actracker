package model

import (
	"log"
	"time"

	"github.com/jasbrake/actracker/db"
)

// Game represents a Server's current game state
type Game struct {
	ID               int
	Key              string    `json:"key"`
	Mode             int       `json:"mode"`
	PlayerCount      int       `json:"player_count" db:"player_count"`
	MinutesRemaining int       `json:"minutes_remaining" db:"minutes_remaining"`
	Map              string    `json:"map" db:"map"`
	Mastermode       int       `json:"mastermode"`
	Password         bool      `json:"password"`
	TimeEnded        time.Time `json:"time_ended" db:"time_ended"`
	Players          []Player  `json:"players"`
}

func (g Game) Save() error {
	statement, err := db.Conn.PrepareNamed(`INSERT INTO game
    (key, mode, player_count, map, mastermode, time_ended)
    VALUES (:key, :mode, :player_count, :map, :mastermode, :time_ended)
    RETURNING id`)
	if err != nil {
		return err
	}

	g.TimeEnded = time.Now().UTC()
	err = statement.Get(&g.ID, g)
	if err != nil {
		return err
	}

	for _, p := range g.Players {
		err := p.SaveForGame(g.ID)
		if err != nil {
			log.Printf("error saving player for game %+v: %s", p, err)
		}
	}
	return nil
}
