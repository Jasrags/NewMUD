package mud

import (
	"strings"

	"github.com/rs/zerolog"
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
