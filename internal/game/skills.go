package game

import "github.com/Jasrags/NewMUD/internal/game/shared"

const (
	SkillsFilepath         = "_data/skills"
	SkillActiveFilepath    = "_data/skills/active"
	SkillGroupsFilepath    = "_data/skills/groups"
	SkillKnowledgeFilepath = "_data/skills/knowledge"
	SkillLanguagesFilepath = "_data/skills/languages"

	SkillTypeActive    SkillType = "Active"
	SkillTypeGroup     SkillType = "Group"
	SkillTypeKnowledge SkillType = "Knowledge"
	SkillTypeLanguage  SkillType = "Language"

	// Skill Categories
	SkillCategoryCombat        SkillCategory = "Combat Active"
	SkillCategoryMagical       SkillCategory = "Magical Active"
	SkillCategoryPhysical      SkillCategory = "Physical Active"
	SkillCategoryPseudoMagical SkillCategory = "Pseudo-Magical Active"
	SkillCategoryResonance     SkillCategory = "Resonance Active"
	SkillCategorySocial        SkillCategory = "Social Active"
	SkillCategoryTechnical     SkillCategory = "Technical Active"
	SkillCategoryVehicle       SkillCategory = "Vehicle Active"
	// Knowledge Skills
	SkillCategoryAcademic     SkillCategory = "Academic"
	SkillCategoryInterest     SkillCategory = "Interest"
	SkillCategoryLanguage     SkillCategory = "Language"
	SkillCategoryProfessional SkillCategory = "Professional"
	SkillCategoryStreet       SkillCategory = "Street"
)

type (
	SkillType     string
	SkillCategory string

	SkillBlueprint struct {
		ID              string            `yaml:"id"`
		Name            string            `yaml:"name"`
		Type            SkillType         `yaml:"type"`
		Description     string            `yaml:"description"`
		IsDefaultable   bool              `yaml:"is_defaultable"`
		LinkedAttribute AttributeType     `yaml:"linked_attribute"`
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
