package game

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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

// Area functions
func (mgr *EntityManager) AddArea(a *Area) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding area",
		slog.String("area_id", a.ID))

	mgr.areas[strings.ToLower(a.ID)] = a
}

func (mgr *EntityManager) GetArea(areaID string) *Area {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Getting area",
		slog.String("area_id", areaID))

	return mgr.areas[strings.ToLower(areaID)]
}

func (mgr *EntityManager) RemoveArea(a *Area) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Removing area",
		slog.String("area_id", a.ID))

	delete(mgr.areas, strings.ToLower(a.ID))
}

// Mob functions
func (mgr *EntityManager) GetAllMobInstances() map[string]*MobInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.mobInstances
}

func (mgr *EntityManager) AddMobInstance(m *MobInstance) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mobInstances[m.InstanceID]; ok {
		slog.Warn("Mob instance already exists",
			slog.String("mob_instance_id", m.InstanceID))
		return
	}

	mgr.mobInstances[m.InstanceID] = m
}

func (mgr *EntityManager) GetMobInstance(id string) *MobInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.mobInstances[id]
}

func (mgr *EntityManager) RemoveMobInstance(m *MobInstance) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mobInstances[m.InstanceID]; !ok {
		slog.Warn("Mob instance not found",
			slog.String("mob_instance_id", m.InstanceID))
		return
	}

	delete(mgr.mobInstances, m.InstanceID)
}

func (mgr *EntityManager) GetAllMobBlueprints() map[string]*MobBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.mobBlueprints
}

func (mgr *EntityManager) AddMobBlueprint(m *MobBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mobBlueprints[m.ID]; ok {
		slog.Warn("Mob blueprint already exists",
			slog.String("mob_blueprint_id", m.ID))
		return
	}

	mgr.mobBlueprints[m.ID] = m
}

func (mgr *EntityManager) RemoveMobBlueprint(m *MobBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mobBlueprints[m.ID]; !ok {
		slog.Warn("Mob blueprint not found",
			slog.String("mob_blueprint_id", m.ID))
		return
	}

	delete(mgr.mobBlueprints, m.ID)
}

func (mgr *EntityManager) GetMobBlueprintByID(id string) *MobBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.mobBlueprints[id]
}

func (mgr *EntityManager) GetMobBlueprintByInstance(mob *MobInstance) *MobBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	bp, ok := mgr.mobBlueprints[mob.BlueprintID]
	if !ok {
		slog.Error("Mob blueprint not found",
			slog.String("mob_blueprint_id", mob.BlueprintID))
		return nil
	}

	return bp
}

func (mgr *EntityManager) CreateMobInstanceFromBlueprintID(id string) *MobInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	bp, ok := mgr.mobBlueprints[id]
	if !ok {
		slog.Error("Mob blueprint not found",
			slog.String("mob_blueprint_id", id))
		return nil
	}

	return mgr.CreateMobInstanceFromBlueprint(bp)
}

func (mgr *EntityManager) CreateMobInstanceFromBlueprint(bp *MobBlueprint) *MobInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	mob := &MobInstance{
		Blueprint:   bp,
		InstanceID:  uuid.New().String(),
		BlueprintID: bp.ID,
		// Dynamic state fields
		CharacterDispositions: make(map[string]string),
		Edge:                  bp.Edge,
		PositionState:         PositionStanding,
		Inventory:             NewInventory(),
		Equipment:             make(map[string]*ItemInstance),
	}

	metatype := mgr.GetMetatype(bp.MetatypeID)
	if metatype == nil {
		slog.Error("Mob metatype not found",
			slog.String("mob_blueprint_id", bp.ID),
			slog.String("mob_metatype_id", bp.MetatypeID))
		return nil
	}

	// Spawn items into the mob's inventory or equipment
	for _, spawn := range bp.Spawns {
		// Check if the spawn is for an item
		if spawn.ItemID != "" {

			// Check if the item is a quality item
			quantity := spawn.Quantity
			if spawn.Quantity == 0 {
				quantity = 1
			}
			// Check if the spawn has a chance
			chance := spawn.Chance
			if spawn.Chance == 0 {
				chance = 100
			}

			for range quantity {
				if !RollChance(chance) {
					continue
				}

				item := mgr.CreateItemInstanceFromBlueprintID(spawn.ItemID)
				if item == nil {
					slog.Error("Item instance not found",
						slog.String("mob_blueprint_id", mob.BlueprintID),
						slog.String("item_blueprint_id", spawn.ItemID))
					continue
				}

				// Equip the item in the specified slot
				if spawn.EquipSlot != "" {
					if _, ok := mob.Equipment[spawn.EquipSlot]; ok {
						slog.Warn("Equip slot already occupied",
							slog.String("mob_blueprint_id", mob.BlueprintID),
							slog.String("equip_slot", spawn.EquipSlot))
						continue
					}
					mob.Equipment[spawn.EquipSlot] = item
				} else {
					// Add the item to the inventory
					mob.Inventory.AddItem(item)
				}
			}
		}
	}

	mob.Blueprint.Metatype = metatype

	return mob
}

