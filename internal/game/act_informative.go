package game

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/charmbracelet/lipgloss"
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
// 		io.WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red\n"))
// 		return
// 	}

// 	io.WriteString(s, cfmt.Sprintf("{{Your current stats:}}::cyan\n"))

// 	attributes := char.Attributes
// 	attributes.Recalculate()

// 	// Helper function to format attributes
// 	formatAttribute := func(name string, attr Attribute[int]) string {
// 		if attr.TotalValue > attr.Base {
// 			return cfmt.Sprintf("{{%-20s}}::white|bold {{%3d}}::cyan{{(}}::white {{%d}}::red{{)}}::white\n", name, attr.Base, attr.TotalValue)
// 		}
// 		return cfmt.Sprintf("{{%-20s}}::white|bold {{%3d}}::cyan\n", name, attr.Base)
// 	}
// 	// Handle float attributes like Essence separately
// 	formatFloatAttribute := func(name string, attr Attribute[float64]) string {
// 		if attr.TotalValue > attr.Base {
// 			return cfmt.Sprintf("{{%-20s}}::white|bold {{%.1f}}::cyan {{(}}::white{{%.1f}}::red{{)}}::white\n", name, attr.Base, attr.TotalValue)
// 		}
// 		return cfmt.Sprintf("{{%-20s}}::white|bold {{%.1f}}::cyan\n", name, attr.Base)
// 	}

// 	io.WriteString(s, formatAttribute("Body", attributes.Body))
// 	io.WriteString(s, formatAttribute("Agility", attributes.Agility))
// 	io.WriteString(s, formatAttribute("Reaction", attributes.Reaction))
// 	io.WriteString(s, formatAttribute("Strength", attributes.Strength))
// 	io.WriteString(s, formatAttribute("Willpower", attributes.Willpower))
// 	io.WriteString(s, formatAttribute("Logic", attributes.Logic))
// 	io.WriteString(s, formatAttribute("Intuition", attributes.Intuition))
// 	io.WriteString(s, formatAttribute("Charisma", attributes.Charisma))
// 	io.WriteString(s, formatAttribute("Edge", attributes.Edge))
// 	io.WriteString(s, formatFloatAttribute("Essence", attributes.Essence))
// 	if attributes.Magic.Base > 0 {
// 		io.WriteString(s, formatAttribute("Magic", attributes.Magic))
// 	}
// 	if attributes.Resonance.Base > 0 {
// 		io.WriteString(s, formatAttribute("Resonance", attributes.Resonance))
// 	}

// 	// Carry weight stats
// 	maxCarryWeight := char.GetLiftCarry()
// 	currentCarryWeight := char.GetCurrentCarryWeight()
// 	io.WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%.2f}}::cyan{{/}}::white{{%d}}::cyan{{kg}}::white\n", "Carry Weight", currentCarryWeight, maxCarryWeight))

// 	// New limits
// 	physicalLimit := char.GetPhysicalLimit()
// 	adjustedPhysicalLimit := char.GetAdjustedPhysicalLimit()
// 	mentalLimit := char.GetMentalLimit()
// 	socialLimit := char.GetSocialLimit()

// 	if adjustedPhysicalLimit < physicalLimit {
// 		io.WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan {{(Adjusted: %d)}}::yellow\n", "Physical Limit", physicalLimit, adjustedPhysicalLimit))
// 	} else {
// 		io.WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan\n", "Physical Limit", physicalLimit))
// 	}

// 	io.WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan\n", "Mental Limit", mentalLimit))
// 	io.WriteString(s, cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan\n", "Social Limit", socialLimit))
// }

// func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
// 	if char == nil {
// 		io.WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red\n"))
// 		return
// 	}

// 	var output strings.Builder

