package mud

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Room represents a room in the game.git s
type Room struct {
	Log         zerolog.Logger
	ID          string
	Title       string
	Description string
	Exits       map[string]*Room
	Players     map[string]*Player
}

func NewRoom(id, title, description string) *Room {
	return &Room{
		Log:         NewDevLogger(),
		ID:          id,
		Title:       title,
		Description: description,
		Exits:       make(map[string]*Room),
		Players:     make(map[string]*Player),
	}
}

func (r *Room) AddPlayer(player *Player) {
	r.Log.Debug().
		Str("player_name", player.Name).
		Str("room_id", r.ID).
		Msg("Add player to room")

	r.Players[player.Name] = player
}

func (r *Room) RemovePlayer(player *Player) {
	r.Log.Debug().
		Str("player_name", player.Name).
		Str("room_id", r.ID).
		Msg("Remove player from room")

	delete(r.Players, player.Name)
}

func setupWorld() *Room {
	room1 := NewRoom("room1", "Small Room", "You are in a small, cozy room. Exits lead north and east.")
	room2 := NewRoom("room2", "Bright Room", "You are in a bright, sunlit room. Exits lead south.")
	room3 := NewRoom("room3", "Dark Room", "You are in a dark, eerie room. Exits lead west.")

	// Connect rooms
	room1.Exits["north"] = room2
	room1.Exits["east"] = room3
	room2.Exits["south"] = room1
	room3.Exits["west"] = room1

	room1.Listen()
	room2.Listen()
	room3.Listen()

	return room1
}

func (r *Room) Broadcast(message string, exclude *Player) {
	for _, player := range r.Players {
		if player != exclude {
			player.Out <- message
		}
	}
}

func (r *Room) Listen() {
	r.Log.Debug().
		Str("room_id", r.ID).
		Msg("Listening for player events")

	// Subscribe to player entrance
	eventBus.Subscribe(EventPlayerEnter, func(player *Player, roomID string) {
		if roomID == r.ID {
			r.Broadcast(fmt.Sprintf("%s enters the room.\r\n", player.Name), player)
			r.AddPlayer(player)
		}
	})

	// Subscribe to player exit
	eventBus.Subscribe(EventPlayerExit, func(player *Player, roomID string) {
		if roomID == r.ID {
			r.Broadcast(fmt.Sprintf("%s leaves the room.\r\n", player.Name), player)
			r.RemovePlayer(player)
		}
	})
}
