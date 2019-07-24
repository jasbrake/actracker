package model

import (
	"sync"
)

// StateManager holds the in-memory state
type StateManager struct {
	sync.RWMutex
	ServersMap map[string]Server
}

// NewStateManager creates a new StateManager
func NewStateManager() *StateManager {
	return &StateManager{
		ServersMap: make(map[string]Server),
	}
}
