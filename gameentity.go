package main

import (
	"sync"

	ee "github.com/vansante/go-event-emitter"
)

type GameEntity struct {
	sync.RWMutex
	Listeners []ee.Listener `yaml:"-"`

	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Room        *Room            `yaml:"-"`
	RoomID      string           `yaml:"room_id"`
	Area        *Area            `yaml:"-"`
	AreaID      string           `yaml:"area_id"`
	Inventory   Inventory        `yaml:"inventory"`
	Equipment   map[string]*Item `yaml:"equipment"`
}

func (e *GameEntity) GetName() string {
	return e.Name
}

func (e *GameEntity) GetID() string {
	return e.ID
}

func (e *GameEntity) SetRoom(room *Room) {
	e.Room = room
	e.RoomID = room.ReferenceID
}
