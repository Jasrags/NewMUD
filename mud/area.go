package mud

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type Area struct {
	Log   zerolog.Logger
	ID    string           `yaml:"id"`
	Title string           `yaml:"title"`
	Rooms map[string]*Room `yaml:"-"`
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
}

func (a *Area) RemoveRoom(id string) {
	a.Log.Debug().
		Str("room_id", id).
		Msg("Removing room")

	delete(a.Rooms, strings.ToLower(id))
}

func (a *Area) AddRoomToMap(room *Room) {
	panic("not implemented")
}

func (a *Area) Update() {
	a.Log.Debug().Msg("Updating area")
}

type AreaManager struct {
	Log             zerolog.Logger
	Areas           map[string]*Area
	placeholderArea *Area
}

func NewAreaManager() *AreaManager {
	return &AreaManager{
		Log:   NewDevLogger(),
		Areas: make(map[string]*Area),
	}
}

func (am *AreaManager) Load() {
	am.Log.Debug().Msg("Loading areas")
	dataPath := "_data/areas"
	files, err := os.ReadDir(dataPath)
	if err != nil {
		am.Log.Error().Err(err).Msg("Failed to read data directory")
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

	am.Log.Debug().
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
