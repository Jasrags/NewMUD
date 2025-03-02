package game

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/i582/cfmt/cmd/cfmt"
)

const (
	MobsFilename = "mobs.yml"

	DispositionFriendly   = "Friendly"
	DispositionNeutral    = "Neutral"
	DispositionAggressive = "Aggressive"
)

type (
	MobBlueprint struct {
	}
	// TODO: Implement mob AI behaviors.
	// TODO: Do we want mobs to be an "instance" that will persist after spawning?
	Mob struct {
		GameEntity            `yaml:",inline"`
		Tags                  []string          `yaml:"tags"`
		ProfessionalRating    int               `yaml:"professional_rating"`
		GeneralDisposition    string            `yaml:"general_disposition"`
		CharacterDispositions map[string]string `yaml:"character_dispositions"`
	}
)

func NewMob() *Mob {
	return &Mob{
		GameEntity:            NewGameEntity(),
		GeneralDisposition:    DispositionNeutral,
		CharacterDispositions: make(map[string]string),
	}
}

func (m *Mob) Init() {
	slog.Debug("Initializing mob",
		slog.String("mob_id", m.ID))
}

func (m *Mob) GetName() string {
	return m.Name
}

func (m *Mob) GetID() string {
	return m.ID
}

func (m *Mob) SetGeneralDisposition(disposition string) {
	m.GeneralDisposition = disposition
}

func (m *Mob) ReactToMessage(sender *Character, message string) {
	// Mobs can "react" based on predefined AI behaviors.
	m.ReactToInteraction(sender, message)
}

func (m *Mob) SetDispositionForCharacter(char *Character, disposition string) {
	m.CharacterDispositions[char.ID] = disposition
}

func (m *Mob) GetDispositionForCharacter(char *Character) string {
	if disposition, exists := m.CharacterDispositions[char.ID]; exists {
		return disposition
	}
	return m.GeneralDisposition // Fallback to general disposition
}

func (m *Mob) ReactToInteraction(sender *Character, message string) {
	if strings.Contains(strings.ToLower(message), "hello") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s says: 'Hello, %s.'}}::green"+CRLF, m.Name, sender.Name), nil)
		}
	} else if strings.Contains(strings.ToLower(message), "attack") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s snarls at %s and prepares to attack!}}::red"+CRLF, m.Name, sender.Name), nil)
		}
	} else {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s looks confused by %s's words.}}::yellow"+CRLF, m.Name, sender.Name), nil)
		}
	}
}

func DescribeMobDisposition(mob *Mob, char *Character) string {
	disposition := mob.GetDispositionForCharacter(char)
	switch disposition {
	case DispositionFriendly:
		return fmt.Sprintf("%s looks at you warmly.", mob.Name)
	case DispositionNeutral:
		return fmt.Sprintf("%s glances at you indifferently.", mob.Name)
	case DispositionAggressive:
		return fmt.Sprintf("%s snarls menacingly at you!", mob.Name)
	default:
		return fmt.Sprintf("%s's demeanor is unreadable.", mob.Name)
	}
}

