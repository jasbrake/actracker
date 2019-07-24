package service

import (
	"time"

	"github.com/jasbrake/actracker/model"
)

// StartSleep creates the sleep service that waits a duration before sending it back to the pinging chan
func StartSleep(c *model.Config) {
	go func() {
		for s := range c.Sleeping {
			t := c.SleepSeconds
			if s.TimeoutCount > 0 {
				t *= s.TimeoutCount
			}
			if t > c.MaxSleepSeconds {
				t = c.MaxSleepSeconds
			}
			go sleepServer(s, time.Duration(t)*time.Second, c.Pinging)
		}
	}()
}

func sleepServer(s *model.Server, t time.Duration, out chan<- *model.Server) {
	time.Sleep(t)
	out <- s
}
