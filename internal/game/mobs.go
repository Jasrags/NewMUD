package game

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
)

const (
	MobsFilename = "mobs.yml"

	DispositionFriendly   MobDisposition = "Friendly"
	DispositionNeutral    MobDisposition = "Neutral"
	DispositionAggressive MobDisposition = "Aggressive"
)

type (
	MobDisposition string

	// TODO: Implement mob AI behaviors.
	// TODO: Do we want mobs to be an "instance" that will persist after spawning?
	Mob struct {
		GameEntity            `yaml:",inline"`
		Tags                  []string                  `yaml:"tags"`
		ProfessionalRating    int                       `yaml:"professional_rating"`
		GeneralDisposition    MobDisposition            `yaml:"general_disposition"`
		CharacterDispositions map[string]MobDisposition `yaml:"character_dispositions"`
	}
)

func NewMob() *Mob {
	return &Mob{
		GameEntity:            NewGameEntity(),
		GeneralDisposition:    DispositionNeutral,
		CharacterDispositions: make(map[string]MobDisposition),
	}
}

func (m *Mob) Init() {
	slog.Debug("Initializing mob",
		slog.String("mob_id", m.ID))
}

func (m *Mob) GetName() string {
	return m.Name
}

func (m *Mob) GetID() string {
	return m.ID
}

func (m *Mob) SetGeneralDisposition(disposition MobDisposition) {
	m.GeneralDisposition = disposition
}

func (m *Mob) ReactToMessage(sender *Character, message string) {
	// Mobs can "react" based on predefined AI behaviors.
	m.ReactToInteraction(sender, message)
}

func (m *Mob) SetDispositionForCharacter(char *Character, disposition MobDisposition) {
	m.CharacterDispositions[char.ID] = disposition
}

func (m *Mob) GetDispositionForCharacter(char *Character) MobDisposition {
	if disposition, exists := m.CharacterDispositions[char.ID]; exists {
		return disposition
	}
	return m.GeneralDisposition // Fallback to general disposition
}

func (m *Mob) ReactToInteraction(sender *Character, message string) {
	if strings.Contains(strings.ToLower(message), "hello") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s says: 'Hello, %s.'}}::green"+CRLF, m.Name, sender.Name), nil)
		}
	} else if strings.Contains(strings.ToLower(message), "attack") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s snarls at %s and prepares to attack!}}::red"+CRLF, m.Name, sender.Name), nil)
		}
	} else {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s looks confused by %s's words.}}::yellow"+CRLF, m.Name, sender.Name), nil)
		}
	}
}

func DescribeMobDisposition(mob *Mob, char *Character) string {
	disposition := mob.GetDispositionForCharacter(char)
	switch disposition {
	case DispositionFriendly:
		return fmt.Sprintf("%s looks at you warmly.", mob.Name)
	case DispositionNeutral:
		return fmt.Sprintf("%s glances at you indifferently.", mob.Name)
	case DispositionAggressive:
		return fmt.Sprintf("%s snarls menacingly at you!", mob.Name)
	default:
		return fmt.Sprintf("%s's demeanor is unreadable.", mob.Name)
	}
}

// RenderMobTable builds a formatted table of a mob's stats.
// It leverages the embedded GameEntity fields from Mob.
func RenderMobTable(mob *Mob) string {
	// Optionally, recalculate attributes if needed.
	mob.Recalculate()

	var builder strings.Builder

	// Header: basic details from GameEntity.
	builder.WriteString(cfmt.Sprintf("{{Name:}}::white|bold {{%s}}::cyan"+CRLF, mob.Name))
	builder.WriteString(cfmt.Sprintf("{{ID:}}::white|bold {{%s}}::cyan"+CRLF, mob.ID))
	builder.WriteString(cfmt.Sprintf("{{Title:}}::white|bold {{%s}}::cyan"+CRLF, mob.Title))
	builder.WriteString(cfmt.Sprintf("{{Description:}}::white|bold {{%s}}::cyan"+CRLF, mob.Description))
	builder.WriteString(cfmt.Sprintf("{{Long Description:}}::white|bold {{%s}}::cyan"+CRLF, mob.LongDescription))
	builder.WriteString(CRLF)

	// Mob-specific data.
	builder.WriteString(cfmt.Sprintf("{{Professional Rating:}}::white|bold {{%d}}::cyan"+CRLF, mob.ProfessionalRating))
	builder.WriteString(cfmt.Sprintf("{{General Disposition:}}::white|bold {{%s}}::cyan"+CRLF, mob.GeneralDisposition))
	builder.WriteString(CRLF)

	// Attributes from the embedded GameEntity.
	builder.WriteString(cfmt.Sprintf("{{Body:}}::white|bold {{%d}}::cyan"+CRLF, mob.Body.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Agility:}}::white|bold {{%d}}::cyan"+CRLF, mob.Agility.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Reaction:}}::white|bold {{%d}}::cyan"+CRLF, mob.Reaction.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Strength:}}::white|bold {{%d}}::cyan"+CRLF, mob.Strength.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Willpower:}}::white|bold {{%d}}::cyan"+CRLF, mob.Willpower.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Logic:}}::white|bold {{%d}}::cyan"+CRLF, mob.Logic.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Intuition:}}::white|bold {{%d}}::cyan"+CRLF, mob.Intuition.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Charisma:}}::white|bold {{%d}}::cyan"+CRLF, mob.Charisma.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Essence:}}::white|bold {{%.1f}}::cyan"+CRLF, mob.Essence.TotalValue))
	if mob.Magic.Base > 0 {
		builder.WriteString(cfmt.Sprintf("{{Magic:}}::white|bold {{%d}}::cyan"+CRLF, mob.Magic.TotalValue))
	}
	if mob.Resonance.Base > 0 {
		builder.WriteString(cfmt.Sprintf("{{Resonance:}}::white|bold {{%d}}::cyan"+CRLF, mob.Resonance.TotalValue))
	}

	return builder.String()
}
