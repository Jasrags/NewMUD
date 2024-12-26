package items

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	Mgr = NewManager()
)

type Manager struct {
	items map[string]*Item `yaml:"items"`
}

func NewManager() *Manager {
	return &Manager{
		items: make(map[string]*Item),
	}
}

func (mgr *Manager) GetByID(id string) *Item {
	return mgr.items[id]
}

func (mgr *Manager) LoadDataFiles() {
	dataFilePath := viper.GetString("data.areas_path")
	itemsFileName := viper.GetString("data.items_file")

	slog.Info("Loading item data files",
		slog.String("datafile_path", dataFilePath),
		slog.String("items_file", itemsFileName))

	files, err := os.ReadDir(dataFilePath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", dataFilePath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if file.IsDir() {
			areaPath := filepath.Join(dataFilePath, file.Name())
			itemsPath := filepath.Join(areaPath, itemsFileName)

			itemsData, err := os.ReadFile(itemsPath)
			if err != nil {
				slog.Error("failed reading items file",
					slog.Any("error", err),
					slog.String("items_path", itemsPath))
				continue
			}

			var items []Item
			if err := yaml.Unmarshal(itemsData, &items); err != nil {
				slog.Error("failed to unmarshal items data",
					slog.Any("error", err),
					slog.String("items_path", itemsPath))
				continue
			}

			for _, item := range items {
				mgr.items[item.ID] = &item
				slog.Debug("Loaded item",
					slog.String("id", item.ID),
					slog.String("name", item.Name))
			}

			slog.Info("Loaded area items",
				slog.Int("count", len(items)),
				slog.String("area_name", file.Name()))
		}
	}
}
