package main

import (
	"io"
	"log/slog"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

/*
Usage:
  - stats
*/
func DoStats(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if char == nil {
		io.WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red\n"))
		return
	}

	io.WriteString(s, cfmt.Sprintf("{{Your current stats:}}::cyan\n"))

	attributes := char.Attributes

	// Display all attributes
	format := "{{%-15s}}::white|bold {{Base: %3d}}::green {{Delta: %+3d}}::yellow {{Total: %3d}}::cyan\n"
	io.WriteString(s, cfmt.Sprintf(format, "Body", attributes.Body.Base, attributes.Body.Delta, attributes.Body.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Agility", attributes.Agility.Base, attributes.Agility.Delta, attributes.Agility.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Reaction", attributes.Reaction.Base, attributes.Reaction.Delta, attributes.Reaction.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Strength", attributes.Strength.Base, attributes.Strength.Delta, attributes.Strength.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Willpower", attributes.Willpower.Base, attributes.Willpower.Delta, attributes.Willpower.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Logic", attributes.Logic.Base, attributes.Logic.Delta, attributes.Logic.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Intuition", attributes.Intuition.Base, attributes.Intuition.Delta, attributes.Intuition.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Charisma", attributes.Charisma.Base, attributes.Charisma.Delta, attributes.Charisma.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Edge", attributes.Edge.Base, attributes.Edge.Delta, attributes.Edge.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Essence", int(attributes.Essence.Base), int(attributes.Essence.Delta), int(attributes.Essence.TotalValue)))
	io.WriteString(s, cfmt.Sprintf(format, "Magic", attributes.Magic.Base, attributes.Magic.Delta, attributes.Magic.TotalValue))
	io.WriteString(s, cfmt.Sprintf(format, "Resonance", attributes.Resonance.Base, attributes.Resonance.Delta, attributes.Resonance.TotalValue))
}

/*
Usage:
  - look
  - look [at] <item|character|direction|mob>
*/
// TODO: This needs work still but it's functional
func DoLook(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if room == nil {
		slog.Error("Character is not in a room",
			slog.String("character_id", char.ID))

		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if len(args) == 0 {
		// No arguments: Look at the room
		io.WriteString(s, RenderRoom(user, char, nil))
		return
	}

	target := strings.Join(args, " ")

	// Check if the target is an item in the room
	if item := room.Inventory.FindItemByName(target); item != nil {
		io.WriteString(s, RenderItemDescription(item))
		return
	}

	// Check if the target is a mob in the room
	if mob := room.FindMobByName(target); mob != nil {
		io.WriteString(s, RenderMobDescription(mob))
		return
	}

	// Check if the target is another character in the room
	if targetChar := room.FindCharacterByName(target); targetChar != nil {
		io.WriteString(s, RenderCharacterDescription(targetChar))
		return
	}

	// Check if the target is a direction
	if room.HasExit(target) {
		io.WriteString(s, RenderExitDescription(target))
		return
	}

	// Target not found
	io.WriteString(s, cfmt.Sprintf("{{You see nothing special about '%s'.}}::yellow\n", target))
}

func SuggestLook(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	if room == nil {
		return suggestions
	}

	switch len(args) {
	case 0:
		// Suggest "at" keyword
		suggestions = append(suggestions, "at")
	case 1:
		if strings.EqualFold(args[0], "at") {
			// Suggest items, mobs, characters, and directions
			for _, i := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(i)
				suggestions = append(suggestions, bp.Name)
			}
			for _, m := range room.Mobs {
				suggestions = append(suggestions, m.Name)
			}
			for _, c := range room.Characters {
				if c.ID != char.ID { // Exclude the player themselves
					suggestions = append(suggestions, char.Name)
				}
			}
			for _, e := range room.Exits {
				suggestions = append(suggestions, e.Direction)
			}
		} else {
			// Suggest items, mobs, characters, and directions directly
			for _, i := range room.Inventory.Items {
				bp := EntityMgr.GetItemBlueprintByInstance(i)
				if strings.HasPrefix(strings.ToLower(bp.Name), strings.ToLower(args[0])) {
					suggestions = append(suggestions, bp.Name)
				}
			}
			for _, m := range room.Mobs {
				if strings.HasPrefix(strings.ToLower(m.Name), strings.ToLower(args[0])) {
					suggestions = append(suggestions, m.Name)
				}
			}
			for _, c := range room.Characters {
				if c.ID != char.ID && strings.HasPrefix(strings.ToLower(char.Name), strings.ToLower(args[0])) {
					suggestions = append(suggestions, char.Name)
				}
			}
			for _, e := range room.Exits {
				if strings.HasPrefix(strings.ToLower(e.Direction), strings.ToLower(args[0])) {
					suggestions = append(suggestions, e.Direction)
				}
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
func DoHelp(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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

func DoInventory(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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

/*
Usage:
  - who
*/
// TODO: Sort all admins to the top of the list
// TODO: Add a CanSee function for characters and have this function use that to determine if a character can see another character in the who list
func DoWho(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
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
  - time
  - time details
*/
func DoTime(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	switch len(args) {
	case 0:
		// Basic time display
		io.WriteString(s, cfmt.Sprintf("{{The current in-game time is %s.}}::cyan\n", gameTime.String()))
	case 1:
		if strings.EqualFold(args[0], "details") {
			// Detailed time information
			hour := gameTime.CurrentHour()
			minute := gameTime.CurrentMinute()
			timeUntilSunrise := calculateTimeUntil(6) // Example sunrise time
			timeUntilSunset := calculateTimeUntil(18) // Example sunset time

			io.WriteString(s, cfmt.Sprintf("{{Current in-game time: %02d:%02d}}::cyan\n", hour, minute))
			io.WriteString(s, cfmt.Sprintf("{{Time until sunrise:}}::green %s\n", formatMinutesAsTime(timeUntilSunrise)))
			io.WriteString(s, cfmt.Sprintf("{{Time until sunset:}}::yellow %s\n", formatMinutesAsTime(timeUntilSunset)))
		} else {
			io.WriteString(s, cfmt.Sprintf("{{Unknown argument '%s'. Usage: time [details]}}::red\n", args[0]))
		}
	default:
		io.WriteString(s, cfmt.Sprintf("{{Invalid usage. Usage: time [details]}}::red\n"))
	}
}
