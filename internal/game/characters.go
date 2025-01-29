package game

import (
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
)

type Interactable interface {
	GetName() string
	GetID() string
	ReactToMessage(sender *Character, message string)
}

type CharacterRole string

const (
	CharacterRoleAdmin  CharacterRole = "admin"
	CharacterRolePlayer CharacterRole = "player"
)

type Character struct {
	GameEntity     `yaml:",inline"`
	User           *Account      `yaml:"-"`
	UserID         string        `yaml:"user_id"`
	Role           CharacterRole `yaml:"role"`
	Prompt         string        `yaml:"prompt"`
	Conn           ssh.Session   `yaml:"-"`
	CreatedAt      time.Time     `yaml:"created_at"`
	UpdatedAt      *time.Time    `yaml:"updated_at"`
	DeletedAt      *time.Time    `yaml:"deleted_at"`
	CommandHistory []string      `yaml:"-"`
}

func NewCharacter() *Character {
	return &Character{
		GameEntity: NewGameEntity(),
		Role:       CharacterRolePlayer,
		CreatedAt:  time.Now(),
	}
}

func (c *Character) Init() {
	slog.Debug("Initializing character",
		slog.String("character_id", c.ID))
}

func (c *Character) Send(msg string) {
	io.WriteString(c.Conn, msg)
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
		io.WriteString(c.Conn, cfmt.Sprintf("{{%s says to you: '%s'}}::green\n", sender.Name, message))
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
		c.Room.Broadcast(cfmt.Sprintf("\n{{%s leaves the room.}}::green\n", c.Name), []string{c.ID})
		c.Room.RemoveCharacter(c)
	}

	c.SetRoom(nextRoom)
	nextRoom.AddCharacter(c)
	c.Room.Broadcast(cfmt.Sprintf("\n{{%s enters the room.}::green}\n", c.Name), []string{c.ID})

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
