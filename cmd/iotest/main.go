package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jasrags/NewMUD/internal/game"
	"github.com/charmbracelet/lipgloss"
)

// const (
// 	ColorBlack = iota
// 	ColorRed
// 	ColorGreen
// 	ColorYellow
// 	ColorBlue
// 	ColorMagenta
// 	ColorCyan
// 	ColorWhite
// 	ColorBrightBlack
// 	ColorBrightRed
// 	ColorBrightGreen
// 	ColorBrightYellow
// 	ColorBrightBlue
// 	ColorBrightMagenta
// 	ColorBrightCyan
// 	ColorBrightWhite
// )

func main() {
	var output strings.Builder
	output.WriteString(RenderCharacterTable())
	// println("Hello, World!")

	// char := &game.Character{
	// 	UserID: "test_user",
	// 	Role:   game.CharacterRolePlayer,
	// 	GameEntity: game.GameEntity{
	// 		Name:            "Street Samurai",
	// 		Title:           "the Uncanny",
	// 		ID:              "ID",
	// 		Metatype:        "Ork",
	// 		Age:             25,
	// 		Sex:             "Male",
	// 		Height:          180,
	// 		Weight:          80,
	// 		Ethnicity:       "White",
	// 		StreetCred:      2,
	// 		Notoriety:       2,
	// 		PublicAwareness: 2,
	// 		Karma:           2,
	// 		TotalKarma:      5,
	// 		Description:     "A street samurai character",
	// 		Attributes:      game.NewAttributes(),
	// 		// Attributes: game.Attributes{
	// 		// 	Body:      game.Attribute[int]{Base: 7},
	// 		// 	Agility:   game.Attribute[int]{Base: 6},
	// 		// 	Reaction:  game.Attribute[int]{Base: 5}, // (7)
	// 		// 	Strength:  game.Attribute[int]{Base: 5},
	// 		// 	Willpower: game.Attribute[int]{Base: 3},
	// 		// 	Logic:     game.Attribute[int]{Base: 2},
	// 		// 	Intuition: game.Attribute[int]{Base: 3},
	// 		// 	Charisma:  game.Attribute[int]{Base: 2},
	// 		// 	Essence:   game.Attribute[float64]{Base: 0.88},
	// 		// 	Magic:     game.Attribute[int]{Base: 0},
	// 		// 	Resonance: game.Attribute[int]{Base: 0},
	// 		// },
	// 		PhysicalDamage: game.PhysicalDamage{
	// 			Current:  0,
	// 			Max:      10,
	// 			Overflow: 0,
	// 		},
	// 		StunDamage: game.StunDamage{
	// 			Current: 0,
	// 			Max:     10,
	// 		},
	// 		Edge: game.Edge{
	// 			Max:       5,
	// 			Available: 5,
	// 		},
	// 		Equipment: map[string]*game.Item{},
	// 	},
	// }

	// // Personal data                Core combat info
	// // 77 38| 25/25/25
	// output.WriteString(cfmt.Sprintf("╭─────────────────────────────────────────────────────────────────────────────╮\n"))
	// output.WriteString(cfmt.Sprintf("│ Name:                                                                       │\n"))
	// output.WriteString(cfmt.Sprintf("│                                                                             │\n"))
	// output.WriteString(cfmt.Sprintf("├--------------------------------------┬--------------------------------------┤\n"))
	// // output.WriteString(cfmt.Sprintf("│ %-36s │ %-36s │\n",
	// // "{{Body}}::white|bold {{3}}::cyan (5)", "{{Essence}}::white|bold 6.0"))
	// output.WriteString(cfmt.Sprintf("│ {{Body}}::white|bold                  {{4}}::cyan ({{6}}::cyan)          │ {{Essence}}::white|bold                  {{4}}::cyan ({{6}}::cyan)       │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Agility}}::white|bold               {{4}}::cyan ({{6}}::cyan)          │                                      │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Reaction}}::white|bold              {{4}}::cyan ({{6}}::cyan)          │ {{Initiative}}::white|bold               {{4}}::cyan ({{6}}::cyan) + {{1d6}}::cyan ({{3d6}}::cyan)       │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Strength}}::white|bold              {{4}}::cyan ({{6}}::cyan)          │ {{Matrix Initiative}}::white|bold        {{4}}::cyan ({{6}}::cyan) + {{1d6}}::cyan ({{3d6}}::cyan)       │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Willpower}}::white|bold             {{4}}::cyan ({{6}}::cyan)          │ {{Astral Initiative}}::white|bold        {{4}}::cyan ({{6}}::cyan) + {{1d6}}::cyan ({{3d6}}::cyan)       │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Logic}}::white|bold                 {{4}}::cyan ({{6}}::cyan)          │ {{Composure}}::white|bold                            │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Intuition}}::white|bold             {{4}}::cyan ({{6}}::cyan)          │ {{Judge Intentions}}::white|bold                     │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Charisma}}::white|bold              {{4}}::cyan ({{6}}::cyan)          │ {{Memory}}::white|bold                               │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Magic/Resonance}}::white|bold       {{4}}::cyan ({{6}}::cyan)          │ {{Lift/Carry}}::white|bold                           │\n"))
	// output.WriteString(cfmt.Sprintf("│                                      │ {{Movement}}::white|bold                             │\n"))
	// output.WriteString(cfmt.Sprintf("│                                      │                                      │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Edge Aval/Max}}::white|bold         {{4}}::cyan/{{5}}::cyan            │                                      │\n"))
	// output.WriteString(cfmt.Sprintf("├-------------------------┬------------┴------------┬-------------------------┤\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Physcial Limit}}::white|bold {{6}}::cyan ({{3}}::red)    │ {{Mental Limit}}::white|bold {{6}}::cyan ({{3}}::red)      │ {{Social Limit}}::white|bold {{6}}::cyan ({{3}}::red)      │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Mental Limit}}::white|bold   {{6}}::cyan ({{3}}::red)    │ {{Mental Limit}}::white|bold {{6}}::cyan ({{3}}::red)      │ {{Social Limit}}::white|bold {{6}}::cyan ({{3}}::red)      │\n"))
	// output.WriteString(cfmt.Sprintf("│ {{Social Limit}}::white|bold   {{6}}::cyan ({{3}}::red)    │ {{Mental Limit}}::white|bold {{6}}::cyan ({{3}}::red)      │ {{Social Limit}}::white|bold {{6}}::cyan ({{3}}::red)      │\n"))
	// output.WriteString(cfmt.Sprintf("│                         │                         │                         │\n"))
	// output.WriteString(cfmt.Sprintf("│                         │                         │                         │\n"))
	// output.WriteString(cfmt.Sprintf("╰─────────────────────────┴─────────────────────────┴─────────────────────────╯\n"))
	// output.WriteString("\n")
	// n := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1, 0, 1).Render
	// v := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1, 0, 1).Render
	// cellStyle := lipgloss.NewStyle().Padding(0, 1, 0, 1).Render

	// t := table.New().Width(80)
	// t.Row(n("Name"))
	// t.Row(n("Body"), n("Essence"))
	// // // t.Row("Bubble Tea", s("Milky"))
	// // // t.Row("Milk Tea", s("Also milky"))
	// // // t.Row("Actual milk", s("Milky as well"))
	// output.WriteString(t.Render())

	// output.WriteString("\n\n")
	// output.WriteString(RenderANSI16Colors())
	// output.WriteString(RenderTable())

	// fmt.Println(t.Render())
	// Attributes                   Conditon monitor

	// Skills                       Qualties

	// IDs, Lifestyle, Currency     Contacts

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

	// Two-column Main Stats Block

	// if char.Attributes.Body.Base < char.Attributes.Body.TotalValue {
	// 	output.WriteString(cfmt.Sprintf(
	// 		"{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%s}}::cyan\n",
	// 		"Body", char.Attributes.Body.TotalValue, "Essence", fmt.Sprintf("%.2f", char.Attributes.Essence.TotalValue)))
	// } else {
	// 	output.WriteString(cfmt.Sprintf(
	// 		"{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%s}}::cyan\n",
	// 		"Body", char.Attributes.Body.TotalValue, "Essence", fmt.Sprintf("%.2f", char.Attributes.Essence.TotalValue)))
	// }
	// output.WriteString(cfmt.Sprintf(
	// 	"{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%s}}::cyan\n",
	// 	"Body", char.Attributes.Body.TotalValue, "Essence", fmt.Sprintf("%.2f", char.Attributes.Essence.TotalValue)))
	// char.Attributes.Body.SetBase(10)
	// char.Attributes.Body.AddDelta(1)

	// char.Attributes.Essence.SetBase(6)
	// char.Attributes.Essence.SubDelta(0.5)

	// char.Attributes.Recalculate()

	// output.WriteString(formatter.MustFormat("{white}{bold}{p0}: {p1} ({p2})", char.Attributes.Essence.Name, char.Attributes.Essence.Base, char.Attributes.Essence.TotalValue))
	// output.WriteString(FormatAttribute(char.Attributes.Body))
	// output.WriteString(FormatAttribute(char.Attributes.Essence))
	// output.WriteString("\n")
	// output.WriteString(FormatAttribute(char.Attributes.Logic))
	// output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan", "Composure", char.GetComposure()))
	// // output.WriteString(FormatAttribute(char.Attributes.Composure))
	// output.WriteString("\n")
	// output.WriteString(cfmt.Sprintf(
	// "{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d (%d)}}::cyan\n",
	// "Agility", char.Attributes.Agility.TotalValue, "Magic/Resonance", char.Attributes.Magic.TotalValue, char.Attributes.Resonance.TotalValue))
	// output.WriteString( cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n", "Reaction", char.Attributes.Reaction.TotalValue, "Initiative", char.Initiative.Base))
	// output.WriteString( cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n", "Strength", char.Attributes.Strength.TotalValue, "Matrix Initiative", char.MatrixInitiative.Base))
	// output.WriteString( cfmt.Sprintf("{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d + 1d6}}::cyan\n", "Willpower", char.Attributes.Willpower.TotalValue, "Astral Initiative", char.AstralInitiative.Base))
	// output.WriteString(cfmt.Sprintf(
	// "{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n",
	// "Logic", char.Attributes.Logic.TotalValue, "Composure", char.GetComposure()))
	// output.WriteString(cfmt.Sprintf(
	// 	"{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n",
	// 	"Intuition", char.Attributes.Intuition.TotalValue, "Judge Intentions", char.GetJudgeIntentions()))
	// output.WriteString(cfmt.Sprintf(
	// 	"{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d}}::cyan\n",
	// 	"Charisma", char.Attributes.Charisma.TotalValue, "Memory", char.GetMemory()))
	// output.WriteString(cfmt.Sprintf(
	// 	"{{%-15s}}::white|bold {{%6d}}::cyan/{{%-6d}}::cyan {{%-20s}}::white|bold {{%.2fkg/%.2fkg}}::cyan\n",
	// 	"Edge", char.Edge.Available, char.Edge.Max, "Lift/Carry", char.GetCurrentCarryWeight(), char.GetLiftCarry()))
	// output.WriteString(cfmt.Sprintf(
	// 	"{{%-20s}}::white|bold {{%d}}::cyan\n\n",
	// 	"Movement", char.GetMovement()))

	// Limits at the Bottom
	// physicalLimit := char.GetPhysicalLimit()
	// adjustedPhysicalLimit := char.GetAdjustedPhysicalLimit()
	// mentalLimit := char.GetMentalLimit()
	// socialLimit := char.GetSocialLimit()

	// output.WriteString(cfmt.Sprintf(
	// 	"Physical Limit: {{%d (%d)}}::cyan  Mental Limit: {{%d (%d)}}::cyan  Social Limit: {{%d (%d)}}::cyan\n",
	// 	physicalLimit, adjustedPhysicalLimit, mentalLimit, char.GetMentalLimit(), socialLimit, char.GetSocialLimit()))

	fmt.Print(output.String())
}

