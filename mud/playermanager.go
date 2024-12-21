package mud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type PlayerManager struct {
	mu      sync.RWMutex
	Log     zerolog.Logger
	Players map[string]*Player
}

// NewPlayerManager creates and initializes a PlayerManager
func NewPlayerManager(l zerolog.Logger) *PlayerManager {
	return &PlayerManager{
		Log:     l,
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

// TODO: We need to add conn to the new player but should that be part of the load or after the load?
func (pm *PlayerManager) LoadPlayer(name string, force bool) {
	pm.Log.Info().
		Str("player_name", name).
		Bool("force", force).
		Msg("Loading player")

	dataPath := viper.GetString("data.players_path")
	filePath := fmt.Sprintf("%s/%s.json", dataPath, strings.ToLower(name))
	file, errOpen := os.Open(filePath)
	if errOpen != nil {
		pm.Log.Error().Err(errOpen).Msg("Failed to open player file")

		return
	}
	defer file.Close()

	player := &Player{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(player); err != nil {
		pm.Log.Error().Err(err).Msg("Failed to decode player file")
		return
	}

	pm.AddPlayer(player)
}

// Save all players
func (pm *PlayerManager) Save() {
	pm.Log.Info().Msg("Saving all players")
	for _, p := range pm.Players {
		p.Save()
	}
}

// Load all players
// func (pm *PlayerManager) Load() {
// pm.Log.Info().Msg("Loading players")
// }