// Item functions
func (mgr *EntityManager) GetAllItemInstances() map[string]*ItemInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.itemInstances
}

func (mgr *EntityManager) AddItemInstance(i *ItemInstance) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.itemInstances[i.InstanceID]; ok {
		slog.Warn("Item instance already exists",
			slog.String("item_instance_id", i.InstanceID))
		return
	}

	mgr.itemInstances[i.InstanceID] = i
}

func (mgr *EntityManager) GetItemInstance(id string) *ItemInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.itemInstances[id]
}

func (mgr *EntityManager) RemoveItemInstance(i *ItemInstance) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.itemInstances[i.InstanceID]; !ok {
		slog.Warn("Item instance not found",
			slog.String("item_instance_id", i.InstanceID))
		return
	}

	delete(mgr.itemInstances, i.InstanceID)
}

func (mgr *EntityManager) GetAllItemBlueprints() map[string]*ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.itemsBlueprints
}

func (mgr *EntityManager) AddItemBlueprint(i *ItemBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding item blueprint",
		slog.String("item_id", i.ID))

	if _, ok := mgr.itemsBlueprints[i.ID]; ok {
		slog.Warn("Item blueprint already exists",
			slog.String("item_id", i.ID))
		return
	}

	mgr.itemsBlueprints[i.ID] = i
}

func (mgr *EntityManager) GetItemBlueprintByID(id string) *ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.itemsBlueprints[id]
}

func (mgr *EntityManager) GetItemBlueprintByInstance(item *ItemInstance) *ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	bp, ok := mgr.itemsBlueprints[item.BlueprintID]
	if !ok {
		slog.Error("Item blueprint not found",
			slog.String("item_blueprint_id", item.BlueprintID))
		return nil
	}

	return bp
}

func (mgr *EntityManager) CreateItemInstanceFromBlueprintID(id string) *ItemInstance {
	mgr.RLock()
	defer mgr.RUnlock()

	bp := mgr.GetItemBlueprintByID(id)
	if bp == nil {
		slog.Error("Item blueprint not found",
			slog.String("item_blueprint_id", id))
		return nil
	}

	return mgr.CreateItemInstanceFromBlueprint(bp)
}

func (mgr *EntityManager) CreateItemInstanceFromBlueprint(bp *ItemBlueprint) *ItemInstance {
	return &ItemInstance{
		Blueprint:   bp,
		BlueprintID: bp.ID,
		InstanceID:  uuid.New().String(),
		Attachments: bp.Attachments,
	}
}

func (mgr *EntityManager) GetItemBlueprint(id string) *ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	bp, ok := mgr.itemsBlueprints[id]
	if !ok {
		slog.Error("Item blueprint not found",
			slog.String("item_blueprint_id", id))
		return nil
	}

	return bp
}

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

// Room Functions
func (mgr *EntityManager) GetAllRooms() map[string]*Room {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.rooms
}

func (mgr *EntityManager) AddRoom(r *Room) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.rooms[r.ID]; ok {
		slog.Warn("Room already exists",
			slog.String("room_id", r.ID))
		return
	}

	mgr.rooms[r.ID] = r
}

func (mgr *EntityManager) GetRoom(id string) *Room {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.rooms[strings.ToLower(id)]
}

func (mgr *EntityManager) RemoveRoom(r *Room) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.rooms, r.ID)
}

// TOOD: Break out room load code into a separate function

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
	mgr.LoadAreasFromFS()

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

