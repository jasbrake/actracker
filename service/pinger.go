package service

import (
	"fmt"
	"strings"

	"github.com/jasbrake/acpinger"
	"github.com/jasbrake/actracker/model"
)

// StartPinger starts the pinging service
func StartPinger(c *model.Config) {
	for i := 0; i < c.PingerCount; i++ {
		go pinger(c.Pinging, c.Updates)
	}
}

func pinger(in <-chan *model.Server, out chan<- *model.Server) {
	for s := range in {
		err := pingServer(s)
		handlePingError(s, err)
		out <- s
	}
}

func pingServer(s *model.Server) error {
	std, err := acpinger.PingStd(s.IP.String(), s.Port, 0)
	if err != nil {
		return err
	}
	s.Description = std.Description
	s.MaxClients = std.MaxClients

	ext, err := acpinger.PingExt(s.IP.String(), s.Port, 0)
	if err != nil {
		return err
	}

	s.TimeoutCount = 0

	players := make([]model.Player, 0)
	if ext.Players != nil {
		for _, p := range ext.Players {
			player := model.NewPlayer(p)
			player.UpdateLocation()
			players = append(players, player)
		}
	}

	s.Game = model.Game{
		Key:              s.Key,
		Mode:             std.Mode,
		PlayerCount:      std.PlayerCount,
		MinutesRemaining: std.MinutesRemaining,
		Map:              std.CurrentMap,
		Mastermode:       std.Mastermode,
		Password:         std.Password,
		Players:          players,
	}
	return nil
}

func handlePingError(s *model.Server, err error) {
	if err != nil {
		e := err.Error()

		if strings.Contains(e, "read: connection refused") ||
			strings.Contains(e, "read: no route to host") {
			s.Dead = true
			return
		}

		if strings.Contains(e, "i/o timeout") {
			s.TimeoutCount++
			return
		}

		fmt.Printf("ping error (%s): %s\n", s.Key, err.Error())
	}
}
