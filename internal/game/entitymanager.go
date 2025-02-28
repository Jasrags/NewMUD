package game

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	EntityMgr = NewEntityManager()
)

type EntityManager struct {
	sync.RWMutex

	areas       map[string]*Area
	items       map[string]*ItemBlueprint
	metatypes   map[string]*Metatype
	mobs        map[string]*Mob
	pregens     map[string]*Pregen
	qualtities  map[string]*QualityBlueprint
	rooms       map[string]*Room
	skills      map[string]*SkillBlueprint
	skillGroups map[string]*SkillGroup
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		areas:       make(map[string]*Area),
		items:       make(map[string]*ItemBlueprint),
		metatypes:   make(map[string]*Metatype),
		mobs:        make(map[string]*Mob),
		pregens:     make(map[string]*Pregen),
		qualtities:  make(map[string]*QualityBlueprint),
		rooms:       make(map[string]*Room),
		skills:      make(map[string]*SkillBlueprint),
		skillGroups: make(map[string]*SkillGroup),
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

// Item functions
func (mgr *EntityManager) GetAllItemBlueprints() map[string]*ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.items
}

func (mgr *EntityManager) AddItemBlueprint(i *ItemBlueprint) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding item blueprint",
		slog.String("item_id", i.ID))

	if _, ok := mgr.items[i.ID]; ok {
		slog.Warn("Item blueprint already exists",
			slog.String("item_id", i.ID))
		return
	}

	mgr.items[i.ID] = i
}

func (mgr *EntityManager) GetItemBlueprintByID(id string) *ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.items[id]
}

func (mgr *EntityManager) GetItemBlueprintByInstance(item *Item) *ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	bp, ok := mgr.items[item.BlueprintID]
	if !ok {
		slog.Error("Item blueprint not found",
			slog.String("item_blueprint_id", item.BlueprintID))
		return nil
	}

	return bp
}

func (mgr *EntityManager) CreateItemInstanceFromBlueprintID(id string) *Item {
	slog.Debug("Creating item instance from blueprint",
		slog.String("item_blueprint_id", id))

	bp, ok := mgr.items[id]
	if !ok {
		slog.Error("Item blueprint not found",
			slog.String("item_blueprint_id", id))
		return nil
	}

	return mgr.CreateItemInstanceFromBlueprint(bp)
}

func (mgr *EntityManager) CreateItemInstanceFromBlueprint(bp *ItemBlueprint) *Item {
	slog.Debug("Creating item instance from blueprint",
		slog.String("item_blueprint_id", bp.ID))

	return &Item{
		InstanceID:  uuid.New().String(),
		BlueprintID: bp.ID,
		Modifiers:   make(map[string]int),
		Attachments: []string{},
	}
}

func (mgr *EntityManager) CreateItemFromBlueprint(bp ItemBlueprint) *Item {
	slog.Debug("Creating item instance from blueprint",
		slog.String("item_blueprint_id", bp.ID))

	return &Item{
		InstanceID:  uuid.New().String(),
		BlueprintID: bp.ID,
		Modifiers:   make(map[string]int),
		Attachments: []string{},
	}
}

func (mgr *EntityManager) GetItemBlueprint(id string) *ItemBlueprint {
	mgr.RLock()
	defer mgr.RUnlock()

	bp, ok := mgr.items[id]
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

// Mob functions
func (mgr *EntityManager) GetAllMobs() map[string]*Mob {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.mobs
}

func (mgr *EntityManager) AddMob(m *Mob) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.mobs[m.ID]; ok {
		slog.Warn("Mob already exists",
			slog.String("mob_id", m.ID))
		return
	}

	mgr.mobs[m.ID] = m
}

func (mgr *EntityManager) GetMob(id string) *Mob {
	mgr.RLock()
	defer mgr.RUnlock()

	for _, mob := range mgr.mobs {
		slog.Debug("Checking mob",
			slog.String("mob_id", mob.ID))
	}

	return mgr.mobs[strings.ToLower(id)]
}

