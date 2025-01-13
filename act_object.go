package main

import (
	"io"
	"strconv"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

/*
Usage:
  - drop [<quantity>] <item>
  - drop all <item>
  - drop all
*/
func DoDrop(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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
func DoGive(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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
func DoGet(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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
