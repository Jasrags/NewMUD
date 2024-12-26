package items

type Type string

const (
	TypeJunk Type = "junk"
)

type Item struct {
	ID          string `yaml:"id"`
	UUID        string `yaml:"uuid"`
	AreaID      string `yaml:"area_id"`
	RoomID      string `yaml:"room_id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        Type   `yaml:"type"`
}

func NewItem() *Item {
	return &Item{}
}
