package game

import (
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
	ee "github.com/vansante/go-event-emitter"
)

const (
	CharactersFilepath = "_data/characters"

	CharacterRoleAdmin  = "admin"
	CharacterRolePlayer = "player"
)

type (
	Karma struct {
		Available int `yaml:"available"` // Karma available to spend
		Total     int `yaml:"total"`     // Total karma earned
	}
	Character struct {
		sync.RWMutex `yaml:"-"`
		Listeners    []ee.Listener `yaml:"-"`

		GameEntity     `yaml:",inline"`
		RoomID         string      `yaml:"room_id"`
		Room           *Room       `yaml:"-"`
		AccountID      string      `yaml:"account_id"`
		Account        *Account    `yaml:"account"`
		PregenID       string      `yaml:"pregen_id"`
		Role           string      `yaml:"role"`
		Prompt         string      `yaml:"prompt"`
		Karma          Karma       `yaml:"karma"`
		Conn           ssh.Session `yaml:"-"`
		CreatedAt      time.Time   `yaml:"created_at"`
		UpdatedAt      *time.Time  `yaml:"updated_at"`
		DeletedAt      *time.Time  `yaml:"deleted_at"`
		CommandHistory []string    `yaml:"-"`
	}
)

func NewCharacter() *Character {
	return &Character{
		GameEntity: NewGameEntity(),
		Role:       CharacterRolePlayer,
		Prompt:     DefaultPrompt,
		CreatedAt:  time.Now(),
	}
}

func (c *Character) Init() {
	slog.Debug("Initializing character",
		slog.String("character_id", c.ID))
}

func (c *Character) Send(msg string) {
	WriteString(c.Conn, msg)
}

func (c *Character) GetName() string {
	return c.Name
}

func (c *Character) GetID() string {
	return c.ID
}

func (c *Character) ReactToMessage(sender *Character, message string) {
	// For Characters, send the message via their session.
	if c.Conn != nil {
		WriteString(c.Conn, cfmt.Sprintf("{{%s says to you: '%s'}}::green"+CRLF, sender.Name, message))
	}
}

func (c *Character) SetRoom(room *Room) {
	c.Room = room
	c.RoomID = room.ReferenceID
}

func (c *Character) MoveToRoom(nextRoom *Room) {
	slog.Debug("Moving character to room",
		slog.String("character_id", c.ID),
		slog.String("room_id", nextRoom.ID))

	if c.Room != nil && c.Room.ID != nextRoom.ID {
		// EventMgr.Publish(EventRoomCharacterLeave, &RoomCharacterLeave{Character: c, Room: c.Room, NextRoom: nextRoom})
		c.Room.Broadcast(cfmt.Sprintf("\n{{%s leaves the room.}}::green"+CRLF, c.Name), []string{c.ID})
		c.Room.RemoveCharacter(c)
	}

	c.SetRoom(nextRoom)
	nextRoom.AddCharacter(c)
	c.Room.Broadcast(cfmt.Sprintf("\n{{%s enters the room.}::green}"+CRLF, c.Name), []string{c.ID})

	// EventMgr.Publish(EventRoomCharacterEnter, &RoomCharacterEnter{Character: c, Room: c.Room, PrevRoom: prevRoom})
	// EventMgr.Publish(EventPlayerEnterRoom, &PlayerEnterRoom{Character: c, Room: c.Room})
}

func (c *Character) Save() error {
	c.Lock()
	defer c.Unlock()

	slog.Debug("Saving character",
		slog.String("character_id", c.ID),
		slog.String("character_name", c.Name))

	dataFilePath := viper.GetString("data.characters_path")

	slog.Debug("Saving character",
		slog.String("character_id", c.ID),
		slog.String("character_name", c.Name))

	t := time.Now()
	c.UpdatedAt = &t

	filePath := filepath.Join(dataFilePath, strings.ToLower(c.Name)+".yml")
	if err := SaveYAML(filePath, c); err != nil {
		slog.Error("failed to save character data",
			slog.Any("error", err))
		return err
	}

	return nil
}

func RenderCharacterTable(char *Character) string {
	metatype := EntityMgr.GetMetatype(char.MetatypeID)

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
					RenderKeyValue("Metatype", metatype.Name), "\t",
				),
				lipgloss.JoinHorizontal(lipgloss.Top,
					RenderKeyValue("Age", "0"), "\t",
					RenderKeyValue("Sex", char.Sex), "\t",
					RenderKeyValue("Height", "0"), "\t",
					RenderKeyValue("Weight", "0"),
				),
				lipgloss.JoinHorizontal(lipgloss.Top,
					cfmt.Sprintf("%s: %d;", "Street Cred:", char.StreetCred),
					cfmt.Sprintf("%s: %d;", "Notoriety:", char.Notoriety),
					cfmt.Sprintf("%s: %d;", "Public Awareness:", char.PublicAwareness),
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
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						cfmt.Sprintf("%-10s %d", "Body:", char.GetBody()),
						cfmt.Sprintf("%-10s %d", "Agility:", char.GetAgility()),
						cfmt.Sprintf("%-10s %d", "Reaction:", char.GetReaction()),
						cfmt.Sprintf("%-10s %d", "Strength:", char.GetStrength()),
						cfmt.Sprintf("%-10s %d", "Willpower:", char.GetWillpower()),
						cfmt.Sprintf("%-10s %d", "Logic:", char.GetLogic()),
						cfmt.Sprintf("%-10s %d", "Intuition:", char.GetIntuition()),
						cfmt.Sprintf("%-10s %d", "Charisma:", char.GetCharisma()),
						cfmt.Sprintf("%-10s %.2f", "Essence:", char.GetEssence()),
						cfmt.Sprintf("%-10s %d", "Magic:", char.GetMagic()),
						cfmt.Sprintf("%-10s %d", "Resonance:", char.GetResonance()),
					),
				),
			),
			// Attributes RIGHT - Derivied attributes
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render(""),
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left,
						cfmt.Sprintf("%-17s %d + 6d%d", "Initiative:", char.GetInitative(), char.GetInitativeDice()),
						cfmt.Sprintf("%-17s %d", "Composure:", char.GetComposure()),
						cfmt.Sprintf("%-17s %d", "Judge Intentions:", char.GetJudgeIntentions()),
						cfmt.Sprintf("%-17s %d", "Memory:", char.GetMemory()),
					),
				),
			),
		),
	)

	return table
}
