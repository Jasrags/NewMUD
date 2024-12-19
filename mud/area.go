package mud

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type AreaFloor struct {
	Log                         zerolog.Logger
	Z, lowX, highX, lowY, highY int
	// Map                         map[int]map[int]*AreaFloor
	Map [][]*Room
}

func NewAreaFloor(z int) *AreaFloor {
	return &AreaFloor{
		Log: NewDevLogger(),
		Z:   z,
	}
}

func (af *AreaFloor) AddRoom(x, y int, room *Room) {
	af.Log.Debug().
		Int("x", x).
		Int("y", y).
		// Str("room_id", room.ID).
		Msg("Adding room to floor")

	if room == nil {
		af.Log.Warn().
			Int("x", x).
			Int("y", y).
			Msg("Room is nil")
		return
	}

	if af.GetRoom(room.Coordinates.X, room.Coordinates.Y) != nil {
		af.Log.Warn().
			Int("x", x).
			Int("y", y).
			Msg("Room already exists")
		return
	}

	if x < af.lowX {
		af.lowX = x
	} else if x > af.highX {
		af.highX = x
	}

	if y < af.lowY {
		af.lowY = y
	} else if y > af.highY {
		af.highY = y
	}

	af.Map[x][y] = room
}

func (af *AreaFloor) GetRoom(x, y int) *Room {
	af.Log.Debug().
		Int("x", x).
		Int("y", y).
		Msg("Getting room from floor")

	if af.Map[x] == nil {
		af.Log.Warn().
			Int("x", x).
			Msg("Column does not exist")
	}

	return af.Map[x][y]
}

func (af *AreaFloor) RemoveRoom(x, y int) {
	af.Log.Debug().
		Int("x", x).
		Int("y", y).
		Msg("Removing room from floor")
}

type Area struct {
	Log   zerolog.Logger
	ID    string           `yaml:"id"`
	Title string           `yaml:"title"`
	Rooms map[string]*Room `yaml:"-"`
	Map   []*AreaFloor     `yaml:"-"`
}

func NewArea() *Area {
	return &Area{
		Log:   NewDevLogger(),
		Rooms: make(map[string]*Room),
	}
}

func (a *Area) GetRoomByID(id string) *Room {
	a.Log.Debug().
		Str("room_id", id).
		Msg("Getting room by ID")

	return a.Rooms[strings.ToLower(id)]
}

func (a *Area) AddRoom(room *Room) {
	a.Log.Debug().
		Str("room_id", room.ID).
		Msg("Adding room")

	a.Rooms[strings.ToLower(room.ID)] = room

	if room.Coordinates != nil {
		a.AddRoomToMap(room)
	}
}

func (a *Area) RemoveRoom(id string) {
	a.Log.Debug().
		Str("room_id", id).
		Msg("Removing room")

	delete(a.Rooms, strings.ToLower(id))
}

func (a *Area) AddRoomToMap(room *Room) {
	a.Log.Debug().
		Str("room_id", room.ID).
		Msg("Adding room to map")

	if room.Coordinates == nil {
		a.Log.Warn().
			Str("room_id", room.ID).
			Msg("Room has no coordinates")

		return
	}

	if a.Map[room.Coordinates.Z] == nil {
		a.Map[room.Coordinates.Z] = NewAreaFloor(room.Coordinates.Z)
	}

	a.Map[room.Coordinates.Z].AddRoom(room.Coordinates.X, room.Coordinates.Y, room)
}

func (a *Area) GetRoomAtCoordinates(z, x, y int) *Room {
	a.Log.Debug().
		Int("z", z).
		Int("x", x).
		Int("y", y).
		Msg("Getting room at coordinates")

	if a.Map[z] == nil {
		a.Log.Warn().
			Int("z", z).
			Msg("Floor does not exist")

		return nil
	}

	return a.Map[z].GetRoom(x, y)
}

func (a *Area) Update() {
	a.Log.Debug().Msg("Updating area")
}

type AreaManager struct {
	Log             zerolog.Logger
	RoomManager     *RoomManager
	Areas           map[string]*Area
	placeholderArea *Area
}

func NewAreaManager(rm *RoomManager) *AreaManager {
	return &AreaManager{
		Log:         NewDevLogger(),
		RoomManager: rm,
		Areas:       make(map[string]*Area),
	}
}

func (am *AreaManager) Load() {
	am.Log.Info().Msg("Loading areas")

	dataPath := "_data/areas"
	files, err := os.ReadDir(dataPath)
	if err != nil {
		am.Log.Fatal().Err(err).Msg("Failed to read data directory")
		return
	}

	for _, file := range files {
		if file.IsDir() {
			areaFilePath := filepath.Join(dataPath, file.Name(), "manifest.yml")
			if _, err := os.Stat(areaFilePath); os.IsNotExist(err) {
				continue
			}

			areaFile, err := os.ReadFile(areaFilePath)
			if err != nil {
				am.Log.Error().Err(err).Msgf("Failed to read area file: %s", areaFilePath)
				continue
			}

			area := NewArea()
			if err := yaml.Unmarshal(areaFile, &area); err != nil {
				am.Log.Error().Err(err).Msgf("Failed to unmarshal area file: %s", areaFilePath)
				continue
			}

			am.AddArea(area)
		}
	}

	am.Log.Info().Msg("Linking rooms to areas")
	for _, room := range am.RoomManager.Rooms {

		area := am.GetAreaByReference(room.AreaID)
		if area == nil {
			am.Log.Warn().
				Str("room_id", room.ID).
				Msg("Room has no area")
			continue
		}

		area.AddRoom(room)
	}

	for _, area := range am.Areas {
		am.Log.Info().
			Str("area_id", area.ID).
			Int("room_count", len(area.Rooms)).
			Msg("Loaded area")
	}

	am.Log.Info().
		Int("area_count", len(am.Areas)).
		Msg("Loaded areas")
}

func (am *AreaManager) AddArea(area *Area) {
	am.Log.Debug().
		Str("area_id", area.ID).
		Msg("Adding area")

	am.Areas[strings.ToLower(area.ID)] = area
}

func (am *AreaManager) GetArea(id string) *Area {
	am.Log.Debug().
		Str("area_id", id).
		Msg("Getting area")

	return am.Areas[strings.ToLower(id)]
}

func (am *AreaManager) GetAreaByReference(referenceID string) *Area {
	am.Log.Debug().
		Str("reference_id", referenceID).
		Msg("Getting area by reference")

	id := strings.Split(referenceID, ":")[0]

	return am.Areas[strings.ToLower(id)]
}

func (am *AreaManager) RemoveArea(id string) {
	am.Log.Debug().
		Str("area_id", id).
		Msg("Removing area")

	delete(am.Areas, strings.ToLower(id))
}

func (am *AreaManager) TickAll() {
	am.Log.Debug().Msg("Ticking all areas")

	for _, area := range am.Areas {
		area.Update()
	}
}

func (am *AreaManager) GetPlaceholderArea() *Area {
	am.Log.Debug().Msg("Getting placeholder area")

	if am.placeholderArea != nil {
		return am.placeholderArea
	}

	// Create a placeholder area
	area := NewArea()
	area.ID = "placeholder"
	area.Title = "Placeholder"

	room := NewRoom()
	room.ID = "placeholder"
	room.Title = "Placeholder"
	room.Description = "This is a placeholder room."

	area.AddRoom(room)

	am.placeholderArea = area

	return am.placeholderArea
}
