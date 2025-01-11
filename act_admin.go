package main

import (
	"io"
	"strconv"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/i582/cfmt/cmd/cfmt"
)

func DoSpawn(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) < 2 {
		io.WriteString(s, cfmt.Sprintf("{{Usage: spawn <item|mob> [<quantity>] <id>}}::yellow\n"))
		return
	}

	entityType := args[0]
	quantity := 1
	entityID := ""

	// Parse optional quantity argument
	if len(args) == 3 {
		if q, err := strconv.Atoi(args[1]); err == nil && q > 0 {
			quantity = q
			entityID = args[2]
		} else {
			io.WriteString(s, cfmt.Sprintf("{{Invalid quantity.}}::yellow\n"))
			return
		}
	} else {
		entityID = args[1]
	}

	switch entityType {
	case "item", "i":
		// Spawn items
		bp := EntityMgr.GetItemBlueprintByID(entityID)
		if bp == nil {
			io.WriteString(s, cfmt.Sprintf("{{Error: No item blueprint named '%s' found.}}::red\n", entityID))
			return
		}

		for i := 0; i < quantity; i++ {
			instance := EntityMgr.CreateItemInstanceFromBlueprint(bp)
			if instance == nil {
				io.WriteString(s, cfmt.Sprintf("{{Error creating item instance for blueprint '%s'.}}::red\n", entityID))
				return
			}
			char.Inventory.AddItem(instance)
		}

		char.Save()
		io.WriteString(s, cfmt.Sprintf("{{You spawn %d %s.}}::green\n", quantity, pluralizer.PluralizeNoun(bp.Name, quantity)))
		room.Broadcast(cfmt.Sprintf("{{%s spawns %d %s.}}::green\n", char.Name, quantity, pluralizer.PluralizeNoun(bp.Name, quantity)), []string{char.ID})

	case "mob", "m":
		// Spawn mobs
		for i := 0; i < quantity; i++ {
			mob := &Mob{
				ID:   uuid.New().String(),
				Name: entityID,
			}
			room.AddMob(mob)
		}

		io.WriteString(s, cfmt.Sprintf("{{You spawn %d mob(s) named '%s'.}}::green\n", quantity, entityID))
		room.Broadcast(cfmt.Sprintf("{{%s spawns %d mob(s) named '%s'.}}::green\n", char.Name, quantity, entityID), []string{char.ID})

	default:
		io.WriteString(s, cfmt.Sprintf("{{Invalid entity type. Usage: spawn <item|mob> [<quantity>] <id>}}::yellow\n"))
	}
}

func SuggestSpawn(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest entity types
		suggestions = append(suggestions, "item", "mob")
	case 1: // Suggest IDs for items or mobs
		if strings.EqualFold(args[0], "item") || strings.EqualFold(args[0], "i") {
			for id := range EntityMgr.items {
				suggestions = append(suggestions, id)
			}
		} else if strings.EqualFold(args[0], "mob") || strings.EqualFold(args[0], "m") {
			// Optionally suggest predefined mob templates if available
		}
	case 2: // Suggest quantity or IDs
		if _, err := strconv.Atoi(args[1]); err == nil {
			if strings.EqualFold(args[0], "item") || strings.EqualFold(args[0], "i") {
				for id := range EntityMgr.items {
					suggestions = append(suggestions, id)
				}
			}
		}
	}

	return suggestions
}
