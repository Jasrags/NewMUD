package game

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
	AttributeEssence   AttributeType = "Essence"
	AttributeMagic     AttributeType = "Magic"
	AttributeResonance AttributeType = "Resonance"
)

type Attributes struct {
	// Base attributes
	Body      Attribute[int]     `yaml:"body"`
	Agility   Attribute[int]     `yaml:"agility"`
	Reaction  Attribute[int]     `yaml:"reaction"`
	Strength  Attribute[int]     `yaml:"strength"`
	Willpower Attribute[int]     `yaml:"willpower"`
	Logic     Attribute[int]     `yaml:"logic"`
	Intuition Attribute[int]     `yaml:"intuition"`
	Charisma  Attribute[int]     `yaml:"charisma"`
	Essence   Attribute[float64] `yaml:"essence"`
	Magic     Attribute[int]     `yaml:"magic"`
	Resonance Attribute[int]     `yaml:"resonance"`
	// Derived attributes
	Initiative      Attribute[int] `yaml:"initiative"`
	InitiativeDice  Attribute[int] `yaml:"initiative_dice"`
	Composure       Attribute[int] `yaml:"composure"`
	JudgeIntentions Attribute[int] `yaml:"judge_intentions"`
	Memory          Attribute[int] `yaml:"memory"`
	Lift            Attribute[int] `yaml:"lift"`
	Carry           Attribute[int] `yaml:"carry"`
	Walk            Attribute[int] `yaml:"walk"`
	Run             Attribute[int] `yaml:"run"`
	Swim            Attribute[int] `yaml:"swim"`
}

func NewAttributes() Attributes {
	return Attributes{
		// Base attributes
		Body:      Attribute[int]{Name: "Body"},
		Agility:   Attribute[int]{Name: "Agility"},
		Reaction:  Attribute[int]{Name: "Reaction"},
		Strength:  Attribute[int]{Name: "Strength"},
		Willpower: Attribute[int]{Name: "Willpower"},
		Logic:     Attribute[int]{Name: "Logic"},
		Intuition: Attribute[int]{Name: "Intuition"},
		Charisma:  Attribute[int]{Name: "Charisma"},
		Essence:   Attribute[float64]{Name: "Essence"},
		Magic:     Attribute[int]{Name: "Magic"},
		Resonance: Attribute[int]{Name: "Resonance"},
		// Derived attributes
		Initiative:      Attribute[int]{Name: "Initiative"},
		InitiativeDice:  Attribute[int]{Name: "Initiative Dice"},
		Composure:       Attribute[int]{Name: "Composure"},
		JudgeIntentions: Attribute[int]{Name: "Judge Intentions"},
		Memory:          Attribute[int]{Name: "Memory"},
		Lift:            Attribute[int]{Name: "Lift"},
		Carry:           Attribute[int]{Name: "Carry"},
		Walk:            Attribute[int]{Name: "Walk"},
		Run:             Attribute[int]{Name: "Run"},
		Swim:            Attribute[int]{Name: "Swim"},
	}
}

func (a *Attributes) Recalculate() {
	// Base attributes
	a.Body.Recalculate()
	a.Agility.Recalculate()
	a.Reaction.Recalculate()
	a.Strength.Recalculate()
	a.Willpower.Recalculate()
	a.Logic.Recalculate()
	a.Intuition.Recalculate()
	a.Charisma.Recalculate()
	a.Essence.Recalculate()
	a.Magic.Recalculate()
	a.Resonance.Recalculate()
	// Derived attributes
	a.Initiative.Recalculate()
	a.InitiativeDice.Recalculate()
	a.Composure.Recalculate()
	a.JudgeIntentions.Recalculate()
	a.Memory.Recalculate()
	a.Lift.Recalculate()
	a.Carry.Recalculate()
	a.Walk.Recalculate()
	a.Run.Recalculate()
	a.Swim.Recalculate()
}

// func (a *Attributes) Reset() {
// 	a.Body.Reset()
// 	a.Agility.Reset()
// 	a.Reaction.Reset()
// 	a.Strength.Reset()
// 	a.Willpower.Reset()
// 	a.Logic.Reset()
// 	a.Intuition.Reset()
// 	a.Charisma.Reset()
// 	a.Essence.Reset()
// 	a.Magic.Reset()
// 	a.Resonance.Reset()
// }

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
