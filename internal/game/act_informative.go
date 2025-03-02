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
// TODO: Change the color of the currenty carry wight when we get closer to max
// func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
// 	if char == nil {
// 		WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red"+CRLF))
// 		return
// 	}

// 	WriteString(s, cfmt.Sprintf("{{Your current stats:}}::cyan"+CRLF))

// 	attributes := char.Attributes
// 	attributes.Recalculate()

// 	// Helper function to format attributes
// 	formatAttribute := func(name string, attr Attribute[int]) string {
// 		if attr.TotalValue > attr.Base {
// 			return cfmt.Sprintf("{{%-20s}}::white|bold {{%3d}}::cyan{{(}}::white {{%d}}::red{{)}}::white"+CRLF, name, attr.Base, attr.TotalValue)
// 		}
// 		return cfmt.Sprintf("{{%-20s}}::white|bold {{%3d}}::cyan"+CRLF, name, attr.Base)
// 	}
// 	// Handle float attributes like Essence separately
// 	formatFloatAttribute := func(name string, attr Attribute[float64]) string {
// 		if attr.TotalValue > attr.Base {
// 			return cfmt.Sprintf("{{%-20s}}::white|bold {{%.1f}}::cyan {{(}}::white{{%.1f}}::red{{)}}::white"+CRLF, name, attr.Base, attr.TotalValue)
// 		}
// 		return cfmt.Sprintf("{{%-20s}}::white|bold {{%.1f}}::cyan"+CRLF, name, attr.Base)
// 	}

// 	WriteString(s, formatAttribute("Body", attributes.Body))
// 	WriteString(s, formatAttribute("Agility", attributes.Agility))
// 	WriteString(s, formatAttribute("Reaction", attributes.Reaction))
// 	WriteString(s, formatAttribute("Strength", attributes.Strength))
// 	WriteString(s, formatAttribute("Willpower", attributes.Willpower))
// 	WriteString(s, formatAttribute("Logic", attributes.Logic))
// 	WriteString(s, formatAttribute("Intuition", attributes.Intuition))
// 	WriteString(s, formatAttribute("Charisma", attributes.Charisma))
// 	WriteString(s, formatAttribute("Edge", attributes.Edge))
// 	WriteString(s, formatFloatAttribute("Essence", attributes.Essence))
// 	if attributes.Magic.Base > 0 {
// 		WriteString(s, formatAttribute("Magic", attributes.Magic))
// 	}
// 	if attributes.Resonance.Base > 0 {
// 		WriteString(s, formatAttribute("Resonance", attributes.Resonance))
// 	}

// 	// Carry weight stats
// 	maxCarryWeight := char.GetLiftCarry()
// 	currentCarryWeight := char.GetCurrentCarryWeight()
// 	WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%.2f}}::cyan{{/}}::white{{%d}}::cyan{{kg}}::white"+CRLF, "Carry Weight", currentCarryWeight, maxCarryWeight))

// 	// New limits
// 	physicalLimit := char.GetPhysicalLimit()
// 	adjustedPhysicalLimit := char.GetAdjustedPhysicalLimit()
// 	mentalLimit := char.GetMentalLimit()
// 	socialLimit := char.GetSocialLimit()

// 	if adjustedPhysicalLimit < physicalLimit {
// 		WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan {{(Adjusted: %d)}}::yellow"+CRLF, "Physical Limit", physicalLimit, adjustedPhysicalLimit))
// 	} else {
// 		WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan"+CRLF, "Physical Limit", physicalLimit))
// 	}

// 	WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan"+CRLF, "Mental Limit", mentalLimit))
// 	WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan"+CRLF, "Social Limit", socialLimit))
// }

// func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
// 	if char == nil {
// 		WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red"+CRLF))
// 		return
// 	}

// 	var output strings.Builder

// 	// Character Info Block
// 	output.WriteString(cfmt.Sprintf("Name: {{%-15s}}::cyan Title: {{%s}}::cyan"+CRLF, char.Name, char.Title))
// 	output.WriteString(cfmt.Sprintf("Metatype: {{%-12s}}::cyan Ethnicity: {{%s}}::cyan"+CRLF, char.Metatype, char.Ethnicity))
// 	output.WriteString(cfmt.Sprintf("Age: {{%-4d}}::cyan Sex: {{%-6s}}::cyan Height: {{%-6s}}::cyan Weight: {{%s}}::cyan"+CRLF,
// 		char.Age, char.Sex, char.Height, char.Weight))
// 	output.WriteString(cfmt.Sprintf("Street Cred: {{%-3d}}::cyan Notoriety: {{%-3d}}::cyan Public Awareness: {{%d}}::cyan"+CRLF,
// 		char.StreetCred, char.Notoriety, char.PublicAwareness))
// 	output.WriteString(cfmt.Sprintf("Karma: {{%-10d}}::cyan Total Karma: {{%d}}::cyan"+CRLF, char.Karma, char.TotalKarma))

