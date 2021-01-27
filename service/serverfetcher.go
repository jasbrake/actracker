package service

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jasbrake/actracker/model"
)

var (
	queued map[string]bool
)

func init() {
	queued = make(map[string]bool)
}

// StartServerFetcher creates the server fetching service.
// It fetches servers from the DB and from the masterserver.
func StartServerFetcher(c *model.Config) {
	go func() {
		for {
			queueServers(c)
			time.Sleep(time.Duration(12) * time.Hour)
		}
	}()
}

func queueServers(c *model.Config) {
	stored, err := model.GetServers()
	if err != nil {
		log.Println(err)
	}
	remote, err := fetchServersFromMS()
	if err != nil {
		log.Println(err)
	}
	servers := append(stored, remote...)

	for _, s := range servers {
		// Queue all servers we don't have
		_, ok := queued[s.Key]
		if !ok {
			queued[s.Key] = true
			c.Pinging <- s
		}
	}
}

func fetchServersFromMS() ([]*model.Server, error) {
	servers := make([]*model.Server, 0)
	resp, err := http.Get("http://ms.cubers.net/retrieve.do?action=list&name=actracker&version=1202&build=68644")
	if err != nil {
		return servers, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return servers, err
	}
	if len(body) == 0 {
		return servers, errors.New("response body empty")
	}

	lines := strings.Split(string(body), "\n")

	for _, l := range lines {
		words := strings.Split(l, " ")
		if words[0] == "addserver" {
			ip := words[1]
			port, err := strconv.Atoi(words[2])
			if err != nil {
				return servers, err
			}
			s := model.NewServer(ip, port)
			s.OnMasterServer = true

			// Only add the server if the IP was validly parsed
			if s.IP.IP != nil {
				servers = append(servers, s)
			}
		}
	}
	return servers, nil
}
