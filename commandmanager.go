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

func (mgr *CommandManager) ParseAndExecute(s ssh.Session, input string, user *User, char *Character, room *Room) {
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
		io.WriteString(s, cfmt.Sprintf("{{You don't have permission to run '%s'.}}::red\n", cmd))
		return
	}

	if command.Suggest != nil {
		suggestions := command.Suggest(input, args, char, room)
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

func SuggestGive(line string, args []string, char *Character, room *Room) []string {
	switch len(args) {
	case 1: // Suggest items for the first argument
		var names []string
		for _, item := range char.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			names = append(names, bp.Name)
		}
		return names
	case 2: // Suggest "to" as the second token
		return []string{"to"}
	case 3: // Suggest characters in the room for the third token
		var names []string
		for _, character := range room.Characters {
			names = append(names, character.Name)
		}
		return names
	default:
		return nil
	}
}
