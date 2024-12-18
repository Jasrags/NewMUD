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
//     ID          string            `yaml:"id"`
//     Title       string            `yaml:"title"`
//     Description string            `yaml:"description"`
//     Exits       map[string]string `yaml:"exits"`
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

			for _, d := range data {
				room := NewRoom()
				room.ID = fmt.Sprintf("%s:%s", areaName, d.ID)
				room.Title = d.Title
				room.Description = d.Description
				room.Coordinates = d.Coordinates

				for _, exit := range d.Exits {
					room.Exits[exit.Direction] = &Exit{
						Room:      room,
						Direction: exit.Direction,
					}
				}
				rm.AddRoom(room)
			}
		}
	}

	rm.Log.Debug().
		Int("room_count", len(rm.Rooms)).
		Msg("Loaded rooms")
}

// func (rm *RoomManager) Load() {
// 	rm.Log.Debug().Msg("Loading rooms")

// 	dataPath := "_data/areas"

// 	files, err := os.ReadDir(dataPath)
// 	if err != nil {
// 		rm.Log.Error().Err(err).Msg("Failed to read data directory")
// 		return
// 	}

// 	for _, file := range files {
// 		if file.IsDir() {
// 			roomFilePath := filepath.Join(dataPath, file.Name(), "rooms.yml")
// 			if _, err := os.Stat(roomFilePath); os.IsNotExist(err) {
// 				continue
// 			}

// 			roomFile, err := os.ReadFile(roomFilePath)
// 			if err != nil {
// 				rm.Log.Error().Err(err).Msgf("Failed to read room file: %s", roomFilePath)
// 				continue
// 			}

// 			var data []RoomData
// 			if err := yaml.Unmarshal(roomFile, &data); err != nil {
// 				rm.Log.Error().Err(err).Msgf("Failed to unmarshal room file: %s", roomFilePath)
// 				continue
// 			}

// 			for _, d := range data {
// 				rm.Log.Debug().
// 					Str("id", d.ID).
// 					Str("title", d.Title).
// 					Str("description", d.Description).
// 					Msg("Adding room")

// 				room := NewRoom()
// 				room.ID = d.ID
// 				room.Title = d.Title
// 				room.Description = d.Description

// 				rm.AddRoom(room)
// 			}

// 			for _, d := range data {
// 				room := rm.GetRoom(d.ID)

// 				for _, exit := range d.Exits {
// 					// rm.Log.Debug().
// 					// 	Str("direction", exit.Direction).
// 					// 	Str("room_id", exit.RoomID).
// 					// 	Msg("Adding exit to room")
// 					room.Exits[exit.Direction] = &Exit{
// 						Room:      room,
// 						Direction: exit.Direction,
// 					}
// 				}
// 			}

// 			// for _, room := range rm.Rooms {
// 			// 	for _, exit := range data {
// 			// 		for _, e := range exit.Exits {
// 			// 			if e.RoomID == room.ID {
// 			// 				rm.Log.Debug().
// 			// 					Str("room_id", room.ID).
// 			// 					Str("exit_room_id", e.RoomID).
// 			// 					Str("direction", e.Direction).
// 			// 					Msg("Adding exit to room")

// 			// 				room.Exits[e.Direction] = &Exit{
// 			// 					Room:      rm.GetRoom(e.RoomID),
// 			// 					Direction: e.Direction,
// 			// 				}
// 			// 			}
// 			// 		}
// 			// 	}
// 			// }
// 		}
// 	}
// }

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

// var Rooms = map[string]*Room{}

// import (
// 	"fmt"

// 	"github.com/rs/zerolog"
// )

// // Room represents a room in the game.git s
// type Room struct {
// 	Log         zerolog.Logger
// 	ID          string
// 	Title       string
// 	Description string
// 	Exits       map[string]*Room
// 	Players     map[string]*Player
// }

// func NewRoom(id, title, description string) *Room {
// 	return &Room{
// 		Log:         NewDevLogger(),
// 		ID:          id,
// 		Title:       title,
// 		Description: description,
// 		Exits:       make(map[string]*Room),
// 		Players:     make(map[string]*Player),
// 	}
// }

// func (r *Room) AddPlayer(player *Player) {
// 	r.Log.Debug().
// 		Str("player_name", player.Name).
// 		Str("room_id", r.ID).
// 		Msg("Add player to room")

// 	r.Players[player.Name] = player
// }

// func (r *Room) RemovePlayer(player *Player) {
// 	r.Log.Debug().
// 		Str("player_name", player.Name).
// 		Str("room_id", r.ID).
// 		Msg("Remove player from room")

// 	delete(r.Players, player.Name)
// }

// func setupWorld() *Room {
// 	room1 := NewRoom("room1", "Small Room", "You are in a small, cozy room. Exits lead north and east.")
// 	room2 := NewRoom("room2", "Bright Room", "You are in a bright, sunlit room. Exits lead south.")
// 	room3 := NewRoom("room3", "Dark Room", "You are in a dark, eerie room. Exits lead west.")

// 	// Connect rooms
// 	room1.Exits["north"] = room2
// 	room1.Exits["east"] = room3
// 	room2.Exits["south"] = room1
// 	room3.Exits["west"] = room1

// 	room1.Listen()
// 	room2.Listen()
// 	room3.Listen()

// 	return room1
// }

// func (r *Room) Broadcast(message string, exclude *Player) {
// 	for _, player := range r.Players {
// 		if player != exclude {
// 			player.Out <- message
// 		}
// 	}
// }

// func (r *Room) Listen() {
// 	r.Log.Debug().
// 		Str("room_id", r.ID).
// 		Msg("Listening for player events")

// 	// Subscribe to player entrance
// 	eventBus.Subscribe(EventPlayerEnter, func(player *Player, roomID string) {
// 		if roomID == r.ID {
// 			r.Broadcast(fmt.Sprintf("%s enters the room.\r\n", player.Name), player)
// 			r.AddPlayer(player)
// 		}
// 	})

// 	// Subscribe to player exit
// 	eventBus.Subscribe(EventPlayerExit, func(player *Player, roomID string) {
// 		if roomID == r.ID {
// 			r.Broadcast(fmt.Sprintf("%s leaves the room.\r\n", player.Name), player)
// 			r.RemovePlayer(player)
// 		}
// 	})
// }