// func FormatAttribute[T int | float64](attribute game.Attribute[T]) string {
// 	if attribute.Base == 0 {
// 		return ""
// 	}
// 	var output strings.Builder
// 	output.WriteString(cfmt.Sprintf("{{%-20s}}::white|bold ", attribute.Name))

// 	// "{{%-20s}}::white|bold {{%-8d}}::cyan {{%-20s}}::white|bold {{%d (%d)}}::cyan\n",

// 	var formattedBaseValue string
// 	switch any(attribute.Base).(type) {
// 	case int:
// 		formattedBaseValue = fmt.Sprintf("%-2d", attribute.TotalValue)
// 	case float64:
// 		formattedBaseValue = fmt.Sprintf("%-2.2f", attribute.TotalValue) // Format float to 2 decimal places
// 	}
// 	output.WriteString(cfmt.Sprintf("{{%s}}::cyan ", formattedBaseValue))

// 	var formattedTotalValue string
// 	switch any(attribute.TotalValue).(type) {
// 	case int:
// 		formattedTotalValue = fmt.Sprintf("%-2d", attribute.TotalValue)
// 	case float64:
// 		formattedTotalValue = fmt.Sprintf("%-2.2f", attribute.TotalValue) // Format float to 2 decimal places
// 	}
// 	output.WriteString(cfmt.Sprintf("{{(%s)}}::cyan ", formattedTotalValue))
// 	// cfmt.Sprint(formattedValue)

