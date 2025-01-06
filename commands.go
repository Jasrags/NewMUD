package main

import (
	"io"
	"log/slog"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
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
			Func:    DoLook,
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
			Func:    DoGet,
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
			Func:    DoGive,
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
			Func:    DoDrop,
		},
		{
			Name:        "help",
			Description: "List available commands",
			Usage: []string{
				"help",
				"help <command>",
			},
			Aliases: []string{"h"},
			Func:    DoHelp,
		},
		{
			Name:        "move",
			Description: "Move to a different room",
			Usage:       []string{"move [direction]"},
			Aliases:     []string{"m", "n", "s", "e", "w", "u", "d", "north", "south", "east", "west", "up", "down"},
			Func:        DoMove,
		},
		{
			Name:        "inventory",
			Description: "List your inventory",
			Usage:       []string{"inventory"},
			Aliases:     []string{"i"},
			Func:        DoInventory,
		},
		{
			Name:        "say",
			Description: "Say something to the room or to a character or mob",
			Usage: []string{
				"say <message>",
				"say @<name> <message>",
			},
			Func: DoSay,
		},
		{
			Name:        "spawn",
			Description: "Spawn an item or mob into the room",
			Usage: []string{
				"spawn item <item>",
				"spawn mob <mob>",
			},
			RequiredRoles: []CharacterRole{CharacterRoleAdmin},
			Func:          DoSpawn,
		},
	}
)

type Command struct {
	Name          string
	Description   string
	Usage         []string
	Aliases       []string
	RequiredRoles []CharacterRole
	IsAdmin       bool
	Func          CommandFunc
}

type CommandFunc func(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room)

/*
Usage:
  - say <message>
  - say @<name> <message>
*/
func DoSay(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Say command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Say what?}}::yellow\n"))
		return
	}

	message := strings.Join(args, " ")

	if strings.HasPrefix(message, "@") {
		// Handle targeted messages
		splitMessage := strings.SplitN(message, " ", 2)
		if len(splitMessage) < 2 {
			io.WriteString(s, cfmt.Sprintf("{{Say what to whom?}}::yellow\n"))
			return
		}

		targetName := splitMessage[0][1:] // Remove '@'
		targetedMessage := splitMessage[1]

		// Find the target in the room
		target := room.FindInteractableByName(targetName)
		if target == nil {
			io.WriteString(s, cfmt.Sprintf("{{No one named '%s' is here.}}::red\n", targetName))
			return
		}

		// Notify the speaker
		io.WriteString(s, cfmt.Sprintf("{{You say to %s: '%s'}}::cyan\n", target.GetName(), targetedMessage))

		// Let the target react to the message
		target.ReactToMessage(char, targetedMessage)

		// Broadcast to the room (excluding speaker and target)
		room.Broadcast(cfmt.Sprintf("{{%s says something to %s.}}::green\n", char.Name, target.GetName()), []string{char.ID, target.GetID()})

	} else {
		// General message to the room
		io.WriteString(s, cfmt.Sprintf("{{You say: '%s'}}::cyan\n", message))
		room.Broadcast(cfmt.Sprintf("{{%s says: '%s'}}::green\n", char.Name, message), []string{char.ID})
	}
}

/*
Usage:
  - help
  - help <command>
*/
func DoHelp(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Help command",
		slog.String("command", cmd),
		slog.Any("args", args))

	uniqueCommands := make(map[string]*Command)
	for _, cmd := range CommandMgr.GetCommands() {
		if CanRunCommand(char, cmd) {
			uniqueCommands[cmd.Name] = cmd
		}
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
			if len(command.Aliases) > 0 {
				builder.WriteString(cfmt.Sprintf("{{Aliases:}}::white|bold %s\n", strings.Join(command.Aliases, ", ")))
			}
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
func DoDrop(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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

	if arg1 == "all" {
		if len(args) < 2 {
			// Drop all items in the inventory
			if len(char.Inventory.Items) == 0 {
				io.WriteString(s, cfmt.Sprintf("{{You have nothing to drop.}}::yellow\n"))
				return
			}

			// Use a copy of the items to safely modify the inventory while iterating
			itemsToDrop := make([]*Item, len(char.Inventory.Items))
			copy(itemsToDrop, char.Inventory.Items)

			for _, item := range itemsToDrop {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp == nil {
					io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
					continue
				}

				char.Inventory.RemoveItem(item)
				char.Save()
				room.Inventory.AddItem(item)
				io.WriteString(s, cfmt.Sprintf("{{You drop %s.}}::green\n", bp.Name))
				room.Broadcast(cfmt.Sprintf("{{%s drops %s.}}::green\n", char.Name, bp.Name), []string{char.ID})
			}
			return
		}

		// Drop all <items>
		query := strings.Join(args[1:], " ")
		singularQuery := Singularize(query)
		matchingItems := SearchInventory(&char.Inventory, singularQuery)

		if len(matchingItems) == 0 {
			io.WriteString(s, cfmt.Sprintf("{{You have no %s to drop.}}::yellow\n", query))
			return
		}

		for _, item := range matchingItems {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp == nil {
				io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
				continue
			}

			char.Inventory.RemoveItem(item)
			char.Save()
			room.Inventory.AddItem(item)
			io.WriteString(s, cfmt.Sprintf("{{You drop %s.}}::green\n", bp.Name))
			room.Broadcast(cfmt.Sprintf("{{%s drops %s.}}::green\n", char.Name, bp.Name), []string{char.ID})
		}
		return
	}

	// Handle single item or numbered items (e.g., "drop rock" or "drop 2 rocks")
	query := strings.Join(args, " ")
	singularQuery := Singularize(query)
	matchingItems := SearchInventory(&char.Inventory, singularQuery)

	if len(matchingItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{You have no %s to drop.}}::yellow\n", query))
		return
	}

	item := matchingItems[0] // Default to the first match if ambiguous
	bp := EntityMgr.GetItemBlueprintByInstance(item)
	if bp == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
		return
	}

	char.Inventory.RemoveItem(item)
	char.Save()
	room.Inventory.AddItem(item)
	io.WriteString(s, cfmt.Sprintf("{{You drop %s.}}::green\n", bp.Name))
	room.Broadcast(cfmt.Sprintf("{{%s drops %s.}}::green\n", char.Name, bp.Name), []string{char.ID})
}

