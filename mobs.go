package main

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type Mob struct {
	sync.RWMutex
	Listeners []ee.Listener `yaml:"-"`

	ID          string `yaml:"id"`
	ReferenceID string `yaml:"reference_id"`
	UUID        string `yaml:"uuid"`
	Area        *Area  `yaml:"-"`
	AreaID      string `yaml:"area_id"`
	Room        *Room  `yaml:"-"`
	RoomID      string `yaml:"room_id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
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
