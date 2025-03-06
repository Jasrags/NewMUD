package game

import (
	"log/slog"

	"github.com/google/uuid"
)

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

	var itemInstance ItemInstance
	itemInstance.InstanceID = uuid.New().String()
	itemInstance.BlueprintID = bp.ID
	itemInstance.Blueprint = bp
	itemInstance.Attachments = bp.Attachments

	return &itemInstance
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
