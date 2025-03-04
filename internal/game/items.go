package game

import (
	"slices"
	"sync"

	ee "github.com/vansante/go-event-emitter"
)

const (
	ItemTagArmor     = "Armor"
	ItemTagJacket    = "Jacket"
	ItemTagSynthetic = "Synthetic"
	ItemTagLeather   = "Leather"

	ItemsFilepath = "data/items"
	ItemsFilename = "items.yml"

	ItemTypeJunk  = "junk"
	ItemTypeKey   = "key"
	ItemTypeArmor = "armor"

	ItemSubtypeNone = "None"

	EquipSlotNone  = "none"
	EquipSlotHead  = "head"
	EquipSlotBody  = "body"
	EquipSlotHands = "hands"
	EquipSlotLegs  = "legs"
)

var (
	EquipSlots = []string{EquipSlotHead, EquipSlotBody, EquipSlotHands, EquipSlotLegs}
)

// TODO: For keys we need subtypes for the different locks they can open.
// TODO: For picks they also need subtypes for the different locks they can pick.
// TODO: For locks they should have a rating of how difficult they are to pick
// TODO: For picks they should have a raiting of how good they are at picking locks
// TODO: Locks should somehow tie into alarm/traps or other events

type (
	ItemBlueprint struct {
		ID          string            `yaml:"id"`
		Name        string            `yaml:"name"`
		Description string            `yaml:"description"`
		Tags        []string          `yaml:"tags"`
		Weight      float64           `yaml:"weight"`
		Value       int               `yaml:"value"`
		BaseStats   map[string]int    `yaml:"base_stats"`
		EquipSlots  []string          `yaml:"equip_slots"`
		Modifiers   map[string]int    `yaml:"modifiers"`
		Attachments map[string]string `yaml:"attachments"`
		Type        string            `yaml:"type"`
		Subtype     string            `yaml:"subtype"`
	}

	// TODO: need to add the weight of attachments to the weight of the item
	ItemInstance struct {
		sync.RWMutex `yaml:"-"`
		Listeners    []ee.Listener `yaml:"-"`

		InstanceID  string         `yaml:"instance_id"`
		BlueprintID string         `yaml:"blueprint_id"`
		Blueprint   *ItemBlueprint `yaml:"-"`
		// Dynamic state fields
		Attachments map[string]string `yaml:"attachments"`
		NestedInv   *Inventory        `yaml:"nested_inventory"`
	}
)

func (ib *ItemBlueprint) HasTags(searchTags ...string) bool {
	for _, searchTag := range searchTags {
		if !slices.Contains(ib.Tags, searchTag) {
			return true
		}
	}

	return false
}
