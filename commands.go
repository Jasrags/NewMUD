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
			Name:        "get",
			Description: "Get an item",
			Aliases:     []string{"g"},
			Func:        Get,
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
	slog.Debug("Help command",
		slog.String("command", cmd),
		slog.Any("args", args))

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

func Get(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Get command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Get what?}}::red\n"))
		return
	}

	arg1 := args[0]

	switch arg1 {
	case "all":
		for _, item := range char.Room.Items {
			char.Room.RemoveItem(item)
			char.AddItem(item)
			io.WriteString(s, cfmt.Sprintf("{{You get %s.}}::green\n", item.Name))
		}
	default:
		io.WriteString(s, cfmt.Sprintf("{{You can't get that.}}::red\n"))
	}

}

func Look(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Look command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	// if no arguments are passed, show the room
	if len(args) == 0 {
		io.WriteString(s, RenderRoom(user, char, nil))
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
		// prevRoom := char.Room
		// char.FromRoom()
		// char.ToRoom(exit.Room)

		// act("$n has arrived.", TRUE, ch, 0,0, TO_ROOM);
		// do_look(ch, "\0",15);
		char.MoveToRoom(exit.Room)
		char.Save()

		io.WriteString(s, cfmt.Sprintf("You move %s.\n\n", dir))
		io.WriteString(s, RenderRoom(user, char, nil))
	} else {
		io.WriteString(s, cfmt.Sprintf("{{You can't go that way.}}::red\n"))
		return
	}
}
