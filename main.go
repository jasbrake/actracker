package main

import (
	"github.com/jasbrake/actracker/api"
	"github.com/jasbrake/actracker/db"
	"github.com/jasbrake/actracker/geoip"
	"github.com/jasbrake/actracker/model"
	"github.com/jasbrake/actracker/service"
)

func main() {
	db.Setup()
	geoip.Setup()
	config := model.InitConfig()
	service.StartSleep(config)
	service.StartUpdate(config)
	service.StartPinger(config)
	service.StartServerFetcher(config)
	api.Start()
}
