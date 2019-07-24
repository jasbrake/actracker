package service

import (
	"fmt"
	"log"

	"github.com/jasbrake/acpinger"
	"github.com/jasbrake/actracker/geoip"
	"github.com/jasbrake/actracker/model"
	"github.com/jasbrake/actracker/state"
)

// StartUpdate starts the update handler service.
// This service determines what to do with server updates from the pinger.
func StartUpdate(c *model.Config) {
	go updates(c.Pinging, c.Sleeping)
}

func updates(in <-chan *model.Server, out chan<- *model.Server) {
	for {
		s := <-in
		if s.TimeoutCount == 0 {
			updateCountry(s)
			err := s.Save()
			if err != nil {
				log.Println(err)
			}

			old, ok := state.GetServer(s.Key)
			if ok {
				handleGameEnd(s, old)
			}
			sCopy := *s
			sCopy.Game.Players = append([]acpinger.Player(nil), s.Game.Players...)
			state.SaveServer(sCopy)
		}
		out <- s
	}
}

func updateCountry(s *model.Server) {
	geo, err := geoip.DB.Country(s.IP.IP)
	if err != nil {
		log.Printf("failed to find GeoIP for %s: %s", s.IP, err)
	} else {
		s.Country = geo.Country.Names["en"]
		s.CountryISO = geo.Country.IsoCode
	}
}

func handleGameEnd(s *model.Server, old model.Server) {
	// Don't save a server when going from an invalid game state to a valid one
	if old.Game.CurrentMap != "" {
		// We consider the game as having ended if either:
		// 1. The new map is different from the old map.
		// OR
		// 2. The new map is the same, but the time remaining has increased
		//    (meaning the same map was voted again).
		if s.Game.CurrentMap != old.Game.CurrentMap ||
			s.Game.MinutesRemaining > old.Game.MinutesRemaining {
			// 1 minute remaining or less and at least 1 player
			if old.Game.MinutesRemaining <= 1 && len(old.Game.Players) > 0 {
				// Save game
				fmt.Printf("Game saved: %+v\n\n%+v\n", old, s)
			}
		}
	}
}
