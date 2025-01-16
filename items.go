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
	ItemTypeJunk ItemType = "junk"
	ItemTypeKey  ItemType = "key"

	ItemSubtypeNone ItemSubtype = "none"

	EquipSlotNone EquipSlot = "none"
	EquipSlotHead EquipSlot = "head"
)

type ItemBlueprint struct {
	ID          string         `yaml:"id"`
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Tags        []string       `yaml:"tags"`
	Weight      float64        `yaml:"weight"`
	Value       int            `yaml:"value"`
	BaseStats   map[string]int `yaml:"base_stats"`
	EquipSlots  []EquipSlot    `yaml:"equip_slots"`
	Type        ItemType       `yaml:"type"`
	Subtype     ItemSubtype    `yaml:"subtype"`
}

// TODO: need to add the weight of attachments to the weight of the item
type Item struct {
	InstanceID  string         `yaml:"instance_id"`
	BlueprintID string         `yaml:"blueprint_id"`
	Modifiers   map[string]int `yaml:"modifiers"`
	Attachments []string       `yaml:"attachments"`
	NestedInv   *Inventory     `yaml:"nested_inventory"`
}

func NewItem(blueprint *ItemBlueprint) *Item {
	return &Item{
		InstanceID:  uuid.New().String(),
		BlueprintID: blueprint.ID,
		Modifiers:   make(map[string]int),
		Attachments: []string{},
	}
}
