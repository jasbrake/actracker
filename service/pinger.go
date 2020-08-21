package service

import (
	"log"
	"net"
	"os"
	"syscall"

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
		logPingError(err)
		out <- s
	}
}

func pingServer(s *model.Server) error {
	std, err := acpinger.PingStd(s.IP.String(), s.Port, 0)
	if err != nil {
		s.TimeoutCount++
		return err
	}
	s.Description = std.Description
	s.MaxClients = std.MaxClients

	ext, err := acpinger.PingExt(s.IP.String(), s.Port, 0)
	if err != nil {
		s.TimeoutCount++
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

func logPingError(err error) {
	if err != nil {
		// Log the error if it isn't a timeout or "connection refused"
		// This is pretty gross but besides just ignoring the error, I'm not sure
		// how best to handle servers being offline or timing out.
		if opErr, ok := err.(*net.OpError); ok {
			if syscallErr, ok := opErr.Err.(*os.SyscallError); ok {
				if syscallErr.Err != syscall.ECONNREFUSED {
					log.Println(syscallErr)
				}
			} else if nerr, ok := err.(net.Error); ok && !nerr.Timeout() {
				log.Println(nerr)
			}
		} else {
			log.Println(err)
		}
	}
}
