package game

import "fmt"

const (
	PreGensFilepath = "_data/pregens"
)

type (
	Pregen struct {
		GameEntityInformation `yaml:",inline"`
		GameEntityStats       `yaml:",inline"`
		GameEntityDynamic     `yaml:",inline"`
	}
)

func NewPregen() *Pregen {
	return &Pregen{
		GameEntityDynamic: NewGameEntityDynamic(),
	}
}

func (p *Pregen) GetSelectionInfo() string {
	return fmt.Sprintf("%s [%s]", p.Title, p.MetatypeID)
}
