package main

import (
	"io"
	"log/slog"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

// TODO: need a manager for this as well

var (
	registeredCommands = []Command{
		{
			Name:        "look",
			Description: "Look around the room",
			Usage: []string{
				"look [item|character|mob|direction]",
			},
			Aliases: []string{"l"},
			Func:    Look,
		},
		{
			Name:        "get",
			Description: "Get an item",
			Usage: []string{
				"get all",
				"get <item>",
				"get <number> <items>",
				"get all <items>",
			},
			Aliases: []string{"g"},
			Func:    Get,
		},
		{
			Name:        "give",
			Description: "Give an item",
			Usage: []string{
				"give <item> [to] <character>",
				"give 2 <items> [to] <character>",
				"give all [to] <character>",
			},
			Aliases: []string{"gi"},
			Func:    Give,
		},
		{
			Name:        "drop",
			Description: "Drop an item",
			Usage: []string{
				"drop all",
				"drop <item>",
				"drop <number> <items>",
				"drop all <items>",
			},
			Aliases: []string{"d"},
			Func:    Drop,
		},
		{
			Name:        "help",
			Description: "List available commands",
			Usage: []string{
				"help",
				"help <command>",
			},
			Aliases: []string{"h"},
			Func:    Help,
		},
		{
			Name:        "move",
			Description: "Move to a different room",
			Usage:       []string{"move [direction]"},
			Aliases:     []string{"m", "n", "s", "e", "w", "u", "d", "north", "south", "east", "west", "up", "down"},
			Func:        Move,
		},
	}
)

type Command struct {
	Name        string
	Description string
	Usage       []string
	Aliases     []string
	IsAdmin     bool
	Func        CommandFunc
}

type CommandFunc func(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room)

/*
Usage:
  - help
  - help <command>
*/
func Help(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Help command",
		slog.String("command", cmd),
		slog.Any("args", args))

	uniqueCommands := make(map[string]*Command)
	for _, cmd := range CommandMgr.GetCommands() {
		uniqueCommands[cmd.Name] = cmd
	}

	var builder strings.Builder
	switch len(args) {
	case 0:
		builder.WriteString(cfmt.Sprintf("{{Available commands:}}::white|bold\n"))
		for _, cmd := range uniqueCommands {
			builder.WriteString(cfmt.Sprintf("{{%s}}::cyan - %s (aliases: %s)\n", cmd.Name, cmd.Description, strings.Join(cmd.Aliases, ", ")))
		}
	case 1:
		if command, ok := uniqueCommands[args[0]]; ok {
			builder.WriteString(cfmt.Sprintf("{{%s}}::cyan\n", strings.ToUpper(command.Name)))
			builder.WriteString(cfmt.Sprintf("{{Description:}}::white|bold %s\n", command.Description))
			builder.WriteString(cfmt.Sprintf("{{Aliases:}}::white|bold %s\n", strings.Join(command.Aliases, ", ")))
			builder.WriteString(cfmt.Sprintf("{{Usage:}}::white|bold\n"))
			for _, usage := range command.Usage {
				builder.WriteString(cfmt.Sprintf("{{  - %s}}::green\n", usage))
			}
		} else {
			builder.WriteString(cfmt.Sprintf("{{Unknown command.}}::red\n"))
		}
	}

	io.WriteString(s, builder.String())
}

/*
Usage:
  - drop all
  - drop <item>
  - drop <number> <items>
  - drop all <items>
*/
func Drop(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Drop command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Drop what?}}::red\n"))
		return
	}

	arg1 := args[0]

	switch arg1 {
	// case "all":
	// for _, item := range char.Items {
	// char.RemoveItem(item)
	// char.Room.AddItem(item)
	// io.WriteString(s, cfmt.Sprintf("{{You drop %s.}}::green\n", item.Spec.Name))
	// char.Room.Broadcast(cfmt.Sprintf("{{%s drops %s.}}::green\n", char.Name, item.Spec.Name), []string{char.ID})
	// }
	default:
		io.WriteString(s, cfmt.Sprintf("{{You can't drop that.}}::red\n"))
	}
}

