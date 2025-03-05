package game

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	EntityMgr = NewEntityManager()
)

type EntityManager struct {
	sync.RWMutex

	areas           map[string]*Area
	itemsBlueprints map[string]*ItemBlueprint
	itemInstances   map[string]*ItemInstance
	metatypes       map[string]*Metatype
	mobBlueprints   map[string]*MobBlueprint
	mobInstances    map[string]*MobInstance
	pregens         map[string]*Pregen
	qualtities      map[string]*QualityBlueprint
	rooms           map[string]*Room
	skills          map[string]*SkillBlueprint
	skillGroups     map[string]*SkillGroup
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		areas:           make(map[string]*Area),
		itemsBlueprints: make(map[string]*ItemBlueprint),
		itemInstances:   make(map[string]*ItemInstance),
		metatypes:       make(map[string]*Metatype),
		mobBlueprints:   make(map[string]*MobBlueprint),
		mobInstances:    make(map[string]*MobInstance),
		pregens:         make(map[string]*Pregen),
		qualtities:      make(map[string]*QualityBlueprint),
		rooms:           make(map[string]*Room),
		skills:          make(map[string]*SkillBlueprint),
		skillGroups:     make(map[string]*SkillGroup),
	}
}

// Generic load function
func (mgr *EntityManager) LoadDataFiles() {
	st := time.Now()
	slog.Info("Loading data files")

	mgr.loadMetatypes()
	mgr.loadPregens()
	mgr.loadQualities()
	mgr.loadSkills(SkillTypeActive, SkillActiveFilepath)
	mgr.loadSkills(SkillTypeKnowledge, SkillKnowledgeFilepath)
	mgr.loadSkills(SkillTypeLanguage, SkillLanguagesFilepath)
	mgr.loadSkillGroups()
	mgr.loadAreasFromFS(os.DirFS(viper.GetString("data.areas_path")))

	took := time.Since(st)
	slog.Info("Loaded data files",
		slog.Duration("took", took),
		slog.Int("areas", len(mgr.areas)),
		slog.Int("items", len(mgr.itemsBlueprints)),
		slog.Int("mobs", len(mgr.mobBlueprints)),
		slog.Int("rooms", len(mgr.rooms)),
		slog.Int("pregens", len(mgr.pregens)),
		slog.Int("qualities", len(mgr.qualtities)),
		slog.Int("skills", len(mgr.skills)),
		slog.Int("skill_groups", len(mgr.skillGroups)),
		slog.Int("metatypes", len(mgr.metatypes)),
	)
}
