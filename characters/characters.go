package characters

type Character struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCharacter() *Character {
	return &Character{}
}
