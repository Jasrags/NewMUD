package mud

import (
	"net"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
)

type Player struct {
	Log    zerolog.Logger `json:"-"`
	Name   string         `json:"name"`
	Conn   net.Conn       `json:"-"`
	RoomID string         `json:"room_id"`
	Room   *Room          `json:"-"`
}

func NewPlayer(conn net.Conn) *Player {
	return &Player{
		Log:  NewDevLogger(),
		Conn: conn,
	}
}

func (p *Player) SetRoom(room *Room) {
	p.Log.Debug().
		Str("player_name", p.Name).
		Str("room_id", room.ID).
		Msg("Set player room")

	p.Room = room
	p.RoomID = room.ID
}

// MoveTo will move the player to the next room and broadcast the player's arrival and departure to the rooms
func (p *Player) MoveTo(nextRoom *Room) {
	p.Log.Debug().
		Str("player_name", p.Name).
		Str("current_room_id", p.Room.ID).
		Str("next_room_id", nextRoom.ID).
		Msg("Move player to room")

	prevRoom := p.Room
	if p.Room != nil && p.Room.ID != nextRoom.ID {
		p.Room.RemovePlayer(p)
	}

	prevRoom.Broadcast(cfmt.Sprintf("\n{{%s}}::green|bold {{has left the room}}::white\n", p.Name), p)

	p.SetRoom(nextRoom)
	nextRoom.AddPlayer(p)

	nextRoom.Broadcast(cfmt.Sprintf("\n{{%s}}::green|bold {{has entered the room}}::white\n", p.Name), p)
}
