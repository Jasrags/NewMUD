package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/exp/rand"
)

func DoLock(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Lock what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])
	if direction == "" {
		WriteString(s, "{{Invalid direction.}}::red"+CRLF)
		return
	}

	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is already locked.}}::yellow"+CRLF, direction)
		return
	}

	if !exit.Door.IsClosed {
		WriteStringF(s, "{{You must close the door to the %s before locking it.}}::yellow"+CRLF, direction)
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
		WriteStringF(s, "{{You don't have the key to lock the door to the %s.}}::red"+CRLF, direction)
		return
	}

	exit.Door.IsLocked = true
	WriteStringF(s, "{{You lock the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s locks the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})
}

func DoUnlock(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Unlock what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if !exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is not locked.}}::yellow"+CRLF, direction)
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
		WriteStringF(s, "{{You don't have the key to unlock the door to the %s.}}::red"+CRLF, direction)
		return
	}

	exit.Door.IsLocked = false
	WriteStringF(s, "{{You unlock the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s unlocks the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})
}

func DoPick(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Pick what?}}::yellow"+CRLF)
		return
	}

	direction := args[0]
	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if !exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is not locked.}}::yellow"+CRLF, direction)
		return
	}

	// if !hasKey {
	//     if exit.Door.PickDifficulty > 0 {
	//         success := AttemptLockPick(char, exit.Door.PickDifficulty)
	//         if success {
	//             exit.Door.IsLocked = false
	//             WriteString(s, "{{You successfully pick the lock on the door to the %s.}}::green"+CRLF, direction)
	//             room.Broadcast("{{%s picks the lock on the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID}
	//             return
	//         } else {
	//             WriteString(s, "{{You fail to pick the lock on the door to the %s.}}::red"+CRLF, direction)
	//             return
	//         }
	//     }

	//     WriteString(s, "{{You don't have the key to unlock the door to the %s.}}::red"+CRLF, direction)
	//     return
	// }

	pickRoll := rand.Intn(100) + 1 // Random roll between 1 and 100
	if pickRoll > exit.Door.PickDifficulty {
		exit.Door.IsLocked = false
		WriteStringF(s, "{{You successfully pick the lock on the door to the %s.}}::green"+CRLF, direction)
		room.Broadcast(cfmt.Sprintf("{{%s picks the lock on the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})
	} else {
		WriteStringF(s, "{{You fail to pick the lock on the door to the %s.}}::red"+CRLF, direction)
	}
}

// This command is for opening closed entities
func DoOpen(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Open what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if !exit.Door.IsClosed {
		WriteStringF(s, "{{The door to the %s is already open.}}::yellow"+CRLF, direction)
		return
	}

	if exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is locked.}}::red"+CRLF, direction)
		return
	}

	exit.Door.IsClosed = false
	WriteStringF(s, "{{You open the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s opens the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s opens from the other side.}}::green"+CRLF, ReverseDirection(direction)), []string{})
	}
}

func DoClose(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Close what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])

	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door.IsClosed {
		WriteStringF(s, "{{The door to the %s is already closed.}}::yellow"+CRLF, direction)
		return
	}

	exit.Door.IsClosed = true
	WriteStringF(s, "{{You close the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s closes the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s closes from the other side.}}::green"+CRLF, ReverseDirection(direction)), []string{})
	}
}

