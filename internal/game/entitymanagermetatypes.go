package game

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Metatype functions
func (mgr *EntityManager) GetMetatype(id string) *Metatype {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.metatypes[strings.ToLower(id)]
}

func (mgr *EntityManager) AddMetatype(m *Metatype) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.metatypes[m.ID] = m
}

func (mgr *EntityManager) RemoveMetatype(m *Metatype) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.metatypes, m.ID)
}

func (mgr *EntityManager) GetMetatypes() map[string]*Metatype {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.metatypes
}

func (mgr *EntityManager) GetMetatypeMenuOptions() map[string]string {
	mgr.RLock()
	defer mgr.RUnlock()

	options := make(map[string]string)
	for _, m := range mgr.metatypes {
		if m.Hidden {
			continue
		}
		options[m.Name] = m.GetSelectionInfo()
	}

	return options
}

func (mgr *EntityManager) loadMetatypes() {
	slog.Info("Loading metatypes")

	st := time.Now()
	files, err := os.ReadDir(MetatypesFilepath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", MetatypesFilepath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if !IsYAMLFile(file.Name()) {
			continue
		}

		var metatype Metatype
		if err := LoadYAML(filepath.Join(MetatypesFilepath, file.Name()), &metatype); err != nil {
			slog.Error("failed to unmarshal metatype data",
				slog.Any("error", err),
				slog.String("metatype_path", filepath.Join(MetatypesFilepath, file.Name())))
			continue
		}
		mgr.AddMetatype(&metatype)

		slog.Debug("Loaded metatype",
			slog.String("metatype_id", metatype.ID))
	}

	took := time.Since(st)
	slog.Debug("Loaded metatypes",
		slog.Duration("took", took),
		slog.Int("count", len(mgr.metatypes)))
}
