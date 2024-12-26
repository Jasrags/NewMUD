package rooms

import (
	"github.com/Jasrags/NewMUD/characters"
	"github.com/Jasrags/NewMUD/exits"
	"github.com/Jasrags/NewMUD/items"
	"github.com/Jasrags/NewMUD/mobs"
)

type Room struct {
	ID           string                  `yaml:"id"`
	AreaID       string                  `yaml:"area_id"`
	Title        string                  `yaml:"title"`
	Description  string                  `yaml:"description"`
	Exits        map[string]exits.Exit   `yaml:"exits"`
	Items        []*items.Item           `yaml:"-"`
	Characters   []*characters.Character `yaml:"-"`
	Mobs         []*mobs.Mob             `yaml:"-"`
	DefaultItems []string                `yaml:"default_items"` // IDs of items to load into the room
	DefaultMobs  []string                `yaml:"default_mobs"`  // IDs of mobs to load into the room
	SpawnedMobs  []*mobs.Mob             `yaml:"-"`             // Mobs that have been spawned into the room
}

func NewRoom() *Room {
	return &Room{
		Exits: make(map[string]exits.Exit),
	}
}

func (r *Room) AddCharacter(c *characters.Character) {
	r.Characters = append(r.Characters, c)
}

func (r *Room) RemoveCharacter(c *characters.Character) {
	for i, char := range r.Characters {
		if char.ID == c.ID {
			r.Characters = append(r.Characters[:i], r.Characters[i+1:]...)
			break
		}
	}
}

func (r *Room) AddMob(m *mobs.Mob) {
	r.Mobs = append(r.Mobs, m)
}

func (r *Room) RemoveMob(m *mobs.Mob) {
	for i, mob := range r.Mobs {
		if mob.ID == m.ID {
			r.Mobs = append(r.Mobs[:i], r.Mobs[i+1:]...)
			break
		}
	}
}

func (r *Room) AddItem(i *items.Item) {
	r.Items = append(r.Items, i)
}

func (r *Room) RemoveItem(i *items.Item) {
	for k, item := range r.Items {
		if item.ID == i.ID {
			r.Items = append(r.Items[:k], r.Items[k+1:]...)
			break
		}
	}
}

// func (r *Room) MoveToRoom(userID, roomID string) {
// 	slog.Debug("Moving user to room",
// 		slog.String("user_id", userID),
// 		slog.String("room_id", roomID))

// get user

// get room

// If the user is already in the room, return

// remove user from current room/area

// add user to new room/area

// }
