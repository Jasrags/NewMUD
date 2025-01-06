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

func (mgr *CommandManager) ParseAndExecute(s ssh.Session, input string, user *User, char *Character, room *Room) {
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

	if command, ok := mgr.commands[cmd]; ok && CanRunCommand(char, command) {
		command.Func(s, cmd, args, user, char, room)
	} else {
		io.WriteString(s, cfmt.Sprintf("{{Unknown command.}}::red\n"))
	}
}

func CanRunCommand(char *Character, cmd *Command) bool {
	slog.Debug("Checking if character can run command",
		slog.String("character_id", char.ID),
		slog.String("command", cmd.Name))

	if len(cmd.RequiredRoles) == 0 {
		return true
	}

	requiredRoles := make(map[CharacterRole]bool)
	for _, role := range cmd.RequiredRoles {
		requiredRoles[role] = true
	}

	if _, ok := requiredRoles[char.Role]; !ok {
		return false
	}

	return true
}
