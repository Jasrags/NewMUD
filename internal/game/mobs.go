package game

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
)

type Disposition string

const (
	DispositionFriendly   Disposition = "Friendly"
	DispositionNeutral    Disposition = "Neutral"
	DispositionAggressive Disposition = "Aggressive"
)

// TODO: Implement mob AI behaviors.
// TODO: Do we want mobs to be an "instance" that will persist after spawning?
type Mob struct {
	GameEntity            `yaml:",inline"`
	GeneralDisposition    Disposition            `yaml:"general_disposition"`
	CharacterDispositions map[string]Disposition `yaml:"character_dispositions"`
}

func NewMob() *Mob {
	return &Mob{
		GameEntity:            NewGameEntity(),
		GeneralDisposition:    DispositionNeutral,
		CharacterDispositions: make(map[string]Disposition),
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

func (m *Mob) SetGeneralDisposition(disposition Disposition) {
	m.GeneralDisposition = disposition
}

func (m *Mob) ReactToMessage(sender *Character, message string) {
	// Mobs can "react" based on predefined AI behaviors.
	m.ReactToInteraction(sender, message)
}

func (m *Mob) SetDispositionForCharacter(char *Character, disposition Disposition) {
	m.CharacterDispositions[char.ID] = disposition
}

func (m *Mob) GetDispositionForCharacter(char *Character) Disposition {
	if disposition, exists := m.CharacterDispositions[char.ID]; exists {
		return disposition
	}
	return m.GeneralDisposition // Fallback to general disposition
}

func (m *Mob) ReactToInteraction(sender *Character, message string) {
	if strings.Contains(strings.ToLower(message), "hello") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s says: 'Hello, %s.'}}::green\n", m.Name, sender.Name), nil)
		}
	} else if strings.Contains(strings.ToLower(message), "attack") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s snarls at %s and prepares to attack!}}::red\n", m.Name, sender.Name), nil)
		}
	} else {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s looks confused by %s's words.}}::yellow\n", m.Name, sender.Name), nil)
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
