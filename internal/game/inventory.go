package game

import (
	"strings"
)

type Inventory struct {
	Items []*Item `yaml:"items"`
}

// NewInventory creates a new inventory
func NewInventory() *Inventory {
	return &Inventory{
		Items: []*Item{},
	}
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

// Find an item by its name
func (inv *Inventory) FindItemByName(name string) *Item {
	for _, item := range inv.Items {
		blueprint := EntityMgr.GetItemBlueprintByInstance(item) // Fetch the blueprint for the item
		if blueprint != nil && strings.EqualFold(blueprint.Name, name) {
			return item
		}
	}
	return nil
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

func (inv *Inventory) Clear() {
	inv.Items = nil
}

func (inv *Inventory) Search(query string) []*Item {
	results := []*Item{}

	if len(inv.Items) == 0 {
		return results
	}

	lowerQuery := strings.ToLower(query)

	for _, item := range inv.Items {
		blueprint := EntityMgr.GetItemBlueprintByInstance(item)
		if blueprint == nil {
			continue
		}

		if strings.Contains(strings.ToLower(blueprint.Name), lowerQuery) {
			results = append(results, item)
			continue
		}

		if matchesTags(blueprint.Tags, lowerQuery) {
			results = append(results, item)
		}
	}

	return results
}

// Helper function to check if any tag matches the query
func matchesTags(tags []string, query string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

func TransferItem(item *Item, from, to Inventory) bool {
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
