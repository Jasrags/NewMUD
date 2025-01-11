package main

import "github.com/google/uuid"

type (
	ItemType    string
	ItemSubtype string
	EquipSlot   string
)

// TODO: For keys we need subtypes for the different locks they can open.
// TODO: For picks they also need subtypes for the different locks they can pick.
// TODO: For locks they should have a rating of how difficult they are to pick
// TODO: For picks they should have a raiting of how good they are at picking locks
// TODO: Locks should somehow tie into alarm/traps or other events

const (
	// General Types
	ItemTypeJunk   ItemType = "junk"
	ItemTypeKey    ItemType = "key"
	ItemTypeFood   ItemType = "food"
	ItemTypeWeapon ItemType = "weapon"
	ItemTypeArmor  ItemType = "armor"

	// Weapon Subtypes
	ItemSubtypeMelee  ItemSubtype = "melee"
	ItemSubtypeRanged ItemSubtype = "ranged"

	// Armor Subtypes
	ItemSubtypeHead  ItemSubtype = "head"
	ItemSubtypeChest ItemSubtype = "chest"
	ItemSubtypeLegs  ItemSubtype = "legs"
)

type ItemBlueprint struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Tags        []string               `yaml:"tags"`
	BaseStats   map[string]int         `yaml:"base_stats"`
	EquipSlots  []EquipSlot            `yaml:"equip_slots"`
	Type        ItemType               `yaml:"type"`
	Subtype     ItemSubtype            `yaml:"subtype"`
	Properties  map[string]interface{} `yaml:"properties"`
}

type Item struct {
	InstanceID  string                 `yaml:"instance_id"`
	BlueprintID string                 `yaml:"blueprint_id"`
	Modifiers   map[string]int         `yaml:"modifiers"`
	Attachments []string               `yaml:"attachments"`
	Properties  map[string]interface{} `yaml:"properties"`
}

// Helper methods for item type checks
func (i *ItemBlueprint) IsWeapon() bool {
	return i.Type == ItemTypeWeapon
}

func (i *ItemBlueprint) IsArmor() bool {
	return i.Type == ItemTypeArmor
}

func (i *ItemBlueprint) IsFood() bool {
	return i.Type == ItemTypeFood
}

func NewItem(blueprint *ItemBlueprint) *Item {
	return &Item{
		InstanceID:  uuid.New().String(),
		BlueprintID: blueprint.ID,
		Modifiers:   make(map[string]int),
		Attachments: []string{},
	}
}
