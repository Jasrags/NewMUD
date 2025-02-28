package game

import "github.com/Jasrags/NewMUD/internal/game/shared"

const (
	SkillsFilepath         = "_data/skills"
	SkillActiveFilepath    = "_data/skills/active"
	SkillGroupsFilepath    = "_data/skills/groups"
	SkillKnowledgeFilepath = "_data/skills/knowledge"
	SkillLanguagesFilepath = "_data/skills/languages"

	SkillTypeActive    = "Active"
	SkillTypeGroup     = "Group"
	SkillTypeKnowledge = "Knowledge"
	SkillTypeLanguage  = "Language"

	// Skill Categories
	SkillCategoryCombat        = "Combat Active"
	SkillCategoryMagical       = "Magical Active"
	SkillCategoryPhysical      = "Physical Active"
	SkillCategoryPseudoMagical = "Pseudo-Magical Active"
	SkillCategoryResonance     = "Resonance Active"
	SkillCategorySocial        = "Social Active"
	SkillCategoryTechnical     = "Technical Active"
	SkillCategoryVehicle       = "Vehicle Active"
	// Knowledge Skills
	SkillCategoryAcademic     = "Academic"
	SkillCategoryInterest     = "Interest"
	SkillCategoryLanguage     = "Language"
	SkillCategoryProfessional = "Professional"
	SkillCategoryStreet       = "Street"
)

type (
	SkillBlueprint struct {
		ID              string            `yaml:"id"`
		Name            string            `yaml:"name"`
		Type            string            `yaml:"type"`
		Description     string            `yaml:"description"`
		IsDefaultable   bool              `yaml:"is_defaultable"`
		LinkedAttribute string            `yaml:"linked_attribute"`
		Specializations []string          `yaml:"specializations"`
		RuleSource      shared.RuleSource `yaml:"rule_source"`
	}
	Skill struct {
		BlueprintID string `yaml:"blueprint_id"`
		// Name           string    `yaml:"name"`
		// Type           SkillType `yaml:"type"`
		Specialization string `yaml:"specialization"`
		Rating         int    `yaml:"rating"`
	}
	SkillGroup struct {
		ID          string            `yaml:"id,omitempty"`
		Name        string            `yaml:"name,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Skills      []string          `yaml:"skills"`
		RuleSource  shared.RuleSource `yaml:"rule_source,omitempty"`
	}
)

func (s *Skill) SetRating(rating int) {
	s.Rating = rating
}

func (s *Skill) SetSpecialization(specialization string) {
	s.Specialization = specialization
}
