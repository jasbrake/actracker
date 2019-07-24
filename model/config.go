package model

import (
	"log"
	"os"
	"strconv"
)

// Config is the main struct to hold all global config and services
type Config struct {
	Pinging         chan *Server
	Sleeping        chan *Server
	Updates         chan *Server
	PingerCount     int
	SleepSeconds    int
	MaxSleepSeconds int
}

// InitConfig initializes the config from the environment or DB
func InitConfig() *Config {
	c := &Config{
		Pinging:  make(chan *Server),
		Sleeping: make(chan *Server),
		Updates:  make(chan *Server),
	}
	c.PingerCount = mustGetenvInt("PINGER_COUNT")
	c.SleepSeconds = mustGetenvInt("SLEEP_SECONDS")
	c.MaxSleepSeconds = mustGetenvInt("MAX_SLEEP_SECONDS")
	return c
}

func mustGetenvInt(v string) int {
	s := os.Getenv(v)
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("could not parse %s (%s) as an int\n", v, s)
	}
	return n
}
