package game

import "fmt"

const (
	PreGensFilepath = "_data/pregens"
)

type (
	PregenSkill struct {
		ID     string `yaml:"id"`
		Rating int    `yaml:"rating"`
	}
	Pregen struct {
		ID          string        `yaml:"id"`
		Title       string        `yaml:"title"`
		Description string        `yaml:"description"`
		MetatypeID  string        `yaml:"metatype_id"`
		Attributes  Attributes    `yaml:"attributes"`
		Skills      []PregenSkill `yaml:"skills"`
		// Qualities   []string   `yaml:"qualities"`
	}
)

func (p *Pregen) GetSelectionInfo() string {
	return fmt.Sprintf("%s [%s]", p.Title, p.MetatypeID)
}
