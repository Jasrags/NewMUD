package commands

import (
	"log/slog"
	"net"

	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
)

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
	commandList = make(map[string]*Command)
)

type Command struct {
	Name        string
	Description string
	Aliases     []string
	IsAdmin     bool
	Func        CommandFunc
}

type CommandFunc func(args []string, user *users.User, room *rooms.Room)

func RegisterCommands() {
	slog.Info("Registering commands")
	for _, command := range registeredCommands {
		slog.Debug("Registering command",
			slog.String("command", command.Name))
		commandList[command.Name] = &command
		for _, alias := range command.Aliases {
			commandList[alias] = &command
		}
	}
}

func GetCommands() map[string]*Command {
	return commandList
}

func ParseAndExecute(conn net.Conn, cmd string, args []string) {
	slog.Debug("Parsing and executing command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if cmd == "" {
		return
	}

	if command, ok := commandList[cmd]; ok {
		command.Func(args, nil, nil)
	} else {
		// io.WriteString(player.Conn, cfmt.Sprintf("{{Unknown command.}}::red\n"))
	}

}
