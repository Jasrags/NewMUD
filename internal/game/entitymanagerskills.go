package game

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// Skill functions
func (mgr *EntityManager) GetSkillBlueprint(id string) *SkillBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.skills[id]
}

func (mgr *EntityManager) AddSkillBlueprint(s *SkillBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	mgr.skills[s.ID] = s
}
func (mgr *EntityManager) RemoveSkillBlueprint(s *SkillBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.skills, s.ID)
}

func (mgr *EntityManager) CreateSkillInstanceFromBlueprintID(id string, rating int, specialization string) *Skill {
	if bp := mgr.GetSkillBlueprint(id); bp != nil {
		return mgr.CreateSkillInstanceFromBlueprint(bp, rating, specialization)
	}

	return nil
}

func (mgr *EntityManager) CreateSkillInstanceFromBlueprint(bp *SkillBlueprint, rating int, specialization string) *Skill {

	return NewSkill(bp, rating, specialization)
}

func (mgr *EntityManager) loadSkills(t string, path string) {
	st := time.Now()
	slog.Info("Loading skills",
		slog.Any("skill_type", t))

	files, err := os.ReadDir(path)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", path),
			slog.Any("error", err))
	}

	for _, file := range files {
		if !IsYAMLFile(file.Name()) {
			continue
		}

		var skill SkillBlueprint
		if err := LoadYAML(filepath.Join(path, file.Name()), &skill); err != nil {
			slog.Error("failed to unmarshal skill data",
				slog.Any("error", err),
				slog.String("skill_path", filepath.Join(path, file.Name())))
			continue
		}
		skill.Type = t
		mgr.skills[skill.ID] = &skill

		slog.Debug("Loaded skill",
			slog.Any("skill_type", skill.Type),
			slog.String("skill_id", skill.ID))
	}

	stook := time.Since(st)
	slog.Info("Loaded skills",
		slog.Duration("took", stook),
		slog.Any("skill_type", t),
		slog.Int("count", len(mgr.skills)))
}

func (mgr *EntityManager) loadSkillGroups() {
	st := time.Now()
	slog.Info("Loading skill groups")

	files, err := os.ReadDir(SkillGroupsFilepath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", SkillGroupsFilepath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if !IsYAMLFile(file.Name()) {
			continue
		}

		var group SkillGroup
		if err := LoadYAML(filepath.Join(SkillGroupsFilepath, file.Name()), &group); err != nil {
			slog.Error("failed to unmarshal skill group data",
				slog.Any("error", err),
				slog.String("skill_group_path", filepath.Join(SkillGroupsFilepath, file.Name())))
			continue
		}
		mgr.skillGroups[group.ID] = &group

		slog.Debug("Loaded skill group",
			slog.String("skill_group_id", group.ID))
	}

	took := time.Since(st)
	slog.Info("Loaded skill groups",
		slog.Duration("took", took),
		slog.Int("count", len(mgr.skillGroups)))
}