// 	// Character Info Block
// 	output.WriteString(cfmt.Sprintf("Name: {{%-15s}}::cyan Title: {{%s}}::cyan\n", char.Name, char.Title))
// 	output.WriteString(cfmt.Sprintf("Metatype: {{%-12s}}::cyan Ethnicity: {{%s}}::cyan\n", char.Metatype, char.Ethnicity))
// 	output.WriteString(cfmt.Sprintf("Age: {{%-4d}}::cyan Sex: {{%-6s}}::cyan Height: {{%-6s}}::cyan Weight: {{%s}}::cyan\n",
// 		char.Age, char.Sex, char.Height, char.Weight))
// 	output.WriteString(cfmt.Sprintf("Street Cred: {{%-3d}}::cyan Notoriety: {{%-3d}}::cyan Public Awareness: {{%d}}::cyan\n",
// 		char.StreetCred, char.Notoriety, char.PublicAwareness))
// 	output.WriteString(cfmt.Sprintf("Karma: {{%-10d}}::cyan Total Karma: {{%d}}::cyan\n\n", char.Karma, char.TotalKarma))

// 	// Damage and Condition Tracking
// 	output.WriteString(cfmt.Sprintf("Physical Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Stun Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Overflow: {{%d}}::cyan\n\n",
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
// 		output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8s}}::cyan {{%-20s}}::white|bold {{%s}}::cyan\n",
// 			stat.LeftName, stat.LeftValue, stat.RightName, stat.RightValue))
// 	}

// 	// Limits at the Bottom
// 	physicalLimit := char.GetPhysicalLimit()
// 	adjustedPhysicalLimit := char.GetAdjustedPhysicalLimit()
// 	mentalLimit := char.GetMentalLimit()
// 	socialLimit := char.GetSocialLimit()

// 	output.WriteString(cfmt.Sprintf("\nPhysical Limit: {{%d (%d)}}::cyan  Mental Limit: {{%d}}::cyan  Social Limit: {{%d}}::cyan\n",
// 		physicalLimit, adjustedPhysicalLimit, mentalLimit, socialLimit))

// 	// Send output to session
// 	io.WriteString(s, output.String())
// }

// func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
// 	if char == nil {
// 		io.WriteString(s, cfmt.Sprintf("{{Error: No character is associated with this session.}}::red\n"))
// 		return
// 	}

// 	var builder strings.Builder

// 	// Character Info Block
// 	builder.WriteString(cfmt.Sprintf("Name: {{%-15s}}::cyan Title: {{%-15s}}::cyan\n", char.Name, char.Title))
// 	builder.WriteString(cfmt.Sprintf("Metatype: {{%-12s}}::cyan Ethnicity: {{%s}}::cyan\n", char.Metatype, char.Ethnicity))
// 	builder.WriteString(cfmt.Sprintf("Age: {{%-4d}}::cyan Sex: {{%-6s}}::cyan Height: {{%-6s}}::cyan Weight: {{%s}}::cyan\n",
// 		char.Age, char.Sex, char.Height, char.Weight))
// 	builder.WriteString(cfmt.Sprintf("Street Cred: {{%-3d}}::cyan Notoriety: {{%-3d}}::cyan Public Awareness: {{%d}}::cyan\n",
// 		char.StreetCred, char.Notoriety, char.PublicAwareness))
// 	builder.WriteString(cfmt.Sprintf("Karma: {{%-10d}}::cyan Total Karma: {{%d}}::cyan\n\n", char.Karma, char.TotalKarma))

// 	// Damage and Condition Tracking
// 	builder.WriteString(cfmt.Sprintf("Physical Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Stun Damage:    {{%2d}}::cyan/{{%2d}}::cyan     Overflow: {{%d}}::cyan\n\n",
// 		char.PhysicalDamage.Current, char.PhysicalDamage.Max, char.StunDamage.Current, char.StunDamage.Max, char.PhysicalDamage.Overflow))

