package main

import (
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
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

func (mgr *CommandManager) RegisterCommand(command Command) {
	slog.Debug("Registering command",
		slog.String("command", command.Name),
		slog.Any("aliases", command.Aliases))
	mgr.commands[command.Name] = &command
	for _, alias := range command.Aliases {
		mgr.commands[alias] = &command
	}
}

func (mgr *CommandManager) GetCommands() map[string]*Command {
	return mgr.commands
}

func (mgr *CommandManager) ParseAndExecute(s ssh.Session, input string, user *Account, char *Character, room *Room) {

	if char != nil && input != "" {
		if strings.HasPrefix(input, "!") {
			historyIndex, err := strconv.Atoi(input[1:])
			if err == nil && historyIndex > 0 && historyIndex <= len(char.CommandHistory) {
				input = char.CommandHistory[historyIndex-1]
				io.WriteString(s, cfmt.Sprintf("{{Re-executing: %s}}::green\n", input))
			} else {
				io.WriteString(s, "{{Invalid history index.}}::red\n")
				return
			}
		}

		char.CommandHistory = append(char.CommandHistory, input)
		maxHistorySize := viper.GetInt("server.max_history_size")
		if len(char.CommandHistory) > maxHistorySize {
			char.CommandHistory = char.CommandHistory[1:] // Remove the oldest entry
		}
	}

	cmd, args := ParseArguments(input)
	if cmd == "" {
		return
	}

	command, ok := mgr.commands[cmd]
	if !ok {
		io.WriteString(s, cfmt.Sprintf("{{Unknown command '%s'. Type 'help' for a list of commands.}}::red\n", cmd))
		return
	}

	if !mgr.CanRunCommand(char, command) {
		io.WriteString(s, cfmt.Sprintf("{{Unknown command '%s'. Type 'help' for a list of commands.}}::red\n", cmd))
		return
	}

	if command.SuggestFunc != nil {
		suggestions := command.SuggestFunc(input, args, char, room)
		if len(suggestions) > 0 {
			io.WriteString(s, cfmt.Sprintf("{{Suggestions:}}::green %s\n", strings.Join(suggestions, ", ")))
		}
	}

	command.Func(s, cmd, args, user, char, room)
}

func (mgr *CommandManager) CanRunCommand(char *Character, cmd *Command) bool {
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

func ParseArguments(input string) (string, []string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	return parts[0], parts[1:]
}
