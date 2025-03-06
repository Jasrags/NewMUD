package game

import (
	"log/slog"
	"strings"
)

// Room Functions
func (mgr *EntityManager) GetAllRooms() map[string]*Room {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.rooms
}

func (mgr *EntityManager) AddRoom(r *Room) {
	mgr.Lock()
	defer mgr.Unlock()

	if _, ok := mgr.rooms[r.ID]; ok {
		slog.Warn("Room already exists",
			slog.String("room_id", r.ID))
		return
	}

	mgr.rooms[r.ID] = r
}

func (mgr *EntityManager) GetRoom(id string) *Room {
	mgr.RLock()
	defer mgr.RUnlock()

	return mgr.rooms[strings.ToLower(id)]
}

func (mgr *EntityManager) RemoveRoom(r *Room) {
	mgr.Lock()
	defer mgr.Unlock()

	delete(mgr.rooms, r.ID)
}