/*
Usage:
  - give <item> [to] <character>
  - give 2 <items> [to] <character>
  - give all [to] <character>
*/
func Give(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Give command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	switch len(args) {
	case 0:
		io.WriteString(s, cfmt.Sprintf("{{Give what?}}::red\n"))
		return
		// case 1:
		// 	io.WriteString(s, cfmt.Sprintf("{{Give to who?}}::red\n"))
		// 	return
		// case 2:
		// 	what := args[0]
		// 	if what == "all" {
		// 		for _, item := range char.Items {
		// 			char.RemoveItem(item)
		// 			char.Room.AddItem(item)
		// 			io.WriteString(s, cfmt.Sprintf("{{You give %s.}}::green\n", item.Spec.Name))
		// 			char.Room.Broadcast(cfmt.Sprintf("{{%s gives %s.}}::green\n", char.Name, item.Spec.Name), []string{char.ID})
		// 		}
		// 		return
	}

	// 	to := args[1]
	// 	for _, c := range char.Room.Characters {
	// 		if strings.EqualFold(c.Name, to) {
	// 			for _, item := range char.Items {
	// 				char.RemoveItem(item)
	// 				c.AddItem(item)
	// 				io.WriteString(s, cfmt.Sprintf("{{You give %s.}}::green\n", item.Spec.Name))
	// 				char.Room.Broadcast(cfmt.Sprintf("{{%s gives %s.}}::green\n", char.Name, item.Spec.Name), []string{char.ID})
	// 			}
	// 			return
	// 		}
	// 	}
	// }

	// arg1 := args[0]
	// arg2 := args[1]
	// if args[1] == nil || args[1] == "" {
	//     io.WriteString(s, cfmt.Sprintf("{{Give to who?}}::red\n"))
	//     return
	// }

	// switch arg1 {
	// case "all":
	// 	for _, item := range char.Items {
	// 		char.RemoveItem(item)
	// 		char.Room.AddItem(item)
	// 		io.WriteString(s, cfmt.Sprintf("{{You give %s.}}::green\n", item.Spec.Name))
	// 		char.Room.Broadcast(cfmt.Sprintf("{{%s gives %s.}}::green\n", char.Name, item.Spec.Name), []string{char.ID})
	// 	}
	// default:
	io.WriteString(s, cfmt.Sprintf("{{You can't give that.}}::red\n"))
	// }
}

/*
Usage:
  - get all
  - get <item>
  - get <number> <items>
  - get all <items>
*/
func Get(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Get command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Get what?}}::red\n"))
		return
	}

	arg1 := args[0]

	switch arg1 {
	// case "all":
	// 	for _, item := range char.Room.Inventory.Items {
	// 		char.Room.RemoveItem(item)
	// 		char.AddItem(item)
	// 		io.WriteString(s, cfmt.Sprintf("{{You get %s.}}::green\n", item.Spec.Name))
	// 		char.Room.Broadcast(cfmt.Sprintf("{{%s gets %s.}}::green\n", char.Name, item.Spec.Name), []string{char.ID})
	// 	}
	default:
		io.WriteString(s, cfmt.Sprintf("{{You can't get that.}}::red\n"))
	}

}

func Look(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Look command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		slog.Error("Character is not in a room",
			slog.String("character_id", char.ID))

		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))

		return
	}

	// if no arguments are passed, show the room
	if len(args) == 0 {
		io.WriteString(s, RenderRoom(user, char, nil))
	} else {
		io.WriteString(s, cfmt.Sprintf("{{Look at what?}}::red\n"))
	}

	// TODO: Support looking at other things, like items, characters, mobs
}

func Move(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Move command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if cmd == "move" && len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Move where?}}::red\n"))
		return
	}

	// Check if the player specified a direction with the move command or used a direction alias
	var dir string
	switch cmd {
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
		slog.Error("Invalid direction",
			slog.String("direction", dir))
	}

	// Check if the exit exists
	if exit, ok := char.Room.Exits[dir]; ok {
		// prevRoom := char.Room
		// char.FromRoom()
		// char.ToRoom(exit.Room)

		// act("$n has arrived.", TRUE, ch, 0,0, TO_ROOM);
		// do_look(ch, "\0",15);
		char.MoveToRoom(exit.Room)
		char.Save()

		io.WriteString(s, cfmt.Sprintf("You move %s.\n\n", dir))
		io.WriteString(s, RenderRoom(user, char, nil))
	} else {
		io.WriteString(s, cfmt.Sprintf("{{You can't go that way.}}::red\n"))
		return
	}
}
