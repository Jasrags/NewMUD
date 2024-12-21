package items

type Item struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewItem() *Item {
	return &Item{}
}
