package mud

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

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
	rm.Log.Info().Msg("Loading rooms")

	dataPath := viper.GetString("data.areas_path")
	files, err := os.ReadDir(dataPath)
	if err != nil {
		rm.Log.Error().Err(err).Msg("Failed to read data directory")
		return
	}
	for _, file := range files {
		if file.IsDir() {
			roomFilePath := filepath.Join(dataPath, file.Name(), viper.GetString("data.rooms_file"))
			if _, err := os.Stat(roomFilePath); os.IsNotExist(err) {
				continue
			}

			areaName := file.Name()
			roomFile, err := os.ReadFile(roomFilePath)
			if err != nil {
				rm.Log.Fatal().Err(err).Msgf("Failed to read room file: %s", roomFilePath)
				continue
			}

			var data []RoomData
			if err := yaml.Unmarshal(roomFile, &data); err != nil {
				rm.Log.Fatal().Err(err).Msgf("Failed to unmarshal room file: %s", roomFilePath)
				continue
			}

			// Build all the rooms prefixed with the area name
			for _, d := range data {
				roomID := CreateEntityRef(areaName, d.ID)
				room := NewRoom()
				room.ID = roomID
				room.Title = d.Title
				room.Description = d.Description
				room.Coordinates = d.Coordinates
				room.AreaID = areaName
				rm.AddRoom(room)
			}

			rm.Log.Info().Msg("Building room exits")
			// Add exits to the rooms
			for _, d := range data {
				room := rm.GetRoom(CreateEntityRef(areaName, d.ID))

				if room == nil {
					rm.Log.Error().
						Str("room_id", fmt.Sprintf("%s:%s", areaName, d.ID)).
						Msg("Exit room not found")
					continue
				}

				for _, exit := range d.Exits {
					exitRoom := rm.GetRoom(exit.RoomID)

					if exitRoom == nil {
						rm.Log.Error().
							Str("room_id", exit.RoomID).
							Msg("Exit room not found")
						continue
					}

					room.Exits[exit.Direction] = &Exit{
						Room:      exitRoom,
						Direction: exit.Direction,
					}
				}
			}
		}
	}

	rm.Log.Info().
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
