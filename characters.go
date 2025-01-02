package main

import (
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
	ee "github.com/vansante/go-event-emitter"
)

type CharacterRole string

const (
	CharacterRoleAdmin  CharacterRole = "admin"
	CharacterRolePlayer CharacterRole = "player"
)

type Character struct {
	sync.RWMutex
	Listeners []ee.Listener `json:"-"`
	Conn      ssh.Session   `json:"-"`

	ID        string        `json:"id"`
	User      *User         `json:"-"`
	UserID    string        `json:"user_id"`
	Name      string        `json:"name"`
	Room      *Room         `json:"-"`
	RoomID    string        `json:"room_id"`
	Area      *Area         `json:"-"`
	AreaID    string        `json:"area_id"`
	Role      CharacterRole `json:"role"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at"`
}

func NewCharacter() *Character {
	return &Character{
		CreatedAt: time.Now(),
		Role:      CharacterRolePlayer,
	}
}

func (c *Character) Init() {
	slog.Debug("Initializing character",
		slog.String("character_id", c.ID))
}

func (c *Character) Send(msg string) {
	slog.Debug("Sending message to character",
		slog.String("character_id", c.ID),
		slog.String("message", msg))

	io.WriteString(c.Conn, msg)
}

func (c *Character) MoveToRoom(nextRoom *Room) {
	c.Lock()
	defer c.Unlock()

	slog.Debug("Moving character to room",
		slog.String("character_id", c.ID),
		slog.String("room_id", nextRoom.ID))

	var prevRoom = c.Room
	if c.Room != nil && c.Room.ID != nextRoom.ID {
		EventMgr.Publish(EventRoomCharacterLeave, &RoomCharacterLeave{Character: c, Room: c.Room, NextRoom: nextRoom})
		c.Room.RemoveCharacter(c)
	}

	c.Room = nextRoom
	c.RoomID = CreateEntityRef(c.AreaID, c.Room.ID)
	nextRoom.AddCharacter(c)

	EventMgr.Publish(EventRoomCharacterEnter, &RoomCharacterEnter{Character: c, Room: c.Room, PrevRoom: prevRoom})
	EventMgr.Publish(EventPlayerEnterRoom, &PlayerEnterRoom{Character: c, Room: c.Room})
}

func (c *Character) Save() error {
	c.Lock()
	defer c.Unlock()

	dataFilePath := viper.GetString("data.characters_path")

	slog.Debug("Saving character",
		slog.String("character_id", c.ID),
		slog.String("character_name", c.Name))

	t := time.Now()
	c.UpdatedAt = &t

	filePath := filepath.Join(dataFilePath, strings.ToLower(c.Name)+".json")
	if err := SaveJSON(filePath, c); err != nil {
		slog.Error("failed to save character data",
			slog.Any("error", err))
		return err
	}
	// filepath := filepath.Join(viper.GetString("data.characters_path"), c.Name+".json")

	return nil
}
