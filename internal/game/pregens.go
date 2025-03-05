package game

import "fmt"

const (
	PreGensFilepath = "_data/pregens"
)

type (
	Pregen struct {
		GameEntity `yaml:",inline"`
	}
)

func (p *Pregen) GetSelectionInfo() string {
	return fmt.Sprintf("%s [%s]", p.Title, p.MetatypeID)
}
