package mud

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// "to the North",
// "to the East",
// "to the South",
// "to the West",
// "up from here",
// "down from here",

type Player struct {
	Log    zerolog.Logger `json:"-"`
	Name   string         `json:"name"`
	Conn   net.Conn       `json:"-"`
	RoomID string         `json:"room_id"`
	Room   *Room          `json:"-"`
	Role   string         `json:"role"`
}

func NewPlayer(l zerolog.Logger, conn net.Conn) *Player {
	return &Player{
		Log:  l,
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
// Fires:
// Room#event:playerLeave
// Room#event:playerEnter
// Player#event:enterRoom
func (p *Player) MoveTo(nextRoom *Room) {
	p.Log.Debug().
		Str("player_name", p.Name).
		Str("next_room_id", nextRoom.ID).
		Msg("Move player to room")

	prevRoom := p.Room
	if p.Room != nil && p.Room.ID != nextRoom.ID {
		p.Room.Emit("playerLeave", p, nextRoom)
		p.Room.RemovePlayer(p)

		prevRoom.Broadcast(cfmt.Sprintf("\n{{%s}}::green|bold {{has left the room}}::white\n", p.Name), p)
	}

	p.SetRoom(nextRoom)
	nextRoom.AddPlayer(p)

	nextRoom.Emit("playerEnter", p, prevRoom)
	p.Emit("enterRoom", nextRoom)

	nextRoom.Broadcast(cfmt.Sprintf("\n{{%s}}::green|bold {{has entered the room}}::white\n", p.Name), p)
}

func (p *Player) Save() {
	p.Log.Debug().
		Str("player_name", p.Name).
		Msg("Save player")

	dataDir := viper.GetString("data.players_path")
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		p.Log.Error().Err(err).Msg("Failed to create data directory")
		return
	}

	filePath := filepath.Join(dataDir, strings.ToLower(p.Name)+".json")
	file, errCreate := os.Create(filePath)
	if errCreate != nil {
		p.Log.Error().Err(errCreate).Msg("Failed to create player file")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(p); err != nil {
		p.Log.Error().Err(err).Msg("Failed to encode player to JSON")
	}
}

func (p *Player) Emit(event string, data ...interface{}) {
	eventName := fmt.Sprintf("Player#event:%s", event)
	p.Log.Debug().
		Str("player_name", p.Name).
		Str("event_name", eventName).
		Msg("Emit event")
}