// RenderMobTable builds a formatted table of a mob's stats.
// It leverages the embedded GameEntity fields from Mob.
func RenderMobTable(mob *Mob) string {
	mob.Recalculate()

	metatype := EntityMgr.GetMetatype(mob.MetatypeID)

	// Define styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")). // Orange
		Align(lipgloss.Center).
		Width(80)

	singleColumnStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")). // White
		Width(80).
		Padding(0, 1)

	doubleColumnStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")). // Cyan
		Width(39).
		Padding(0, 1)

		// Character Info
		// characterName := headerStyle.Render("Ork Thug Lieutenant")
		// characterDetails := singleColumnStyle.Render(cfmt.Sprintf("Metatype %s; %s; Age: %d; Height: %dcm; Weight: %dkg; Street Cred: %d; Notoriety: %d; Public Awareness: %d",
		// metatype.Name, mob.Sex, mob.Age, mob.Height, mob.Weight, mob.StreetCred, mob.Notoriety, mob.PublicAwareness))

		// Attributes Header (Single Column spanning both double columns)
		// attributesHeader := headerStyle.Render("Attributes")

		// Double-column attributes
		// attributes := lipgloss.JoinHorizontal(lipgloss.Top,
		// doubleColumnStyle.Render(fmt.Sprintf("Body: %d", mob.Body.TotalValue)),
		// doubleColumnStyle.Render(fmt.Sprintf("Professional Rating: %d", mob.ProfessionalRating)),
		// )

		// Another row of double-column attributes
		// attributesRow2 := lipgloss.JoinHorizontal(lipgloss.Top,
		// doubleColumnStyle.Render("Agility: 3"),
		// doubleColumnStyle.Render("Magic: 1"),
		// )

		// Collect skills into a slice
	// var skillsBlock []string
	// for _, skill := range mob.Skills {
	// 	bp := EntityMgr.GetSkillBlueprint(skill.BlueprintID)
	// 	skillsBlock = append(skillsBlock, fmt.Sprintf("%s (%d)", bp.Name, skill.Rating))
	// }

	// skillsDisplay := FormatTwoColumnBlock(skillsBlock, rightBlock, 24)

	// Render full character sheet
	characterSheet := lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("Ork Thug Lieutenant"),
		singleColumnStyle.Render(cfmt.Sprintf("Metatype %s; %s; Age: %d; Height: %dcm; Weight: %dkg; Street Cred: %d; Notoriety: %d; Public Awareness: %d",
			metatype.Name, mob.Sex, mob.Age, mob.Height, mob.Weight, mob.StreetCred, mob.Notoriety, mob.PublicAwareness)),
		headerStyle.Render("Attributes"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Body:", mob.Body.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Professional Rating:", mob.ProfessionalRating)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Agility:", mob.Agility.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %.1f", "Essence:", mob.Essence.TotalValue)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Reaction:", mob.Reaction.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d/%d", "Edge:", mob.Edge.Available, mob.Edge.Max)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Strength:", mob.Strength.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d+%dd6", "Initiative:",
				mob.GetInitative(), mob.GetInitativeDice())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Willpower:", mob.Willpower.TotalValue)),
			doubleColumnStyle.Render(headerStyle.Width(39).Render("Inheret Limits")),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Logic:", mob.Logic.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Physical Limit:", mob.GetPhysicalLimit())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Intuition:", mob.Intuition.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Mental Limit:", mob.GetMentalLimit())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Charisma:", mob.Charisma.TotalValue)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Social Limit:", mob.GetSocialLimit())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Magic:", mob.Magic.TotalValue)),
			""),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Resonance:", mob.Resonance.TotalValue)),
			""),
		headerStyle.Render("Movement"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%dm/%dm/+%d Land Movement", 8, 16, 2)),
			doubleColumnStyle.Render(fmt.Sprintf("%dm/+%d Swimming", 4, 1)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			headerStyle.Width(39).Render("Active Skills"),
			headerStyle.Width(39).Render("Knowledge Skills"),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%d [%s] %s %d", 7, "A", "Blades", 3)),
			"",
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			headerStyle.Width(39).Render("Attribute-Only Tests"),
			headerStyle.Width(39).Render("Toxin Resistances"),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Composure:", mob.GetComposure())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Contact:", 7, 7)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Judge Intentions:", mob.GetJudgeIntentions())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Ingestion:", 7, 7)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Lifting & Carrying:", (mob.Strength.TotalValue+mob.Body.TotalValue))),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Inhalation:", 7, 7)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Memory:", mob.GetMemory())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Injection:", 7, 7)),
		),
		headerStyle.Width(80).Render("Addiction Resistance"),
		singleColumnStyle.Render(fmt.Sprintf("%s %d", "Resist Physical Addiction:", 7)),
		singleColumnStyle.Render(fmt.Sprintf("%s %d", "Resist Psychological Addiction:", 6)),
		headerStyle.Width(80).Render("Damage Resistances"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Armor:", mob.GetArmorValue())),
			"",
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Acid Proection:", mob.GetAcidResistance())),
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Electricity Protection:", mob.GetElectricityResistance())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Cold Proection:", mob.GetColdResistance())),
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Fire Protection:", mob.GetFireResistance())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Falling Proection:", mob.GetFallingResistance())),
			doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 7, "Fatigue Resistance:", mob.GetFatigueResistance())),
		),
		headerStyle.Width(80).Render("Metatype Abilities"),
		singleColumnStyle.Render(fmt.Sprintf("%s", "Enhanced Senses: Low-Light Vision")),
		// Edge Pool
		// Defenses
		// Damage
		// Inventory

		// Nuyen
	)

	// TODO: temp display of inventory
	for i, item := range mob.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp == nil {
			continue
		}
		if i == 0 {
			characterSheet += headerStyle.Render("Inventory") + CRLF
		}
		characterSheet += fmt.Sprintf("%s %s\n", bp.Name, item.InstanceID)
	}

	return characterSheet
}

