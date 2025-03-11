package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInventory(t *testing.T) {
	inv := NewInventory()
	assert.NotNil(t, inv, "Expected new inventory to be created")
	assert.Equal(t, 0, len(inv.Items), "Expected new inventory to have no items")
}
func TestAddItem(t *testing.T) {
	inv := NewInventory()
	item := &ItemInstance{InstanceID: "item1"}

	inv.Add(item)
	assert.Equal(t, 1, len(inv.Items), "Expected inventory to have 1 item after adding")
	assert.Equal(t, item, inv.Items[0], "Expected the added item to be in the inventory")
}

// func TestRemoveItem(t *testing.T) {
// 	item1 := &ItemInstance{InstanceID: "item1"}
// 	item2 := &ItemInstance{InstanceID: "item2"}
// 	inv := NewInventory()
// 	inv.Add(item1)
// 	inv.Add(item2)

// 	// Test removing an existing item
// 	removed := inv.Remove(item1)
// 	assert.True(t, removed, "Expected item1 to be removed")
// 	assert.Equal(t, 1, len(inv.Items), "Expected inventory to have 1 item after removal")
// 	assert.Equal(t, item2, inv.Items[0], "Expected item2 to be the remaining item")

// 	// Test removing a non-existing item
// 	removed = inv.Remove(item1)
// 	assert.False(t, removed, "Expected removal of non-existing item to return false")
// 	assert.Equal(t, 1, len(inv.Items), "Expected inventory to still have 1 item after failed removal")
// }

// func TestFindItemByName(t *testing.T) {
// 	itemBP1 := &ItemBlueprint{ID: "item1", Name: "Sword"}
// 	itemBP2 := &ItemBlueprint{ID: "item2", Name: "Shield"}
// 	EntityMgr := NewEntityManager()
// 	EntityMgr.AddItemBlueprint(itemBP1)
// 	EntityMgr.AddItemBlueprint(itemBP2)

// 	item1 := &Item{InstanceID: "item1", BlueprintID: itemBP1.ID}
// 	item2 := &Item{InstanceID: "item2", BlueprintID: itemBP2.ID}
// 	inv := NewInventory()
// 	inv.AddItem(item1)
// 	inv.AddItem(item2)

// 	// Test finding an existing item by name
// 	foundItem := inv.FindItemByName("Sword")
// 	assert.NotNil(t, foundItem, "Expected to find item with name 'Sword'")
// 	assert.Equal(t, item1, foundItem, "Expected to find item1 with name 'Sword'")

// 	// Test finding an item with different case
// 	foundItem = inv.FindItemByName("sword")
// 	assert.NotNil(t, foundItem, "Expected to find item with name 'sword' (case insensitive)")
// 	assert.Equal(t, item1, foundItem, "Expected to find item1 with name 'sword'")

// 	// Test finding a non-existing item by name
// 	foundItem = inv.FindItemByName("Bow")
// 	assert.Nil(t, foundItem, "Expected not to find item with name 'Bow'")
// }
