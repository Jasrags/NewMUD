package main

type AttributeType string

const (
	AttributeBody      AttributeType = "Body"
	AttributeAgility   AttributeType = "Agility"
	AttributeReaction  AttributeType = "Reaction"
	AttributeStrength  AttributeType = "Strength"
	AttributeWillpower AttributeType = "Willpower"
	AttributeLogic     AttributeType = "Logic"
	AttributeIntuition AttributeType = "Intuition"
	AttributeCharisma  AttributeType = "Charisma"
	AttributeEdge      AttributeType = "Edge"
	AttributeEssence   AttributeType = "Essence"
	AttributeMagic     AttributeType = "Magic"
	AttributeResonance AttributeType = "Resonance"
)

type Attributes struct {
	Body      Attribute[int]     `yaml:"body"`
	Agility   Attribute[int]     `yaml:"agility"`
	Reaction  Attribute[int]     `yaml:"reaction"`
	Strength  Attribute[int]     `yaml:"strength"`
	Willpower Attribute[int]     `yaml:"willpower"`
	Logic     Attribute[int]     `yaml:"logic"`
	Intuition Attribute[int]     `yaml:"intuition"`
	Charisma  Attribute[int]     `yaml:"charisma"`
	Edge      Attribute[int]     `yaml:"edge"`
	Essence   Attribute[float64] `yaml:"essence"`
	Magic     Attribute[int]     `yaml:"magic"`
	Resonance Attribute[int]     `yaml:"resonance"`
}

func (a *Attributes) Recalculate() {
	a.Body.Recalculate()
	a.Agility.Recalculate()
	a.Reaction.Recalculate()
	a.Strength.Recalculate()
	a.Willpower.Recalculate()
	a.Logic.Recalculate()
	a.Intuition.Recalculate()
	a.Charisma.Recalculate()
	a.Edge.Recalculate()
	a.Essence.Recalculate()
	a.Magic.Recalculate()
	a.Resonance.Recalculate()
}

func (a *Attributes) Reset() {
	a.Body.Reset()
	a.Agility.Reset()
	a.Reaction.Reset()
	a.Strength.Reset()
	a.Willpower.Reset()
	a.Logic.Reset()
	a.Intuition.Reset()
	a.Charisma.Reset()
	a.Edge.Reset()
	a.Essence.Reset()
	a.Magic.Reset()
	a.Resonance.Reset()
}

type AttributeT[T int | float64] interface{}

type Attribute[T int | float64] struct {
	Name       string `yaml:"name"`
	Base       T      `yaml:"base"`
	Delta      T      `yaml:"delta"`
	TotalValue T      `yaml:"total_value"`
}

func NewAttribute[T int | float64](name AttributeType, base T) *Attribute[T] {
	return &Attribute[T]{
		Base: base,
	}
}

func (a *Attribute[T]) SetBase(value T) {
	a.Base = value
	a.Recalculate()
}

func (a *Attribute[T]) AddDelta(value T) {
	a.Delta += value
	a.Recalculate()
}

func (a *Attribute[T]) SubDelta(value T) {
	a.Delta -= value
	a.Recalculate()
}

func (a *Attribute[T]) Recalculate() {
	a.TotalValue = a.Base + a.Delta
}

func (a *Attribute[T]) Reset() {
	a.Base = 0
	a.Delta = 0
	a.TotalValue = 0
}
