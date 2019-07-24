package state

import (
	"github.com/jasbrake/actracker/model"
)

var (
	// Manager is responsible for handling the in-memory server states
	Manager *model.StateManager
)

func init() {
	Manager = model.NewStateManager()
}

// GetActiveServers returns all servers with a valid game
func GetActiveServers() []model.Server {
	servers := make([]model.Server, 0)
	Manager.RLock()
	for _, v := range Manager.ServersMap {
		if v.Game.PlayerCount > 0 {
			servers = append(servers, v)
		}
	}
	Manager.RUnlock()
	return servers
}

// GetServer retrieves a specific server from the state by key
func GetServer(key string) (server model.Server, ok bool) {
	Manager.RLock()
	s, ok := Manager.ServersMap[key]
	Manager.RUnlock()
	return s, ok
}

// SaveServer saves or overwrites a server in the state
func SaveServer(s model.Server) {
	Manager.Lock()
	Manager.ServersMap[s.Key] = s
	Manager.Unlock()
}

// DeleteServer removes a server from the state
func DeleteServer(key string) {
	Manager.Lock()
	delete(Manager.ServersMap, key)
	Manager.Unlock()
}
