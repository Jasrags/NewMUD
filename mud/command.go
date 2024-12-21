package mud

import (
	"fmt"
	"io"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
)

// CommandHandler is a function type for handling commands
type CommandHandler func(ctx *GameContext, player *Player, command string, args []string)

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
		// Log:         l,
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

func NewCommandManager(l zerolog.Logger) *CommandManager {
	return &CommandManager{
		Log:      l,
		Commands: make(map[string]*Command),
	}
}

func (cm *CommandManager) Load() {
	cm.Log.Info().Msg("Load commands")

	for _, cmd := range commands {
		cm.RegisterCommand(cmd)
	}

	cm.Log.Info().
		Int("num_commands", len(cm.Commands)).
		Msg("Commands loaded")
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

func (cm *CommandManager) ParseAndExecute(gc *GameContext, input string, player *Player) {
	cm.Log.Debug().
		Str("input", input).
		Str("player_name", player.Name).
		Msg("Parse and execute command")

	if len(input) == 0 {
		cm.Log.Debug().Msg("Empty input")
		return
	}
	parts := strings.Fields(strings.ToLower(strings.TrimRight(input, "\n")))
	commandName := parts[0]
	args := parts[1:]

	cm.Log.Debug().
		Str("command_name", commandName).
		Strs("args", args).
		Msg("Command name")

	if cmd, exists := cm.Commands[commandName]; exists {
		cmd.Execute(gc, player, commandName, args)
	} else {
		io.WriteString(player.Conn, cfmt.Sprintf("{{Unknown command.}}::red\n"))
	}
}

var commands = []*Command{
	NewCommand(
		"say",
		"Say something to the room",
		[]string{"s"},
		func(ctx *GameContext, player *Player, commandName string, args []string) {
			ctx.Log.Debug().Msg("Say command")

			if player.Room == nil {
				io.WriteString(player.Conn, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
				return
			}

			if len(args) == 0 {
				io.WriteString(player.Conn, cfmt.Sprintf("{{You must specify something to say.}}::red\n"))
				return
			}

			message := strings.Join(args, " ")
			for _, p := range player.Room.Players {
				if p != player {
					io.WriteString(p.Conn, cfmt.Sprintf("{{%s says:}}::cyan %s\n", player.Name, message))
				} else {
					io.WriteString(p.Conn, cfmt.Sprintf("{{You say:}}::cyan %s\n", message))
				}
			}
		}),
	// TODO: This needs to display the main command and show it's aliases
	NewCommand(
		"help",
		"List all available commands",
		[]string{"h"},
		func(ctx *GameContext, player *Player, commandName string, args []string) {
			ctx.Log.Debug().Msg("Help command")

			uniqueCommands := make(map[string]*Command)
			for _, cmd := range ctx.CommandManager.Commands {
				uniqueCommands[cmd.Name] = cmd
			}

			var builder strings.Builder
			builder.WriteString(cfmt.Sprintf("{{Available commands:}}::white|bold\n"))
			for _, cmd := range uniqueCommands {
				builder.WriteString(cfmt.Sprintf("{{%s}}::cyan - %s (aliases: %s)\n", cmd.Name, cmd.Description, strings.Join(cmd.Aliases, ", ")))
			}
			io.WriteString(player.Conn, builder.String())
		}),
	NewCommand(
		"look",
		"Look around the room",
		[]string{"l"},
		func(ctx *GameContext, player *Player, commandName string, args []string) {
			ctx.Log.Debug().Msg("Look command")

			if player.Room == nil {
				io.WriteString(player.Conn, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
				return
			}

			io.WriteString(player.Conn, RenderRoom(player, player.Room))
		}),
	NewCommand(
		"move",
		"Move to another room",
		[]string{"m",
			"n", "north",
			"s", "south",
			"e", "east",
			"w", "west",
			"u", "up",
			"d", "down",
		},
		func(ctx *GameContext, player *Player, commandName string, args []string) {
			ctx.Log.Debug().
				Str("player_name", player.Name).
				Str("command_name", commandName).
				Strs("args", args).
				Msg("Move command")

				// Check if the player is in a room
			if player.Room == nil {
				io.WriteString(player.Conn, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
				return
			}

			// Check if the player specified a direction with the move command
			if commandName == "move" && len(args) == 0 {
				io.WriteString(player.Conn, cfmt.Sprintf("{{You must specify a direction.}}::red\n"))
				return
			}

			// Check if the player specified a direction with the move command or used a direction alias
			var dir string
			switch commandName {
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
				dir = args[0]

			}

			// Check if the exit exists
			if exit, ok := player.Room.Exits[dir]; ok {
				player.MoveTo(exit.Room)
				io.WriteString(player.Conn, cfmt.Sprintf("You move %s.\n\n", dir))
			} else {
				io.WriteString(player.Conn, cfmt.Sprintf("{{You can't go that way.}}::red\n"))
			}

			io.WriteString(player.Conn, RenderRoom(player, player.Room))
		}),
	NewCommand(
		"quit",
		"Quit the game",
		[]string{"q"},
		func(ctx *GameContext, player *Player, commandName string, args []string) {
			ctx.Log.Debug().Msg("Quit command")

			io.WriteString(player.Conn, cfmt.Sprintf("Goodbye!\n"))
			fmt.Println("Player disconnected:", player.Conn.RemoteAddr())
			player.Conn.Close()
		}),
}
