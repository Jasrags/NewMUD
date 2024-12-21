package characters

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	manager = NewManager()
)

type Manager struct {
	characters map[string]*Character
}

func NewManager() *Manager {
	return &Manager{
		characters: make(map[string]*Character),
	}
}

func LoadDataFiles() {
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
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				slog.Error("failed reading file",
					slog.String("file", file.Name()),
					slog.Any("error", err))
			}

			var c Character
			if err := json.Unmarshal(fileContent, &c); err != nil {
				slog.Error("failed to unmarshal character data",
					slog.Any("error", err),
					slog.String("file", file.Name()))
			}

			manager.characters[strings.ToLower(c.Name)] = &c

			slog.Debug("Loaded character",
				slog.String("id", c.ID),
				slog.String("name", c.Name))
		}
	}

	slog.Info("Loaded characters",
		slog.Int("count", len(manager.characters)))
}
