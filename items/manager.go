package items

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	manager = NewManager()
)

type Manager struct {
	items map[string]*Item `yaml:"items"`
}

func NewManager() *Manager {
	return &Manager{
		items: make(map[string]*Item),
	}
}

func GetByID(id string) *Item {
	return manager.items[id]
}

func LoadDataFiles() {
	dataFilePath := viper.GetString("data.areas_path")
	// manifestFileName := viper.GetString("data.manifest_file")
	itemsFileName := viper.GetString("data.items_file")

	slog.Info("Loading item data files",
		slog.String("datafile_path", dataFilePath),
		// slog.String("manifest_file", manifestFileName),
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
			// manifestPath := filepath.Join(areaPath, manifestFileName)
			itemsPath := filepath.Join(areaPath, itemsFileName)

			// Load area manifest
			// manifestData, errReadFile := os.ReadFile(manifestPath)
			// if errReadFile != nil {
			// 	slog.Error("failed reading manifest file",
			// 		slog.Any("error", errReadFile),
			// 		slog.String("area_path", areaPath))
			// 	continue
			// }

			// var area Area
			// if err := yaml.Unmarshal(manifestData, &area); err != nil {
			// 	slog.Error("failed to unmarshal manifest data",
			// 		slog.Any("error", err),
			// 		slog.String("area_path", areaPath))
			// 	continue
			// }

			// slog.Info("Loaded area manifest",
			// 	slog.String("area_name", file.Name()))

			// Add area to roomManager
			// manager.areas[area.ID] = &area

			// Load rooms
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

			// Add rooms to roomManager
			for _, item := range items {
				manager.items[item.ID] = &item
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
