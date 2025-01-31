package game

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
		ID                  string     `yaml:"id"`
		PointCost           int        `yaml:"point_cost"`
		Name                string     `yaml:"name"`
		Category            string     `yaml:"category"`
		Description         string     `yaml:"description"`
		Attributes          Attributes `yaml:"attributes"`
		Hidden              bool       `yaml:"hidden"`
		Qualities           []string   `yaml:"qualities"`
		QualityRestrictions []string   `yaml:"quality_restrictions"`
		RuleSource          string     `yaml:"rule_source"`
	}
)
