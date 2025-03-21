package game

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/i582/cfmt/cmd/cfmt"
	ee "github.com/vansante/go-event-emitter"
)

const (
	MobsFilename = "mobs.yml"

	DispositionFriendly   = "Friendly"
	DispositionNeutral    = "Neutral"
	DispositionAggressive = "Aggressive"
)

type (
	MobSpawns struct {
		ItemID    string `yaml:"item_id"`
		Quantity  int    `yaml:"quantity"`
		Chance    int    `yaml:"chance"`
		EquipSlot string `yaml:"equip_slot"`
	}
	MobBlueprint struct {
		GameEntityInformation `yaml:",inline"`
		GameEntityStats       `yaml:",inline"`
		Spawns                []MobSpawns `yaml:"spawns"`
	}
	MobInstance struct {
		sync.RWMutex `yaml:"-"`
		Listeners    []ee.Listener `yaml:"-"`

		GameEntityDynamic `yaml:",inline"`

		InstanceID  string        `yaml:"instance_id"`
		BlueprintID string        `yaml:"blueprint_id"`
		Blueprint   *MobBlueprint `yaml:"-"`
		RoomID      string        `yaml:"room_id"`
		Room        *Room         `yaml:"-"`
	}
	// TODO: Implement mob AI behaviors.
)

func (m *MobInstance) GetArmorValue() int {
	var totalValue int

	modifiers := m.GetAllModifiers()
	if value, ok := modifiers["armor_value"]; ok {
		totalValue += value
	}

	return m.GetBody() + totalValue
}

func (m *MobInstance) GetBody() int {
	return m.Blueprint.Body
}

func (m *MobInstance) GetAgility() int {
	return m.Blueprint.Agility
}

func (m *MobInstance) GetReaction() int {
	return m.Blueprint.Reaction
}

func (m *MobInstance) GetStrength() int {
	return m.Blueprint.Strength
}

func (m *MobInstance) GetWillpower() int {
	return m.Blueprint.Willpower
}

func (m *MobInstance) GetLogic() int {
	return m.Blueprint.Logic
}

func (m *MobInstance) GetIntuition() int {
	return m.Blueprint.Intuition
}

func (m *MobInstance) GetCharisma() int {
	return m.Blueprint.Charisma
}

func (m *MobInstance) GetEssence() float64 {
	return m.Blueprint.Essence
}

func (m *MobInstance) GetMagic() int {
	return m.Blueprint.Magic
}

func (m *MobInstance) GetResonance() int {
	return m.Blueprint.Resonance
}

func (m *MobInstance) ReactToMessage(sender *Character, message string) {
	// Mobs can "react" based on predefined AI behaviors.
	m.ReactToInteraction(sender, message)
}

func (m *MobInstance) ReactToInteraction(sender *Character, message string) {
	if strings.Contains(strings.ToLower(message), "hello") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s says: 'Hello, %s.'}}::green"+CRLF, m.Blueprint.Name, sender.Name), nil)
		}
	} else if strings.Contains(strings.ToLower(message), "attack") {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s snarls at %s and prepares to attack!}}::red"+CRLF, m.Blueprint.Name, sender.Name), nil)
		}
	} else {
		room := sender.Room
		if room != nil {
			room.Broadcast(cfmt.Sprintf("{{%s looks confused by %s's words.}}::yellow"+CRLF, m.Blueprint.Name, sender.Name), nil)
		}
	}
}

func DescribeMobDisposition(mob *MobInstance, char *Character) string {
	switch mob.CharacterDispositions[char.ID] {
	case DispositionFriendly:
		return fmt.Sprintf("%s looks at you warmly.", mob.Blueprint.Name)
	case DispositionNeutral:
		return fmt.Sprintf("%s glances at you indifferently.", mob.Blueprint.Name)
	case DispositionAggressive:
		return fmt.Sprintf("%s snarls menacingly at you!", mob.Blueprint.Name)
	default:
		return fmt.Sprintf("%s's demeanor is unreadable.", mob.Blueprint.Name)
	}
}

