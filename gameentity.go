package main

import (
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type GameEntity struct {
	sync.RWMutex `yaml:"-"`
	Listeners    []ee.Listener `yaml:"-"`

	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Attributes  Attributes       `yaml:"attributes"`
	Room        *Room            `yaml:"-"`
	RoomID      string           `yaml:"room_id"`
	Area        *Area            `yaml:"-"`
	AreaID      string           `yaml:"area_id"`
	Inventory   Inventory        `yaml:"inventory"`
	Equipment   map[string]*Item `yaml:"equipment"`
}

func NewGameEntity() GameEntity {
	return GameEntity{
		ID: uuid.New().String(),
		// Attributes: NewAttributes(),
		// Inventory:  NewInventory(),
		Equipment: make(map[string]*Item),
		Listeners: make([]ee.Listener, 0),
	}
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

func (e *GameEntity) GetLiftCarry() int {
	baseCarryWeight := 10
	return (e.Attributes.Strength.TotalValue + e.Attributes.Body.TotalValue) * baseCarryWeight
}

func (e *GameEntity) GetCurrentCarryWeight() float64 {
	totalWeight := 0.0

	for _, item := range e.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp != nil {
			totalWeight += bp.Weight
		}
	}

	return totalWeight
}
