package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	EntityMgr = NewEntityManager()
)

type EntityManager struct {
	sync.RWMutex

	areas map[string]*Area
	items map[string]*ItemBlueprint
	mobs  map[string]*Mob
	rooms map[string]*Room
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		areas: make(map[string]*Area),
		items: make(map[string]*ItemBlueprint),
		mobs:  make(map[string]*Mob),
		rooms: make(map[string]*Room),
	}
}

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

	return mgr.mobs[strings.ToLower(id)]
}

func (mgr *EntityManager) RemoveMob(m *Mob) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.mobs, m.ID)
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

func (mgr *EntityManager) LoadDataFiles() {
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
		slog.Int("areas_count", len(mgr.areas)),
		slog.Int("rooms_count", len(mgr.rooms)),
		slog.Int("items_count", len(mgr.items)),
		slog.Int("mobs_count", len(mgr.mobs)))
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
