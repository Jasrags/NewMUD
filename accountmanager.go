package main

// import (
// 	"log/slog"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/spf13/viper"
// )

// var (
// 	AccountMgr = NewAccountManager()
// )

// type AccountManager struct {
// 	sync.RWMutex

// 	online      map[string]*Account
// 	accounts    map[string]*Account
// 	bannedNames map[string]bool
// }

// func NewAccountManager() *AccountManager {
// 	um := &AccountManager{
// 		online:      make(map[string]*Account),
// 		accounts:    make(map[string]*Account),
// 		bannedNames: make(map[string]bool),
// 	}

// 	return um
// }

// func (mgr *AccountManager) AddAccount(u *Account) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	mgr.accounts[strings.ToLower(u.Username)] = u
// }

// func (mgr *AccountManager) GetAccountByID(userID string) *Account {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	return mgr.accounts[userID]
// }

// func (mgr *AccountManager) GetByUsername(username string) *Account {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	return mgr.accounts[strings.ToLower(username)]
// }
// func (mgr *AccountManager) RemoveAccount(u *Account) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	delete(mgr.accounts, strings.ToLower(u.Username))
// }

// func (mgr *AccountManager) SetOnline(u *Account) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	t := time.Now()
// 	u.LastLoginAt = &t
// 	u.Save()

// 	mgr.online[u.ID] = u
// }

// func (mgr *AccountManager) SetOffline(u *Account) {
// 	mgr.Lock()
// 	defer mgr.Unlock()

// 	delete(mgr.online, u.ID)
// }

// func (mgr *AccountManager) Exists(username string) bool {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	return mgr.accounts[strings.ToLower(username)] != nil
// }

// func (mgr *AccountManager) IsBannedName(name string) bool {
// 	mgr.RLock()
// 	defer mgr.RUnlock()

// 	return mgr.bannedNames[strings.ToLower(name)]
// }

// func (mgr *AccountManager) LoadDataFiles() {
// 	dataFilePath := viper.GetString("data.accounts_path")
// 	bannedNames := viper.GetStringSlice("banned_names")

// 	// Load banned names
// 	slog.Info("Loading banned user names")
// 	for _, name := range bannedNames {
// 		mgr.bannedNames[strings.ToLower(name)] = true
// 	}

// 	slog.Info("Loading user data files",
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

// 			var u Account
// 			if err := LoadYAML(filePath, &u); err != nil {
// 				slog.Error("failed to unmarshal user data",
// 					slog.Any("error", err),
// 					slog.String("file", file.Name()))
// 			}

// 			mgr.AddAccount(&u)

// 			slog.Debug("Loaded user",
// 				slog.String("id", u.ID),
// 				slog.String("username", u.Username))
// 		}
// 	}

// 	slog.Info("Loaded users",
// 		slog.Int("count", len(mgr.accounts)))
// }