/*
Usage:
  - give <item> [to] <character>
  - give 2 <items> [to] <character>
  - give all <item> [to] <character>
*/
// TODO: Fix this to work with new inventory system
func DoGive(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Give command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) < 2 {
		io.WriteString(s, cfmt.Sprintf("{{Give what to who?}}::red\n"))
		return
	}

	// Parse the command arguments
	what := args[0]
	recipientName := args[1]

	if len(args) > 2 && (recipientName == "to" || args[2] == "to") {
		recipientName = args[len(args)-1]
	}

	// Find the recipient in the room
	var recipient *Character
	for _, r := range room.Characters {
		if strings.EqualFold(r.Name, recipientName) {
			recipient = r
			break
		}
	}

	if recipient == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no one named '%s' here.}}::red\n", recipientName))
		return
	}

	switch what {
	case "all":
		if len(args) == 2 {
			// Give all items to the recipient
			if len(char.Inventory.Items) == 0 {
				io.WriteString(s, cfmt.Sprintf("{{You have nothing to give.}}::yellow\n"))
				return
			}

			itemsToGive := make([]*Item, len(char.Inventory.Items))
			copy(itemsToGive, char.Inventory.Items)

			for _, item := range itemsToGive {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp == nil {
					io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
					continue
				}

				char.Inventory.RemoveItem(item)
				char.Save()
				recipient.Inventory.AddItem(item)
				recipient.Save()
				io.WriteString(s, cfmt.Sprintf("{{You give %s to %s.}}::green\n", bp.Name, recipient.Name))
				room.Broadcast(cfmt.Sprintf("{{%s gives %s to %s.}}::green\n", char.Name, bp.Name, recipient.Name), []string{char.ID})
			}
			return
		}

		// Give all <items> to the recipient
		query := strings.Join(args[1:len(args)-1], " ")
		singularQuery := Singularize(query)
		matchingItems := SearchInventory(&char.Inventory, singularQuery)

		if len(matchingItems) == 0 {
			io.WriteString(s, cfmt.Sprintf("{{You have no %s to give.}}::yellow\n", query))
			return
		}

		for _, item := range matchingItems {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp == nil {
				io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
				continue
			}

			char.Inventory.RemoveItem(item)
			char.Save()
			recipient.Inventory.AddItem(item)
			recipient.Save()
			io.WriteString(s, cfmt.Sprintf("{{You give %s to %s.}}::green\n", bp.Name, recipient.Name))
			room.Broadcast(cfmt.Sprintf("{{%s gives %s to %s.}}::green\n", char.Name, bp.Name, recipient.Name), []string{char.ID})
		}
		return

	default:
		// Give single or numbered items
		query := what
		singularQuery := Singularize(query)
		matchingItems := SearchInventory(&char.Inventory, singularQuery)

		if len(matchingItems) == 0 {
			io.WriteString(s, cfmt.Sprintf("{{You have no %s to give.}}::yellow\n", query))
			return
		}

		item := matchingItems[0] // Default to the first match if ambiguous
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp == nil {
			io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
			return
		}

		char.Inventory.RemoveItem(item)
		char.Save()
		recipient.Inventory.AddItem(item)
		recipient.Save()
		io.WriteString(s, cfmt.Sprintf("{{You give %s to %s.}}::green\n", bp.Name, recipient.Name))
		room.Broadcast(cfmt.Sprintf("{{%s gives %s to %s.}}::green\n", char.Name, bp.Name, recipient.Name), []string{char.ID})
	}
}

