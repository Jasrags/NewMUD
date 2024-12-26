package rooms

type Area struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

func NewArea() *Area {
	return &Area{}
}
