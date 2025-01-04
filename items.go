package main

import (
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type DefaultItem struct {
	ID               string `yaml:"id"`
	RespawnChance    int    `yaml:"respawn_chance"`
	MaxLoad          int    `yaml:"max_load"`
	ReplaceOnRespawn bool   `yaml:"replace_on_respawn"`
	Quantity         int    `yaml:"quantity"`
}

type Type string

const (
	TypeJunk Type = "junk"
)

type Item struct {
	sync.RWMutex `yaml:"-"`
	Listeners    []ee.Listener `yaml:"-"`

	ID          string `yaml:"id"`
	ReferenceID string `yaml:"reference_id"`
	Area        *Area  `yaml:"-"`
	AreaID      string `yaml:"area_id"`
	Room        *Room  `yaml:"-"`
	RoomID      string `yaml:"room_id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        Type   `yaml:"type"`
}

func NewItem() *Item {
	return &Item{
		ID: uuid.New().String(),
	}
}
