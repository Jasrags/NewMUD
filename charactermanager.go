package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	CharacterMgr = NewCharacterManager()
)

type CharacterManager struct {
	sync.RWMutex
	characters map[string]*Character
}

func NewCharacterManager() *CharacterManager {
	return &CharacterManager{
		characters: make(map[string]*Character),
	}
}

func (mgr *CharacterManager) AddCharacter(c *Character) {
	slog.Debug("Adding character",
		slog.String("character_id", c.ID))

	mgr.Lock()
	defer mgr.Unlock()

	mgr.characters[strings.ToLower(c.Name)] = c
}

func (mgr *CharacterManager) GetCharacterByName(name string) *Character {
	slog.Debug("Getting character by name",
		slog.String("character_name", name))

	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.characters[strings.ToLower(name)]
}

func (mgr *CharacterManager) RemoveCharacter(c *Character) {
	slog.Debug("Removing character",
		slog.String("character_name", c.Name))

	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.characters, strings.ToLower(c.Name))
}

func (mgr *CharacterManager) LoadDataFiles() {
	dataFilePath := viper.GetString("data.characters_path")
	slog.Info("Loading character data files",
		slog.String("datafile_path", dataFilePath))

	files, err := os.ReadDir(dataFilePath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", dataFilePath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(dataFilePath, file.Name())

			var c Character
			if err := LoadJSON(filePath, &c); err != nil {
				slog.Error("failed to load character data",
					slog.Any("error", err),
					slog.String("file", file.Name()))
			}

			mgr.AddCharacter(&c)

			slog.Debug("Loaded character",
				slog.String("id", c.ID),
				slog.String("name", c.Name))
		}
	}

	slog.Info("Loaded characters",
		slog.Int("count", len(mgr.characters)))
}
