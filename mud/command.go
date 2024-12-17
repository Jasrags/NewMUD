package mud

import (
	"fmt"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
)

func RenderRoom(room *Room) string {
	var sb strings.Builder
	sb.WriteString(cfmt.Sprintf("{{%s}}::green|bold\r\n", room.Title))
	sb.WriteString(cfmt.Sprintf("{{%s}}::white|bold\r\n", room.Description))
	sb.WriteString(cfmt.Sprint("{{Exits:}}::yellow|bold\r\n"))

	if len(room.Exits) == 0 {
		sb.WriteString(" {{- None.}}::yellow|bold\r\n")
	} else {
		for direction, _ := range room.Exits {
			sb.WriteString(cfmt.Sprintf(" - {{%s}}::yellow\r\n", direction))
			// player.Out <- cfmt.Sprintf("{{%s}}::white\r\n", room.Title)
		}
	}
	sb.WriteString(cfmt.Sprint("\r\n"))

	return sb.String()
}

var (
	lookCommand = &Command{
		Name:        "look",
		Description: "Look around the room or at something specific.",
		Execute: func(player *Player, args string) {
			if args == "" {
				player.Out <- RenderRoom(player.Room)
			} else {
				player.Out <- cfmt.Sprintf("You don't see any '%s' here.\r\n", args)
			}
		},
	}
	moveCommand = &Command{
		Name:        "move",
		Description: "Move to another room.",
		Execute: func(player *Player, args string) {
			direction := strings.ToLower(args)
			nextRoom, exists := player.Room.Exits[direction]
			if !exists {
				player.Out <- "You can't go that way.\r\n"
				return
			}

			// Emit exit event
			eventBus.Publish(EventPlayerExit, player, player.Room.ID)

			// Update player's room
			player.Room = nextRoom

			// Emit entrance event
			eventBus.Publish(EventPlayerEnter, player, nextRoom.ID)

			// Send room description to the player
			player.Out <- fmt.Sprintf("You move %s.\n%s\r\n", direction, nextRoom.Description)
		},
	}
)

// Command represents a game command.
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Execute     func(player *Player, args string)
}

type CommandParser struct {
	Log      zerolog.Logger
	commands map[string]*Command
}

func NewCommandParser() *CommandParser {
	return &CommandParser{
		Log:      NewDevLogger(),
		commands: make(map[string]*Command),
	}
}

func (cp *CommandParser) RegisterCommand(cmd *Command) {
	cp.Log.Debug().
		Str("command", cmd.Name).
		Strs("aliases", cmd.Aliases).
		Msg("Register command")

	cp.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		cp.commands[alias] = cmd
	}
}

func (cp *CommandParser) ParseAndExecute(input string, player *Player) {
	input = strings.TrimSpace(input)
	cp.Log.Debug().
		Str("input", input).
		Str("player_name", player.Name).
		Msg("Parse and execute command")

	parts := strings.SplitN(input, " ", 2)
	commandName := strings.ToLower(parts[0])
	args := ""
	cp.Log.Debug().
		Str("command_name", commandName).
		Str("args", args).
		Msg("Command name")
	if len(parts) > 1 {
		args = parts[1]
	}

	if cmd, exists := cp.commands[commandName]; exists {
		cmd.Execute(player, args)
		cp.Log.Debug().Str("command_name", commandName).Msg("Command executed")
	} else {
		player.Out <- "Unknown command."
		cp.Log.Debug().Str("command_name", commandName).Msg("Unknown command")
	}
}