// 	// return cfmt.Sprintf("{{%-16s}}::white|bold %v (%v)",
// 	// attribute.Name, attribute.Base, attribute.Delta, formattedValue)
// 	// var formattedValue string
// 	// switch any(attribute.TotalValue).(type) {
// 	// case int:
// 	// 	formattedValue = fmt.Sprintf("%d", attribute.TotalValue)
// 	// case float64:
// 	// 	formattedValue = fmt.Sprintf("%.2f", attribute.TotalValue) // Format float to 2 decimal places
// 	// }

// 	// if attribute.Base < attribute.TotalValue {
// 	// 	return formatter.MustFormat("{white}{bold}{p0} {p1} ({p2})", attribute.Name, attribute.Base, attribute.TotalValue)

// 	// return cfmt.Sprintf("{{%-20s}}::white|bold {{%-2d}}::cyan ({{%-2d}}::cyan)", attribute.Name, attribute.Base, attribute.TotalValue)
// 	// } else {
// 	// return formatter.MustFormat("{white}{bold}{p0} {p1}", attribute.Name, attribute.Base, attribute.TotalValue)
// 	// return cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan", attribute.Name, attribute.Base)
// 	// }

// 	return output.String()
// }

// func FormatAttribute(attribute game.Attribute[int]) string {
// 	if base == 0 {
// 		return ""
// 	}