/*
Usage:
  - get all
  - get <item>
  - get <number> <items>
  - get all <items>
*/
func DoGet(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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

	if arg1 == "all" {
		if len(args) < 2 {
			io.WriteString(s, cfmt.Sprintf("{{Get all what?}}::red\n"))
			return
		}

		// Combine remaining args into the query
		query := strings.Join(args[1:], " ")

		// Handle plural and singular forms
		singularQuery := Singularize(query)
		matchingItems := SearchInventory(&room.Inventory, singularQuery)

		if len(matchingItems) == 0 {
			io.WriteString(s, cfmt.Sprintf("{{There are no %s here.}}::yellow\n", query))
			return
		}

		for _, item := range matchingItems {
			bp := EntityMgr.GetItemBlueprintByInstance(item) // Fetch the blueprint
			if bp == nil {
				io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
				continue
			}

			room.Inventory.RemoveItem(item)
			char.Inventory.AddItem(item)
			char.Save()
			io.WriteString(s, cfmt.Sprintf("{{You get %s.}}::green\n", bp.Name))
			room.Broadcast(cfmt.Sprintf("{{%s gets %s.}}::green\n", char.Name, bp.Name), []string{char.ID})
		}
		return
	}

	// Handle single item search (e.g., "get rock")
	query := strings.Join(args, " ")
	singularQuery := Singularize(query)
	matchingItems := SearchInventory(&room.Inventory, singularQuery)

	if len(matchingItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{There is no %s here.}}::yellow\n", query))
		return
	}

	item := matchingItems[0]                         // Default to the first match if ambiguous
	bp := EntityMgr.GetItemBlueprintByInstance(item) // Fetch the blueprint
	if bp == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
		return
	}

	room.Inventory.RemoveItem(item)
	char.Inventory.AddItem(item)
	char.Save()
	io.WriteString(s, cfmt.Sprintf("{{You get %s.}}::green\n", bp.Name))
	room.Broadcast(cfmt.Sprintf("{{%s gets %s.}}::green\n", char.Name, bp.Name), []string{char.ID})
}

/*
Usage:
  - look
  - look [at] <item|character|direction|mob>
*/
func DoLook(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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

/*
Usage:
  - move <north,n,south,s,east,e,west,w,up,u,down,d>
  - <north,n,south,s,east,e,west,w,up,u,down,d>
*/
func DoMove(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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

func DoSpawn(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Spawn command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) < 2 {
		io.WriteString(s, cfmt.Sprintf("{{Usage: spawn <item|mob> <name>}}::yellow\n"))
		return
	}

	entityType := args[0]
	entityName := strings.Join(args[1:], " ")

	switch entityType {
	case "i":
		// Spawn an item into the character inventory
		bp := EntityMgr.GetItemBlueprintByID(entityName)
		i := EntityMgr.CreateItemInstanceFromBlueprint(bp)
		if i == nil {
			io.WriteString(s, cfmt.Sprintf("{{Error: No item blueprint named '%s' found.}}::red\n", entityName))
			return
		}

		char.Inventory.AddItem(i)
		// room.Inventory.AddItem(i)
		io.WriteString(s, cfmt.Sprintf("{{You spawn a %s.}}::green\n", bp.Name))
		room.Broadcast(cfmt.Sprintf("{{%s spawns a %s.}}::green\n", char.Name, bp.Name), []string{char.ID})

	case "m":
		// Spawn a mob into the room
		mob := &Mob{
			ID:   uuid.New().String(),
			Name: entityName,
		}

		room.AddMob(mob)
		io.WriteString(s, cfmt.Sprintf("{{You spawn a mob named %s.}}::green\n", entityName))
		room.Broadcast(cfmt.Sprintf("{{%s spawns a mob named %s.}}::green\n", char.Name, entityName), []string{char.ID})

	default:
		io.WriteString(s, cfmt.Sprintf("{{Invalid entity type. Usage: spawn <item|mob> <name>}}::yellow\n"))
	}
}

func DoInventory(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Inventory command",
		slog.String("command", cmd),
		slog.Any("args", args))

	if char == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red\n"))
		return
	}

	if len(char.Inventory.Items) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{You are not carrying anything.}}::yellow\n"))
		return
	}

	io.WriteString(s, cfmt.Sprintf("{{You are carrying:}}::cyan\n"))
	itemCounts := make(map[string]int)

	// Count items based on their blueprint name
	for _, item := range char.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp == nil {
			io.WriteString(s, cfmt.Sprintf("{{Error: Unable to retrieve item blueprint.}}::red\n"))
			continue
		}
		itemCounts[bp.Name]++
	}

	// Display the inventory
	for name, count := range itemCounts {
		io.WriteString(s, cfmt.Sprintf("- %s\n",
			pluralizer.PluralizeNounPhrase(name, count)))
	}
}
