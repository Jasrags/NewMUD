package game

import (
	"slices"
	"strings"
	"sync"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/wordwrap"
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

func (i *ItemInstance) FormatListItem() string {
	var sb strings.Builder
	sb.WriteString(cfmt.Sprintf("%s", i.Blueprint.Name))
	return sb.String()
}

func (i *ItemInstance) FormatDetailed() string {
	var sb strings.Builder
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Name:", i.Blueprint.Name))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Description:", i.Blueprint.Description))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %.2f"+CRLF, "Weight:", i.Blueprint.Weight))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %d"+CRLF, "Value:", i.Blueprint.Value))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Type:", i.Blueprint.Type))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Subtype:", i.Blueprint.Subtype))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Tags:", strings.Join(i.Blueprint.Tags, ", ")))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Equip Slots:", strings.Join(i.Blueprint.EquipSlots, ", ")))
	// sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Base Stats:", i.Blueprint.BaseStats))
	// sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Modifiers:", i.Blueprint.Modifiers))
	// sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Attachments:", i.Blueprint.Attachments))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Instance ID:", i.InstanceID))

	return wordwrap.String(sb.String(), 80)
}

func (ib *ItemBlueprint) HasTags(searchTags ...string) bool {
	for _, searchTag := range searchTags {
		if !slices.Contains(ib.Tags, searchTag) {
			return true
		}
	}

	return false
}