// 	if base < total {
// 		return cfmt.Sprintf("{{%-20s}}::white|bold {{%-2d}}::cyan ({{%-2d}}::cyan)", label, base, total)
// 	} else {
// 		return cfmt.Sprintf("{{%-20s}}::white|bold {{%d}}::cyan", label, base)
// 	}
// }

func RenderANSI16Colors() string {
	// ANSI-16 colors with their respective codes and names
	colors := []struct {
		code int
		name string
	}{
		{0, "Black"}, {1, "Red"}, {2, "Green"}, {3, "Yellow"},
		{4, "Blue"}, {5, "Magenta"}, {6, "Cyan"}, {7, "White"},
		{8, "Bright Black"}, {9, "Bright Red"}, {10, "Bright Green"},
		{11, "Bright Yellow"}, {12, "Bright Blue"}, {13, "Bright Magenta"},
		{14, "Bright Cyan"}, {15, "Bright White"},
	}

	// Render styles for each color
	var output string
	for _, c := range colors {
		// Foreground style
		fgStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("%d", c.code))).
			Padding(0, 2)

		fgStyleBold := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(fmt.Sprintf("%d", c.code))).
			Padding(0, 2)

		// Background style
		bgStyle := lipgloss.NewStyle().
			Background(lipgloss.Color(fmt.Sprintf("%d", c.code))).
			Padding(0, 2)

		// Combine styles
		output += fmt.Sprintf(
			"%s  %s  %s %s\n",
			fgStyle.Render(fmt.Sprintf("[%d]", c.code)),
			fgStyle.Render(c.name),
			fgStyleBold.Render(c.name),
			bgStyle.Render("   "),
		)
	}

	return output
}

