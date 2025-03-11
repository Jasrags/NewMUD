package game

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

/*
Usage:
  - stats
*/
func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
	// If arguments are provided, assume the user is requesting stats for another character.
	if len(args) > 0 {
		// // Only admins can view other characters' stats.
		// if char.Role != CharacterRoleAdmin {
		//     WriteString(s, "{{You are not authorized to view other characters' stats.}}::red"+CRLF)
		//     return
		// }

		// Join the args to form the target character name.
		targetName := args[0]
		targetChar := CharacterMgr.GetCharacterByName(targetName)
		if targetChar == nil {
			WriteString(s, fmt.Sprintf("{{Character '%s' not found.}}::red"+CRLF, targetName))
			return
		}

		char = targetChar

		// // Render and display the stats for the target character.
		// WriteString(s, RenderCharacterTable(targetChar))
		// WriteString(s, CRLF)
		// return
	}

	// No target specified; display the current character's stats.
	WriteString(s, RenderCharacterTable(char))
	WriteString(s, CRLF)
}

/*
Usage:
  - look
  - look [at] <item|character|direction|mob>
*/
// TODO: This needs work still but it's functional
func DoLook(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) == 0 {
		// No arguments: Look at the room
		WriteString(s, RenderRoom(user, char, nil))
		return
	}

	target := strings.Join(args, " ")

	// Check if the target is an item in the room
	if item := room.Inventory.FindItemByName(target); item != nil {
		WriteString(s, RenderItemDescription(item))
		return
	}

	// Check if the target is a mob in the room
	if mob := room.FindMobByName(target); mob != nil {
		WriteString(s, RenderMobDescription(mob))
		return
	}

	// Check if the target is another character in the room
	if targetChar := room.FindCharacterByName(target); targetChar != nil {
		WriteString(s, RenderCharacterDescription(targetChar))
		return
	}

	// Check if the target is a direction
	if room.HasExit(target) {
		WriteString(s, RenderExitDescription(target))
		return
	}

	// Target not found
	WriteString(s, cfmt.Sprintf("{{You see nothing special about '%s'.}}::yellow"+CRLF, target))
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
			for _, m := range room.MobInstances {
				suggestions = append(suggestions, m.Blueprint.Name)
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
			for _, m := range room.MobInstances {
				if strings.HasPrefix(strings.ToLower(m.Blueprint.Name), strings.ToLower(args[0])) {
					suggestions = append(suggestions, m.Blueprint.Name)
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
	// Retrieve all registered commands
	commands := CommandMgr.GetCommands()

	// Check if specific command help is requested
	if len(args) == 1 {
		commandName := args[0]
		command, exists := commands[commandName]
		if !exists || (!CanSeeCommand(char, command)) {
			WriteStringF(s, "{{Unknown command '%s'. Type 'help' for a list of commands.}}::red"+CRLF, commandName)
			return
		}

		var builder strings.Builder
		builder.WriteString(cfmt.Sprintf("{{%s}}::cyan"+CRLF, strings.ToUpper(command.Name)))
		builder.WriteString(cfmt.Sprintf("{{Description:}}::white|bold %s"+CRLF, command.Description))
		if len(command.Aliases) > 0 {
			builder.WriteString(cfmt.Sprintf("{{Aliases:}}::white|bold %s"+CRLF, strings.Join(command.Aliases, ", ")))
		}
		builder.WriteString("{{Usage:}}::white|bold" + CRLF)
		for _, usage := range command.Usage {
			builder.WriteString(cfmt.Sprintf("  - {{%s}}::green"+CRLF, usage))
		}
		WriteString(s, builder.String())
		return
	}

	// Organize commands by category while preventing duplicate entries
	categorizedCommands := make(map[CommandCategory][]*Command)
	displayedCommands := make(map[string]bool) // Tracks displayed commands

	for _, command := range commands {
		// Skip commands that the character should not see
		if !CanSeeCommand(char, command) {
			continue
		}

		// Ensure we only add the primary command (not aliases) to the output
		if _, exists := displayedCommands[command.Name]; exists {
			continue // Skip duplicates
		}
		displayedCommands[command.Name] = true // Mark as displayed

		category := command.CommandCategory
		categorizedCommands[category] = append(categorizedCommands[category], command)
	}

	// Generate the help output grouped by category
	var builder strings.Builder
	builder.WriteString("{{Available Commands:}}::white|bold" + CRLF)

	for category, cmdList := range categorizedCommands {
		// Category header
		builder.WriteString(cfmt.Sprintf("\n{{%s}}::yellow|bold"+CRLF, category))

		// List commands under the category
		for _, cmd := range cmdList {
			aliases := ""
			if len(cmd.Aliases) > 0 {
				aliases = fmt.Sprintf(" (aliases: %s)", strings.Join(cmd.Aliases, ", "))
			}
			builder.WriteString(cfmt.Sprintf("  {{%-10s}}::cyan - %s%s"+CRLF, cmd.Name, cmd.Description, aliases))
		}
	}

	WriteString(s, builder.String())
}

func DoInventory(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(char.Inventory.Items) == 0 {
		WriteString(s, cfmt.Sprintf("{{You are not carrying anything.}}::yellow"+CRLF))
		return
	}

	WriteString(s, cfmt.Sprintf("{{You are carrying:}}::cyan"+CRLF))
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
		WriteString(s, cfmt.Sprintf("- %s"+CRLF,
			pluralizer.PluralizeNounPhrase(name, count)))
	}
}

/*
Usage:
  - equipment
*/
func DoEquipment(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	// Display equipped items for each supported slot.
	WriteString(s, cfmt.Sprintf("{{Equipped Items:}}::cyan"+CRLF))
	for _, slot := range EquipSlots {
		if item, exists := char.Equipment.Slots[slot]; exists && item != nil {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp != nil {
				WriteStringF(s, "{{%-6s}}::cyan %s"+CRLF, slot+":", bp.Name)
			} else {
				WriteStringF(s, "{{%-6s}}::cyan {{<unknown>}}::red"+CRLF, slot+":")
			}
		} else {
			WriteStringF(s, "{{%-6s}}::cyan <empty>"+CRLF, slot+":")
		}
	}
}

/*
Usage:
  - who
*/
// TODO: Sort all admins to the top of the list
// TODO: Add a CanSee function for characters and have this function use that to determine if a character can see another character in the who list
func DoWho(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	activeCharacters := CharacterMgr.GetOnlineCharacters()

	if len(activeCharacters) == 0 {
		WriteString(s, cfmt.Sprintf("{{No one else is in the game right now.}}::yellow"+CRLF))
		return
	}

	WriteString(s, cfmt.Sprintf("{{Players currently in the game:}}::green"+CRLF))

	for _, activeChar := range activeCharacters {
		color := "cyan"
		if activeChar.Role == CharacterRoleAdmin {
			color = "yellow"
		}

		if activeChar.Title == "" {
			activeChar.Title = "the Basic"
		}

		// Display character title and name
		WriteString(s, cfmt.Sprintf("{{%s - %s}}::%s"+CRLF, activeChar.Name, activeChar.Title, color))
	}
}

func DoPrompt(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	// If no arguments, display current prompt
	if len(args) == 0 {
		WriteStringF(s, "{{Your current prompt:}}::cyan %s"+CRLF, char.Prompt)
		WriteString(s, "{{Use 'prompt <new format>' to set a custom prompt.}}::yellow"+CRLF)
		WriteString(s, "{{Available Macros:}}::green {{time}}, {{hp}}, {{gold}}, {{stamina}} "+CRLF)
		return
	}

	// Collect user input
	newPrompt := strings.Join(args, " ")

	// Validate prompt
	if !ValidatePrompt(newPrompt) {
		placeholders := []string{}
		for n := range promptPlaceholders {
			placeholders = append(placeholders, n)
		}
		WriteString(s, "{{Invalid prompt format! Please use only supported macros.}}::red"+CRLF)
		WriteStringF(s, "{{Available Macros:}}::green %s"+CRLF, strings.Join(placeholders, ", "))
		return
	}

	// Save new prompt
	char.Prompt = newPrompt
	char.Save()

	WriteStringF(s, "{{Prompt updated successfully! New prompt:}}::green %s"+CRLF, newPrompt)
}

// ValidatePrompt ensures that only allowed macros (from promptPlaceholders) are used
func ValidatePrompt(prompt string) bool {
	re := regexp.MustCompile(`{{[^{}]+}}`)
	matches := re.FindAllString(prompt, -1)

	for _, match := range matches {
		if _, exists := promptPlaceholders[match]; !exists {
			return false
		}
	}
	return true
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
		WriteStringF(s, cfmt.Sprintf("{{The current in-game time is %s.}}::cyan"+CRLF, GameTimeMgr.GetFormattedTime()))
	case 1:
		if strings.EqualFold(args[0], "details") {
			// Detailed time information
			hour := GameTimeMgr.CurrentHour()
			minute := GameTimeMgr.CurrentMinute()
			timeUntilSunrise := calculateTimeUntil(6) // Example sunrise time
			timeUntilSunset := calculateTimeUntil(18) // Example sunset time

			WriteStringF(s, "{{Current in-game time: %02d:%02d}}::cyan"+CRLF, hour, minute)
			WriteStringF(s, "{{Time until sunrise:}}::green %s"+CRLF, formatMinutesAsTime(timeUntilSunrise))
			WriteStringF(s, "{{Time until sunset:}}::yellow %s"+CRLF, formatMinutesAsTime(timeUntilSunset))
		} else {
			WriteStringF(s, "{{Unknown argument '%s'. Usage: time [details]}}::red"+CRLF, args[0])
		}
	default:
		WriteString(s, "{{Invalid usage. Usage: time [details]}}::red"+CRLF)
	}
}

func DoHistory(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(char.CommandHistory) == 0 {
		WriteString(s, "{{No command history available.}}::yellow"+CRLF)
		return
	}

	if len(args) > 0 {
		search := strings.ToLower(args[0])
		WriteStringF(s, "{{Search results for '%s':}}::green"+CRLF, search)
		found := false
		for i, entry := range char.CommandHistory {
			if strings.Contains(strings.ToLower(entry), search) {
				WriteStringF(s, "{{%d: %s}}::cyan"+CRLF, i+1, entry)
				found = true
			}
		}
		if !found {
			WriteStringF(s, "{{No history entries found for '%s'.}}::red"+CRLF, search)
		}
		return
	}

	WriteString(s, "{{Command history:}}::green"+CRLF)
	for i, entry := range char.CommandHistory {
		WriteStringF(s, "{{%d: %s}}::cyan"+CRLF, i+1, entry)
	}
}
