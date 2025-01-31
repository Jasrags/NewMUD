package game

import "github.com/Jasrags/NewMUD/internal/game/shared"

const (
	QualitiesFilepath = "_data/qualities"

	QualityTypePositive QualityType = "Positive"
	TypeNegative        QualityType = "Negative"
)

type (
	QualityType string

	QualityBlueprint struct {
		ID          string            `yaml:"id"`
		Type        QualityType       `yaml:"type"`
		Name        string            `yaml:"name"`
		Description string            `yaml:"description"`
		Modifiers   []string          `yaml:"modifiers"`
		Cost        int               `yaml:"cost"`
		RuleSource  shared.RuleSource `yaml:"rule_source"`
	}

	Quality struct {
		BlueprintID string      `yaml:"blueprint_id"`
		Type        QualityType `yaml:"type"`
		Name        string      `yaml:"name"`
		Rating      int         `yaml:"rating"`
	}
)

func (q *Quality) SetRating(rating int) {
	q.Rating = rating
}
