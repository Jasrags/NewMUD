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
	ItemsFilepath = "data/items"
	ItemsFilename = "items.yml"

	ItemTagArmor     = "Armor"
	ItemTagJacket    = "Jacket"
	ItemTagSynthetic = "Synthetic"
	ItemTagLeather   = "Leather"
	ItemTagWeapon    = "Weapon"

	ItemTypeJunk = "junk"
	ItemTypeKey  = "key"

	ItemSubtypeNone  = "None"
	ItemSubtypeMelee = "Melee"

	EquipSlotNone    = "none"
	EquipSlotHead    = "head"
	EquipSlotBody    = "body"
	EquipSlotHands   = "hands"
	EquipSlotLegs    = "legs"
	EquipSlotWeapon  = "weapon"
	EquipSlotOffhand = "offhand"

	WeaponDamagePhysical = "Physical"
	WeaponDamageStun     = "Stun"

	LegalityTypeLegal      = "Legal"
	LegalityTypeRestricted = "Restricted"
	LegalityTypeForbidden  = "Forbidden"

	// Item types
	ItemTypeArmor  = "Armor"
	ItemTypeWeapon = "Weapon"

	// Armor categories
	ItemCategoryArmor    = "Armor"
	ItemCategoryClothing = "Clothing"

	// Weapon categories
	ItemCategoryBlades = "Blades"
	ItemCategoryClubs  = "Clubs"

	MountPointUnderBarrel = "Under-Barrel"
	MountPointBarrel      = "Barrel"
	MountPointStock       = "Stock"
	MountPointTop         = "Top"
	MountPointSide        = "Side"
	MountPointInternal    = "Internal"

	WeaponRangedReloadBreakAction        = "b"
	WeaponRangedReloadDetachableMagazine = "c"
	WeaponRangedReloadDrum               = "d"
	WeaponRangedReloadMuzzleLoader       = "ml"
	WeaponRangedReloadInternalMagazine   = "m"
	WeaponRangedReloadCylinder           = "cy"
	WeaponRangedReloadBelt               = "belt"

	WeaponFiringModeSingleShot    = "Single-Shot"
	WeaponFiringModeSemiAutomatic = "Semi-Automatic"
	WeaponFiringModeBurstFire     = "Burst Fire"
	WeaponFiringModeLongBurst     = "Long Burst"
	WeaponFiringModeFullAuto      = "Full Auto"
)

var (
	EquipSlots = []string{
		EquipSlotHead,
		EquipSlotBody,
		EquipSlotHands,
		EquipSlotLegs,
		EquipSlotWeapon,
		EquipSlotOffhand,
	}
)

// TODO: For keys we need subtypes for the different locks they can open.
// TODO: For picks they also need subtypes for the different locks they can pick.
// TODO: For locks they should have a rating of how difficult they are to pick
// TODO: For picks they should have a raiting of how good they are at picking locks
// TODO: Locks should somehow tie into alarm/traps or other events

type (
	Damage struct {
		Attribute string `yaml:"attribute,omitempty"`
		Type      string `yaml:"type,omitempty"`
		Value     int    `yaml:"value,omitempty"`
	}
	ItemBlueprint struct {
		ID           string         `yaml:"id"`
		Hide         bool           `yaml:"hide"`
		Type         string         `yaml:"type"`
		Category     string         `yaml:"category"`
		Name         string         `yaml:"name"`
		Description  string         `yaml:"description"`
		Availability int            `yaml:"availability,omitempty"`
		Legality     string         `yaml:"legality,omitempty"`
		Cost         int            `yaml:"cost,omitempty"`
		Weight       float64        `yaml:"weight"`
		Tags         []string       `yaml:"tags"`
		EquipSlots   []string       `yaml:"equip_slots"`
		Modifiers    map[string]int `yaml:"modifiers,omitempty"`
		// Needed?
		BaseStats          map[string]int    `yaml:"base_stats"`
		AllowedAttachments []string          `yaml:"allowed_attachments,omitempty"`
		Attachments        map[string]string `yaml:"attachments,omitempty"`
		// Armor
		ArmorValue    int `yaml:"armor_value,omitempty"`
		ArmorCapacity int `yaml:"armor_capacity,omitempty"`
		GearCapacity  int `yaml:"gear_capacity,omitempty"`
		MaxRating     int `yaml:"max_rating,omitempty"`
		// Weapons
		Conceal          int      `yaml:"conceal,omitempty"`
		Accuracy         int      `yaml:"accuracy,omitempty"`
		Reach            int      `yaml:"reach,omitempty"`
		Damage           Damage   `yaml:"damage,omitempty"`
		ArmorPenetration int      `yaml:"armor_penetration,omitempty"`
		FireModes        []string `yaml:"fire_modes,omitempty"`
		Recoil           int      `yaml:"recoil,omitempty"`
		AmmoCapacity     int      `yaml:"ammo_capacity,omitempty"`
		AmmoTypes        []string `yaml:"ammo_type,omitempty"`
		ReloadType       string   `yaml:"reload_type,omitempty"`
	}

	// TODO: need to add the weight of attachments to the weight of the item
	ItemInstance struct {
		sync.RWMutex `yaml:"-"`
		Listeners    []ee.Listener `yaml:"-"`

		InstanceID  string         `yaml:"instance_id"`
		BlueprintID string         `yaml:"blueprint_id"`
		Blueprint   *ItemBlueprint `yaml:"-"`
		// Dynamic state fields
		Attachments map[string]string `yaml:"attachments,omitempty"`
		NestedInv   *Inventory        `yaml:"nested_inventory,omitempty"`

		// Weapons
		SelectedFireMode string `yaml:"selected_fire_mode,omitempty"`
		AmmoCount        int    `yaml:"ammo_count,omitempty"`
		AmmoType         string `yaml:"ammo_type,omitempty"`
		// Armor
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
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %d"+CRLF, "Value:", i.Blueprint.Cost))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Type:", i.Blueprint.Type))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Subtype:", i.Blueprint.Category))
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
