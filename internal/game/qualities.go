package game

import (
	"strings"

	"github.com/Jasrags/NewMUD/internal/game/shared"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/wordwrap"
)

const (
	QualitiesFilepath = "_data/qualities"

	QualityTypePositive QualityType = "Positive"
	QualityTypeNegative QualityType = "Negative"
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
		BlueprintID string            `yaml:"blueprint_id"`
		Blueprint   *QualityBlueprint `yaml:"-"`
		Rating      int               `yaml:"rating"`
	}
)

func NewQuality(bp *QualityBlueprint, rating int) *Quality {
	return &Quality{
		BlueprintID: bp.ID,
		Blueprint:   bp,
		Rating:      rating,
	}
}

func (q *Quality) FormatListItem() string {
	nameColor := "green"
	if q.Blueprint.Type == QualityTypeNegative {
		nameColor = "red"
	}

	var sb strings.Builder
	if q.Rating != 0 {
		sb.WriteString(cfmt.Sprintf("{{%s}}::%s %d", q.Blueprint.Name+":", nameColor, q.Rating))
	} else {
		sb.WriteString(cfmt.Sprintf("{{%s}}::%s", q.Blueprint.Name, nameColor))
	}

	return sb.String()
}

func (q *Quality) FormatDetailed() string {
	var sb strings.Builder

	nameColor := "green"
	if q.Blueprint.Type == QualityTypeNegative {
		nameColor = "red"
	}

	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Name:", q.Blueprint.Name))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold {{%s}}::%s"+CRLF, "Type:", q.Blueprint.Type, nameColor))
	if len(q.Blueprint.Modifiers) > 0 {
		sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Modifiers:", strings.Join(q.Blueprint.Modifiers, ", ")))
	}
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %d"+CRLF, "Cost:", q.Blueprint.Cost))
	if q.Rating != 0 {
		sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %d"+CRLF, "Rating:", q.Rating))
	}
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Description:", q.Blueprint.Description))
	sb.WriteString(cfmt.Sprintf("{{%-12s}}::white|bold %s"+CRLF, "Rule Source:", q.Blueprint.RuleSource))

	return wordwrap.String(sb.String(), 80)
}
