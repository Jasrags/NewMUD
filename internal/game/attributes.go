package game

type AttributeType string

// TODO: Revisit this, the Derivied attributes are not really attributes, they are calculated from the base attributes and should be stored in the character, not in the attributes.
const (
	AttributeBody            AttributeType = "Body"
	AttributeAgility         AttributeType = "Agility"
	AttributeReaction        AttributeType = "Reaction"
	AttributeStrength        AttributeType = "Strength"
	AttributeWillpower       AttributeType = "Willpower"
	AttributeLogic           AttributeType = "Logic"
	AttributeIntuition       AttributeType = "Intuition"
	AttributeCharisma        AttributeType = "Charisma"
	AttributeEssence         AttributeType = "Essence"
	AttributeMagic           AttributeType = "Magic"
	AttributeResonance       AttributeType = "Resonance"
	AttributeInitiative      AttributeType = "Initiative"
	AttributeInitiativeDice  AttributeType = "Initiative Dice"
	AttributeComposure       AttributeType = "Composure"
	AttributeJudgeIntentions AttributeType = "Judge Intentions"
	AttributeMemory          AttributeType = "Memory"
	AttributeLift            AttributeType = "Lift"
	AttributeCarry           AttributeType = "Carry"
	AttributeWalk            AttributeType = "Walk"
	AttributeRun             AttributeType = "Run"
	AttributeSwim            AttributeType = "Swim"
)

type (
	AttributeT[T int | float64] interface{}

	Attribute[T int | float64] struct {
		Name       string `yaml:"name"`
		Base       T      `yaml:"base"`
		Delta      T      `yaml:"delta"`
		TotalValue T      `yaml:"total_value"`
		Min        T      `yaml:"min"`
		Max        T      `yaml:"max"`
		AugMax     T      `yaml:"aug_max"`
	}

	// Attributes struct {

	// 	// Derived attributes
	// 	Initiative      Attribute[int] `yaml:"initiative"`
	// 	InitiativeDice  Attribute[int] `yaml:"initiative_dice"`
	// 	Composure       Attribute[int] `yaml:"composure"`
	// 	JudgeIntentions Attribute[int] `yaml:"judge_intentions"`
	// 	Memory          Attribute[int] `yaml:"memory"`
	// 	Lift            Attribute[int] `yaml:"lift"`
	// 	Carry           Attribute[int] `yaml:"carry"`
	// 	Walk            Attribute[int] `yaml:"walk"`
	// 	Run             Attribute[int] `yaml:"run"`
	// 	Swim            Attribute[int] `yaml:"swim"`
	// }
)

func NewAttribute[T int | float64](name AttributeType, base T) *Attribute[T] {
	return &Attribute[T]{
		Base: base,
	}
}

func (a *Attribute[T]) SetBase(value T) {
	a.Base = value
	a.Recalculate()
}

func (a *Attribute[T]) AddBase(value T) {
	a.Base += value
	a.Recalculate()
}

func (a *Attribute[T]) SubBase(value T) {
	a.Base -= value
	a.Recalculate()
}

func (a *Attribute[T]) SetDelta(value T) {
	a.Delta = value
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

func (a *Attribute[T]) SetMin(value T) {
	a.Min = value
}

func (a *Attribute[T]) SetMax(value T) {
	a.Max = value
}

func (a *Attribute[T]) SetAugMax(value T) {
	a.AugMax = value
}

func (a *Attribute[T]) Recalculate() {
	a.TotalValue = a.Base + a.Delta
}

func (a *Attribute[T]) Reset() {
	a.Base = 0
	a.Delta = 0
	a.TotalValue = 0
	a.Min = 0
	a.Max = 0
	a.AugMax = 0
}
