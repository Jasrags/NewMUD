package game

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
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

// func (mgr *CommandManager) ParseAndExecute(s ssh.Session, input string, user *Account, char *Character, room *Room) {

// 	if char != nil && input != "" {
// 		if strings.HasPrefix(input, "!") {
// 			historyIndex, err := strconv.Atoi(input[1:])
// 			if err == nil && historyIndex > 0 && historyIndex <= len(char.CommandHistory) {
// 				input = char.CommandHistory[historyIndex-1]
// 				WriteStringF(s, "{{Re-executing: %s}}::green"+CRLF, input)
// 			} else {
// 				WriteString(s, "{{Invalid history index.}}::red"+CRLF)
// 				return
// 			}
// 		}

// 		char.CommandHistory = append(char.CommandHistory, input)
// 		maxHistorySize := viper.GetInt("server.max_history_size")
// 		if len(char.CommandHistory) > maxHistorySize {
// 			char.CommandHistory = char.CommandHistory[1:] // Remove the oldest entry
// 		}
// 	}

// 	cmd, args := ParseArguments(input)
// 	if cmd == "" {
// 		return
// 	}

// 	command, ok := mgr.commands[cmd]
// 	if !ok {
// 		WriteStringF(s, "{{Unknown command '%s'. Type 'help' for a list of commands.}}::red"+CRLF, cmd)
// 		return
// 	}

// 	if !mgr.CanRunCommand(char, command) {
// 		WriteStringF(s, "{{Unknown command '%s'. Type 'help' for a list of commands.}}::red"+CRLF, cmd)
// 		return
// 	}

// 	if command.SuggestFunc != nil {
// 		suggestions := command.SuggestFunc(input, args, char, room)
// 		if len(suggestions) > 0 {
// 			WriteStringF(s, "{{Suggestions:}}::green %s"+CRLF, strings.Join(suggestions, ", "))
// 		}
// 	}

// 	command.Func(s, cmd, args, user, char, room)
// }

func (mgr *CommandManager) ParseAndExecute(s ssh.Session, input string, user *Account, char *Character, room *Room) {
	// Global check: Ensure a character is associated with the session before proceeding.
	if user == nil {
		WriteString(s, "{{Error: No user is associated with this session.}}::red"+CRLF)
		slog.Error("No user is associated with this session",
			slog.String("session_id", s.Context().SessionID()))
		return
	}
	if char == nil {
		WriteString(s, "{{Error: No character is associated with this session.}}::red"+CRLF)
		slog.Error("No character is associated with this session",
			slog.String("session_id", s.Context().SessionID()))
		return
	}
	if room == nil {
		WriteString(s, "{{Error: No room is associated with this session}}::red"+CRLF)
		slog.Error("Error: No room is associated with this session",
			slog.String("session_id", s.Context().SessionID()))
		return
	}

	// Handle command history recall (e.g., "!3" to repeat the 3rd command)
	if input != "" && strings.HasPrefix(input, "!") {
		historyIndex, err := strconv.Atoi(input[1:])
		if err == nil && historyIndex > 0 && historyIndex <= len(char.CommandHistory) {
			input = char.CommandHistory[historyIndex-1]
			WriteStringF(s, "{{Re-executing: %s}}::green"+CRLF, input)
		} else {
			WriteString(s, "{{Invalid history index.}}::red"+CRLF)
			return
		}
	}

	// Store the command in the character's command history
	char.CommandHistory = append(char.CommandHistory, input)
	maxHistorySize := viper.GetInt("server.max_history_size")
	if len(char.CommandHistory) > maxHistorySize {
		char.CommandHistory = char.CommandHistory[1:] // Remove the oldest entry
	}

	// Parse command and arguments
	cmd, args := ParseArguments(input)
	if cmd == "" {
		return
	}

	// Look up the command
	command, ok := mgr.commands[cmd]
	if !ok {
		WriteStringF(s, "{{Unknown command '%s'. Type 'help' for a list of commands.}}::red"+CRLF, cmd)
		return
	}

	// Check if the character is allowed to run this command
	if !mgr.CanRunCommand(char, command) {
		WriteStringF(s, "{{Unknown command '%s'. Type 'help' for a list of commands.}}::red"+CRLF, cmd)
		return
	}

	// Provide command suggestions if applicable
	if command.SuggestFunc != nil {
		suggestions := command.SuggestFunc(input, args, char, room)
		if len(suggestions) > 0 {
			WriteStringF(s, "{{Suggestions:}}::green %s"+CRLF, strings.Join(suggestions, ", "))
		}
	}

	// Execute the command
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
