package game

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Quality functions
func (mgr *EntityManager) GetQualityBlueprint(id string) *QualityBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.qualtities[id]
}

func (mgr *EntityManager) AddQualityBlueprint(q *QualityBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.qualtities[q.ID] = q
}

func (mgr *EntityManager) RemoveQualityBlueprint(q *QualityBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.qualtities, q.ID)
}

func (mgr *EntityManager) CreateQualityFromBlueprintID(id string, rating int) *Quality {
	if bp := mgr.GetQualityBlueprint(id); bp != nil {
		return mgr.CreateQualityFromBlueprint(bp, rating)
	}

	return nil
}

func (mgr *EntityManager) CreateQualityFromBlueprint(bp *QualityBlueprint, rating int) *Quality {
	return NewQuality(bp, rating)
}

func (mgr *EntityManager) loadQualities() {
	slog.Info("Loading qualities")

	st := time.Now()

	files, err := os.ReadDir(QualitiesFilepath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", QualitiesFilepath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if !IsYAMLFile(file.Name()) {
			continue
		}

		var quality QualityBlueprint
		if err := LoadYAML(filepath.Join(QualitiesFilepath, file.Name()), &quality); err != nil {
			slog.Error("failed to unmarshal quality data",
				slog.Any("error", err),
				slog.String("quality_path", filepath.Join(QualitiesFilepath, file.Name())))
			continue
		}
		mgr.qualtities[quality.ID] = &quality

		slog.Debug("Loaded quality",
			slog.String("quality_id", quality.ID))
	}

	took := time.Since(st)
	slog.Info("Loaded qualities",
		slog.Duration("took", took),
		slog.Int("count", len(mgr.qualtities)))
}
