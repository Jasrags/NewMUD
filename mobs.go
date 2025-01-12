package main

import (
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/i582/cfmt/cmd/cfmt"
)

// TODO: Implement mob AI behaviors.
// TODO: Do we want mobs to be an "instance" that will persist after spawning?
type Mob struct {
	GameEntity `yaml:",inline"`
}

func NewMob() *Mob {
	return &Mob{
		GameEntity: GameEntity{
			ID:        uuid.New().String(),
			Equipment: make(map[string]*Item),
		},
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

func (m *Mob) ReactToMessage(sender *Character, message string) {
	// Mobs can "react" based on predefined AI behaviors.
	m.ReactToInteraction(sender, message)
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
