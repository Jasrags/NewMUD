package rooms

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	Mgr = NewManager()
)

type Manager struct {
	rooms map[string]*Room
	areas map[string]*Area
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
		areas: make(map[string]*Area),
	}
}

// CreateEntityRef creates an entity reference from an area and ID.
func CreateEntityRef(area, id string) string {
	return strings.ToLower(fmt.Sprintf("%s:%s", area, id))
}

// ParseEntityRef parses an entity reference into its area and ID parts.
func ParseEntityRef(entityRef string) (area, id string) {
	parts := strings.Split(strings.ToLower(entityRef), ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func (mgr *Manager) AddRoom(r *Room) {
	slog.Debug("Adding room",
		slog.String("area_id", r.AreaID),
		slog.String("room_id", r.ID))

	mgr.rooms[CreateEntityRef(r.AreaID, r.ID)] = r
}

func (mgr *Manager) GetRoom(entityRef string) *Room {
	slog.Debug("Getting room",
		slog.String("entity_ref", entityRef))

	return mgr.rooms[entityRef]
}

func (mgr *Manager) RemoveRoom(r *Room) {
	slog.Debug("Removing room",
		slog.String("area_id", r.AreaID),
		slog.String("room_id", r.ID))

	delete(mgr.rooms, CreateEntityRef(r.AreaID, r.ID))
}

func (mgr *Manager) LoadDataFiles() {
	dataFilePath := viper.GetString("data.areas_path")
	manifestFileName := viper.GetString("data.manifest_file")
	roomsFileName := viper.GetString("data.rooms_file")

	slog.Info("Loading room data files",
		slog.String("datafile_path", dataFilePath),
		slog.String("manifest_file", manifestFileName),
		slog.String("rooms_file", roomsFileName))

	files, err := os.ReadDir(dataFilePath)
	if err != nil {
		slog.Error("failed reading directory",
			slog.String("datafile_path", dataFilePath),
			slog.Any("error", err))
	}

	for _, file := range files {
		if file.IsDir() {
			areaPath := filepath.Join(dataFilePath, file.Name())
			manifestPath := filepath.Join(areaPath, manifestFileName)
			roomsPath := filepath.Join(areaPath, roomsFileName)

			// Load area manifest
			manifestData, errReadFile := os.ReadFile(manifestPath)
			if errReadFile != nil {
				slog.Error("failed reading manifest file",
					slog.Any("error", errReadFile),
					slog.String("area_path", areaPath))
				continue
			}

			var area Area
			if err := yaml.Unmarshal(manifestData, &area); err != nil {
				slog.Error("failed to unmarshal manifest data",
					slog.Any("error", err),
					slog.String("area_path", areaPath))
				continue
			}

			slog.Info("Loaded area manifest",
				slog.String("area_name", file.Name()))

			// Add area to roomManager
			mgr.areas[area.ID] = &area

			// Load rooms
			roomsData, err := os.ReadFile(roomsPath)
			if err != nil {
				slog.Error("failed reading rooms file",
					slog.Any("error", err),
					slog.String("rooms_path", roomsPath))
				continue
			}

			var rooms []Room
			if err := yaml.Unmarshal(roomsData, &rooms); err != nil {
				slog.Error("failed to unmarshal rooms data",
					slog.Any("error", err),
					slog.String("rooms_path", roomsPath))
				continue
			}

			// Add rooms to roomManager
			for _, room := range rooms {
				room.AreaID = area.ID
				mgr.AddRoom(&room)
				// mgr.rooms[CreateEntityRef(room.AreaID, room.ID)] = &room
				// slog.Debug("Loaded room",
				// 	slog.String("area_id", room.AreaID),
				// 	slog.String("room_id", room.ID))
			}

			slog.Info("Loaded area rooms",
				slog.Int("count", len(rooms)),
				slog.String("area_name", file.Name()))
		}
	}
}
