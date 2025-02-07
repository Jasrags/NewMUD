package game

import (
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
)

const (
	MetatypesFilepath = "_data/metatypes"

	RacialTraitLowLightVision                    RacialTrait = "Low-Light Vision"
	RacialTraitThermographicVision               RacialTrait = "Thermographic Vision"
	RacialTrait2DicForPathogenAndToxinResistance RacialTrait = "+2 dice for pathogen and toxin resistance"
	RacialTrait20PercentIncreasedLifestyleCost   RacialTrait = "+20% increased Lifestyle cost"
	RacialTrait1Reach                            RacialTrait = "+1 Reach"
	RacialTrait1DermalArmor                      RacialTrait = "+1 dermal armor"
	RacialTraitDoubleLifestyleCosts              RacialTrait = "Double Lifestyle costs"

	MetatypeCategoryMetahuman    MetatypeCategory = "Metahuman"
	MetatypeCategoryMetavariant  MetatypeCategory = "Metavariant"
	MetatypeCategoryMetasapient  MetatypeCategory = "Metasapient"
	MetatypeCategoryShapeshifter MetatypeCategory = "Shapeshifter"

	MetatypeNameHuman MetatypeName = "Human"
	MetatypeNameElf   MetatypeName = "Elf"
	MetatypeNameDwarf MetatypeName = "Dwarf"
	MetatypeNameOrk   MetatypeName = "Ork"
	MetatypeNameTroll MetatypeName = "Troll"
)

type (
	MetatypeName     string
	RacialTrait      string
	MetatypeCategory string

	Metatype struct {
		ID          string `yaml:"id"`
		PointCost   int    `yaml:"point_cost"`
		Name        string `yaml:"name"`
		Category    string `yaml:"category"`
		Description string `yaml:"description"`

		Body      Attribute[int]     `yaml:"body"`
		Agility   Attribute[int]     `yaml:"agility"`
		Reaction  Attribute[int]     `yaml:"reaction"`
		Strength  Attribute[int]     `yaml:"strength"`
		Willpower Attribute[int]     `yaml:"willpower"`
		Logic     Attribute[int]     `yaml:"logic"`
		Intuition Attribute[int]     `yaml:"intuition"`
		Charisma  Attribute[int]     `yaml:"charisma"`
		Essence   Attribute[float64] `yaml:"essence"`
		Magic     Attribute[int]     `yaml:"magic"`
		Resonance Attribute[int]     `yaml:"resonance"`

		// Attributes          Attributes `yaml:"attributes"`
		Hidden              bool     `yaml:"hidden"`
		Qualities           []string `yaml:"qualities"`
		QualityRestrictions []string `yaml:"quality_restrictions"`
		RuleSource          string   `yaml:"rule_source"`
	}
)

func (m *Metatype) GetSelectionInfo() string {
	var output strings.Builder
	output.WriteString(cfmt.Sprintf("{{Name:}}::white|bold %s (%s)"+CRLF, m.Name, m.Category))
	output.WriteString(cfmt.Sprintf("{{Description:}}::white|bold %s"+CRLF, m.Description))

	if len(m.Qualities) > 0 {
		output.WriteString("{{Qualities:}}::white|bold " + CRLF)
		for _, qualityID := range m.Qualities {
			quality := EntityMgr.GetQualityBlueprint(qualityID)
			output.WriteString(cfmt.Sprintf("  {{-}}::white|bold {{%s}}::cyan"+CRLF, quality.Name))
		}
	}

	if len(m.QualityRestrictions) > 0 {
		output.WriteString("{{Qualities Restrictions:}}::white|bold " + CRLF)
		for _, qualityID := range m.QualityRestrictions {
			quality := EntityMgr.GetQualityBlueprint(qualityID)
			output.WriteString(cfmt.Sprintf("  {{-}}::white|bold {{%s}}::red"+CRLF, quality.Name))
		}
	}
	return output.String()
}
