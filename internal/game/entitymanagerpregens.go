package game

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Pregen functions
func (mgr *EntityManager) AddPregen(p *Pregen) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.pregens[p.ID] = p
}

func (mgr *EntityManager) GetPregen(id string) *Pregen {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.pregens[id]
}

func (mgr *EntityManager) RemovePregen(p *Pregen) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.pregens, p.ID)
}

func (mgr *EntityManager) GetPregens() map[string]*Pregen {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.pregens
}

func (mgr *EntityManager) GetPregenMenuOptions() map[string]string {
	mgr.RLock()
	defer mgr.RUnlock()

	options := make(map[string]string)
	for _, p := range mgr.pregens {
		options[p.ID] = p.GetSelectionInfo()
	}

	return options
}

func (mgr *EntityManager) loadPregens() {
	slog.Info("Loading pregens")

	st := time.Now()

	files, err := os.ReadDir(PreGensFilepath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", PreGensFilepath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if !IsYAMLFile(file.Name()) {
			continue
		}

		var pregen Pregen
		if err := LoadYAML(filepath.Join(PreGensFilepath, file.Name()), &pregen); err != nil {
			slog.Error("failed to unmarshal pregen data",
				slog.Any("error", err),
				slog.String("pregens_path", filepath.Join(PreGensFilepath, file.Name())))
			continue
		}
		mgr.pregens[pregen.ID] = &pregen

		slog.Debug("Loaded pregen",
			slog.String("pregen_id", pregen.ID))
	}

	took := time.Since(st)
	slog.Info("Loaded pregens",
		slog.Duration("took", took),
		slog.Int("count", len(mgr.pregens)))
}
