package main

type Inventory struct {
	Items []*Item `yaml:"items"`
}

// Add an item to the inventory
func (inv *Inventory) AddItem(item *Item) {
	inv.Items = append(inv.Items, item)
}

// Remove an item from the inventory
func (inv *Inventory) RemoveItem(item *Item) bool {
	for i, existingItem := range inv.Items {
		if existingItem == item {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return true
		}
	}
	return false
}

// Find an item by its instance ID
func (inv *Inventory) FindItemByID(instanceID string) *Item {
	for _, item := range inv.Items {
		if item.InstanceID == instanceID {
			return item
		}
	}
	return nil
}

func TransferItem(item *Item, from, to *Inventory) bool {
	if from.RemoveItem(item) {
		to.AddItem(item)
		return true
	}
	return false
}

// Combine base stats and modifiers for a given item instance
func GetCombinedStats(instance *Item, em *EntityManager) map[string]int {
	blueprint := em.GetBlueprint(instance)
	if blueprint == nil {
		return nil
	}

	combinedStats := make(map[string]int)
	for key, value := range blueprint.BaseStats {
		combinedStats[key] = value
	}
	for key, value := range instance.Modifiers {
		combinedStats[key] += value
	}
	return combinedStats
}
