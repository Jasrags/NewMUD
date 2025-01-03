package main

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type Exit struct {
	Room      *Room  `yaml:"-"`
	RoomID    string `yaml:"room_id"`
	Direction string `yaml:"direction"`
}

type Corrdinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

// TODO: Add Doors and Locks
type Room struct {
	sync.RWMutex
	Listeners []ee.Listener `yaml:"-"`

	ID           string          `yaml:"id"`
	ReferenceID  string          `yaml:"reference_id"`
	UUID         string          `yaml:"uuid"`
	AreaID       string          `yaml:"area_id"`
	Area         *Area           `yaml:"-"`
	Title        string          `yaml:"title"`
	Description  string          `yaml:"description"`
	Exits        map[string]Exit `yaml:"exits"`
	Corrdinates  *Corrdinates    `yaml:"corrdinates"`
	Items        []*Item         `yaml:"-"`
	Characters   []*Character    `yaml:"-"`
	Mobs         []*Mob          `yaml:"-"`
	DefaultItems []string        `yaml:"default_items"` // IDs of items to load into the room
	DefaultMobs  []string        `yaml:"default_mobs"`  // IDs of mobs to load into the room
	SpawnedMobs  []*Mob          `yaml:"-"`             // Mobs that have been spawned into the room
}

func NewRoom() *Room {
	return &Room{
		UUID:  uuid.New().String(),
		Exits: make(map[string]Exit),
	}
}

// func (r *Room) Init() {
// 	slog.Debug("Initializing room",
// 		slog.String("room_id", r.ID))

// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomCharacterEnter, r.onRoomCharacterEnter))
// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomCharacterLeave, r.onRoomCharacterLeave))
// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomMobEnter, r.onRoomMobEnter))
// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomMobLeave, r.onRoomMobLeave))
// 	// r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomSpawn, r.onRoomSpawn))
// 	// r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomUpdate, r.onRoomUpdate))
// }

func (r *Room) GetExits() {
	r.RLock()
	defer r.RUnlock()

	// adjacents := map[string]Corrdinates{
	// 	"north":     {X: r.Corrdinates.X, Y: r.Corrdinates.Y + 1, Z: r.Corrdinates.Z},
	// 	"south":     {X: r.Corrdinates.X, Y: r.Corrdinates.Y - 1, Z: r.Corrdinates.Z},
	// 	"east":      {X: r.Corrdinates.X + 1, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z},
	// 	"west":      {X: r.Corrdinates.X - 1, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z},
	// 	"up":        {X: r.Corrdinates.X, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z + 1},
	// 	"down":      {X: r.Corrdinates.X, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z - 1},
	// 	"northeast": {X: r.Corrdinates.X + 1, Y: r.Corrdinates.Y + 1, Z: r.Corrdinates.Z},
	// 	"northwest": {X: r.Corrdinates.X - 1, Y: r.Corrdinates.Y + 1, Z: r.Corrdinates.Z},
	// 	"southeast": {X: r.Corrdinates.X + 1, Y: r.Corrdinates.Y - 1, Z: r.Corrdinates.Z},
	// 	"southwest": {X: r.Corrdinates.X - 1, Y: r.Corrdinates.Y - 1, Z: r.Corrdinates.Z},
	// }

	// var exits []string
	// for direction := range r.Exits {
	// 	exits = append(exits, direction)
	// }
}

func (r *Room) AddCharacter(c *Character) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Adding character to room",
		slog.String("room_id", r.ID),
		slog.String("character_id", c.ID))

	r.Characters = append(r.Characters, c)
}

func (r *Room) RemoveCharacter(c *Character) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Removing character from room",
		slog.String("room_id", r.ID),
		slog.String("character_id", c.ID))

	for i, char := range r.Characters {
		if char.ID == c.ID {
			r.Characters = append(r.Characters[:i], r.Characters[i+1:]...)
			break
		}
	}
}

func (r *Room) AddMob(m *Mob) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Adding mob to room",
		slog.String("room_id", r.ID),
		slog.String("mob_id", m.ID))

	r.Mobs = append(r.Mobs, m)
}

func (r *Room) RemoveMob(m *Mob) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Removing mob from room",
		slog.String("room_id", r.ID),
		slog.String("mob_id", m.ID))

	for i, mob := range r.Mobs {
		if mob.ID == m.ID {
			r.Mobs = append(r.Mobs[:i], r.Mobs[i+1:]...)
			break
		}
	}
}

func (r *Room) AddItem(i *Item) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Adding item to room",
		slog.String("room_id", r.ID),
		slog.String("item_id", i.ID))

	r.Items = append(r.Items, i)
}

func (r *Room) RemoveItem(i *Item) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Removing item from room",
		slog.String("room_id", r.ID),
		slog.String("item_id", i.ID))

	for k, item := range r.Items {
		if item.ID == i.ID {
			r.Items = append(r.Items[:k], r.Items[k+1:]...)
			break
		}
	}
}

func (r *Room) Broadcast(msg string, excludeIDs []string) {
	slog.Debug("Broadcasting message to room",
		slog.String("room_id", r.ID),
		slog.String("message", msg),
		slog.Any("exclude_ids", excludeIDs))

	for _, char := range r.Characters {
		slog.Debug("Broadcasting message to character",
			slog.String("character_id", char.ID),
			slog.String("message", msg))

		for _, excludeID := range excludeIDs {
			if char.ID != excludeID {
				char.Send(msg)
			}
		}
	}
}

// // Event functions
// func (r *Room) onRoomCharacterEnter(arguments ...interface{}) {
// 	slog.Debug("Room character enter event",
// 		slog.String("room_id", r.ID),
// 		slog.Any("args", arguments))

// 	arg := arguments[0].(*RoomCharacterEnter)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A character has entered the room", []string{arg.Character.ID})
// }

// func (r *Room) onRoomCharacterLeave(arguments ...interface{}) {
// 	slog.Debug("Room character leave event",
// 		slog.String("room_id", r.ID))

// 	arg := arguments[0].(*RoomCharacterLeave)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A character has left the room", []string{arg.Character.ID})
// }

// func (r *Room) onRoomMobEnter(arguments ...interface{}) {
// 	slog.Debug("Room mob enter event",
// 		slog.String("room_id", r.ID))

// 	arg := arguments[0].(*RoomMobEnter)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A mob has entered the room", []string{arg.Mob.ID})

// }

// func (r *Room) onRoomMobLeave(arguments ...interface{}) {
// 	slog.Debug("Room mob leave event",
// 		slog.String("room_id", r.ID))

// 	arg := arguments[0].(*RoomMobLeave)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A mob has left the room", []string{arg.Mob.ID})
// }
