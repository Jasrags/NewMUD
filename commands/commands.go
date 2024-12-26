package commands

import (
	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/gliderlabs/ssh"
)

// TODO: need a manager for this as well

var (
	registeredCommands = []Command{
		{
			Name:        "look",
			Description: "Look around the room",
			Aliases:     []string{"l"},
			Func:        Look,
		},
		{
			Name:        "help",
			Description: "List available commands",
			Aliases:     []string{"h"},
			Func:        Help,
		},
	}
)

type Command struct {
	Name        string
	Description string
	Aliases     []string
	IsAdmin     bool
	Func        CommandFunc
}

type CommandFunc func(s ssh.Session, args []string, user *users.User, room *rooms.Room)
