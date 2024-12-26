package users

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	Mgr = NewManager()
)

type Manager struct {
	onlineUsers map[string]*User
	users       map[string]*User
}

func NewManager() *Manager {
	return &Manager{
		onlineUsers: make(map[string]*User),
		users:       make(map[string]*User),
	}
}

func (mgr *Manager) GetByUsername(username string) *User {
	slog.Debug("Getting user by username",
		slog.String("username", username))

	return mgr.users[strings.ToLower(username)]
}

func (mgr *Manager) SetOnline(u *User) {
	slog.Debug("Setting user online",
		slog.String("id", u.ID),
		slog.String("username", u.Username))
	t := time.Now()
	u.LastLoginAt = &t
	u.Save()

	mgr.onlineUsers[u.ID] = u
}

func (mgr *Manager) SetOffline(u *User) {
	slog.Debug("Setting user offline",
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	delete(mgr.onlineUsers, u.ID)
}

func (mgr *Manager) LoadDataFiles() {
	dataFilePath := viper.GetString("data.users_path")

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
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				slog.Error("failed reading file",
					slog.String("file", file.Name()),
					slog.Any("error", err))
			}

			var u User
			if err := json.Unmarshal(fileContent, &u); err != nil {
				slog.Error("failed to unmarshal user data",
					slog.Any("error", err),
					slog.String("file", file.Name()))
			}

			mgr.users[strings.ToLower(u.Username)] = &u

			slog.Debug("Loaded user",
				slog.String("id", u.ID),
				slog.String("username", u.Username))
		}
	}

	slog.Info("Loaded users",
		slog.Int("count", len(mgr.users)))
}