// 	// Two-column Main Stats Block
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%.2f}}::cyan\n",
// 		"Body", char.Attributes.Body.TotalValue, "Essence", char.Attributes.Essence.TotalValue))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d (%d)}}::cyan\n",
// 		"Agility", char.Attributes.Agility.TotalValue, "Magic/Resonance", char.Attributes.Magic.TotalValue, char.Attributes.Resonance.TotalValue))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n",
// 		"Reaction", char.Attributes.Reaction.TotalValue, "Initiative", char.Initiative.Base))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n",
// 		"Strength", char.Attributes.Strength.TotalValue, "Matrix Initiative", char.MatrixInitiative.Base))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n",
// 		"Willpower", char.Attributes.Willpower.TotalValue, "Composure", char.GetComposure()))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n",
// 		"Logic", char.Attributes.Logic.TotalValue, "Judge Intentions", char.GetJudgeIntentions()))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n",
// 		"Intuition", char.Attributes.Intuition.TotalValue, "Memory", char.GetMemory()))
// 	builder.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%.2fkg/%.2fkg}}::cyan\n",
// 		"Edge", char.EdgePoints, "Lift/Carry", char.GetCurrentCarryWeight(), char.GetLiftCarry()))

// 	// Limits at the Bottom
// 	builder.WriteString("\n")
// 	builder.WriteString(cfmt.Sprintf("Physical Limit: {{%d (%d)}}::cyan  Mental Limit: {{%d (%d)}}::cyan  Social Limit: {{%d (%d)}}::cyan\n",
// 		char.GetPhysicalLimit(), char.GetAdjustedPhysicalLimit(),
// 		char.GetMentalLimit(), char.GetMentalLimitAdjusted(),
// 		char.GetSocialLimit(), char.GetSocialLimitAdjusted()))

// 	// Write everything to the session
// 	io.WriteString(s, builder.String())
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

func RenderCharacterTable(char *Character) string {
	char.Recalculate()
	table := lipgloss.JoinVertical(lipgloss.Left,
		// Personal Data
		headerStyle.Render("Personal Data"),
		singleColumnStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Top,
					RenderKeyValue("Name", char.Name), "\t",
					RenderKeyValue("Title", char.Title),
				),
				lipgloss.JoinHorizontal(lipgloss.Top,
					RenderKeyValue("Metatype", char.Metatype), "\t",
					RenderKeyValue("Ethnicity", char.Ethnicity),
				),
				lipgloss.JoinHorizontal(lipgloss.Top,
					RenderKeyValue("Age", "0"), "\t",
					RenderKeyValue("Sex", char.Sex), "\t",
					RenderKeyValue("Height", "0"), "\t",
					RenderKeyValue("Weight", "0"),
				),
				lipgloss.JoinHorizontal(lipgloss.Top,
					RenderKeyValue("Street Cred", "0"), "\t",
					RenderKeyValue("Notoriety", "0"), "\t",
					RenderKeyValue("Public Awareness", "0"),
				),
				lipgloss.JoinHorizontal(lipgloss.Top,
					RenderKeyValue("Karma", "0"), "\t",
					RenderKeyValue("Total Karma", "0"),
				),
			),
		),
		// Attributes doble column
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render("Attributes"),
				// Attributes - LEFT - Base attributes
				// Formats:
				// Reaction   5  (7)
				// Essence    6.00
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						RenderAttribute(char.Attributes.Body),      // 5  (7)
						RenderAttribute(char.Attributes.Agility),   // 5  (7)
						RenderAttribute(char.Attributes.Reaction),  // 5  (7)
						RenderAttribute(char.Attributes.Strength),  // 5  (7)
						RenderAttribute(char.Attributes.Willpower), // 5  (7)
						RenderAttribute(char.Attributes.Logic),     // 5  (7)
						RenderAttribute(char.Attributes.Intuition), // 5  (7)
						RenderAttribute(char.Attributes.Charisma),  // 5  (7)
						RenderAttribute(char.Attributes.Essence),   // 5  (7)
						RenderAttribute(char.Attributes.Magic),     // 5  (7)
						RenderAttribute(char.Attributes.Resonance), // Essence    6.00
						// strs...,
					),
				),
			),
			// Attributes RIGHT - Derivied attributes
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render(""),
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						RenderAttribute(char.Attributes.Initiative), // Initiative 10 (12) + 1d6 (2d6)
						RenderAttribute(char.Attributes.InitiativeDice),
						RenderAttribute(char.Attributes.Composure),       // 5  (7)
						RenderAttribute(char.Attributes.JudgeIntentions), // 5  (7)
						RenderAttribute(char.Attributes.Memory),          // 5  (7)
						RenderAttribute(char.Attributes.Lift),            // 5  (7)
						RenderAttribute(char.Attributes.Carry),           // 5  (7)
						RenderAttribute(char.Attributes.Walk),            // 5  (7)
						RenderAttribute(char.Attributes.Run),             // 5  (7)
						RenderAttribute(char.Attributes.Swim),            // 5  (7)
						"",
					),
				),
			),
		),
	)

	return table
}

