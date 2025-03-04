package game

import (
	"slices"
	"strings"
)

const ()

type (
	Inventory struct {
		Items []*ItemInstance `yaml:"items"`
	}
)

// NewInventory creates a new inventory
func NewInventory() Inventory {
	return Inventory{
		Items: []*ItemInstance{},
	}
}

// Add an item to the inventory
func (inv *Inventory) AddItem(item *ItemInstance) {
	inv.Items = append(inv.Items, item)
}

// Remove an item from the inventory
func (inv *Inventory) RemoveItem(item *ItemInstance) bool {
	for i, existingItem := range inv.Items {
		if existingItem == item {
			inv.Items = slices.Delete(inv.Items, i, i+1)
			return true
		}
	}
	return false
}

// Find an item by its name
func (inv *Inventory) FindItemByName(name string) *ItemInstance {
	for _, item := range inv.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp != nil && strings.EqualFold(bp.Name, name) {
			return item
		}
	}
	return nil
}

// Find an item by its instance ID
func (inv *Inventory) FindItemByID(instanceID string) *ItemInstance {
	for _, item := range inv.Items {
		if item.InstanceID == instanceID {
			return item
		}
	}
	return nil
}

func (inv *Inventory) FindItemByTags(tags ...string) []*ItemInstance {
	results := []*ItemInstance{}

	for _, item := range inv.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp == nil {
			continue
		}

		if bp.HasTags(tags...) {
			results = append(results, item)
		}
	}

	return results
}

func (inv *Inventory) Clear() {
	inv.Items = nil
}

func (inv *Inventory) Search(query string) []*ItemInstance {
	results := []*ItemInstance{}

	if len(inv.Items) == 0 {
		return results
	}

	lowerQuery := strings.ToLower(query)

	for _, item := range inv.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp == nil {
			continue
		}

		if strings.Contains(strings.ToLower(bp.Name), lowerQuery) {
			results = append(results, item)
			continue
		}

		if matchesTags(bp.Tags, lowerQuery) {
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

func TransferItem(item *ItemInstance, from, to Inventory) bool {
	if from.RemoveItem(item) {
		to.AddItem(item)
		return true
	}
	return false
}

// Combine base stats and modifiers for a given item instance
func GetCombinedStats(i *ItemInstance, em *EntityManager) map[string]int {
	bp := em.GetItemBlueprintByInstance(i)
	if bp == nil {
		return nil
	}

	combinedStats := make(map[string]int)
	for key, value := range bp.BaseStats {
		combinedStats[key] = value
	}
	for key, value := range i.Blueprint.Modifiers {
		combinedStats[key] += value
	}
	return combinedStats
}