// 	// Damage and Condition Tracking
// 	output.WriteString(cfmt.Sprintf("Physical Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Stun Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Overflow: {{%d}}::cyan"+CRLF,
// 		char.PhysicalDamage.Current, char.PhysicalDamage.Max,
// 		char.StunDamage.Current, char.StunDamage.Max,
// 		char.PhysicalDamage.Overflow))

// 	// Two-column Main Stats Block
// 	stats := []struct {
// 		LeftName  string
// 		LeftValue string
// 		RightName string
// 		RightValue string
// 	}{
// 		{"Body", fmt.Sprintf("%d", char.Attributes.Body.TotalValue), "Essence", fmt.Sprintf("%.2f", char.Attributes.Essence.TotalValue)},
// 		{"Agility", fmt.Sprintf("%d", char.Attributes.Agility.TotalValue), "Magic/Resonance", fmt.Sprintf("%d (%d)", char.Attributes.Magic.TotalValue, char.Attributes.Resonance.TotalValue)},
// 		{"Reaction", fmt.Sprintf("%d", char.Attributes.Reaction.TotalValue), "Initiative", fmt.Sprintf("%d + 1d6", char.Initiative.Base)},
// 		{"Strength", fmt.Sprintf("%d", char.Attributes.Strength.TotalValue), "Matrix Initiative", fmt.Sprintf("%d + 1d6", char.MatrixInitiative.Base)},
// 		{"Willpower", fmt.Sprintf("%d", char.Attributes.Willpower.TotalValue), "Astral Initiative", fmt.Sprintf("%d + 1d6", char.AstralInitiative.Base)},
// 		{"Logic", fmt.Sprintf("%d", char.Attributes.Logic.TotalValue), "Composure", fmt.Sprintf("%d", char.GetComposure())},
// 		{"Intuition", fmt.Sprintf("%d", char.Attributes.Intuition.TotalValue), "Judge Intentions", fmt.Sprintf("%d", char.GetJudgeIntentions())},
// 		{"Charisma", fmt.Sprintf("%d", char.Attributes.Charisma.TotalValue), "Memory", fmt.Sprintf("%d", char.GetMemory())},
// 		{"Edge", fmt.Sprintf("%d/%d", char.EdgePoints, char.Attributes.Edge.TotalValue), "Lift/Carry", fmt.Sprintf("%.2fkg/%.2fkg", char.GetCurrentCarryWeight(), char.GetLiftCarry())},
// 		{"Edge Points", fmt.Sprintf("%d", char.EdgePoints), "Movement", fmt.Sprintf("%d", char.GetMovement())},
// 	}

// 	for _, stat := range stats {
// 		output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8s}}::cyan {{%-20s}}::white|bold {{%s}}::cyan"+CRLF,
// 			stat.LeftName, stat.LeftValue, stat.RightName, stat.RightValue))
// 	}

// 	// Limits at the Bottom
// 	physicalLimit := char.GetPhysicalLimit()
// 	adjustedPhysicalLimit := char.GetAdjustedPhysicalLimit()
// 	mentalLimit := char.GetMentalLimit()
// 	socialLimit := char.GetSocialLimit()

// 	output.WriteString(cfmt.Sprintf("\nPhysical Limit: {{%d (%d)}}::cyan  Mental Limit: {{%d}}::cyan  Social Limit: {{%d}}::cyan"+CRLF,
// 		physicalLimit, adjustedPhysicalLimit, mentalLimit, socialLimit))

// 	// Send output to session
// 	WriteString(s, output.String())
// }

// func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
// 	if char == nil {
// 		WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red"+CRLF))
// 		return
// 	}

// 	var builder strings.Builder

// 	// Character Info Block
// 	builder.WriteString(cfmt.Sprintf("Name: {{%-15s}}::cyan Title: {{%-15s}}::cyan"+CRLF, char.Name, char.Title))
// 	builder.WriteString(cfmt.Sprintf("Metatype: {{%-12s}}::cyan Ethnicity: {{%s}}::cyan"+CRLF, char.Metatype, char.Ethnicity))
// 	builder.WriteString(cfmt.Sprintf("Age: {{%-4d}}::cyan Sex: {{%-6s}}::cyan Height: {{%-6s}}::cyan Weight: {{%s}}::cyan"+CRLF,
// 		char.Age, char.Sex, char.Height, char.Weight))
// 	builder.WriteString(cfmt.Sprintf("Street Cred: {{%-3d}}::cyan Notoriety: {{%-3d}}::cyan Public Awareness: {{%d}}::cyan"+CRLF,
// 		char.StreetCred, char.Notoriety, char.PublicAwareness))
// 	builder.WriteString(cfmt.Sprintf("Karma: {{%-10d}}::cyan Total Karma: {{%d}}::cyan"+CRLF, char.Karma, char.TotalKarma))