func FormatTwoColumnBlock(leftItems []string, rightItems []string, colWidth int) string {
	var leftColumn []string
	var rightColumn []string

	// Ensure both columns have the same number of rows
	maxRows := max(len(leftItems), len(rightItems))

	// Fill left column with skills
	for i := 0; i < maxRows; i++ {
		if i < len(leftItems) {
			leftColumn = append(leftColumn, leftItems[i])
		} else {
			leftColumn = append(leftColumn, "") // Empty row for alignment
		}
	}

	// Fill right column with other dynamic content
	for i := 0; i < maxRows; i++ {
		if i < len(rightItems) {
			rightColumn = append(rightColumn, rightItems[i])
		} else {
			rightColumn = append(rightColumn, "") // Empty row for alignment
		}
	}

	// Define Lipgloss styles
	leftStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")). // Gold for skills
		Width(colWidth).
		Padding(0, 1)

	rightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEFA")). // Light Blue for right column
		Width(colWidth).
		Padding(0, 1)

	// Format the output
	var formattedRows []string
	for i := 0; i < maxRows; i++ {
		leftText := leftStyle.Render(leftColumn[i])
		rightText := rightStyle.Render(rightColumn[i])

		formattedRows = append(formattedRows, lipgloss.JoinHorizontal(lipgloss.Top, leftText, rightText))
	}

	return lipgloss.JoinVertical(lipgloss.Left, formattedRows...)
}

// intAttributeStr := "{{%-10s}}::white|bold {{%-2d}}::cyan" + CRLF
// // floatAttributeStr := "{{%-10s}}::white|bold {{%.1f}}::cyan" + CRLF
// // strAttributeStr := "{{%-10s}}::white|bold {{%-2s}}::cyan" + CRLF

// metatype := EntityMgr.GetMetatype(mob.MetatypeID)
// var builder strings.Builder

// // Characater details
// // builder.WriteString(cfmt.Sprintf("{{%s}}::white|bold [%s]"+CRLF, mob.Name, mob.ID))
// // builder.WriteString(cfmt.Sprintf("Metatype %s; %s; Age: %d; Height: %dcm; Weight: %dkg; Street Cred: %d; Notoriety: %d; Public Awareness: %d"+CRLF,
// // metatype.Name, mob.Sex, mob.Age, mob.Height, mob.Weight, mob.StreetCred, mob.Notoriety, mob.PublicAwareness))
// // builder.WriteString(CRLF)

// builder.WriteString(
// 	lipgloss.JoinVertical(lipgloss.Top,
// 		cfmt.Sprintf("{{%s}}::white|bold [%s]"+CRLF, mob.Name, mob.ID),
// 		cfmt.Sprintf("Metatype %s; %s; Age: %d; Height: %dcm; Weight: %dkg; Street Cred: %d; Notoriety: %d; Public Awareness: %d",
// 			metatype.Name, mob.Sex, mob.Age, mob.Height, mob.Weight, mob.StreetCred, mob.Notoriety, mob.PublicAwareness),
// 		cfmt.Sprintf("{{ATTRIBUTES}}::white|bold"),
// 		lipgloss.JoinHorizontal(lipgloss.Left,
// 			lipgloss.JoinVertical(lipgloss.Top,
// 				cfmt.Sprintf(intAttributeStr, "Body:", mob.Body.TotalValue),
// 			),
// 			lipgloss.JoinVertical(lipgloss.Top,
// 				cfmt.Sprintf(intAttributeStr, "Professional Rating:", mob.ProfessionalRating),
// 			),
// 		),
// 		// lipgloss.JoinHorizontal(lipgloss.Left,
// 		// cfmt.Sprintf(intAttributeStr, "Professional Rating:", mob.ProfessionalRating),
// 		// ),
// 		// "bottom",
// 	),
// )
// // 	// Attributes
// //     lipgloss.JoinHorizontal(
// //         lipgloss.Top,
// //         libgloss.Joinvertical(
// //             lipgloss.Left,
// //             "ONE",
// //     )
// // )
// // )
// // builder.WriteString(cfmt.Sprintf("{{ATTRIBUTES}}::white|bold" + CRLF))
// // Attributes from the embedded GameEntity.
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Body:", mob.Body.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Agility:", mob.Agility.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Reaction:", mob.Reaction.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Strength:", mob.Strength.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Willpower:", mob.Willpower.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Logic:", mob.Logic.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Intuition:", mob.Intuition.TotalValue))
// // builder.WriteString(cfmt.Sprintf(intAttributeStr, "Charisma:", mob.Charisma.TotalValue))
// // builder.WriteString(cfmt.Sprintf(floatAttributeStr, "Essence:", mob.Essence.TotalValue))
// // if mob.Magic.Base > 0 {
// // 	builder.WriteString(cfmt.Sprintf(intAttributeStr, "Magic:", mob.Magic.TotalValue))
// // }
// // if mob.Resonance.Base > 0 {
// // 	builder.WriteString(cfmt.Sprintf(intAttributeStr, "Resonance:", mob.Resonance.TotalValue))
// // }
// builder.WriteString(CRLF)
// // builder.WriteString(cfmt.Sprintf(strAttributeStr, "ID:", mob.ID))
// // builder.WriteString(cfmt.Sprintf(strAttributeStr, "Name:", mob.Name))
// // builder.WriteString(cfmt.Sprintf(strAttributeStr, "Title:", mob.Title))
// // builder.WriteString(cfmt.Sprintf(strAttributeStr, "Description:", mob.Description))
// // builder.WriteString(cfmt.Sprintf(strAttributeStr, "Long Description:", mob.LongDescription))
// // builder.WriteString(CRLF)