/*
Usage:
  - drop [<quantity>] <item>
  - drop all <item>
  - drop all
*/
func DoDrop(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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
			room.Inventory.Add(item)
			char.Inventory.Remove(item)
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
			room.Inventory.Add(item)
			char.Inventory.Remove(item)
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
		room.Inventory.Add(item)
		char.Inventory.Remove(item)
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
	givenItems := []*ItemInstance{}
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
		char.Inventory.Remove(item)
		recipient.Inventory.Add(item)
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
	pickedItems := []*ItemInstance{}
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
		room.Inventory.Remove(item)
		char.Inventory.Add(item)
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

/*
DoEquip will equip an item from the character's inventory

If there are multiple item matches it will prompt the user to select one via a menu.
If the item has no EquipSlots it will not be equipped and an error will be shown.
If the item has multiple EquipSlots it will attempt to use the listed slots in order.
If the Equipment slots are already occupied it will not be equipped and an error will be shown.
If the item can be equiped it will be removed from the inventory and added to the Equipment map and a message will be shown.
The EquipItem function should be used to contain the logic for equipping an item so it can be reused.

Usage:
  - equip <item> [slot]
*/
// func DoEquip(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
// 	if len(args) == 0 {
// 		WriteString(s, "{{Usage: equip <item name> [index]}}::yellow"+CRLF)
// 		return
// 	}

// 	// Check if the last argument is an index number.
// 	var indexProvided bool
// 	var selectedIndex int
// 	searchArgs := args
// 	lastArg := args[len(args)-1]
// 	if i, err := strconv.Atoi(lastArg); err == nil {
// 		indexProvided = true
// 		selectedIndex = i
// 		searchArgs = args[:len(args)-1]
// 	}
// 	searchTerm := strings.Join(searchArgs, " ")

// 	// Get all inventory items that partially match the search term.
// 	matches := char.Inventory.Search(searchTerm)
// 	if len(matches) == 0 {
// 		WriteStringF(s, "{{No items found matching '%s' in your inventory.}}::red"+CRLF, searchTerm)
// 		return
// 	}

// 	if len(matches) > 1 {
// 		// TODO: handle one item
// 	}

// 	var options []MenuOption
// 	for _, m := range matches {
// 		options = append(options, MenuOption{
// 			DisplayText: fmt.Sprintf("%s - %s ", m.Blueprint.Name, m.InstanceID),
// 			Value:       m.InstanceID,
// 		})
// 	}

// 	chosen, err := PromptForMenu(s, "Please select a item:", options)
// 	if err != nil {
// 		WriteString(s, "{{Error receiving input.}}::red"+CRLF)
// 		return
// 	}

// 	selectedItem := char.Inventory.FindItemByID(chosen)
// 	if selectedItem == nil {
// 		WriteString(s, fmt.Sprintf("{{No item found with ID '%s'.}}::red"+CRLF, chosen))
// 		return
// 	}

// 	// If multiple matches exist and no index was provided, list them.
// 	if len(matches) > 1 && !indexProvided {
// 		var builder strings.Builder
// 		builder.WriteString(fmt.Sprintf("{{Multiple items found matching '%s':}}::yellow"+CRLF, searchTerm))
// 		for i, item := range matches {
// 			bp := EntityMgr.GetItemBlueprintByInstance(item)
// 			builder.WriteString(fmt.Sprintf("  %d) %s"+CRLF, i+1, bp.Name))
// 		}
// 		builder.WriteString("{{Please re-run the command with the desired index.}}::yellow" + CRLF)
// 		WriteString(s, builder.String())
// 		return
// 	}

// 	// TODO: support multiple matches via the menu prompt

// 	// Select the appropriate item.
// 	var chosenItem *ItemInstance
// 	if len(matches) == 1 {
// 		chosenItem = matches[0]
// 	} else {
// 		if selectedIndex < 1 || selectedIndex > len(matches) {
// 			WriteStringF(s, "{{Invalid index. There are %d matching items for '%s'.}}::red"+CRLF, len(matches), searchTerm)
// 			return
// 		}
// 		chosenItem = matches[selectedIndex-1]
// 	}

// 	bp := EntityMgr.GetItemBlueprintByInstance(chosenItem)
// 	if bp == nil {
// 		WriteString(s, "{{Item blueprint not found.}}::red"+CRLF)
// 		return
// 	}

// 	// Ensure the item is equippable.
// 	if len(bp.EquipSlots) == 0 || bp.EquipSlots[0] == EquipSlotNone {
// 		WriteString(s, "{{That item cannot be equipped.}}::red"+CRLF)
// 		return
// 	}
// 	slot := bp.EquipSlots[0]

// 	// Check if the slot is already occupied.
// 	if equippedItem, exists := char.Equipment[slot]; exists && equippedItem != nil {
// 		equippedBP := EntityMgr.GetItemBlueprintByInstance(equippedItem)
// 		WriteStringF(s, "{{The %s slot is already occupied by %s. Please unequip it first.}}::yellow"+CRLF, slot, equippedBP.Name)
// 		return
// 	}

// 	// Remove the item from inventory and equip it.
// 	char.Inventory.RemoveItem(chosenItem)
// 	char.Equipment[slot] = chosenItem
// 	char.Save()
// 	WriteStringF(s, "{{You have equipped %s in the %s slot.}}::green"+CRLF, bp.Name, slot)
// }

// func EquipItem(char *Character, item *ItemInstance) {
// 	if item == nil {
// 		return
// 	}

// 	if len(item.Blueprint.EquipSlots) == 0 || item.Blueprint.EquipSlots[0] == EquipSlotNone {
// 		return
// 	}

// 	slot := item.Blueprint.EquipSlots[0]
// 	char.Equipment[slot] = item
// }

func EquipItem(char *Character, item *ItemInstance) error {
	// bp := EntityMgr.GetItemBlueprintByInstance(item)
	// if bp == nil {
	// 	return fmt.Errorf("item blueprint not found")
	// }
	// Check if the item is equippable.
	if len(item.Blueprint.EquipSlots) == 0 || item.Blueprint.EquipSlots[0] == EquipSlotNone {
		return fmt.Errorf("that item cannot be equipped")
	}

	// Use the first valid slot.
	slot := item.Blueprint.EquipSlots[0]

	// Check if the slot is already occupied.
	if equippedItem, exists := char.Equipment.Slots[slot]; exists && equippedItem != nil {
		// equippedBP := EntityMgr.GetItemBlueprintByInstance(equippedItem)
		return fmt.Errorf("the %s slot is already occupied by %s; please unequip it first", slot, equippedItem.Blueprint.Name)
	}
	// char.Inventory.RemoveItem(item)
	// i :=     inventory.RemoveItem(item)
	// if
	// Remove the item from inventory.

	if item := char.Inventory.Remove(item); item != nil {
		char.Equipment.Equip(slot, item)

		if err := char.Save(); err != nil {
			return fmt.Errorf("failed to save your changes: %w", err)
		}
		return nil
	}
	return fmt.Errorf("failed to remove the item from your inventory")
}

func DoEquip(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) == 0 {
		WriteString(s, "{{Usage: equip <item name> [index]}}::yellow"+CRLF)
		return
	}

	// Check if the last argument is an index number.
	var indexProvided bool
	var selectedIndex int
	searchArgs := args
	lastArg := args[len(args)-1]
	if i, err := strconv.Atoi(lastArg); err == nil {
		indexProvided = true
		selectedIndex = i
		searchArgs = args[:len(args)-1]
	}
	searchTerm := strings.Join(searchArgs, " ")

	// Search the actor's inventory for matching items.
	matches := char.Inventory.Search(searchTerm)
	if len(matches) == 0 {
		WriteStringF(s, "{{No items found matching '%s' in your inventory.}}::red"+CRLF, searchTerm)
		return
	}

	// If multiple items match and no index was provided, prompt the user.
	if len(matches) > 1 && !indexProvided {
		var options []MenuOption
		for _, m := range matches {
			options = append(options, MenuOption{
				DisplayText: fmt.Sprintf("%s - %s", m.Blueprint.Name, m.InstanceID),
				Value:       m.InstanceID,
			})
		}
		chosen, err := PromptForMenu(s, "Please select an item:", options)
		if err != nil {
			WriteString(s, "{{Error receiving input.}}::red"+CRLF)
			return
		}

		selectedItem := char.Inventory.FindItemByID(chosen)
		if selectedItem == nil {
			WriteString(s, fmt.Sprintf("{{No item found with ID '%s'.}}::red"+CRLF, chosen))
			return
		}
		if err := EquipItem(char, selectedItem); err != nil {
			WriteStringF(s, "{{%s}}::red"+CRLF, err.Error())
			return
		}
		bp := EntityMgr.GetItemBlueprintByInstance(selectedItem)
		WriteStringF(s, "{{You have equipped %s in the %s slot.}}::green"+CRLF, bp.Name, bp.EquipSlots[0])
		return
	}

	// If only one match or an index was provided.
	var chosenItem *ItemInstance
	if len(matches) == 1 {
		chosenItem = matches[0]
	} else {
		if selectedIndex < 1 || selectedIndex > len(matches) {
			WriteStringF(s, "{{Invalid index. There are %d matching items for '%s'.}}::red"+CRLF, len(matches), searchTerm)
			return
		}
		chosenItem = matches[selectedIndex-1]
	}

	// Check if the item is equippable.
	// bp := EntityMgr.GetItemBlueprintByInstance(chosenItem)
	// if bp == nil {
	// WriteString(s, "{{Item blueprint not found.}}::red"+CRLF)
	// return
	// }
	if len(chosenItem.Blueprint.EquipSlots) == 0 || chosenItem.Blueprint.EquipSlots[0] == EquipSlotNone {
		WriteString(s, "{{That item cannot be equipped.}}::red"+CRLF)
		return
	}

	// Now, try equipping the item using the shared helper.
	if err := EquipItem(char, chosenItem); err != nil {
		WriteStringF(s, "{{%s}}::red"+CRLF, err.Error())
		return
	}
	WriteStringF(s, "{{You have equipped %s in the %s slot.}}::green"+CRLF, chosenItem.Blueprint.Name, chosenItem.Blueprint.EquipSlots[0])
}

