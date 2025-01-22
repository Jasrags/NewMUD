package main

// import (
// 	"log/slog"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"

// 	"github.com/spf13/viper"
// )

// var (
// 	CharacterMgr = NewCharacterManager()
// )

// type CharacterManager struct {
// 	sync.RWMutex

// 	characters       map[string]*Character
// 	onlineCharacters map[string]*Character
// 	bannedNames      map[string]bool
// }

// func NewCharacterManager() *CharacterManager {
// 	return &CharacterManager{
// 		characters:       make(map[string]*Character),
// 		onlineCharacters: make(map[string]*Character),
// 		bannedNames:      make(map[string]bool),
// 	}
// }

// func (mgr *CharacterManager) GetOnlineCharacters() map[string]*Character {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	return mgr.onlineCharacters
// }

// func (mgr *CharacterManager) SetCharacterOnline(c *Character) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	slog.Debug("Setting character online",
// 		slog.String("character_id", c.ID))

// 	mgr.onlineCharacters[strings.ToLower(c.Name)] = c
// }

// func (mgr *CharacterManager) SetCharacterOffline(c *Character) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	slog.Debug("Setting character offline",
// 		slog.String("character_id", c.ID))

// 	delete(mgr.onlineCharacters, strings.ToLower(c.Name))
// }

// func (mgr *CharacterManager) AddCharacter(c *Character) {
// 	slog.Debug("Adding character",
// 		slog.String("character_id", c.ID))

// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	mgr.characters[strings.ToLower(c.Name)] = c
// }

// func (mgr *CharacterManager) GetCharacterByName(name string) *Character {
// 	slog.Debug("Getting character by name",
// 		slog.String("character_name", name))

// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	return mgr.characters[strings.ToLower(name)]
// }

// func (mgr *CharacterManager) RemoveCharacter(c *Character) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	slog.Debug("Removing character",
// 		slog.String("character_name", c.Name))

// 	delete(mgr.characters, strings.ToLower(c.Name))
// }
// func (mgr *CharacterManager) Exists(name string) bool {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	slog.Debug("Checking if character exists",
// 		slog.String("character_name", name))

// 	return mgr.characters[strings.ToLower(name)] != nil
// }

// func (mgr *CharacterManager) IsBannedName(name string) bool {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	slog.Debug("Checking if character name is banned",
// 		slog.String("character_name", name))

// 	return mgr.bannedNames[strings.ToLower(name)]
// }

// func (mgr *CharacterManager) LoadDataFiles() {
// 	dataFilePath := viper.GetString("data.characters_path")
// 	bannedNames := viper.GetStringSlice("banned_names")

// 	// Load banned names
// 	slog.Info("Loading banned character names")
// 	for _, name := range bannedNames {
// 		mgr.bannedNames[strings.ToLower(name)] = true
// 	}

// 	slog.Info("Loading character data files",
// 		slog.String("datafile_path", dataFilePath))

// 	files, err := os.ReadDir(dataFilePath)
// 	if err != nil {
// 		slog.Error("failed reading directory",
// 			slog.String("datafile_path", dataFilePath),
// 			slog.Any("error", err))
// 	}

// 	for _, file := range files {
// 		if filepath.Ext(file.Name()) == ".yml" {
// 			filePath := filepath.Join(dataFilePath, file.Name())

// 			var c Character
// 			if err := LoadYAML(filePath, &c); err != nil {
// 				slog.Error("failed to load character data",
// 					slog.Any("error", err),
// 					slog.String("file", file.Name()))
// 			}

// 			mgr.AddCharacter(&c)

// 			slog.Debug("Loaded character",
// 				slog.String("id", c.ID),
// 				slog.String("name", c.Name))
// 		}
// 	}

// 	slog.Info("Loaded characters",
// 		slog.Int("count", len(mgr.characters)))
// }
