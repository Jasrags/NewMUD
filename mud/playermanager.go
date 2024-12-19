package mud

import (
	"sync"

	"github.com/rs/zerolog"
)

type PlayerManager struct {
	mu      sync.RWMutex
	Log     zerolog.Logger
	Players map[string]*Player // Keyed by player name
}

// NewPlayerManager creates and initializes a PlayerManager
func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		Log:     NewDevLogger(),
		Players: make(map[string]*Player),
	}
}

// AddPlayer adds a player to the manager
func (pm *PlayerManager) AddPlayer(player *Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.Players[player.Name] = player
}

// RemovePlayer removes a player from the manager by name
func (pm *PlayerManager) RemovePlayer(name string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.Players, name)
}

// GetPlayer retrieves a player by name
func (pm *PlayerManager) GetPlayer(name string) *Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.Players[name]
}

// GetAllPlayers returns a list of all players
func (pm *PlayerManager) GetAllPlayers() []*Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var allPlayers []*Player
	for _, player := range pm.Players {
		allPlayers = append(allPlayers, player)
	}
	return allPlayers
}
