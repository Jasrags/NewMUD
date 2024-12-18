package mud

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type RoomData struct {
	ID          string       `yaml:"id"`
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Coordinates *Coordinates `yaml:"coordinates"`
	Exits       []ExitData   `yaml:"exits"`
}

type ExitData struct {
	RoomID    string `yaml:"room_id"`
	Direction string `yaml:"direction"`
	// Inferred  bool   `yaml:"-"`
}

// Room represents a game room
type Room struct {
	Log         zerolog.Logger
	ID          string
	Title       string
	Description string
	Coordinates *Coordinates
	Players     map[string]*Player
	Exits       map[string]*Exit
	Area        *Area
	// TODO: Add doors
}

type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Exit struct {
	Room      *Room
	Direction string
	// Inferred  bool
}

// type RoomData struct {
// 	ID          string            `yaml:"id"`
// 	Title       string            `yaml:"title"`
// 	Description string            `yaml:"description"`
// 	Coordinates *Coordinates      `yaml:"coordinates"`
// 	Exits       map[string]string `yaml:"exits"`
// }

func NewRoom() *Room {
	return &Room{
		Log:     NewDevLogger(),
		Players: make(map[string]*Player),
		Exits:   make(map[string]*Exit),
	}
}

func (r *Room) AddPlayer(player *Player) {
	r.Players[player.Name] = player
}

func (r *Room) RemovePlayer(player *Player) {
	delete(r.Players, player.Name)
}

type RoomManager struct {
	Log   zerolog.Logger
	Rooms map[string]*Room
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Log:   NewDevLogger(),
		Rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) Load() {
	rm.Log.Debug().Msg("Loading rooms")

	dataPath := "_data/areas"
	files, err := os.ReadDir(dataPath)
	if err != nil {
		rm.Log.Error().Err(err).Msg("Failed to read data directory")
		return
	}
	for _, file := range files {
		if file.IsDir() {
			roomFilePath := filepath.Join(dataPath, file.Name(), "rooms.yml")
			if _, err := os.Stat(roomFilePath); os.IsNotExist(err) {
				continue
			}

			areaName := file.Name()
			roomFile, err := os.ReadFile(roomFilePath)
			if err != nil {
				rm.Log.Error().Err(err).Msgf("Failed to read room file: %s", roomFilePath)
				continue
			}

			var data []RoomData
			if err := yaml.Unmarshal(roomFile, &data); err != nil {
				rm.Log.Error().Err(err).Msgf("Failed to unmarshal room file: %s", roomFilePath)
				continue
			}

			// Build all the rooms prefixed with the area name
			for _, d := range data {
				room := NewRoom()
				room.ID = fmt.Sprintf("%s:%s", areaName, d.ID)
				room.Title = d.Title
				room.Description = d.Description
				room.Coordinates = d.Coordinates
				rm.AddRoom(room)
			}

			// Add exits to the rooms
			for _, d := range data {
				room := rm.GetRoom(fmt.Sprintf("%s:%s", areaName, d.ID))
				for _, exit := range d.Exits {
					room.Exits[exit.Direction] = &Exit{
						Room:      rm.GetRoom(exit.RoomID),
						Direction: exit.Direction,
					}
				}
			}
		}
	}

	rm.Log.Debug().
		Int("room_count", len(rm.Rooms)).
		Msg("Loaded rooms")
}

func (rm *RoomManager) AddRoom(room *Room) {
	rm.Log.Debug().
		Str("room_id", room.ID).
		Msg("Adding room")

	rm.Rooms[strings.ToLower(room.ID)] = room
}

func (rm *RoomManager) GetRoom(entityRef string) *Room {
	rm.Log.Debug().
		Str("entity_ref", entityRef).
		Msg("Getting room")

	return rm.Rooms[strings.ToLower(entityRef)]
}

func (rm *RoomManager) RemoveRoom(id string) {
	rm.Log.Debug().
		Str("room_id", id).
		Msg("Removing room")

	delete(rm.Rooms, strings.ToLower(id))
}
