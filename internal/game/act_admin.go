package game

import (
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

func DoSpawn(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if room == nil {
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

	if len(args) < 2 {
		WriteString(s, "{{Usage: spawn <item|mob> <name>}}::yellow"+CRLF)
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
			WriteStringF(s, "{{Error: No item blueprint named '%s' found.}}::red"+CRLF, entityName)
			return
		}

		char.Inventory.AddItem(i)
		// room.Inventory.AddItem(i)
		WriteStringF(s, "{{You spawn a %s.}}::green"+CRLF, bp.Name)
		room.Broadcast(cfmt.Sprintf("{{%s spawns a %s.}}::green"+CRLF, char.Name, bp.Name), []string{char.ID})

	case "m":
		// Spawn a mob into the room
		mob := NewMob()
		mob.Name = entityName

		room.AddMob(mob)
		WriteStringF(s, "{{You spawn a mob named %s.}}::green"+CRLF, entityName)
		room.Broadcast(cfmt.Sprintf("{{%s spawns a mob named %s.}}::green"+CRLF, char.Name, entityName), []string{char.ID})

	default:
		WriteString(s, "{{Invalid entity type. Usage: spawn <item|mob> <name>}}::yellow"+CRLF)
	}
}
