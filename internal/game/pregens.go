package game

const (
	PreGensFilepath = "_data/pregens"
)

type (
	PreGen struct {
		ID          string     `yaml:"id"`
		Title       string     `yaml:"title"`
		Description string     `yaml:"description"`
		MetatypeID  string     `yaml:"metatype_id"`
		Attributes  Attributes `yaml:"attributes"`
		Skills      []string   `yaml:"skills"`
		Qualities   []string   `yaml:"qualities"`
	}
)
