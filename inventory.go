package main

import "strings"

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

func SearchInventory(inv *Inventory, query string) []*Item {
	var results []*Item

	for _, item := range inv.Items {
		blueprint := EntityMgr.GetItemBlueprintByInstance(item) // Assume this fetches the blueprint for the item
		if blueprint == nil {
			continue
		}

		// Match against Name
		if strings.Contains(strings.ToLower(blueprint.Name), strings.ToLower(query)) {
			results = append(results, item)
			continue
		}

		// Match against Tags
		for _, tag := range blueprint.Tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(query)) {
				results = append(results, item)
				break
			}
		}
	}

	return results
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
	blueprint := em.GetItemBlueprintByInstance(instance)
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
