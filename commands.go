package main

import (
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/exp/rand"
)

// TODO: We need a RP consistent way to communicate directly with other individuals or groups of individuals I.E. for shadowrun it could be via comlinks and some group or party system

type SuggestFunc func(line string, args []string, char *Character, room *Room) []string

type Command struct {
	Name          string
	Description   string
	Usage         []string
	Aliases       []string
	RequiredRoles []CharacterRole
	Func          CommandFunc
	SuggestFunc   SuggestFunc // Optional suggestion logic
}

type CommandFunc func(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room)

func DoPick(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Pick what?}}::yellow\n"))
		return
	}

	direction := args[0]
	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if !exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is not locked.}}::yellow\n", direction))
		return
	}

	// if !hasKey {
	//     if exit.Door.PickDifficulty > 0 {
	//         success := AttemptLockPick(char, exit.Door.PickDifficulty)
	//         if success {
	//             exit.Door.IsLocked = false
	//             io.WriteString(s, cfmt.Sprintf("{{You successfully pick the lock on the door to the %s.}}::green\n", direction))
	//             room.Broadcast(cfmt.Sprintf("{{%s picks the lock on the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
	//             return
	//         } else {
	//             io.WriteString(s, cfmt.Sprintf("{{You fail to pick the lock on the door to the %s.}}::red\n", direction))
	//             return
	//         }
	//     }

	//     io.WriteString(s, cfmt.Sprintf("{{You don't have the key to unlock the door to the %s.}}::red\n", direction))
	//     return
	// }

	pickRoll := rand.Intn(100) + 1 // Random roll between 1 and 100
	if pickRoll > exit.Door.PickDifficulty {
		exit.Door.IsLocked = false
		io.WriteString(s, cfmt.Sprintf("{{You successfully pick the lock on the door to the %s.}}::green\n", direction))
		room.Broadcast(cfmt.Sprintf("{{%s picks the lock on the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
	} else {
		io.WriteString(s, cfmt.Sprintf("{{You fail to pick the lock on the door to the %s.}}::red\n", direction))
	}
}

func DoClose(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Close what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])

	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if exit.Door.IsClosed {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is already closed.}}::yellow\n", direction))
		return
	}

	exit.Door.IsClosed = true
	io.WriteString(s, cfmt.Sprintf("{{You close the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s closes the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s closes from the other side.}}::green\n", ReverseDirection(direction)), []string{})
	}
}

func DoLock(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Lock what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])
	if direction == "" {
		io.WriteString(s, cfmt.Sprintf("{{Invalid direction.}}::red\n"))
		return
	}

	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is already locked.}}::yellow\n", direction))
		return
	}

	if !exit.Door.IsClosed {
		io.WriteString(s, cfmt.Sprintf("{{You must close the door to the %s before locking it.}}::yellow\n", direction))
		return
	}

	validKeys := make(map[string]bool)
	for _, key := range exit.Door.KeyIDs {
		validKeys[key] = true
	}

	hasKey := false
	for _, item := range char.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp.Type == ItemTypeKey && validKeys[bp.ID] {
			hasKey = true
			break
		}
	}

	if !hasKey {
		io.WriteString(s, cfmt.Sprintf("{{You don't have the key to lock the door to the %s.}}::red\n", direction))
		return
	}

	exit.Door.IsLocked = true
	io.WriteString(s, cfmt.Sprintf("{{You lock the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s locks the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
}

func DoUnlock(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Unlock what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if !exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is not locked.}}::yellow\n", direction))
		return
	}

	validKeys := make(map[string]bool)
	for _, key := range exit.Door.KeyIDs {
		validKeys[key] = true
	}

	// Check if character has the correct key
	hasKey := false
	for _, item := range char.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp.Type == ItemTypeKey && validKeys[bp.ID] {
			hasKey = true
			break
		}
	}

	if !hasKey {
		io.WriteString(s, cfmt.Sprintf("{{You don't have the key to unlock the door to the %s.}}::red\n", direction))
		return
	}

	exit.Door.IsLocked = false
	io.WriteString(s, cfmt.Sprintf("{{You unlock the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s unlocks the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
}

// This command is for opening closed entities
func DoOpen(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Open what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if !exit.Door.IsClosed {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is already open.}}::yellow\n", direction))
		return
	}

	if exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is locked.}}::red\n", direction))
		return
	}

	exit.Door.IsClosed = false
	io.WriteString(s, cfmt.Sprintf("{{You open the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s opens the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s opens from the other side.}}::green\n", ReverseDirection(direction)), []string{})
	}
}

