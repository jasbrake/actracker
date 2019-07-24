package model

import (
	"github.com/jasbrake/acpinger"
)

// Game represents a Server's current game state
type Game struct {
	ID               int
	Key              string            `json:"key"`
	Mode             int               `json:"mode"`
	PlayerCount      int               `json:"player_count"`
	MinutesRemaining int               `json:"minutes_remaining"`
	CurrentMap       string            `json:"current_map"`
	Mastermode       int               `json:"mastermode"`
	Password         bool              `json:"password"`
	Players          []acpinger.Player `json:"players"`
}
