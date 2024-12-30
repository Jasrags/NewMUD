package main

import (
	"io"
	"log/slog"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

var (
	CommandMgr = NewCommandManager()
)

type CommandManager struct {
	commands map[string]*Command
}

func NewCommandManager() *CommandManager {
	return &CommandManager{
		commands: make(map[string]*Command),
	}
}

func (mgr *CommandManager) RegisterCommands() {
	slog.Info("Registering commands")
	for _, command := range registeredCommands {
		slog.Debug("Registering command",
			slog.String("command", command.Name))
		mgr.commands[command.Name] = &command
		for _, alias := range command.Aliases {
			mgr.commands[alias] = &command
		}
	}
}

func (mgr *CommandManager) GetCommands() map[string]*Command {
	return mgr.commands
}

func (mgr *CommandManager) ParseAndExecute(s ssh.Session, input string, user *User, room *Room) {
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

	if command, ok := mgr.commands[cmd]; ok {
		command.Func(s, cmd, args, user, user.ActiveCharacter.Room)
	} else {
		io.WriteString(s, cfmt.Sprintf("{{Unknown command.}}::red\n"))
	}
}
