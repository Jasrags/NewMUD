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
*/
// TODO: overall for communication commands we need to log messages to a database with time, to/from, and message.
// TODO: need to implement a block/unblock function for preventing messages from certain users
func DoSay(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{What do you want to say?}}::red\n"))
		return
	}

	message := strings.Join(args, " ")

	// Broadcast message to the room
	room.Broadcast(cfmt.Sprintf("{{%s says: \"%s\"}}::green\n", char.Name, message), []string{char.ID})

	// Message the player
	io.WriteString(s, cfmt.Sprintf("{{You say: \"%s\"}}::green\n", message))
}

func DoTell(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) < 2 {
		io.WriteString(s, cfmt.Sprintf("{{Usage: tell <username> <message>.}}::red\n"))
		return
	}

	recipientName := args[0]
	message := strings.Join(args[1:], " ")

	var recipient *Character
	for _, r := range room.Characters {
		if strings.EqualFold(r.Name, recipientName) {
			recipient = r
			break
		}
	}

	if recipient == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no one named '%s' here.}}::yellow\n", recipientName))
		return
	}

	// Message the recipient
	io.WriteString(recipient.Conn, cfmt.Sprintf("{{%s tells you: \"%s\"}}::cyan\n", char.Name, message))

	// Message the sender
	io.WriteString(s, cfmt.Sprintf("{{You tell %s: \"%s\"}}::green\n", recipient.Name, message))

	// Message the room (excluding sender and recipient)
	room.Broadcast(cfmt.Sprintf("{{%s tells %s something privately.}}::green\n", char.Name, recipient.Name), []string{char.ID, recipient.ID})
}

func SuggestTell(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest names of characters in the room
		for _, r := range room.Characters {
			if !strings.EqualFold(r.Name, char.Name) { // Exclude self
				suggestions = append(suggestions, r.Name)
			}
		}
	case 1: // Suggest partial names
		for _, r := range room.Characters {
			if !strings.EqualFold(r.Name, char.Name) && strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(args[0])) {
				suggestions = append(suggestions, r.Name)
			}
		}
	}

	return suggestions
}

/*
Usage:
  - help
  - help <command>
*/
func DoHelp(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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
  - drop [<quantity>] <item>
  - drop all <item>
  - drop all
*/
func DoDrop(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Drop what?}}::red\n"))
		return
	}

	quantity := 1
	itemQuery := ""

	// Parse arguments
	if strings.EqualFold(args[0], "all") {
		if len(args) == 1 {
			// Usage: "drop all" (all items in inventory)
			quantity = -1 // Indicate "all"
			itemQuery = "all"
		} else {
			// Usage: "drop all <item>"
			itemQuery = strings.Join(args[1:], " ")
			quantity = -1
		}
	} else if len(args) > 1 {
		// Usage: "drop [<quantity>] <item>"
		if parsedQuantity, err := strconv.Atoi(args[0]); err == nil {
			quantity = parsedQuantity
			itemQuery = strings.Join(args[1:], " ")
		} else {
			itemQuery = strings.Join(args, " ")
		}
	} else {
		// Usage: "drop <item>"
		itemQuery = args[0]
	}

	if itemQuery == "all" {
		// Handle "drop all" (drop all items in inventory)
		if len(char.Inventory.Items) == 0 {
			io.WriteString(s, cfmt.Sprintf("{{You have no items to drop.}}::yellow\n"))
			return
		}

		for _, item := range char.Inventory.Items {
			room.Inventory.AddItem(item)
		}

		count := len(char.Inventory.Items)
		char.Inventory.Clear()
		char.Save()

		// Messaging
		io.WriteString(s, cfmt.Sprintf("{{You drop %d items.}}::green\n", count))
		room.Broadcast(cfmt.Sprintf("{{%s drops all items.}}::green\n", char.Name), []string{char.ID})
		return
	}

	// Handle "drop all <item>" or "drop [<quantity>] <item>"
	singularQuery := Singularize(itemQuery)
	matchingItems := SearchInventory(&char.Inventory, singularQuery)
	if len(matchingItems) == 0 {
		if strings.EqualFold(itemQuery, "all") {
			io.WriteString(s, cfmt.Sprintf("{{You have no items to drop.}}::yellow\n"))
		} else {
			io.WriteString(s, cfmt.Sprintf("{{You have no %s to drop.}}::yellow\n", itemQuery))
		}
		return
	}

	if quantity == -1 { // Handle "all <item>"
		quantity = len(matchingItems)
	}

	if quantity > len(matchingItems) {
		io.WriteString(s, cfmt.Sprintf("{{You do not have %d %s.}}::yellow\n", quantity, itemQuery))
		return
	}

	// Transfer items
	itemsToDrop := matchingItems[:quantity]
	for _, item := range itemsToDrop {
		char.Inventory.RemoveItem(item)
		room.Inventory.AddItem(item)
	}
	char.Save()

	// Use the first item's blueprint to format the message
	bp := EntityMgr.GetItemBlueprintByInstance(itemsToDrop[0])
	if bp == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
		return
	}

	itemName := bp.Name
	if quantity > 1 {
		itemName = pluralizer.PluralizeNoun(itemName, quantity)
	}

	// Messaging
	io.WriteString(s, cfmt.Sprintf("{{You drop %d %s.}}::green\n", quantity, itemName))
	room.Broadcast(cfmt.Sprintf("{{%s drops %d %s.}}::green\n", char.Name, quantity, itemName), []string{char.ID})
}