/*
Usage:
  - who
*/
// TODO: Sort all admins to the top of the list
// TODO: Add a CanSee function for characters and have this function use that to determine if a character can see another character in the who list
func DoWho(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Who command",
		slog.String("command", cmd),
		slog.Any("args", args))

	// Simulated global list of active characters
	activeCharacters := CharacterMgr.GetOnlineCharacters()

	if len(activeCharacters) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{No one else is in the game right now.}}::yellow\n"))
		return
	}

	io.WriteString(s, cfmt.Sprintf("{{Players currently in the game:}}::green\n"))

	for _, activeChar := range activeCharacters {
		color := "cyan"
		if activeChar.Role == CharacterRoleAdmin {
			color = "yellow"
		}

		if activeChar.Title == "" {
			activeChar.Title = "the Basic"
		}

		// Display character title and name
		io.WriteString(s, cfmt.Sprintf("{{%s - %s}}::%s\n", activeChar.Name, activeChar.Title, color))
	}
}

/*
Usage:
  - say <message>
  - say @<name> <message>
*/
// TODO: overall for communication commands we need to log messages to a database with time, to/from, and message.
// TODO: need to implement a block/unblock function for preventing messages from certain users
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
		if CommandMgr.CanRunCommand(char, cmd) {
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
  - give <character> [<quantity>] <item>
*/
func DoGive(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	slog.Debug("Give command",
		slog.String("command", cmd),
		slog.Any("args", args),
		slog.String("character_id", char.ID),
		slog.String("character_name", char.Name))

	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) < 2 {
		io.WriteString(s, cfmt.Sprintf("{{Usage: give <character> [<quantity>] <item>.}}::red\n"))
		return
	}

	// Parse recipient
	recipientName := args[0]
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

	// Parse item and quantity
	quantity := 1
	itemArgs := args[1:]
	if len(itemArgs) > 1 {
		// Check if the second argument is numeric
		if parsedQuantity, err := strconv.Atoi(itemArgs[0]); err == nil {
			quantity = parsedQuantity
			itemArgs = itemArgs[1:] // Remove quantity from arguments
		}
	}

	// Ensure the item query is provided
	if len(itemArgs) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{You must specify an item to give.}}::yellow\n"))
		return
	}

	itemQuery := strings.Join(itemArgs, " ")
	singularQuery := Singularize(itemQuery)

	// Search inventory
	matchingItems := SearchInventory(&char.Inventory, singularQuery)
	if len(matchingItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{You have no %s to give.}}::yellow\n", itemQuery))
		return
	}

	if quantity > len(matchingItems) {
		io.WriteString(s, cfmt.Sprintf("{{You do not have %d %s.}}::yellow\n", quantity, itemQuery))
		return
	}

	// Transfer items
	itemsToGive := matchingItems[:quantity]
	for _, item := range itemsToGive {
		char.Inventory.RemoveItem(item)
		recipient.Inventory.AddItem(item)
	}
	char.Save()
	recipient.Save()

	// Get item blueprint for consistent naming
	bp := EntityMgr.GetItemBlueprintByInstance(itemsToGive[0])
	if bp == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
		return
	}

	itemName := bp.Name
	if quantity > 1 {
		itemName = pluralizer.PluralizeNoun(itemName, quantity)
	}

	// Message the giver
	io.WriteString(s, cfmt.Sprintf("{{You give %s %d %s.}}::green\n", recipient.Name, quantity, itemName))

	// Message the room
	room.Broadcast(cfmt.Sprintf("{{%s gives %s %d %s.}}::green\n", char.Name, recipient.Name, quantity, itemName), []string{char.ID, recipient.ID})

	// Message the recipient
	io.WriteString(recipient.Conn, cfmt.Sprintf("{{%s gives you %d %s.}}::cyan\n", char.Name, quantity, itemName))
}

func SuggestGive(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest character names for the first argument
		for _, r := range room.Characters {
			if !strings.EqualFold(r.Name, char.Name) { // Exclude self
				suggestions = append(suggestions, r.Name)
			}
		}
	case 1: // Suggest numeric quantities or items
		if parsedQuantity, err := strconv.Atoi(args[0]); err == nil {
			suggestions = append(suggestions, strconv.Itoa(parsedQuantity+1), strconv.Itoa(parsedQuantity+2))
		} else {
			for _, item := range char.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp == nil {
					continue
				}
				suggestions = append(suggestions, bp.Name)
			}
		}
	case 2: // Suggest item names if quantity is already provided
		for _, item := range char.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp == nil {
				continue
			}
			suggestions = append(suggestions, bp.Name)
		}
	default:
		// No suggestions for further arguments
	}

	return suggestions
	// return FilterSuggestions(line, suggestions)
}

func FilterSuggestions(input string, suggestions []string) []string {
	matches := []string{}
	lowerInput := strings.ToLower(input) // Normalize input to lowercase for case-insensitive matching
	for _, suggestion := range suggestions {
		if strings.HasPrefix(strings.ToLower(suggestion), lowerInput) {
			matches = append(matches, suggestion)
		}
	}
	return matches
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

	dir := ParseDirection(cmd)

	// Check if the exit exists
	if exit, ok := char.Room.Exits[dir]; ok {
		if exit.Door != nil && exit.Door.IsClosed {
			io.WriteString(s, cfmt.Sprintf("{{The door to the %s is closed.}}::red\n", dir))
			return
		}

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
