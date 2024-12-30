package main

import (
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type Type string

const (
	TypeJunk Type = "junk"
)

type Item struct {
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
	Type        Type   `yaml:"type"`
}

func NewItem() *Item {
	return &Item{
		UUID: uuid.New().String(),
	}
}
