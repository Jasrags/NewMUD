package game

import (
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
)

const (
	CharactersFilepath = "_data/characters"

	CharacterRoleAdmin  CharacterRole = "admin"
	CharacterRolePlayer CharacterRole = "player"
)

type (
	CharacterRole string

	Interactable interface {
		GetName() string
		GetID() string
		ReactToMessage(sender *Character, message string)
	}

	Character struct {
		GameEntity     `yaml:",inline"`
		AccountID      string        `yaml:"account_id"`
		PregenID       string        `yaml:"pregen_id"`
		Role           CharacterRole `yaml:"role"`
		Prompt         string        `yaml:"prompt"`
		Conn           ssh.Session   `yaml:"-"`
		CreatedAt      time.Time     `yaml:"created_at"`
		UpdatedAt      *time.Time    `yaml:"updated_at"`
		DeletedAt      *time.Time    `yaml:"deleted_at"`
		CommandHistory []string      `yaml:"-"`
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

// func (c *Character) FromRoom() {
// 	c.Lock()
// 	defer c.Unlock()

// 	slog.Debug("Removing character from room",
// 		slog.String("character_id", c.ID))

// 	if c.Room == nil {
// 		slog.Error("Character has no room",
// 			slog.String("character_id", c.ID))
// 		return
// 	}

// 	c.Room.RemoveCharacter(c)
// 	c.Room = nil
// 	c.RoomID = ""
// }

// func (c *Character) ToRoom(nextRoom *Room) {
// 	c.Lock()
// 	defer c.Unlock()

// 	slog.Debug("Moving character to room",
// 		slog.String("character_id", c.ID),
// 		slog.String("room_id", nextRoom.ID))

// 	c.Room = nextRoom
// 	c.RoomID = CreateEntityRef(c.AreaID, c.Room.ID)
// 	c.Room.AddCharacter(c)
// }

func (c *Character) SetRoom(room *Room) {
	slog.Debug("Setting character room",
		slog.String("character_id", c.ID),
		slog.String("room_reference_id", room.ReferenceID))

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
					RenderKeyValue("Metatype", metatype.Name), "\t",
					// RenderKeyValue("Ethnicity", char.Ethnicity),
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
						RenderAttribute(char.Body),      // 5  (7)
						RenderAttribute(char.Agility),   // 5  (7)
						RenderAttribute(char.Reaction),  // 5  (7)
						RenderAttribute(char.Strength),  // 5  (7)
						RenderAttribute(char.Willpower), // 5  (7)
						RenderAttribute(char.Logic),     // 5  (7)
						RenderAttribute(char.Intuition), // 5  (7)
						RenderAttribute(char.Charisma),  // 5  (7)
						RenderAttribute(char.Essence),   // 5  (7)
						RenderAttribute(char.Magic),     // 5  (7)
						RenderAttribute(char.Resonance), // Essence    6.00
						// strs...,
					),
				),
			),
			// Attributes RIGHT - Derivied attributes
			lipgloss.JoinVertical(lipgloss.Left,
				headerStyle.Render(""),
				dualColumnStyle.Render(
					lipgloss.JoinVertical(lipgloss.Left), // RenderAttribute(char.Attributes.Initiative), // Initiative 10 (12) + 1d6 (2d6)
					// 			RenderAttribute(char.Attributes.InitiativeDice),
					// 			RenderAttribute(char.Attributes.Composure),       // 5  (7)
					// 			RenderAttribute(char.Attributes.JudgeIntentions), // 5  (7)
					// 			RenderAttribute(char.Attributes.Memory),          // 5  (7)
					// 			RenderAttribute(char.Attributes.Lift),            // 5  (7)
					// 			RenderAttribute(char.Attributes.Carry),           // 5  (7)
					// 			RenderAttribute(char.Attributes.Walk),            // 5  (7)
					// 			RenderAttribute(char.Attributes.Run),             // 5  (7)
					// 			RenderAttribute(char.Attributes.Swim),            // 5  (7)
					// 			"",
				),
			),
		),
	)

	return table
}