// {0, "Black"}, {1, "Red"}, {2, "Green"}, {3, "Yellow"},
//
//	{4, "Blue"}, {5, "Magenta"}, {6, "Cyan"}, {7, "White"},
//
// │ Name: Shadow Walker                                                            │
var (
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1, 0, 1)

		// Table formatting
	singleColumnStyle = borderStyle.Width(80)
	dualColumnStyle   = borderStyle.Width(39)

	// Text styles
	attrNameStyle   = lipgloss.NewStyle().Bold(true).Foreground(white)
	attrValueStyle  = lipgloss.NewStyle().Bold(false).Foreground(cyan)
	attrPosModStyle = lipgloss.NewStyle().Bold(true).Foreground(green)
	attrNegModStyle = lipgloss.NewStyle().Bold(true).Foreground(red)

	// Colors
	black   = lipgloss.Color("0")
	red     = lipgloss.Color("1")
	green   = lipgloss.Color("2")
	yellow  = lipgloss.Color("3")
	blue    = lipgloss.Color("4")
	magenta = lipgloss.Color("5")
	cyan    = lipgloss.Color("6")
	white   = lipgloss.Color("7")
)

var attrs = game.NewAttributes()

// RenderAttributes renders the attributes of a character.
func RenderAttributes(attrs game.Attributes) []string {
	var strs []string

	strs = append(strs, RenderAttribute(attrs.Body))
	strs = append(strs, RenderAttribute(attrs.Agility))
	strs = append(strs, RenderAttribute(attrs.Reaction))
	strs = append(strs, RenderAttribute(attrs.Strength))
	strs = append(strs, RenderAttribute(attrs.Willpower))
	strs = append(strs, RenderAttribute(attrs.Logic))
	strs = append(strs, RenderAttribute(attrs.Intuition))
	strs = append(strs, RenderAttribute(attrs.Charisma))
	strs = append(strs, RenderAttribute(attrs.Essence))
	strs = append(strs, RenderAttribute(attrs.Magic))
	strs = append(strs, RenderAttribute(attrs.Resonance))

	// Now we want to remove all the empty strings from the slice
	for i, s := range strs {
		if s == "" {
			strs = append(strs[:i], strs[i+1:]...)
		}

	}

	return strs
}

