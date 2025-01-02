package main

import (
	"io"
	"log/slog"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
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
		{
			Name:        "move",
			Description: "Move to a different room",
			Aliases:     []string{"m", "n", "s", "e", "w", "u", "d", "north", "south", "east", "west", "up", "down"},
			Func:        Move,
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

type CommandFunc func(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room)

func Help(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Help command")

	uniqueCommands := make(map[string]*Command)
	for _, cmd := range CommandMgr.GetCommands() {
		uniqueCommands[cmd.Name] = cmd
	}

	var builder strings.Builder
	builder.WriteString(cfmt.Sprintf("{{Available commands:}}::white|bold\n"))
	for _, cmd := range uniqueCommands {
		builder.WriteString(cfmt.Sprintf("{{%s}}::cyan - %s (aliases: %s)\n", cmd.Name, cmd.Description, strings.Join(cmd.Aliases, ", ")))
	}

	io.WriteString(s, builder.String())
}

func Look(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Look command")

	// if no arguments are passed, show the room
	if len(args) == 0 {
		io.WriteString(s, RenderRoom(user, char, room))
	} else {
		io.WriteString(s, cfmt.Sprintf("{{Look at what?}}::red\n"))
	}
	// TODO: Support looking at other things, like items, characters, mobs
}

func Move(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Move command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if cmd == "move" && len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Move where?}}::red\n"))
		return
	}

	// Check if the player specified a direction with the move command or used a direction alias
	var dir string
	switch cmd {
	case "n", "north":
		dir = "north"
	case "s", "south":
		dir = "south"
	case "e", "east":
		dir = "east"
	case "w", "west":
		dir = "west"
	case "u", "up":
		dir = "up"
	case "d", "down":
		dir = "down"
	default:
		slog.Error("Invalid direction",
			slog.String("direction", dir))
	}

	// Check if the exit exists
	if exit, ok := char.Room.Exits[dir]; ok {
		char.MoveToRoom(exit.Room)
		char.Save()
		io.WriteString(s, cfmt.Sprintf("You move %s.\n\n", dir))
		io.WriteString(s, RenderRoom(user, char, room))
	} else {
		io.WriteString(s, cfmt.Sprintf("{{You can't go that way.}}::red\n"))
		return
	}
}
