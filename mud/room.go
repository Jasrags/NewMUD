package mud

// Room represents a room in the game.
type Room struct {
	ID          string
	Title       string
	Description string
	Exits       map[string]*Room
	Players     map[string]*Player
}

func (r *Room) AddPlayer(player *Player) {
	r.Players[player.Name] = player
}

func (r *Room) RemovePlayer(player *Player) {
	delete(r.Players, player.Name)
}

func setupWorld() *Room {
	room1 := &Room{
		ID:          "room1",
		Title:       "Small Room",
		Description: "You are in a small, cozy room. Exits lead north and east.",
		Exits:       make(map[string]*Room),
	}
	room2 := &Room{
		ID:          "room2",
		Title:       "Bright Room",
		Description: "You are in a bright, sunlit room. Exits lead south.",
		Exits:       make(map[string]*Room),
	}
	room3 := &Room{
		ID:          "room3",
		Title:       "Dark Room",
		Description: "You are in a dark, eerie room. Exits lead west.",
		Exits:       make(map[string]*Room),
	}

	// Connect rooms
	room1.Exits["north"] = room2
	room1.Exits["east"] = room3
	room2.Exits["south"] = room1
	room3.Exits["west"] = room1

	return room1
}
