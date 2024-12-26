package mobs

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
	mobs map[string]*Mob
}

func NewManager() *Manager {
	return &Manager{
		mobs: make(map[string]*Mob),
	}
}

func (mgr *Manager) LoadDataFiles() {
	dataFilePath := viper.GetString("data.areas_path")
	mobsFileName := viper.GetString("data.mobs_file")

	slog.Info("Loading mob data files",
		slog.String("datafile_path", dataFilePath),
		slog.String("mobs_file", mobsFileName))

	files, err := os.ReadDir(dataFilePath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", dataFilePath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if file.IsDir() {
			areaPath := filepath.Join(dataFilePath, file.Name())
			mobsPath := filepath.Join(areaPath, mobsFileName)

			npcsData, err := os.ReadFile(mobsPath)
			if err != nil {
				slog.Error("failed reading mobs file",
					slog.Any("error", err),
					slog.String("mobs_path", mobsPath))
				continue
			}

			var mobs []Mob
			if err := yaml.Unmarshal(npcsData, &mobs); err != nil {
				slog.Error("failed to unmarshal mobs data",
					slog.Any("error", err),
					slog.String("mobs_path", mobsPath))
				continue
			}

			for _, mob := range mobs {
				mgr.mobs[mob.ID] = &mob
				slog.Debug("Loaded mob",
					slog.String("id", mob.ID),
					slog.String("name", mob.Name))
			}

			slog.Info("Loaded area mobs",
				slog.Int("count", len(mobs)),
				slog.String("area_name", file.Name()))
		}
	}
}
