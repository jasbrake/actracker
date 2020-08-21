package geoip

import (
	"log"
	"os"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	DB *geoip2.Reader
)

// Setup initializes the GeoIP DB
func Setup() {
	geoipDBPath := os.Getenv("GEOIP_DB_PATH")
	var err error
	DB, err = geoip2.Open(geoipDBPath)
	if err != nil {
		log.Fatalf("Failed to read GeoIP DB: %s", err)
	}
}