func DoUnequip(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) == 0 {
		WriteString(s, "{{Usage: unequip <item name or slot> [index]}}::yellow"+CRLF)
		return
	}

	// If the argument exactly matches a valid equip slot, unequip from that slot.
	slotArg := strings.ToLower(args[0])
	validSlots := map[string]bool{
		"head": true, "body": true, "hands": true, "legs": true,
	}
	if validSlots[slotArg] {
		if item, exists := char.Equipment.Slots[slotArg]; exists && item != nil {
			delete(char.Equipment.Slots, slotArg)
			char.Inventory.Add(item)
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			WriteStringF(s, "{{You have unequipped %s from the %s slot.}}::green"+CRLF, bp.Name, slotArg)
		} else {
			WriteStringF(s, "{{The %s slot is already empty.}}::yellow"+CRLF, slotArg)
		}
		return
	}

	// Otherwise, treat the input as a partial match for equipped items.
	var indexProvided bool
	var selectedIndex int
	searchArgs := args
	lastArg := args[len(args)-1]
	if i, err := strconv.Atoi(lastArg); err == nil {
		indexProvided = true
		selectedIndex = i
		searchArgs = args[:len(args)-1]
	}
	searchTerm := strings.Join(searchArgs, " ")

	// Search through equipped items for matches.
	type equippedMatch struct {
		slot string
		item *ItemInstance
	}
	var matches []equippedMatch
	for slot, item := range char.Equipment.Slots {
		if item == nil {
			continue
		}
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp == nil {
			continue
		}
		if strings.Contains(strings.ToLower(bp.Name), strings.ToLower(searchTerm)) {
			matches = append(matches, equippedMatch{slot, item})
		}
	}

	if len(matches) == 0 {
		WriteStringF(s, "{{No equipped items found matching '%s'.}}::red"+CRLF, searchTerm)
		return
	}

	if len(matches) > 1 && !indexProvided {
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("{{Multiple equipped items found matching '%s':}}::yellow"+CRLF, searchTerm))
		for i, m := range matches {
			bp := EntityMgr.GetItemBlueprintByInstance(m.item)
			builder.WriteString(fmt.Sprintf("  %d) [%s slot] %s"+CRLF, i+1, m.slot, bp.Name))
		}
		builder.WriteString("{{Please re-run the command with the desired index.}}::yellow" + CRLF)
		WriteString(s, builder.String())
		return
	}

	var chosenMatch equippedMatch
	if len(matches) == 1 {
		chosenMatch = matches[0]
	} else {
		if selectedIndex < 1 || selectedIndex > len(matches) {
			WriteStringF(s, "{{Invalid index. There are %d matching equipped items for '%s'.}}::red"+CRLF, len(matches), searchTerm)
			return
		}
		chosenMatch = matches[selectedIndex-1]
	}

	// Unequip the chosen item.
	item := char.Equipment.Unequip(chosenMatch.slot)
	// delete(char.Equipment, chosenMatch.slot)
	char.Inventory.Add(item)
	char.Save()
	bp := EntityMgr.GetItemBlueprintByInstance(chosenMatch.item)
	WriteStringF(s, "{{You have unequipped %s from the %s slot.}}::green"+CRLF, bp.Name, chosenMatch.slot)
}
