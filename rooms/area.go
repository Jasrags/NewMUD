package rooms

type Area struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
}

func NewArea() *Area {
	return &Area{}
}
