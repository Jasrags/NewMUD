package game

import (
	"log/slog"

	"github.com/google/uuid"
)

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

	var mob MobInstance
	mob.InstanceID = uuid.New().String()
	mob.Blueprint = bp
	mob.BlueprintID = bp.ID
	mob.GameEntityDynamic = NewGameEntityDynamic()

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
					if _, ok := mob.Equipment.Slots[spawn.EquipSlot]; ok {
						slog.Warn("Equip slot already occupied",
							slog.String("mob_blueprint_id", mob.BlueprintID),
							slog.String("equip_slot", spawn.EquipSlot))
						continue
					}
					mob.Equipment.Slots[spawn.EquipSlot] = item
				} else {
					// Add the item to the inventory
					mob.Inventory.Add(item)
				}
			}
		}
	}

	return &mob
}
