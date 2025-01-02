package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	UserMgr = NewUserManager()
)

type UserManager struct {
	sync.RWMutex

	onlineUsers map[string]*User
	users       map[string]*User
	bannedNames map[string]bool
}

func NewUserManager() *UserManager {
	um := &UserManager{
		onlineUsers: make(map[string]*User),
		users:       make(map[string]*User),
		bannedNames: make(map[string]bool),
	}

	return um
}

func (mgr *UserManager) AddUser(u *User) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding user",
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	mgr.users[strings.ToLower(u.Username)] = u
}

func (mgr *UserManager) GetUserByID(userID string) *User {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Getting user by ID",
		slog.String("id", userID))

	return mgr.users[userID]
}

func (mgr *UserManager) GetByUsername(username string) *User {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Getting user by username",
		slog.String("username", username))

	return mgr.users[strings.ToLower(username)]
}
func (mgr *UserManager) RemoveUser(u *User) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Removing user",
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	delete(mgr.users, strings.ToLower(u.Username))
}

func (mgr *UserManager) SetOnline(u *User) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Setting user online",
		slog.String("id", u.ID),
		slog.String("username", u.Username))
	t := time.Now()
	u.LastLoginAt = &t
	u.Save()

	mgr.onlineUsers[u.ID] = u
}

func (mgr *UserManager) SetOffline(u *User) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Setting user offline",
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	delete(mgr.onlineUsers, u.ID)
}

func (mgr *UserManager) Exists(username string) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Checking if user exists",
		slog.String("username", username))

	return mgr.users[strings.ToLower(username)] != nil
}

func (mgr *UserManager) IsBannedName(name string) bool {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Checking if user name is banned",
		slog.String("username", name))

	return mgr.bannedNames[strings.ToLower(name)]
}

func (mgr *UserManager) LoadDataFiles() {
	dataFilePath := viper.GetString("data.users_path")
	bannedNames := viper.GetStringSlice("banned_names")

	// Load banned names
	slog.Info("Loading banned user names")
	for _, name := range bannedNames {
		mgr.bannedNames[strings.ToLower(name)] = true
	}

	slog.Info("Loading user data files",
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

			var u User
			if err := LoadJSON(filePath, &u); err != nil {
				slog.Error("failed to unmarshal user data",
					slog.Any("error", err),
					slog.String("file", file.Name()))
			}

			mgr.AddUser(&u)

			slog.Debug("Loaded user",
				slog.String("id", u.ID),
				slog.String("username", u.Username))
		}
	}

	slog.Info("Loaded users",
		slog.Int("count", len(mgr.users)))
}
