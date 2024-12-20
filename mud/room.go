package mud

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
)

type Direction string

const (
	DirectionNorth Direction = "north"
	DirectionSouth Direction = "south"
	DirectionEast  Direction = "east"
	DirectionWest  Direction = "west"
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
)

// RoomData represents file data for a room
type RoomData struct {
	ID          string       `yaml:"id"`
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Coordinates *Coordinates `yaml:"coordinates"`
	Exits       []ExitData   `yaml:"exits"`
}

// ExitData represents file data for an exit
type ExitData struct {
	RoomID    string `yaml:"room_id"`
	Direction string `yaml:"direction"`
	// Inferred  bool   `yaml:"-"`
}

// Room represents a game room
type Room struct {
	Log         zerolog.Logger
	ID          string
	Title       string
	Description string
	Coordinates *Coordinates
	Players     map[string]*Player
	Exits       map[string]*Exit
	Area        *Area
	AreaID      string
	// TODO: Add doors
}

// Coordinates represents 3D coordinates of the room in a map
type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

// Exit represents a room exit
type Exit struct {
	Room      *Room
	Direction string
	// Inferred  bool
}

func NewRoom() *Room {
	r := &Room{
		Log:     NewDevLogger(),
		Players: make(map[string]*Player),
		Exits:   make(map[string]*Exit),
	}

	return r
}

// func (r *Room) onPlayerEnter(player *Player) {
// 	r.Log.Debug().
// 		Str("player_name", player.Name).
// 		Str("room_id", r.ID).
// 		Msg("Player has entered the room")
// }

func (r *Room) AddPlayer(player *Player) {
	r.Players[player.Name] = player
}

func (r *Room) RemovePlayer(player *Player) {
	delete(r.Players, player.Name)
}

func (r *Room) Broadcast(message string, exclude *Player) {
	for _, player := range r.Players {
		if exclude != nil && player.Name == exclude.Name {
			continue
		}
		io.WriteString(player.Conn, message)
	}
}

func (r *Room) Emit(event string, data ...interface{}) {
	eventName := fmt.Sprintf("Room#event:%s", event)
	r.Log.Debug().
		Str("room_id", r.ID).
		Str("event_name", eventName).
		Msg("Emit event")
}
