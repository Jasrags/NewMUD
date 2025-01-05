package main

import "github.com/google/uuid"

type (
	ItemType    string
	ItemSubtype string
	EquipSlot   string
)

const (
	ItempTypeJunk ItemType = "junk"

	ItemSubtypeNone ItemSubtype = "none"

	EquipSlotNone EquipSlot = "none"
	EquipSlotHead EquipSlot = "head"
)

type ItemBlueprint struct {
	ID          string         `yaml:"id"`
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Tags        []string       `yaml:"tags"`
	BaseStats   map[string]int `yaml:"base_stats"`
	EquipSlots  []EquipSlot    `yaml:"equip_slots"`
	Type        ItemType       `yaml:"type"`
	Subtype     ItemSubtype    `yaml:"subtype"`
}

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
