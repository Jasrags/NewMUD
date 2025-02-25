package game

import (
	"fmt"
	"strconv"
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

// findMobsByName searches the current room's mobs and returns all instances
// that match the provided name (case-insensitive).
func findMobsByName(room *Room, name string) []*Mob {
	var matches []*Mob
	room.RLock()
	defer room.RUnlock()

	for _, mob := range room.Mobs {
		if strings.EqualFold(mob.Name, name) {
			matches = append(matches, mob)
		}
	}
	return matches
}

// RenderMobTable builds a formatted table of a mob's stats.
// It leverages the embedded GameEntity fields from Mob.
func RenderMobTable(mob *Mob) string {
	// Optionally, recalculate attributes if needed.
	mob.Recalculate()

	var builder strings.Builder

	// Header: basic details from GameEntity.
	builder.WriteString(cfmt.Sprintf("{{Name:}}::white|bold {{%s}}::cyan"+CRLF, mob.Name))
	builder.WriteString(cfmt.Sprintf("{{ID:}}::white|bold {{%s}}::cyan"+CRLF, mob.ID))
	builder.WriteString(cfmt.Sprintf("{{Title:}}::white|bold {{%s}}::cyan"+CRLF, mob.Title))
	builder.WriteString(cfmt.Sprintf("{{Description:}}::white|bold {{%s}}::cyan"+CRLF, mob.Description))
	builder.WriteString(cfmt.Sprintf("{{Long Description:}}::white|bold {{%s}}::cyan"+CRLF, mob.LongDescription))
	builder.WriteString(CRLF)

	// Mob-specific data.
	builder.WriteString(cfmt.Sprintf("{{Professional Rating:}}::white|bold {{%d}}::cyan"+CRLF, mob.ProfessionalRating))
	builder.WriteString(cfmt.Sprintf("{{General Disposition:}}::white|bold {{%s}}::cyan"+CRLF, mob.GeneralDisposition))
	builder.WriteString(CRLF)

	// Attributes from the embedded GameEntity.
	builder.WriteString(cfmt.Sprintf("{{Body:}}::white|bold {{%d}}::cyan"+CRLF, mob.Body.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Agility:}}::white|bold {{%d}}::cyan"+CRLF, mob.Agility.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Reaction:}}::white|bold {{%d}}::cyan"+CRLF, mob.Reaction.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Strength:}}::white|bold {{%d}}::cyan"+CRLF, mob.Strength.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Willpower:}}::white|bold {{%d}}::cyan"+CRLF, mob.Willpower.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Logic:}}::white|bold {{%d}}::cyan"+CRLF, mob.Logic.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Intuition:}}::white|bold {{%d}}::cyan"+CRLF, mob.Intuition.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Charisma:}}::white|bold {{%d}}::cyan"+CRLF, mob.Charisma.TotalValue))
	builder.WriteString(cfmt.Sprintf("{{Essence:}}::white|bold {{%.1f}}::cyan"+CRLF, mob.Essence.TotalValue))
	if mob.Magic.Base > 0 {
		builder.WriteString(cfmt.Sprintf("{{Magic:}}::white|bold {{%d}}::cyan"+CRLF, mob.Magic.TotalValue))
	}
	if mob.Resonance.Base > 0 {
		builder.WriteString(cfmt.Sprintf("{{Resonance:}}::white|bold {{%d}}::cyan"+CRLF, mob.Resonance.TotalValue))
	}

	return builder.String()
}

// DoMobStats is an admin-only command that displays the stats for a specific mob
// in the current room. Usage: mobstats <mob_name> [index]
// If multiple mobs match the given name and no index is provided,
// a list is shown so the admin can re-run the command with an index.
func DoMobStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
	// We require at least one argument (the mob name).
	if len(args) == 0 {
		WriteString(s, "{{Usage: mobstats <mob_name> [index]}}::yellow"+CRLF)
		return
	}

	var mobName string
	indexProvided := false
	var mobIndex int

	// If more than one argument is provided, try to parse the last one as an index.
	if len(args) > 1 {
		if i, err := strconv.Atoi(args[len(args)-1]); err == nil {
			indexProvided = true
			mobIndex = i
			// The mob name is everything except the last argument.
			mobName = strings.Join(args[:len(args)-1], " ")
		} else {
			// Otherwise, treat all arguments as part of the mob name.
			mobName = strings.Join(args, " ")
		}
	} else {
		mobName = args[0]
	}

	// Find all matching mobs in the current room.
	matches := findMobsByName(room, mobName)
	if len(matches) == 0 {
		WriteString(s, fmt.Sprintf("{{No mob found matching '%s' in this room.}}::red"+CRLF, mobName))
		return
	}

	if !indexProvided {
		// If exactly one match exists, display its stats.
		if len(matches) == 1 {
			WriteString(s, RenderMobTable(matches[0]))
			WriteString(s, CRLF)
			return
		}

		// Multiple matches found; list them and instruct the admin.
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("{{Multiple mobs found matching '%s':}}::yellow"+CRLF, mobName))
		for i, m := range matches {
			// Provide a brief summary for each mob (e.g. Name and Title).
			builder.WriteString(fmt.Sprintf("  %d) %s - %s"+CRLF, i+1, m.Name, m.Title))
		}
		builder.WriteString("{{Please re-run the command with the desired index.}}::yellow" + CRLF)
		WriteString(s, builder.String())
		return
	}

	// If an index was provided, validate it.
	if mobIndex < 1 || mobIndex > len(matches) {
		WriteString(s, fmt.Sprintf("{{Invalid mob index. There are %d mobs matching '%s'.}}::red"+CRLF, len(matches), mobName))
		return
	}

	// Display the stats for the chosen mob.
	chosenMob := matches[mobIndex-1]
	WriteString(s, RenderMobTable(chosenMob))
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

// hasAnyTag returns true if filterTags is empty or if any tag in filterTags is found (case-insensitive)
// in the object's tags.
func hasAnyTag(objectTags []string, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, ft := range filterTags {
		for _, ot := range objectTags {
			if strings.EqualFold(ot, ft) {
				return true
			}
		}
	}
	return false
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
			if hasAnyTag(r.Tags, filterTags) {
				outputBuilder.WriteString(cfmt.Sprintf("ID: {{%-35s}}::white|bold  Title: {{%-35s}}::white|bold  Tags: {{%s}}::white|bold"+CRLF, r.ID, r.Title, strings.Join(r.Tags, ", ")))
			}
		}
		if outputBuilder.Len() == 0 {
			outputBuilder.WriteString(cfmt.Sprintf("{{No rooms found with the specified tags.}}::red" + CRLF))
		}
	case "items", "i":
		items := EntityMgr.GetAllItemBlueprints()
		for _, it := range items {
			if hasAnyTag(it.Tags, filterTags) {
				outputBuilder.WriteString(cfmt.Sprintf("ID: {{%-20s}}::white|bold  Name: {{%-20s}}::white|bold  Tags: {{%s}}::white|bold"+CRLF, it.ID, it.Name, strings.Join(it.Tags, ", ")))
			}
		}
		if outputBuilder.Len() == 0 {
			outputBuilder.WriteString(cfmt.Sprintf("{{No items found with the specified tags.}}::red" + CRLF))
		}
	case "mobs", "m":
		mobs := EntityMgr.GetAllMobs()
		for _, m := range mobs {
			if hasAnyTag(m.Tags, filterTags) {
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
