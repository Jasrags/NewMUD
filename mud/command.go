package mud

import (
	"fmt"
	"io"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
)

// CommandHandler is a function type for handling commands
type CommandHandler func(player *Player, args []string)

// CommandMap maps command names to their handlers
var CommandMap = map[string]CommandHandler{}

type Command struct {
	Log         zerolog.Logger
	Name        string
	Description string
	Aliases     []string
	Execute     CommandHandler
}

func NewCommand(name, description string, aliases []string, execute CommandHandler) *Command {
	return &Command{
		Log:         NewDevLogger(),
		Name:        name,
		Description: description,
		Aliases:     aliases,
		Execute:     execute,
	}
}

type CommandManager struct {
	Log      zerolog.Logger
	Commands map[string]*Command
}

func NewCommandManager() *CommandManager {
	return &CommandManager{
		Log:      NewDevLogger(),
		Commands: make(map[string]*Command),
	}
}

func (cm *CommandManager) RegisterCommand(cmd *Command) {
	cm.Log.Debug().
		Str("command", cmd.Name).
		Strs("aliases", cmd.Aliases).
		Msg("Register command")

	cm.Commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		cm.Commands[alias] = cmd
	}
}

func (cm *CommandManager) ParseAndExecute(input string, player *Player) {
	cm.Log.Debug().
		Str("input", input).
		Str("player_name", player.Name).
		Msg("Parse and execute command")

	if len(input) == 0 {
		cm.Log.Debug().Msg("Empty input")
		return
	}
	input = strings.ToLower(input)
	input = strings.TrimRight(input, "\n")
	parts := strings.Fields(input)
	commandName := parts[0]
	args := parts[1:]

	cm.Log.Debug().
		Strs("parts", parts).
		Str("command_name", commandName).
		Strs("args", args).
		Msg("Command name")

	if cmd, exists := cm.Commands[commandName]; exists {
		cmd.Execute(player, args)
	} else {
		io.WriteString(player.Conn, cfmt.Sprintf("{{Unknown command.}}::red\n"))
	}
}

var commands = []*Command{
	NewCommand("look", "Look around the room", []string{"l"}, func(player *Player, args []string) {
		room := Rooms[player.RoomID]
		io.WriteString(player.Conn, cfmt.Sprintf("{{%s}}::green|bold\n", room.Title))
		io.WriteString(player.Conn, cfmt.Sprintf("{{%s}}::white\n", room.Description))

		if len(room.Exits) == 0 {
			io.WriteString(player.Conn, cfmt.Sprint("{{There are no exits.}}::red\n"))
		} else {
			io.WriteString(player.Conn, cfmt.Sprint("{{Exits:}}::yellow|bold\n"))
			for direction, _ := range room.Exits {
				io.WriteString(player.Conn, cfmt.Sprintf("{{ - %s}}::yellow\n", direction))
			}
		}
	}),
	NewCommand("move", "Move to another room", []string{"m"}, func(player *Player, args []string) {
		if len(args) == 0 {
			io.WriteString(player.Conn, cfmt.Sprintf("{{You must specify a direction.}}::red\n"))
			return
		}
		dir := args[0]

		room := Rooms[player.RoomID]
		if nextRoomID, exists := room.Exits[dir]; exists {
			player.RoomID = nextRoomID
			nextRoom := Rooms[nextRoomID]
			io.WriteString(player.Conn, cfmt.Sprintf("You move %s.\n%s\n", dir, nextRoom.Description))
		} else {
			io.WriteString(player.Conn, cfmt.Sprintf("{{You can't go that way.}}::red\n"))
		}
	}),
	NewCommand("quit", "Quit the game", []string{"q"}, func(player *Player, args []string) {
		io.WriteString(player.Conn, cfmt.Sprintf("Goodbye!\n"))
		fmt.Println("Player disconnected:", player.Conn.RemoteAddr())
		player.Conn.Close()
	}),
}

func registerCommands(cm *CommandManager) {
	for _, cmd := range commands {
		cm.RegisterCommand(cmd)
	}
}
