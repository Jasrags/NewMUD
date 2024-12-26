package characters

import "time"

type Character struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	RoomID    string    `json:"room_id"`
	AreaID    string    `json:"area_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCharacter() *Character {
	return &Character{
		CreatedAt: time.Now(),
	}
}
