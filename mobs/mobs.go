package mobs

type Mob struct {
	ID          string `yaml:"id"`
	UUID        string `yaml:"uuid"`
	AreaID      string `yaml:"area_id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func NewMob() *Mob {
	return &Mob{}
}
