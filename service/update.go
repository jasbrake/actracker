package service

import (
	"log"

	"github.com/jasbrake/actracker/model"
	"github.com/jasbrake/actracker/state"
)

// StartUpdate starts the update handler service.
// This service determines what to do with server updates from the pinger.
func StartUpdate(c *model.Config) {
	go updates(c.Updates, c.Sleeping)
}

func updates(in <-chan *model.Server, out chan<- *model.Server) {
	for {
		s := <-in
		if s.TimeoutCount == 0 {
			s.UpdateLocation()
			err := s.Save()
			if err != nil {
				log.Println(err)
			}

			old, ok := state.GetServer(s.Key)
			if ok {
				handleGameEnd(s, old)
			}
			state.SaveServer(*s)
		}
		out <- s
	}
}

func handleGameEnd(s *model.Server, old model.Server) {
	// Don't save a server when going from an invalid game state to a valid one
	if old.Game.Map != "" {
		// We consider the game as having ended if either:
		// 1. The new map is different from the old map.
		// OR
		// 2. The new map is the same, but the time remaining has increased
		//    (meaning the same map was voted again).
		if s.Game.Map != old.Game.Map ||
			s.Game.MinutesRemaining > old.Game.MinutesRemaining {
			// 1 minute remaining or less and at least 1 player
			if old.Game.MinutesRemaining <= 1 && len(old.Game.Players) > 0 {
				err := old.Game.Save()
				if err != nil {
					log.Printf("could not save game %+v: %s\n", old.Game, err)
				}
			}
		}
	}
}