// RenderMobTable builds a formatted table of a mob's stats.
// It leverages the embedded GameEntity fields from Mob.
func RenderMobTable(mob *MobInstance) string {
	metatype := EntityMgr.GetMetatype(mob.Blueprint.MetatypeID)

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
		headerStyle.Render(cfmt.Sprintf("%s", mob.Blueprint.Name)),
		singleColumnStyle.Render(cfmt.Sprintf("ID %s; InstanceID %s;", mob.BlueprintID, mob.InstanceID)),
		singleColumnStyle.Render(cfmt.Sprintf("Metatype %s; %s; Age: %d; Height: %dcm; Weight: %dkg; Street Cred: %d; Notoriety: %d; Public Awareness: %d",
			metatype.Name, mob.Blueprint.Sex, mob.Blueprint.Age, mob.Blueprint.Height, mob.Blueprint.Weight, mob.Blueprint.StreetCred, mob.Blueprint.Notoriety, mob.Blueprint.PublicAwareness)),
		headerStyle.Render("Attributes"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Body:", mob.Blueprint.Body)),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Professional Rating:", mob.Blueprint.ProfessionalRating)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Agility:", mob.GetAgility())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %.1f", "Essence:", mob.GetEssence())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Reaction:", mob.GetReaction())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Edge:", mob.Edge)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Strength:", mob.GetStrength())),
			// doubleColumnStyle.Render(fmt.Sprintf("%s %d+%dd6", "Initiative:",
			// mob.GetInitative(), mob.GetInitativeDice())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Willpower:", mob.GetWillpower())),
			doubleColumnStyle.Render(headerStyle.Width(39).Render("Inheret Limits")),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Logic:", mob.GetLogic())),
			// doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Physical Limit:", mob.GetPhysicalLimit())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Intuition:", mob.GetIntuition())),
			// doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Mental Limit:", mob.GetMentalLimit())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Charisma:", mob.GetCharisma())),
			// doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Social Limit:", mob.GetSocialLimit())),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Magic:", mob.GetMagic())),
			""),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Resonance:", mob.GetResonance())),
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
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Composure:", mob.Blueprint.GetComposure())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Contact:", 7, 7)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Judge Intentions:", mob.Blueprint.GetJudgeIntentions())),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Ingestion:", 7, 7)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Lifting & Carrying:", (mob.Blueprint.GetStrength()+mob.Blueprint.GetBody()))),
			doubleColumnStyle.Render(fmt.Sprintf("%s %d %d", "Inhalation:", 7, 7)),
		),
		lipgloss.JoinHorizontal(lipgloss.Top,
			doubleColumnStyle.Render(fmt.Sprintf("%s %d", "Memory:", mob.Blueprint.GetMemory())),
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
		lipgloss.JoinHorizontal(lipgloss.Top),
		doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Acid Proection:", mob.Blueprint.GetAcidResistance())),
		doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Electricity Protection:", mob.Blueprint.GetElectricityResistance())),

		lipgloss.JoinHorizontal(lipgloss.Top),
		doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Cold Proection:", mob.Blueprint.GetColdResistance())),
		doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Fire Protection:", mob.Blueprint.GetFireResistance())),

		lipgloss.JoinHorizontal(lipgloss.Top),
		doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 4, "Falling Proection:", mob.Blueprint.GetFallingResistance())),
		doubleColumnStyle.Render(fmt.Sprintf("%d %s %d", 7, "Fatigue Resistance:", mob.Blueprint.GetFatigueResistance())),

		headerStyle.Width(80).Render("Metatype Abilities"),
		singleColumnStyle.Render(fmt.Sprintf("%s", "Enhanced Senses: Low-Light Vision")),
		// Edge Pool
		// Defenses
		// Damage
		// Inventory
		// Nuyen
	)

	// TODO: temp display of inventory
	characterSheet += headerStyle.Render("Inventory") + CRLF
	for _, item := range mob.Inventory.Items {
		characterSheet += fmt.Sprintf("%s %s"+CRLF, item.Blueprint.Name, item.InstanceID)
	}
	characterSheet += headerStyle.Render("Equipment") + CRLF
	for slot, item := range mob.Equipment.Slots {
		characterSheet += fmt.Sprintf("%s %s %s"+CRLF, slot, item.Blueprint.Name, item.InstanceID)
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
