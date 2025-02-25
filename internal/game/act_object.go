package game

import (
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
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

	if len(args) == 0 {
		WriteString(s, "{{Drop what?}}::red"+CRLF)
		return
	}

	quantity := 1
	itemQuery := ""

	// Parse quantity and item query
	if len(args) > 1 {
		if parsedQuantity, err := strconv.Atoi(args[0]); err == nil {
			quantity = parsedQuantity
			itemQuery = strings.Join(args[1:], " ")
		} else {
			itemQuery = strings.Join(args, " ")
		}
	} else {
		itemQuery = args[0]
	}

	if strings.EqualFold(itemQuery, "all") {
		// Drop all items
		if len(char.Inventory.Items) == 0 {
			WriteString(s, "{{You do not seem to have anything.}}::yellow"+CRLF)
			return
		}

		droppedItems := make(map[string]int)

		for i := 0; i < len(char.Inventory.Items); {
			item := char.Inventory.Items[i]
			blueprint := EntityMgr.GetItemBlueprintByInstance(item)
			droppedItems[blueprint.Name]++
			room.Inventory.AddItem(item)
			char.Inventory.RemoveItem(item)
		}

		for itemName, count := range droppedItems {
			WriteStringF(s, "{{You drop %d %s.}}::green"+CRLF, count, pluralizer.PluralizeNoun(itemName, count))
			room.Broadcast(cfmt.Sprintf("{{%s drops %d %s.}}::green"+CRLF, char.Name, count, pluralizer.PluralizeNoun(itemName, count)), []string{char.ID})
		}
		return
	}

	if strings.HasPrefix(strings.ToLower(itemQuery), "all ") {
		// Drop all items matching a prefix
		prefix := strings.ToLower(strings.TrimPrefix(strings.ToLower(itemQuery), "all "))
		found := false
		droppedItems := make(map[string]int)

		for i := 0; i < len(char.Inventory.Items); {
			item := char.Inventory.Items[i]
			blueprint := EntityMgr.GetItemBlueprintByInstance(item)
			if !strings.HasPrefix(strings.ToLower(blueprint.Name), prefix) {
				i++
				continue
			}

			droppedItems[blueprint.Name]++
			room.Inventory.AddItem(item)
			char.Inventory.RemoveItem(item)
			found = true
		}

		for itemName, count := range droppedItems {
			WriteStringF(s, "{{You drop %d %s.}}::green"+CRLF, count, pluralizer.PluralizeNoun(itemName, count))
			room.Broadcast(cfmt.Sprintf("{{%s drops %d %s.}}::green"+CRLF, char.Name, count, pluralizer.PluralizeNoun(itemName, count)), []string{char.ID})
		}

		if !found {
			WriteStringF(s, "{{You do not have any items matching '%s'.}}::yellow"+CRLF, prefix)
		}
		return
	}

	// Drop a specific item or quantity of an item
	items := char.Inventory.Search(itemQuery)
	if len(items) == 0 {
		WriteString(s, "{{You do not have that item.}}::yellow"+CRLF)
		return
	}

	if quantity > len(items) {
		blueprint := EntityMgr.GetItemBlueprintByInstance(items[0])
		WriteStringF(s, "{{You do not have %d %s.}}::yellow"+CRLF, quantity, pluralizer.PluralizeNoun(blueprint.Name, quantity))
		return
	}

	droppedItems := make(map[string]int)
	for _, item := range items[:quantity] {
		itemName := EntityMgr.GetItemBlueprintByInstance(item).Name
		droppedItems[itemName]++
		room.Inventory.AddItem(item)
		char.Inventory.RemoveItem(item)
	}

	for itemName, count := range droppedItems {
		WriteStringF(s, "{{You drop %d %s.}}::green"+CRLF, count, pluralizer.PluralizeNoun(itemName, count))
		room.Broadcast(cfmt.Sprintf("{{%s drops %d %s.}}::green"+CRLF, char.Name, count, pluralizer.PluralizeNoun(itemName, count)), []string{char.ID})
	}
}

