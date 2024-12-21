package mud

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type AreaManager struct {
	Log             zerolog.Logger
	RoomManager     *RoomManager
	Areas           map[string]*Area
	placeholderArea *Area
}

func NewAreaManager(l zerolog.Logger, rm *RoomManager) *AreaManager {
	return &AreaManager{
		Log:         l,
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

			area := NewArea(am.Log)
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
	area := NewArea(am.Log)
	area.ID = "placeholder"
	area.Title = "Placeholder"

	room := NewRoom(am.Log)
	room.ID = "placeholder"
	room.Title = "Placeholder"
	room.Description = "This is a placeholder room."

	area.AddRoom(room)

	am.placeholderArea = area

	return am.placeholderArea
}