// // Composure, Judge Intentions, Lifting & Carrying, Memory
// builder.WriteString(cfmt.Sprintf("{{Composure:}}::white|bold {{%d}}::cyan"+CRLF, mob.GetComposure()))
// builder.WriteString(cfmt.Sprintf("{{Judge Intentions:}}::white|bold {{%d}}::cyan"+CRLF, mob.GetJudgeIntentions()))
// builder.WriteString(cfmt.Sprintf("{{Lifting & Carrying:}}::white|bold {{%d}}::cyan"+CRLF, (mob.Strength.TotalValue + mob.Body.TotalValue)))
// builder.WriteString(cfmt.Sprintf("{{Memory:}}::white|bold {{%d}}::cyan"+CRLF, mob.GetMemory()))
// builder.WriteString(CRLF)

// // Initiative
// builder.WriteString(cfmt.Sprintf("{{Initiative:}}::white|bold {{%d + 1d6}}::cyan"+CRLF, (mob.Reaction.TotalValue + mob.Intuition.TotalValue)))
// builder.WriteString(CRLF)

// // Edge
// builder.WriteString(cfmt.Sprintf("{{Edge:}}::white|bold {{%d/%d}}::cyan"+CRLF, mob.Edge.Available, mob.Edge.Max))

// // Limits
// builder.WriteString(cfmt.Sprintf("{{Limits:}}::white|bold Physical %-2d Mental %-2d Social %-2d"+CRLF,
// 	mob.GetPhysicalLimit(), mob.GetMentalLimit(), mob.GetSocialLimit()))

// // Condition monitors
// builder.WriteString(cfmt.Sprintf("{{Condition:}}::white|bold Physical %2d/%-2d Stun %2d/%-2d Overflow %2d/%-2d"+CRLF,
// 	0, mob.GetPhysicalConditionMax(), 0, mob.GetStunConditionMax(), 0, mob.GetOverflowConditionMax()))

// // Mob-specific data.
// builder.WriteString(cfmt.Sprintf("{{Professional Rating:}}::white|bold {{%d}}::cyan"+CRLF, mob.ProfessionalRating))
// builder.WriteString(cfmt.Sprintf("{{General Disposition:}}::white|bold {{%s}}::cyan"+CRLF, mob.GeneralDisposition))
// builder.WriteString(CRLF)

// // Skills
// builder.WriteString(cfmt.Sprintf("{{Skills:}}::white|bold" + CRLF))
// for _, skill := range mob.Skills {
// 	bp := EntityMgr.GetSkillBlueprint(skill.BlueprintID)
// 	builder.WriteString(cfmt.Sprintf("  - %s: (%d)"+CRLF, bp.Name, skill.Rating))
// }

// return wordwrap.String(builder.String(), 80)
// }