func SuggestDrop(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest "all" or items
		suggestions = append(suggestions, "all")
		for _, item := range char.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp != nil {
				suggestions = append(suggestions, bp.Name)
			}
		}
	case 1: // Suggest items or "all <item>"
		if strings.EqualFold(args[0], "all") {
			for _, item := range char.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil {
					suggestions = append(suggestions, "all "+bp.Name)
				}
			}
		} else {
			suggestions = append(suggestions, "all")
			for _, item := range char.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil {
					suggestions = append(suggestions, bp.Name)
				}
			}
		}
	}

	return suggestions
}

/*
Usage:
  - give <character> [<quantity>] <item>
*/
func DoGive(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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
}

/*
Usage:
  - get [<quantity>] <item>
  - get all <item>
  - get all
*/
func DoGet(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Get what?}}::red\n"))
		return
	}

	quantity := 1
	itemQuery := ""

	// Parse arguments
	if strings.EqualFold(args[0], "all") {
		if len(args) == 1 {
			// Usage: "get all" (all items in the room)
			quantity = -1 // Indicate "all"
			itemQuery = "all"
		} else {
			// Usage: "get all <item>"
			itemQuery = strings.Join(args[1:], " ")
			quantity = -1
		}
	} else if len(args) > 1 {
		// Usage: "get [<quantity>] <item>"
		if parsedQuantity, err := strconv.Atoi(args[0]); err == nil {
			quantity = parsedQuantity
			itemQuery = strings.Join(args[1:], " ")
		} else {
			itemQuery = strings.Join(args, " ")
		}
	} else {
		// Usage: "get <item>"
		itemQuery = args[0]
	}

	if itemQuery == "all" {
		// Handle "get all" (retrieve all items in the room)
		if len(room.Inventory.Items) == 0 {
			io.WriteString(s, cfmt.Sprintf("{{There are no items here to get.}}::yellow\n"))
			return
		}

		for _, item := range room.Inventory.Items {
			char.Inventory.AddItem(item)
		}

		count := len(room.Inventory.Items)
		room.Inventory.Clear()
		char.Save()

		// Messaging
		io.WriteString(s, cfmt.Sprintf("{{You get %d items.}}::green\n", count))
		room.Broadcast(cfmt.Sprintf("{{%s gets all items.}}::green\n", char.Name), []string{char.ID})
		return
	}

	// Handle "get all <item>" or "get [<quantity>] <item>"
	singularQuery := Singularize(itemQuery)
	matchingItems := SearchInventory(&room.Inventory, singularQuery)
	if len(matchingItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{There are no %s here.}}::yellow\n", itemQuery))
		return
	}

	if quantity == -1 { // Handle "all <item>"
		quantity = len(matchingItems)
	}

	if quantity > len(matchingItems) {
		io.WriteString(s, cfmt.Sprintf("{{There are not %d %s here.}}::yellow\n", quantity, itemQuery))
		return
	}

	// Transfer items
	itemsToGet := matchingItems[:quantity]
	for _, item := range itemsToGet {
		room.Inventory.RemoveItem(item)
		char.Inventory.AddItem(item)
	}
	char.Save()

	// Use the first item's blueprint to format the message
	bp := EntityMgr.GetItemBlueprintByInstance(itemsToGet[0])
	if bp == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error retrieving item blueprint.}}::red\n"))
		return
	}

	itemName := bp.Name
	if quantity > 1 {
		itemName = pluralizer.PluralizeNoun(itemName, quantity)
	}

	// Messaging
	io.WriteString(s, cfmt.Sprintf("{{You get %d %s.}}::green\n", quantity, itemName))
	room.Broadcast(cfmt.Sprintf("{{%s gets %d %s.}}::green\n", char.Name, quantity, itemName), []string{char.ID})
}

func SuggestGet(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest "all" or items
		suggestions = append(suggestions, "all")
		for _, item := range room.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp != nil {
				suggestions = append(suggestions, bp.Name)
			}
		}
	case 1: // Suggest items or "all <item>"
		if strings.EqualFold(args[0], "all") {
			for _, item := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil {
					suggestions = append(suggestions, "all "+bp.Name)
				}
			}
		} else {
			suggestions = append(suggestions, "all")
			for _, item := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil {
					suggestions = append(suggestions, bp.Name)
				}
			}
		}
	default: // No suggestions for further arguments
	}

	return suggestions
}

/*
Usage:
  - look
  - look [at] <item|character|direction|mob>
*/
func DoLook(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
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
