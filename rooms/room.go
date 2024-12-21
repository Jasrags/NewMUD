package rooms

import (
	"log/slog"

	"github.com/Jasrags/NewMUD/exits"
)

type Room struct {
	ID          string                `yaml:"id"`
	AreaID      string                `yaml:"area_id"`
	Title       string                `yaml:"title"`
	Description string                `yaml:"description"`
	Exits       map[string]exits.Exit `yaml:"exits"`
}

func NewRoom() *Room {
	return &Room{
		Exits: make(map[string]exits.Exit),
	}
}

func MoveToRoom(userID, roomID string) {
	slog.Debug("Moving user to room",
		slog.String("user_id", userID),
		slog.String("room_id", roomID))

	// get user

	// get room

	// If the user is already in the room, return

	// remove user from current room/area

	// add user to new room/area

}
