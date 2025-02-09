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
