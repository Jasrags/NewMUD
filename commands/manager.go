package commands

import (
	"io"
	"log/slog"
	"strings"

	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

var (
	Mgr = NewManager()
)

type Manager struct {
	commandList map[string]*Command
}

func NewManager() *Manager {
	return &Manager{
		commandList: make(map[string]*Command),
	}
}

func (mgr *Manager) RegisterCommands() {
	slog.Info("Registering commands")
	for _, command := range registeredCommands {
		slog.Debug("Registering command",
			slog.String("command", command.Name))
		mgr.commandList[command.Name] = &command
		for _, alias := range command.Aliases {
			mgr.commandList[alias] = &command
		}
	}
}

func (mgr *Manager) GetCommands() map[string]*Command {
	return mgr.commandList
}

func (mgr *Manager) ParseAndExecute(s ssh.Session, input string, user *users.User, room *rooms.Room) {
	slog.Debug("Parsing and executing command",
		slog.String("input", input))

	if input == "" {
		return
	}

	parts := strings.Fields(input)
	cmd := parts[0]
	args := parts[1:]

	slog.Debug("Command name",
		slog.String("command_name", cmd),
		slog.Any("args", args))

	if command, ok := mgr.commandList[cmd]; ok {
		command.Func(s, args, nil, nil)
	} else {
		io.WriteString(s, cfmt.Sprintf("{{Unknown command.}}::red\n"))
	}
}