func (mgr *EntityManager) LoadAreasFromFS() {
	start := time.Now()
	basePath := viper.GetString("data.areas_path")
	areasFS := os.DirFS(basePath)

	slog.Info("Loading areas",
		slog.String("datafile_path", basePath))

	areaDirs, err := fs.ReadDir(areasFS, ".")
	if err != nil {
		slog.Error("failed reading areas directory", "error", err)
		return
	}

	for _, d := range areaDirs {
		if !d.IsDir() {
			continue
		}

		areaFS, err := fs.Sub(areasFS, d.Name())
		if err != nil {
			slog.Error("failed to create sub FS for area", "area", d.Name(), "error", err)
			continue
		}

		// Load the area manifest
		manifestBytes, err := fs.ReadFile(areaFS, "manifest.yml")
		if err != nil {
			slog.Error("failed reading manifest", "area", d.Name(), "error", err)
			continue
		}
		var area Area
		if err := yaml.Unmarshal(manifestBytes, &area); err != nil {
			slog.Error("failed to unmarshal manifest data", "area", d.Name(), "error", err)
			continue
		}
		slog.Info("Loaded area", "area", d.Name())
		mgr.AddArea(&area)

		// Load rooms
		loadFilesFromDir(areaFS, "rooms", func(data []byte) {
			var room Room
			if err := yaml.Unmarshal(data, &room); err != nil {
				slog.Error("failed to unmarshal room data", "area", d.Name(), "error", err)
				return
			}
			room.MobInstances = make(map[string]*MobInstance)
			room.Characters = make(map[string]*Character)
			mgr.AddRoom(&room)
		})

		// Load items (ItemBlueprints)
		loadFilesFromDir(areaFS, "items", func(data []byte) {
			var item ItemBlueprint
			if err := yaml.Unmarshal(data, &item); err != nil {
				slog.Error("failed to unmarshal item data", "area", d.Name(), "error", err)
				return
			}
			mgr.AddItemBlueprint(&item)
		})

		// Load mobs (MobBlueprints)
		loadFilesFromDir(areaFS, "mobs", func(data []byte) {
			var mob MobBlueprint
			if err := yaml.Unmarshal(data, &mob); err != nil {
				slog.Error("failed to unmarshal mob data", "area", d.Name(), "error", err)
				return
			}

			// TODO: Add items into the mobs inventory
			// TODO: Add items to the mobs equipment

			mgr.AddMobBlueprint(&mob)
			// mgr.AddMob(&mob)
		})
	}

	mgr.BuildRooms()
	slog.Info("Loaded areas", "duration", time.Since(start))
}

func (mgr *EntityManager) BuildRooms() {
	slog.Info("Building rooms")

	for _, room := range mgr.rooms {
		// Build exits
		for dir, exit := range room.Exits {
			exit.Room = mgr.GetRoom(exit.RoomID)

			if exit.Room == nil {
				slog.Warn("Exit room not found",
					slog.String("room_id", room.ID),
					slog.String("exit_dir", dir),
					slog.String("exit_room_id", exit.RoomID))
				// TODO: Do we need to remove the exit from the room?
				continue
			}

			if exit.Door != nil {
				exit.Room.Exits[ReverseDirection(dir)].Door = exit.Door
			}
		}

		// Loop over the mobs we need to spawn into the room
		for _, spawn := range room.Spawns {
			quantity := spawn.Quantity
			if spawn.Quantity == 0 {
				quantity = 1
			}
			chance := spawn.Chance
			if spawn.Chance == 0 {
				chance = 100
			}

			if spawn.ItemID != "" {
				// Spawn an item into the room
				bp := mgr.GetItemBlueprintByID(spawn.ItemID)
				if bp == nil {
					slog.Warn("Item blueprint not found",
						slog.String("room_id", room.ID),
						slog.String("item_id", spawn.ItemID))
					continue
				}

				for range quantity {
					if !RollChance(chance) {
						continue
					}

					i := mgr.CreateItemInstanceFromBlueprint(bp)
					if i == nil {
						slog.Warn("Item instance not found",
							slog.String("room_id", room.ID),
							slog.String("item_id", spawn.ItemID))
						continue
					}
					room.Inventory.AddItem(i)
				}
			} else if spawn.MobID != "" {
				bp := mgr.GetMobBlueprintByID(spawn.MobID)
				if bp == nil {
					slog.Warn("Mob blueprint not found",
						slog.String("room_id", room.ID),
						slog.String("mob_blueprint_id", spawn.MobID))
					continue
				}

				for range quantity {
					if !RollChance(chance) {
						continue
					}

					mob := mgr.CreateMobInstanceFromBlueprint(bp)
					if mob == nil {
						slog.Warn("Mob not found",
							slog.String("room_id", room.ID),
							slog.String("mob_id", spawn.MobID))
						continue
					}

					room.AddMobInstance(mob)
				}
			}
		}
	}
}
