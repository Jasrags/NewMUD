package game

import (
	"strings"

	"github.com/Jasrags/NewMUD/internal/game/shared"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/wordwrap"
)

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

	SkillCategoryCombat        = "Combat Active"
	SkillCategoryMagical       = "Magical Active"
	SkillCategoryPhysical      = "Physical Active"
	SkillCategoryPseudoMagical = "Pseudo-Magical Active"
	SkillCategoryResonance     = "Resonance Active"
	SkillCategorySocial        = "Social Active"
	SkillCategoryTechnical     = "Technical Active"
	SkillCategoryVehicle       = "Vehicle Active"

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
		IsDefaultable   bool              `yaml:"is_defaultable,omitempty"`
		LinkedAttribute string            `yaml:"linked_attribute"`
		Specializations []string          `yaml:"specializations"`
		RuleSource      shared.RuleSource `yaml:"rule_source"`
	}
	Skill struct {
		BlueprintID    string          `yaml:"blueprint_id"`
		Blueprint      *SkillBlueprint `yaml:"-"`
		Specialization string          `yaml:"specialization,omitempty"`
		Rating         int             `yaml:"rating"`
	}
	SkillGroup struct {
		ID          string            `yaml:"id,omitempty"`
		Name        string            `yaml:"name,omitempty"`
		Description string            `yaml:"description,omitempty"`
		Skills      []string          `yaml:"skills"`
		RuleSource  shared.RuleSource `yaml:"rule_source,omitempty"`
	}
)

func NewSkill(bp *SkillBlueprint, rating int, specialization string) *Skill {
	return &Skill{
		BlueprintID:    bp.ID,
		Blueprint:      bp,
		Specialization: specialization,
		Rating:         rating,
	}
}

func (s *Skill) FormatListItem() string {
	var sb strings.Builder
	sb.WriteString(cfmt.Sprintf("%s %d", s.Blueprint.Name+":", s.Rating))
	// TODO: Specializations should pull in the item blueprint and display the name
	if s.Specialization != "" {
		sb.WriteString(cfmt.Sprintf(" (%s)", s.Specialization))
	}

	return sb.String()
}

func (s *Skill) FormatDetailed() string {
	var sb strings.Builder
	sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %s"+CRLF, "Name:", s.Blueprint.Name))
	sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %s"+CRLF, "Type:", s.Blueprint.Type))
	sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %s"+CRLF, "Attribute:", s.Blueprint.LinkedAttribute))
	sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %s"+CRLF, "Description:", s.Blueprint.Description))
	sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %d"+CRLF, "Rating:", s.Rating))
	// TODO: Specializations should pull in the item blueprint and display the name
	if s.Specialization != "" {
		sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %s"+CRLF, "Specialization:", s.Specialization))
	}
	sb.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold %v"+CRLF, "Defaultable:", s.Blueprint.IsDefaultable))

	return wordwrap.String(sb.String(), 80)
}