func SuggestDrop(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	if len(args) == 0 {
		// Suggest "all" or item names
		suggestions = append(suggestions, "all")
		for _, item := range char.Inventory.Items {
			blueprint := EntityMgr.GetItemBlueprintByInstance(item)
			suggestions = append(suggestions, blueprint.Name)
		}
	} else if len(args) == 1 {
		// Suggest items based on partial input
		for _, item := range char.Inventory.Items {
			blueprint := EntityMgr.GetItemBlueprintByInstance(item)
			if strings.HasPrefix(strings.ToLower(blueprint.Name), strings.ToLower(args[0])) {
				suggestions = append(suggestions, blueprint.Name)
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
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

	if len(args) < 2 {
		WriteString(s, "{{Usage: give <character> [<quantity>] <item>.}}::red"+CRLF)
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
		WriteStringF(s, "{{There is no one named '%s' here.}}::red"+CRLF, recipientName)
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
	matchingItems := char.Inventory.Search(singularQuery)
	if len(matchingItems) == 0 {
		WriteStringF(s, "{{You have no %s to give.}}::yellow"+CRLF, itemQuery)
		return
	}

	if quantity > len(matchingItems) {
		WriteStringF(s, "{{You do not have %d %s.}}::yellow"+CRLF, quantity, itemQuery)
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
		WriteStringF(s, "{{%s cannot carry any more weight.}}::yellow"+CRLF, recipient.Name)
		return
	}

	itemName := pluralizer.PluralizeNoun(EntityMgr.GetItemBlueprintByInstance(givenItems[0]).Name, len(givenItems))
	WriteStringF(s, "{{You give %s %d %s.}}::green"+CRLF, recipient.Name, len(givenItems), itemName)
	room.Broadcast(cfmt.Sprintf("{{%s gives %s %d %s.}}::green"+CRLF, char.Name, recipient.Name, len(givenItems), itemName), []string{char.ID, recipient.ID})

	if len(givenItems) < quantity {
		WriteStringF(s, "{{%s could not take all items due to weight limits.}}::yellow"+CRLF, recipient.Name)
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
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

	if len(args) == 0 {
		WriteString(s, "{{Get what?}}::red"+CRLF)
		return
	}

	quantity := 1
	itemQuery := ""

	// Parse arguments
	if strings.EqualFold(args[0], "all") {
		if len(args) == 1 {
			quantity = -1 // "get all"
			itemQuery = "all"
		} else {
			itemQuery = strings.Join(args[1:], " ") // "get all <item>"
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
		itemQuery = args[0] // "get <item>"
	}

	singularQuery := Singularize(itemQuery)
	matchingItems := room.Inventory.Search(singularQuery)

	// If no items match the query, inform the user
	if len(matchingItems) == 0 {
		if itemQuery == "all" {
			WriteString(s, "{{There are no items left here to get.}}::yellow"+CRLF)
		} else {
			WriteStringF(s, "{{There are no %s here.}}::yellow"+CRLF, itemQuery)
		}
		return
	}

	if quantity == -1 { // "get all <item>"
		quantity = len(matchingItems)
	}

	if quantity > len(matchingItems) {
		WriteStringF(s, "{{There are not %d %s here.}}::yellow"+CRLF, quantity, itemQuery)
		return
	}

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

	// If no items could be picked up, send an appropriate message
	if len(pickedItems) == 0 {
		if itemQuery == "all" {
			WriteString(s, "{{You cannot carry any more items here due to weight or capacity limits.}}::yellow"+CRLF)
		} else {
			WriteStringF(s, "{{You cannot carry any %s due to weight or capacity limits.}}::yellow"+CRLF, itemQuery)
		}
		return
	}

	// Transfer items to the character's inventory
	for _, item := range pickedItems {
		room.Inventory.RemoveItem(item)
		char.Inventory.AddItem(item)
	}

	// Remove picked items from `matchingItems` to update for subsequent "get all" commands
	matchingItems = matchingItems[len(pickedItems):]

	char.Save()

	// Inform the user about the items they successfully picked up
	itemName := pluralizer.PluralizeNoun(EntityMgr.GetItemBlueprintByInstance(pickedItems[0]).Name, len(pickedItems))
	WriteStringF(s, "{{You get %d %s.}}::green"+CRLF, len(pickedItems), itemName)
	room.Broadcast(cfmt.Sprintf("{{%s picks up %d %s.}}::green"+CRLF, char.Name, len(pickedItems), itemName), []string{char.ID})

	// Additional feedback for partial pickups
	if itemQuery == "all" && len(matchingItems) > 0 {
		WriteString(s, "{{You left some items behind due to weight limits.}}::yellow"+CRLF)
	} else if len(pickedItems) < quantity {
		WriteStringF(s, "{{You left some %s behind due to weight limits.}}::yellow"+CRLF, itemQuery)
	}
}
func SuggestGet(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0:
		// Suggest "all" or individual items
		suggestions = append(suggestions, "all")
		seenItems := map[string]bool{}
		for _, item := range room.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp != nil && !seenItems[bp.Name] {
				seenItems[bp.Name] = true
				suggestions = append(suggestions, bp.Name)
			}
		}
	case 1:
		if strings.EqualFold(args[0], "all") {
			// Suggest "all <item>" for items in the room
			seenItems := map[string]bool{}
			for _, item := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil && !seenItems[bp.Name] {
					seenItems[bp.Name] = true
					suggestions = append(suggestions, "all "+bp.Name)
				}
			}
		} else {
			// Suggest individual items matching the query
			seenItems := map[string]bool{}
			for _, item := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(item)
				if bp != nil && !seenItems[bp.Name] {
					seenItems[bp.Name] = true
					suggestions = append(suggestions, bp.Name)
				}
			}
		}
	default:
		// No further suggestions for additional arguments
	}

	return suggestions
}
