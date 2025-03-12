package game

import (
	"io/fs"
	"log/slog"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

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

func (mgr *EntityManager) loadAreasFromFS(areasFS fs.FS) {
	start := time.Now()

	slog.Info("Loading areas")

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
			var modBlueprint MobBlueprint
			if err := yaml.Unmarshal(data, &modBlueprint); err != nil {
				slog.Error("failed to unmarshal mob data", "area", d.Name(), "error", err)
				return
			}

			// Add pointer to metatype
			metatype := mgr.GetMetatype(modBlueprint.MetatypeID)
			if metatype == nil {
				slog.Warn("Metatype not found",
					slog.String("mob_blueprint_id", modBlueprint.ID),
					slog.String("metatype_id", modBlueprint.MetatypeID))
				return
			}

			modBlueprint.Metatype = metatype

			mgr.AddMobBlueprint(&modBlueprint)
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
					room.Inventory.Add(i)
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