// RenderAttribute renders a single attribute for display.
func RenderAttribute[T int | float64](attr game.Attribute[T]) string {
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

// RenderCharacterTable renders the entire character table.
func RenderCharacterTable() string {
	attrs.Body.SetBase(10)
	attrs.Body.AddDelta(2)
	attrs.Agility.SetBase(4)
	attrs.Reaction.SetBase(4)
	attrs.Strength.SetBase(4)
	attrs.Strength.SubDelta(2)
	attrs.Willpower.SetBase(4)
	attrs.Logic.SetBase(4)
	attrs.Intuition.SetBase(4)
	attrs.Charisma.SetBase(4)
	attrs.Essence.SetBase(6.0)
	attrs.Essence.SubDelta(0.5)
	// attrs.Magic.SetBase(4)
	attrs.Resonance.SetBase(4)
	attrs.Recalculate()

	strs := RenderAttributes(attrs)
	table := lipgloss.JoinVertical(lipgloss.Left,
		// Personal Data
		singleColumnStyle.Render("Name: Shadow Walker"),
		// Metatype, Ethnicity, Age, Sex, Height, Weight, Street Cred, Notoriety, Public Awareness, Karma, Total Karma
		lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left,
				// Attributes left
				// All true attributes
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						strs...,
					),
				),
			),
			// Attributes right
			// Derivied attributes - Iniitiative, Matrix Initiative, Astral Initiative, Composure, Judge Intentions, Memory, Lift/Carry, Movement
			lipgloss.JoinVertical(lipgloss.Left,
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						attrNameStyle.Render("Initiative"),
						attrNameStyle.Render("Matrix Initiative"),
						attrNameStyle.Render("Astral Initiative"),
						attrNameStyle.Render("Composure"),
						attrNameStyle.Render("Judge Intentions"),
						attrNameStyle.Render("Memory"),
						attrNameStyle.Render("Lift/Carry"),
						attrNameStyle.Render("Movement"),
					),
				),
			),
		),
	)

	return table

	// table := lipgloss.JoinVertical(lipgloss.Left,
	// 	singleColumnStyle.Render("Name: Shadow Walker"),

	// 	lipgloss.JoinHorizontal(lipgloss.Top,
	// 		dualColumnStyle.Render(
	// 			// dualColumnStyle.Render(RenderAttributes(attrs)),
	// 			// output.WriteString(RenderAttribute(attrs.Body))
	// 			RenderAttribute(attrs.Agility),
	// 			RenderAttribute(attrs.Reaction),
	// 			RenderAttribute(attrs.Strength),
	// 			RenderAttribute(attrs.Willpower),
	// 			RenderAttribute(attrs.Logic),
	// 			RenderAttribute(attrs.Intuition),
	// 			RenderAttribute(attrs.Charisma),
	// 			RenderAttribute(attrs.Essence),
	// 			RenderAttribute(attrs.Magic),
	// 			RenderAttribute(attrs.Resonance),
	// 		),
	// 		dualColumnStyle.Render("empty"),
	// 	),
	// 	// dualColumnStyle.Render("empty"),
	// )

	// // Styles for table components
	// headerStyle := lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderBottom(true).
	// 	Width(60).
	// 	Padding(0, 1).
	// 	Bold(true)

	// dualColumnStyle := lipgloss.NewStyle().
	// 	Width(30).
	// 	Padding(0, 1)

	// divider := lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderBottom(true).
	// 	Render("────────────────────────────────────────────")

	// // Header row
	// header := headerStyle.Render("Name: Shadow Walker")

	// // Dual-column rows
	// bodyColumn := dualColumnStyle.Render("Body: 4 (6)")
	// essenceColumn := dualColumnStyle.Render("Essence: 4 (6)")

	// agilityColumn := dualColumnStyle.Render("Agility: 4 (6)")
	// emptyColumn := dualColumnStyle.Render("") // Blank column

	// dualRow1 := lipgloss.JoinHorizontal(lipgloss.Top, bodyColumn, essenceColumn)
	// dualRow2 := lipgloss.JoinHorizontal(lipgloss.Top, agilityColumn, emptyColumn)

	// // Combine everything into a single table
	// table := lipgloss.JoinVertical(lipgloss.Left,
	// 	header,
	// 	divider,
	// 	dualRow1,
	// 	dualRow2,
	// )

	return table
}