// 	// Damage and Condition Tracking
// 	builder.WriteString(cfmt.Sprintf("Physical Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Stun Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Overflow: {{%d}}::cyan"+CRLF,
// 		char.PhysicalDamage.Current, char.PhysicalDamage.Max, char.StunDamage.Current, char.StunDamage.Max, char.PhysicalDamage.Overflow))

// 	// Two-column Main Stats Block
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%.2f}}::cyan"+CRLF,
// 		"Body", char.Attributes.Body.TotalValue, "Essence", char.Attributes.Essence.TotalValue))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d (%d)}}::cyan"+CRLF,
// 		"Agility", char.Attributes.Agility.TotalValue, "Magic/Resonance", char.Attributes.Magic.TotalValue, char.Attributes.Resonance.TotalValue))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan"+CRLF,
// 		"Reaction", char.Attributes.Reaction.TotalValue, "Initiative", char.Initiative.Base))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan"+CRLF,
// 		"Strength", char.Attributes.Strength.TotalValue, "Matrix Initiative", char.MatrixInitiative.Base))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan"+CRLF,
// 		"Willpower", char.Attributes.Willpower.TotalValue, "Composure", char.GetComposure()))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan"+CRLF,
// 		"Logic", char.Attributes.Logic.TotalValue, "Judge Intentions", char.GetJudgeIntentions()))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan"+CRLF,
// 		"Intuition", char.Attributes.Intuition.TotalValue, "Memory", char.GetMemory()))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%.2fkg/%.2fkg}}::cyan"+CRLF,
// 		"Edge", char.EdgePoints, "Lift/Carry", char.GetCurrentCarryWeight(), char.GetLiftCarry()))

// 	// Limits at the Bottom
// 	builder.WriteString(""+CRLF)
// 	builder.WriteString(cfmt.Sprintf("Physical Limit: {{%d (%d)}}::cyan  Mental Limit: {{%d (%d)}}::cyan  Social Limit: {{%d (%d)}}::cyan"+CRLF,
// 		char.GetPhysicalLimit(), char.GetAdjustedPhysicalLimit(),
// 		char.GetMentalLimit(), char.GetMentalLimitAdjusted(),
// 		char.GetSocialLimit(), char.GetSocialLimitAdjusted()))

// 	// Write everything to the session
// 	WriteString(s, builder.String())
// }

// // FormatColumn formats a single column with dynamic width and data type.
// func FormatColumn(label string, value interface{}, width int) string {
// 	switch v := value.(type) {
// 	case int:
// 		return cfmt.Sprintf("%-*s %d", width, label, v)
// 	case float64:
// 		return cfmt.Sprintf("%-*s %.2f", width, label, v)
// 	case string:
// 		return cfmt.Sprintf("%-*s %s", width, label, v)
// 	default:
// 		return cfmt.Sprintf("%-*s %v", width, label, v) // Fallback for other types
// 	}
// }

// func FormatSingleColumn(label string, value interface{}) string {
// 	return FormatColumn(label, value, 20)
// }

// func FormatDoubleColumn(label1 string, value1 interface{}, label2 string, value2 interface{}) string {
// 	return cfmt.Sprintf("%-20s %-8v %-20s %-8v",
// 		FormatColumn(label1, value1, 20),
// 		FormatColumn(label2, value2, 20))
// }

// func FormatTripleColumn(label1 string, value1 interface{}, label2 string, value2 interface{}, label3 string, value3 interface{}) string {
// 	return fmt.Sprintf("%-26s %-26s %-26s",
// 		FormatColumn(label1, value1, 20),
// 		FormatColumn(label2, value2, 20),
// 		FormatColumn(label3, value3, 20))
// }

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

// func DoPrompt(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
// 	if char == nil {
// 		WriteString(s, cfmt.Sprint("{{Error: No character is associated with this session.}}::red"+CRLF))
// 		return
// 	}

// 	// If no arguments, display current prompt
// 	if len(args) == 0 {
// 		WriteString(s, cfmt.Sprintf("{{Your current prompt:}}::cyan \"%s\""+CRLF, char.Prompt))
// 		WriteString(s, cfmt.Sprint("{{Use 'prompt <new format>' to set a custom prompt.}}::yellow"+CRLF))
// 		return
// 	}

// 	// Set a new custom prompt
// 	newPrompt := strings.Join(args, " ")
// 	char.Prompt = newPrompt
// 	char.Save()

// 	WriteString(s, cfmt.Sprintf("{{Prompt updated successfully!\nNew prompt:}}::green \"%s\""+CRLF, newPrompt))
// }

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
