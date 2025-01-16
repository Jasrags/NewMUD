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
	matchingItems := SearchInventory(char.Inventory, singularQuery)
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
		if parsedQuantity, err := strconv.Atoi(itemArgs[0]); err == nil {
			quantity = parsedQuantity
			itemArgs = itemArgs[1:]
		}
	}

	itemQuery := strings.Join(itemArgs, " ")
	singularQuery := Singularize(itemQuery)

	// Search inventory
	matchingItems := SearchInventory(char.Inventory, singularQuery)
	if len(matchingItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{You have no %s to give.}}::yellow\n", itemQuery))
		return
	}

	if quantity > len(matchingItems) {
		io.WriteString(s, cfmt.Sprintf("{{You do not have %d %s.}}::yellow\n", quantity, itemQuery))
		return
	}

	// Calculate recipient's capacity and transfer items
	remainingCapacity := float64(recipient.GetLiftCarry()) - recipient.GetCurrentCarryWeight()
	givenItems := []*Item{}
	totalWeight := 0.0

	for _, item := range matchingItems[:quantity] {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp != nil && bp.Weight <= remainingCapacity {
			totalWeight += bp.Weight
			if totalWeight > remainingCapacity {
				break
			}
			givenItems = append(givenItems, item)
			remainingCapacity -= bp.Weight
		}
	}

	// Transfer given items
	for _, item := range givenItems {
		char.Inventory.RemoveItem(item)
		recipient.Inventory.AddItem(item)
	}

	char.Save()
	recipient.Save()

	if len(givenItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{%s cannot carry any more weight.}}::yellow\n", recipient.Name))
		return
	}

	itemName := pluralizer.PluralizeNoun(EntityMgr.GetItemBlueprintByInstance(givenItems[0]).Name, len(givenItems))
	io.WriteString(s, cfmt.Sprintf("{{You give %s %d %s.}}::green\n", recipient.Name, len(givenItems), itemName))
	room.Broadcast(cfmt.Sprintf("{{%s gives %s %d %s.}}::green\n", char.Name, recipient.Name, len(givenItems), itemName), []string{char.ID, recipient.ID})

	if len(givenItems) < quantity {
		io.WriteString(s, cfmt.Sprintf("{{%s could not take all items due to weight limits.}}::yellow\n", recipient.Name))
	}
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
			// "get all"
			quantity = -1
			itemQuery = "all"
		} else {
			// "get all <item>"
			itemQuery = strings.Join(args[1:], " ")
			quantity = -1
		}
	} else if len(args) > 1 {
		if parsedQuantity, err := strconv.Atoi(args[0]); err == nil {
			quantity = parsedQuantity
			itemQuery = strings.Join(args[1:], " ")
		} else {
			itemQuery = strings.Join(args, " ")
		}
	} else {
		itemQuery = args[0]
	}

	// Search for items
	singularQuery := Singularize(itemQuery)
	matchingItems := SearchInventory(room.Inventory, singularQuery)
	if len(matchingItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{There are no %s here.}}::yellow\n", itemQuery))
		return
	}

	if quantity == -1 { // "get all <item>"
		quantity = len(matchingItems)
	}

	if quantity > len(matchingItems) {
		io.WriteString(s, cfmt.Sprintf("{{There are not %d %s here.}}::yellow\n", quantity, itemQuery))
		return
	}

	// Calculate carry capacity and transfer items
	remainingCapacity := float64(char.GetLiftCarry()) - char.GetCurrentCarryWeight()
	pickedItems := []*Item{}
	totalWeight := 0.0

	for _, item := range matchingItems[:quantity] {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp != nil && bp.Weight <= remainingCapacity {
			totalWeight += bp.Weight
			if totalWeight > remainingCapacity {
				break
			}
			pickedItems = append(pickedItems, item)
			remainingCapacity -= bp.Weight
		}
	}

	// Transfer picked items
	for _, item := range pickedItems {
		room.Inventory.RemoveItem(item)
		char.Inventory.AddItem(item)
	}

	char.Save()

	if len(pickedItems) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{You cannot carry any more weight.}}::yellow\n"))
		return
	}

	itemName := pluralizer.PluralizeNoun(EntityMgr.GetItemBlueprintByInstance(pickedItems[0]).Name, len(pickedItems))
	io.WriteString(s, cfmt.Sprintf("{{You get %d %s.}}::green\n", len(pickedItems), itemName))
	room.Broadcast(cfmt.Sprintf("{{%s picks up %d %s.}}::green\n", char.Name, len(pickedItems), itemName), []string{char.ID})

	// Specifically handle "get all" feedback
	if itemQuery == "all" && len(pickedItems) < len(room.Inventory.Items) {
		io.WriteString(s, cfmt.Sprintf("{{You left some items behind due to weight limits.}}::yellow\n"))
	} else if len(pickedItems) < quantity {
		io.WriteString(s, cfmt.Sprintf("{{You left some %s behind due to weight limits.}}::yellow\n", itemQuery))
	}
}

func SuggestGet(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest "all" or individual items
		suggestions = append(suggestions, "all")
		seenItems := map[string]bool{}
		for _, item := range room.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp != nil && !seenItems[bp.Name] {
				seenItems[bp.Name] = true
				suggestions = append(suggestions, bp.Name)
			}
		}
	case 1: // Handle "all" case explicitly
		if strings.EqualFold(args[0], "all") {
			suggestions = append(suggestions, "all")
			seenItems := map[string]bool{}
			for _, item := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil && !seenItems[bp.Name] {
					seenItems[bp.Name] = true
					suggestions = append(suggestions, "all "+bp.Name)
				}
			}
		} else { // Suggest individual items
			seenItems := map[string]bool{}
			for _, item := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil && !seenItems[bp.Name] {
					seenItems[bp.Name] = true
					suggestions = append(suggestions, bp.Name)
				}
			}
		}
	}

	return suggestions
}
