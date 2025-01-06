package main

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/i582/cfmt/cmd/cfmt"
	ee "github.com/vansante/go-event-emitter"
)

type Mob struct {
	sync.RWMutex
	Listeners []ee.Listener `yaml:"-"`

	ID          string           `yaml:"id"`
	ReferenceID string           `yaml:"reference_id"`
	UUID        string           `yaml:"uuid"`
	Area        *Area            `yaml:"-"`
	AreaID      string           `yaml:"area_id"`
	Room        *Room            `yaml:"-"`
	RoomID      string           `yaml:"room_id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Inventory   Inventory        `yaml:"inventory"`
	Equipment   map[string]*Item `yaml:"equipment"`
}

func NewMob() *Mob {
	return &Mob{
		UUID: uuid.New().String(),
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
