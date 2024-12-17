package mud

import (
	"net"
)

type Player struct {
	Name   string
	Conn   net.Conn
	RoomID string
}

func NewPlayer(name string, conn net.Conn) *Player {
	return &Player{
		Name:   name,
		Conn:   conn,
		RoomID: "room1",
	}
}

// type Player struct {
// 	Log  zerolog.Logger
// 	Name string
// 	Room *Room
// 	Out  chan string // For sending messages to the player
// 	Conn net.Conn
// }

// func NewPlayer(name string, conn net.Conn) *Player {
// 	return &Player{
// 		Log:  NewDevLogger(),
// 		Name: name,
// 		Out:  make(chan string),
// 		Conn: conn,
// 	}
// }

// // TODO: make the enter and exit messages work properly
// func (p *Player) MoveTo(nextRoom *Room) {
// 	p.Log.Debug().
// 		Str("player_name", p.Name).
// 		// Str("current_room_id", p.Room.ID).
// 		Str("next_room_id", nextRoom.ID).
// 		Msg("Move player to room")

// 	// prevRoom := p.Room

// 	if p.Room != nil && p.Room.ID != nextRoom.ID {
// 		p.Room.RemovePlayer(p)
// 	}

// 	// for _, player := range prevRoom.Players {
// 	// 	if player.Name != p.Name {
// 	// 		player.Out <- p.Name + " has left the room.\r\n"
// 	// 	}
// 	// }

// 	p.Room = nextRoom
// 	nextRoom.AddPlayer(p)
// 	p.Out <- "You have left the room.\r\n"
// }
