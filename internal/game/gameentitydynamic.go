package game

import (
	"log/slog"
	"slices"
	"strings"
)

const (
	PositionStanding    = "Standing"
	PositionSitting     = "Sitting"
	PositionKneeling    = "Kneeling"
	PositionLying       = "Lying"
	PositionProne       = "Prone"
	PositionCrouching   = "Crouching"
	PositionResting     = "Resting"
	PositionSleeping    = "Sleeping"
	PositionUnconscious = "Unconscious"
)

type GameEntityDynamic struct {
	Edge                  int                 `yaml:"edge"`
	PhysicalDamage        int                 `yaml:"physical_damage,omitempty"`
	StunDamage            int                 `yaml:"stun_damage,omitempty"`
	OverflowDamage        int                 `yaml:"overflow_damage,omitempty"`
	PositionState         string              `yaml:"position_state"`
	Inventory             Inventory           `yaml:"inventory,omitempty"`
	Equipment             Equipment           `yaml:"equipment,omitempty"`
	Qualtities            map[string]*Quality `yaml:"qualities,omitempty"`
	Skills                map[string]*Skill   `yaml:"skills,omitempty"`
	CharacterDispositions map[string]string   `yaml:"character_dispositions,omitempty"`
	Labels                []string            `yaml:"labels,omitempty"`
}

func NewGameEntityDynamic() GameEntityDynamic {
	var ged GameEntityDynamic

	ged.Inventory = NewInventory()
	ged.PositionState = PositionStanding
	ged.Equipment = NewEquipment()
	ged.Qualtities = make(map[string]*Quality)
	ged.Skills = make(map[string]*Skill)
	ged.Labels = make([]string, 0)
	ged.CharacterDispositions = make(map[string]string)

	return ged
}

func (ged *GameEntityDynamic) GetAllModifiers() map[string]int {
	modifiers := make(map[string]int)

	// TODO: add support for effects from spells/damage/etc
	for _, item := range ged.Equipment.Slots {
		for key, value := range item.Blueprint.Modifiers {
			slog.Info("Item Modifier",
				slog.String("key", key),
				slog.Int("value", value))
			modifiers[key] += value
		}
	}

	for _, quality := range ged.Qualtities {
		for key, value := range quality.Blueprint.Modifiers {
			slog.Info("Quality Modifier",
				slog.String("key", key),
				slog.Int("value", value))
			modifiers[key] += value
		}
	}

	return modifiers
}

func (ged *GameEntityDynamic) GetQuality(qualityID string) *Quality {
	quality, ok := ged.Qualtities[qualityID]
	if !ok {
		return nil
	}

	return quality
}

func (ged *GameEntityDynamic) FormatQualtities() string {
	qualities := make([]string, len(ged.Qualtities))
	for _, quality := range ged.Qualtities {
		qualities = append(qualities, quality.FormatListItem())
	}
	return strings.Join(qualities, "\n")
}

func (ged *GameEntityDynamic) GetSkill(skillID string) *Skill {
	skill, ok := ged.Skills[skillID]
	if !ok {
		return nil
	}

	return skill
}

// TODO: Add support for adding/removing labels

func (ged *GameEntityDynamic) GetEdge() int {
	return ged.Edge
}

func (ged *GameEntityDynamic) GetCharacterDisposition(characterID string) string {
	disposition, ok := ged.CharacterDispositions[characterID]
	if !ok {
		return DispositionNeutral
	}
	return disposition
}

func (ged *GameEntityDynamic) GetCurrentCarryWeight() float64 {
	weight := 0.0
	for _, item := range ged.Inventory.Items {
		weight += item.Blueprint.Weight
	}
	return weight
}

func (ged *GameEntityDynamic) AddLabel(label string) {
	lowerLabel := strings.ToLower(label)
	if !slices.Contains(ged.Labels, lowerLabel) {
		ged.Labels = append(ged.Labels, lowerLabel)
	}
}

func (ged *GameEntityDynamic) HasLabels(labels ...string) bool {
	for _, label := range labels {
		if found := slices.Contains(ged.Labels, strings.ToLower(label)); found {
			return true
		}
	}

	return false
}

func (ged *GameEntityDynamic) RemoveLabel(label string) {
	for i, l := range ged.Labels {
		if l == strings.ToLower(label) {
			ged.Labels = slices.Delete(ged.Labels, i, i+1)
			return
		}
	}
}

// Validate the game entity dynamic
// TODO: Implement validation
func (ged *GameEntityDynamic) Validate() error {
	return nil
}