func (mgr *EntityManager) RemoveMob(m *Mob) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.mobs, m.ID)
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

	slog.Info("Loading entity data files")
	dataFilePath := viper.GetString("data.areas_path")
	manifestFileName := viper.GetString("data.manifest_file")
	roomsFileName := viper.GetString("data.rooms_file")
	itemsFileName := viper.GetString("data.items_file")
	mobsFileName := viper.GetString("data.mobs_file")

	slog.Info("Loading entity data files",
		slog.String("datafile_path", dataFilePath),
		slog.String("manifest_file", manifestFileName),
		slog.String("rooms_file", roomsFileName),
		slog.String("items_file", itemsFileName),
		slog.String("mobs_file", mobsFileName))

	mgr.loadMetatypes()
	mgr.loadPregens()
	mgr.loadQualities()
	mgr.loadSkills(SkillTypeActive, SkillActiveFilepath)
	mgr.loadSkills(SkillTypeKnowledge, SkillKnowledgeFilepath)
	mgr.loadSkills(SkillTypeLanguage, SkillLanguagesFilepath)
	mgr.loadSkillGroups()

	files, err := os.ReadDir(dataFilePath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", dataFilePath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if file.IsDir() {
			areaPath := filepath.Join(dataFilePath, file.Name())
			manifestPath := filepath.Join(areaPath, manifestFileName)
			roomsPath := filepath.Join(areaPath, roomsFileName)
			itemsPath := filepath.Join(areaPath, itemsFileName)
			mobsPath := filepath.Join(areaPath, mobsFileName)

			// Load area
			var area Area
			if err := LoadYAML(manifestPath, &area); err != nil {
				slog.Error("failed to unmarshal manifest data",
					slog.Any("error", err),
					slog.String("area_path", areaPath))
				continue
			}

			slog.Info("Loaded area",
				slog.String("area_name", file.Name()))

			mgr.AddArea(&area)

			// Load rooms
			if FileExists(roomsPath) {
				slog.Info("Loading rooms",
					slog.String("path", roomsPath),
					slog.String("area_id", area.ID))

				var rooms []Room
				if err := LoadYAML(roomsPath, &rooms); err != nil {
					slog.Error("failed to unmarshal rooms data",
						slog.Any("error", err),
						slog.String("rooms_path", roomsPath))
					continue
				}

				for i := range rooms {
					slog.Debug("Adding room",
						slog.String("room_id", rooms[i].ID))
					mgr.AddRoom(&rooms[i])
				}
			}

			// Load items
			if FileExists(itemsPath) {
				slog.Info("Loading items",
					slog.String("path", itemsPath),
					slog.String("area_id", area.ID))

				var items []ItemBlueprint
				if err := LoadYAML(itemsPath, &items); err != nil {
					slog.Error("failed to unmarshal item data",
						slog.Any("error", err),
						slog.String("items_path", itemsPath))
					continue
				}

				for i := range items {
					mgr.AddItemBlueprint(&items[i])
				}
			}

			// Load mobs
			if FileExists(mobsPath) {
				slog.Info("Loading mobs",
					slog.String("path", mobsPath),
					slog.String("area_id", area.ID))

				var mobs []Mob
				if err := LoadYAML(mobsPath, &mobs); err != nil {
					slog.Error("failed to unmarshal mobs data",
						slog.Any("error", err),
						slog.String("mobs_path", mobsPath))
					continue
				}

				for i := range mobs {
					mgr.AddMob(&mobs[i])
				}
			}
		}
	}

	mgr.BuildRooms()

	slog.Info("Loaded entities",
		slog.Duration("took", time.Since(st)))
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

		// Spawn default items
		// TODO: Support for respawn_chance, max_load, replace_on_respawn, quantity
		for _, di := range room.DefaultItems {
			i := EntityMgr.CreateItemInstanceFromBlueprintID(di.ID)
			room.Inventory.AddItem(i)
		}

		for _, dm := range room.DefaultMobs {
			m := EntityMgr.GetMob(dm.ID)
			if m == nil {
				slog.Warn("Mob not found",
					slog.String("mob_id", dm.ID))
				continue
			}

			room.AddMob(m)
		}
	}
}
