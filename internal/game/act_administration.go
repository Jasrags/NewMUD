package game

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

func DoSpawn(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 2 {
		WriteString(s, "{{Usage: spawn <item|mob> <name>}}::yellow"+CRLF)
		return
	}

	entityType := args[0]
	entityName := strings.Join(args[1:], " ")

	switch entityType {
	case "i":
		// Spawn an item into the character inventory
		item := EntityMgr.CreateItemInstanceFromBlueprintID(entityName)
		if item == nil {
			WriteStringF(s, "{{Error: No item blueprint named '%s' found.}}::red"+CRLF, entityName)
			return
		}

		char.Inventory.Add(item)

		WriteStringF(s, "{{You spawn a %s.}}::green"+CRLF, item.Blueprint.Name)
		room.Broadcast(cfmt.Sprintf("{{%s spawns a %s.}}::green"+CRLF, char.Name, item.Blueprint.Name), []string{char.ID})
		char.Save()

	case "m":
		// Spawn a mob into the room
		mob := EntityMgr.CreateMobInstanceFromBlueprintID(entityName)
		if mob == nil {
			WriteStringF(s, "{{Error: No mob blueprint named '%s' found.}}::red"+CRLF, entityName)
			return
		}

		room.AddMobInstance(mob)

		WriteStringF(s, "{{You spawn a mob named %s.}}::green"+CRLF, entityName)
		room.Broadcast(cfmt.Sprintf("{{%s spawns a mob named %s.}}::green"+CRLF, char.Name, entityName), []string{char.ID})

	default:
		WriteString(s, "{{Invalid entity type. Usage: spawn <item|mob> <name>}}::yellow"+CRLF)
	}
}

func DoMobStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
	if len(args) == 0 {
		WriteString(s, "{{Usage: mobstats <mob_name>}}::yellow"+CRLF)
		return
	}

	// Join arguments to form the search term
	mobName := strings.Join(args, " ")

	// Find mobs with partial name matching
	matches := room.FindMobsByPartialName(mobName)

	if len(matches) == 0 {
		WriteString(s, fmt.Sprintf("{{No mob found matching '%s' in this room.}}::red"+CRLF, mobName))
		return
	}

	// If exactly one match is found, show stats immediately
	if len(matches) == 1 {
		WriteString(s, RenderMobTable(matches[0]))
		WriteString(s, CRLF)
		return
	}

	// If we have multiple options prompt to select one
	var options []MenuOption
	for _, m := range matches {
		options = append(options, MenuOption{
			DisplayText: fmt.Sprintf("%s - %s ", m.Blueprint.Name, m.InstanceID),
			Value:       m.InstanceID,
		})
	}

	chosen, err := PromptForMenu(s, "Please select a mob:", options)
	if err != nil {
		WriteString(s, "{{Error receiving input.}}::red"+CRLF)
		return
	}

	selectedMob := room.FindMobByInstanceID(chosen)
	if selectedMob == nil {
		WriteString(s, fmt.Sprintf("{{No mob found with ID '%s'.}}::red"+CRLF, chosen))
		return
	}

	WriteString(s, RenderMobTable(selectedMob))
	WriteString(s, CRLF)
}

func DoGoto(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	// Ensure an argument is provided
	if len(args) == 0 {
		WriteString(s, "{{Usage: goto <room_id|character_name>}}::yellow"+CRLF)
		return
	}

	target := args[0]

	// Check if the target is a room ID
	newRoom := EntityMgr.GetRoom(target)
	if newRoom != nil {
		if char.Room == newRoom {
			WriteString(s, "{{You are already in that room.}}::yellow"+CRLF)
			return
		}

		char.MoveToRoom(newRoom)
		WriteStringF(s, "{{You teleport to room '%s'.}}::green"+CRLF, newRoom.ID)
		WriteString(s, RenderRoom(user, char, nil))
		return
	}

	// Check if the target is a character
	targetChar := CharacterMgr.GetCharacterByName(target)
	if targetChar != nil && targetChar.Room != nil {
		if char.Room == targetChar.Room {
			WriteString(s, "{{You are already in the same room as that character.}}::yellow"+CRLF)
			return
		}

		char.MoveToRoom(targetChar.Room)
		WriteStringF(s, "{{You teleport to %s.}}::green"+CRLF, targetChar.Name)
		WriteString(s, RenderRoom(user, char, nil))
		return
	}

	// If no valid target is found, return an error
	WriteStringF(s, "{{No room or character found matching '%s'.}}::red"+CRLF, target)
}

// DoList implements the admin "list" command.
// Usage: list <mobs|items|rooms> [tags]
// e.g. "list rooms shopping", "list items weapon,gun", "list mobs thug"
func DoList(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	// Validate arguments.
	if len(args) < 1 {
		WriteString(s, cfmt.Sprintf("{{Usage: list <mobs|items|rooms> [tags]}}::yellow"+CRLF))
		return
	}

	category := strings.ToLower(args[0])
	var filterTags []string
	if len(args) > 1 {
		// Split the tags by comma and trim any extra spaces.
		filterTags = strings.Split(args[1], ",")
		for i, tag := range filterTags {
			filterTags[i] = strings.TrimSpace(tag)
		}
	}

	var outputBuilder strings.Builder

	switch category {
	case "rooms", "r":
		rooms := EntityMgr.GetAllRooms()
		for _, r := range rooms {
			if HasAnyTag(r.Tags, filterTags) {
				outputBuilder.WriteString(cfmt.Sprintf("ID: {{%-35s}}::white|bold  Title: {{%-35s}}::white|bold  Tags: {{%s}}::white|bold"+CRLF, r.ID, r.Title, strings.Join(r.Tags, ", ")))
			}
		}
		if outputBuilder.Len() == 0 {
			outputBuilder.WriteString(cfmt.Sprintf("{{No rooms found with the specified tags.}}::red" + CRLF))
		}
	case "items", "i":
		items := EntityMgr.GetAllItemBlueprints()
		for _, it := range items {
			if HasAnyTag(it.Tags, filterTags) {
				outputBuilder.WriteString(cfmt.Sprintf("ID: {{%-20s}}::white|bold  Name: {{%-20s}}::white|bold  Tags: {{%s}}::white|bold"+CRLF, it.ID, it.Name, strings.Join(it.Tags, ", ")))
			}
		}
		if outputBuilder.Len() == 0 {
			outputBuilder.WriteString(cfmt.Sprintf("{{No items found with the specified tags.}}::red" + CRLF))
		}
	case "mobs", "m":
		mobs := EntityMgr.GetAllMobBlueprints()
		for _, m := range mobs {
			if HasAnyTag(m.Tags, filterTags) {
				outputBuilder.WriteString(cfmt.Sprintf("ID: {{%-20s}}::white|bold  Name: {{%-20s}}::white|bold  Tags: {{%s}}::white|bold"+CRLF, m.ID, m.Name, strings.Join(m.Tags, ", ")))
			}
		}
		if outputBuilder.Len() == 0 {
			outputBuilder.WriteString(cfmt.Sprintf("{{No mobs found with the specified tags.}}::red" + CRLF))
		}
	default:
		WriteString(s, cfmt.Sprintf("{{Unknown category '%s'. Valid categories are: mobs, items, rooms.}}::red"+CRLF, args[0]))
		return
	}

	WriteString(s, outputBuilder.String())
}
