package main

import (
	"io"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

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
		mob := NewMob()
		mob.Name = entityName
		// mob := &Mob{
		// ID:   uuid.New().String(),
		// Name: entityName,
		// }

		room.AddMob(mob)
		io.WriteString(s, cfmt.Sprintf("{{You spawn a mob named %s.}}::green\n", entityName))
		room.Broadcast(cfmt.Sprintf("{{%s spawns a mob named %s.}}::green\n", char.Name, entityName), []string{char.ID})

	default:
		io.WriteString(s, cfmt.Sprintf("{{Invalid entity type. Usage: spawn <item|mob> <name>}}::yellow\n"))
	}
}
