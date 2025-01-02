package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	EntityMgr = NewEntityManager()
)

type EntityManager struct {
	sync.RWMutex

	areas map[string]*Area
	items map[string]*Item
	mobs  map[string]*Mob
	rooms map[string]*Room
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		areas: make(map[string]*Area),
		items: make(map[string]*Item),
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

func (mgr *EntityManager) AddItem(i *Item) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding item",
		slog.String("area_id", i.AreaID),
		slog.String("item_id", i.ID))

	mgr.items[CreateEntityRef(i.AreaID, i.ID)] = i
}

func (mgr *EntityManager) GetItem(referenceID string) *Item {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Getting item",
		slog.String("item_reference_id", referenceID))

	return mgr.items[strings.ToLower(referenceID)]
}

func (mgr *EntityManager) RemoveItem(i *Item) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Removing item",
		slog.String("area_id", i.AreaID),
		slog.String("item_id", i.ID))

	delete(mgr.items, CreateEntityRef(i.AreaID, i.ID))
}

func (mgr *EntityManager) AddMob(m *Mob) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding mob",
		slog.String("area_id", m.AreaID),
		slog.String("mob_id", m.ID))

	mgr.mobs[CreateEntityRef(m.AreaID, m.ID)] = m
}

func (mgr *EntityManager) GetMob(referenceID string) *Mob {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Getting mob",
		slog.String("mob_reference_id", referenceID))

	return mgr.mobs[strings.ToLower(referenceID)]
}

func (mgr *EntityManager) RemoveMob(m *Mob) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Removing mob",
		slog.String("area_id", m.AreaID),
		slog.String("mob_id", m.ID))

	delete(mgr.mobs, CreateEntityRef(m.AreaID, m.ID))
}

func (mgr *EntityManager) AddRoom(r *Room) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Adding room",
		slog.String("area_id", r.AreaID),
		slog.String("room_id", r.ID))

	mgr.rooms[CreateEntityRef(r.AreaID, r.ID)] = r
	r.Init()
}

func (mgr *EntityManager) GetRoom(referenceID string) *Room {
	mgr.RLock()
	defer mgr.RUnlock()

	slog.Debug("Getting room",
		slog.String("room_reference_id", referenceID))

	return mgr.rooms[strings.ToLower(referenceID)]
}

func (mgr *EntityManager) RemoveRoom(r *Room) {
	mgr.Lock()
	defer mgr.Unlock()

	slog.Debug("Removing room",
		slog.String("room_id", r.ID))

	delete(mgr.rooms, CreateEntityRef(r.AreaID, r.ID))
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
				var rooms []Room
				if err := LoadYAML(roomsPath, &rooms); err != nil {
					slog.Error("failed to unmarshal rooms data",
						slog.Any("error", err),
						slog.String("rooms_path", roomsPath))
					continue
				}

				for i := range rooms {
					room := &rooms[i]
					// room.Init()
					room.ReferenceID = CreateEntityRef(area.ID, room.ID)
					// room.Area = &area
					room.AreaID = area.ID
					mgr.AddRoom(room)
				}
			}

			// Load items
			if FileExists(itemsPath) {
				var items []Item
				if err := LoadYAML(itemsPath, &items); err != nil {
					slog.Error("failed to unmarshal items data",
						slog.Any("error", err),
						slog.String("items_path", itemsPath))
					continue
				}

				for i := range items {
					item := &items[i]
					// item.Init()
					item.ReferenceID = CreateEntityRef(area.ID, item.ID)
					// item.Area = &area
					item.AreaID = area.ID
					mgr.AddItem(item)
				}
			}

			// Load mobs
			if FileExists(mobsPath) {
				var mobs []Mob
				if err := LoadYAML(mobsPath, &mobs); err != nil {
					slog.Error("failed to unmarshal mobs data",
						slog.Any("error", err),
						slog.String("mobs_path", mobsPath))
					continue
				}

				for i := range mobs {
					mob := &mobs[i]
					// mob.Init()
					mob.ReferenceID = CreateEntityRef(area.ID, mob.ID)
					// mob.Area = &area
					mob.AreaID = area.ID
					mgr.AddMob(mob)
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
	// mgr.Lock()
	// defer mgr.Unlock()

	slog.Info("Building room exits")
	for _, room := range mgr.rooms {
		for dir, exit := range room.Exits {
			exit.Room = mgr.GetRoom(exit.RoomID)
			room.Exits[dir] = exit
		}
		// r := mgr.rooms[id]
		// for dir, _ := range r.Exits {
		// e := r.Exits[dir]
		// e.Room = mgr.GetRoom(e.RoomID)
		// exit.Room = mgr.GetRoom(exit.RoomID)
		// exit.Room = mgr.GetRoom(exit.RoomID)
		// room.Exits[exit.Direction] = mgr.GetRoom(exit.ReferenceID)
		// }
	}
}