// RenderAttribute renders a single attribute for display.
func RenderAttribute[T int | float64](attr Attribute[T]) string {
	var output strings.Builder

	if attr.Base == 0 {
		return ""
	}

	output.WriteString(attrNameStyle.Render(fmt.Sprintf("%-10s", attr.Name)))
	output.WriteString(attrValueStyle.Render(fmt.Sprintf(" %-2v", renderValue(attr.Base))))
	if attr.TotalValue != attr.Base {
		style := attrPosModStyle
		if attr.TotalValue < attr.Base {
			style = attrNegModStyle
		}
		output.WriteString(style.Render(fmt.Sprintf(" (%v)", renderValue(attr.TotalValue))))
	}

	return output.String()
}

// renderValue formats the value of an attribute for display.
func renderValue[T int | float64](value T) string {
	switch v := any(value).(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
func RenderKeyValue(key, value string) string {
	return fmt.Sprintf("%s: %s", attrNameStyle.Render(key), attrTextValueStyle.Render(value))
}

func DoStats(s ssh.Session, cmd string, args []string, acct *Account, char *Character, room *Room) {
	if char == nil {
		io.WriteString(s, cfmt.Sprint("{{Error: No character is associated with this session.}}::red\n"))
		return
	}

	io.WriteString(s, RenderCharacterTable(char))
	io.WriteString(s, "\n")
}

// return RenderCharacterTable(char)

// )
// // Single Column
// fmt.Printf("%-20s %d\n", "Movement", 5)

// // Double Column
// fmt.Printf("%-20s %-8d %-20s %.2f\n", "Body", 5, "Essence", 5.60)

// // Triple Column
// fmt.Printf("%-20s %-8s %-20s %-8s %-15s %d\n",
// 	"Physical Damage", "0/11",
// 	"Stun Damage", "0/10",
// 	"Overflow", 0)

// var output strings.Builder

// Character Info Block
// output.WriteString(FormatDoubleColumn("Name:", char.Name, "Title:", char.Title))
// output.WriteString(cfmt.Sprintf(
// 	"{{%-20s}}::white|bold {{%-8s}}::cyan {{%-20s}}::white|bold {{%-8s}}::cyan\n",
// 	"Name:", char.Name, "Title:", char.Title))
// output.WriteString(cfmt.Sprintf(
// 	"Metatype: {{%-12s}}::cyan Ethnicity: {{%s}}::cyan\n",
// 	char.Metatype, char.Ethnicity))
// output.WriteString(cfmt.Sprintf(
// 	"Age: {{%-4d}}::cyan Sex: {{%-6s}}::cyan Height: {{%-6d}}::cyan Weight: {{%d}}::cyan\n",
// 	char.Age, char.Sex, char.Height, char.Weight))
// output.WriteString(cfmt.Sprintf("Street Cred: {{%-3d}}::cyan Notoriety: {{%-3d}}::cyan Public Awareness: {{%d}}::cyan\n",
// 	char.StreetCred, char.Notoriety, char.PublicAwareness))
// output.WriteString(cfmt.Sprintf("Karma: {{%-10d}}::cyan Total Karma: {{%d}}::cyan\n\n", char.Karma, char.TotalKarma))

// // Damage and Condition Tracking
// output.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold {{%6d}}::cyan/{{%-6d}}::cyan {{%-15s}}::white|bold {{%6d}}::cyan/{{%-6d}}::cyan {{%-15s}}::white|bold {{%d}}::cyan\n\n",
// 	"Physical Damage", char.PhysicalDamage.Current, char.GetPhysicalConditionMax(), "Stun Damage", char.StunDamage.Current, char.GetStunConditionMax(), "Overflow", char.PhysicalDamage.Overflow))

// // Two-column Main Stats Block
// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%s}}::cyan\n", "Body", char.Attributes.Body.TotalValue, "Essence", fmt.Sprintf("%.2f", char.Attributes.Essence.TotalValue)))
// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d (%d)}}::cyan\n", "Agility", char.Attributes.Agility.TotalValue, "Magic/Resonance", char.Attributes.Magic.TotalValue, char.Attributes.Resonance.TotalValue))
// // output.WriteString( cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n", "Reaction", char.Attributes.Reaction.TotalValue, "Initiative", char.Initiative.Base))
// // output.WriteString( cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n", "Strength", char.Attributes.Strength.TotalValue, "Matrix Initiative", char.MatrixInitiative.Base))
// // output.WriteString( cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n", "Willpower", char.Attributes.Willpower.TotalValue, "Astral Initiative", char.AstralInitiative.Base))
// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n", "Logic", char.Attributes.Logic.TotalValue, "Composure", char.GetComposure()))
// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n", "Intuition", char.Attributes.Intuition.TotalValue, "Judge Intentions", char.GetJudgeIntentions()))
// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n", "Charisma", char.Attributes.Charisma.TotalValue, "Memory", char.GetMemory()))
// output.WriteString(cfmt.Sprintf("{{%-15s}}::white|bold {{%6d}}::cyan/{{%-6d}}::cyan {{%-20s}}::white|bold {{%.2fkg/%.2fkg}}::cyan\n", "Edge", char.Edge.Available, char.Edge.Max, "Lift/Carry", char.GetCurrentCarryWeight(), char.GetLiftCarry()))
// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan\n\n", "Movement", char.GetMovement()))

// // Limits at the Bottom
// physicalLimit := char.GetPhysicalLimit()
// adjustedPhysicalLimit := char.GetAdjustedPhysicalLimit()
// mentalLimit := char.GetMentalLimit()
// socialLimit := char.GetSocialLimit()

// output.WriteString(cfmt.Sprintf("Physical Limit: {{%d (%d)}}::cyan  Mental Limit: {{%d (%d)}}::cyan  Social Limit: {{%d (%d)}}::cyan\n",
// 	physicalLimit, adjustedPhysicalLimit, mentalLimit, char.GetMentalLimit(), socialLimit, char.GetSocialLimit()))

// 	io.WriteString(s, output.String())
// }

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
		io.WriteString(s, cfmt.Sprintf("{{The current in-game time is %s.}}::cyan\n", GameTimeMgr.GetFormattedTime()))
	case 1:
		if strings.EqualFold(args[0], "details") {
			// Detailed time information
			hour := GameTimeMgr.CurrentHour()
			minute := GameTimeMgr.CurrentMinute()
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

func DoHistory(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if char == nil || len(char.CommandHistory) == 0 {
		io.WriteString(s, "{{No command history available.}}::yellow\n")
		return
	}

	if len(args) > 0 {
		search := strings.ToLower(args[0])
		io.WriteString(s, cfmt.Sprintf("{{Search results for '%s':}}::green\n", search))
		found := false
		for i, entry := range char.CommandHistory {
			if strings.Contains(strings.ToLower(entry), search) {
				io.WriteString(s, cfmt.Sprintf("{{%d: %s}}::cyan\n", i+1, entry))
				found = true
			}
		}
		if !found {
			io.WriteString(s, cfmt.Sprintf("{{No history entries found for '%s'.}}::red\n", search))
		}
		return
	}

	io.WriteString(s, "{{Command history:}}::green\n")
	for i, entry := range char.CommandHistory {
		io.WriteString(s, cfmt.Sprintf("{{%d: %s}}::cyan\n", i+1, entry))
	}
}
