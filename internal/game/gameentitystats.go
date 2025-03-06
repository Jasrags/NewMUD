package game

type GameEntityStats struct {
	Body      int     `yaml:"body"`
	Agility   int     `yaml:"agility"`
	Reaction  int     `yaml:"reaction"`
	Strength  int     `yaml:"strength"`
	Willpower int     `yaml:"willpower"`
	Logic     int     `yaml:"logic"`
	Intuition int     `yaml:"intuition"`
	Charisma  int     `yaml:"charisma"`
	Essence   float64 `yaml:"essence"`
	Magic     int     `yaml:"magic"`
	Resonance int     `yaml:"resonance"`
}

func (ges *GameEntityStats) GetBody() int {
	return ges.Body
}

func (ges *GameEntityStats) GetAgility() int {
	return ges.Agility
}

func (ges *GameEntityStats) GetReaction() int {
	return ges.Reaction
}

func (ges *GameEntityStats) GetStrength() int {
	return ges.Strength
}

func (ges *GameEntityStats) GetWillpower() int {
	return ges.Willpower
}

func (ges *GameEntityStats) GetLogic() int {
	return ges.Logic
}

func (ges *GameEntityStats) GetIntuition() int {
	return ges.Intuition
}

func (ges *GameEntityStats) GetCharisma() int {
	return ges.Charisma
}

func (ges *GameEntityStats) GetEssence() float64 {
	return ges.Essence
}

func (ges *GameEntityStats) GetMagic() int {
	return ges.Magic
}

func (ges *GameEntityStats) GetResonance() int {
	return ges.Resonance
}

// ATTRIBUTE-ONLY TESTS

// GetComposure calculates and returns the Composure of the character.
// Formula: (WIL + CHA)
func (ges *GameEntityStats) GetComposure() int {
	return ges.GetWillpower() + ges.GetCharisma()
}

// GetJudgeIntentions calculates and returns the Judge Intentions of the character.
// Formula: (INT + CHA)
func (ges *GameEntityStats) GetJudgeIntentions() int {
	return ges.GetIntuition() + ges.GetCharisma()
}

// GetMemory calculates and returns the Memory of the character.
// Formula: (LOG + WIL)
func (ges *GameEntityStats) GetMemory() int {
	return ges.GetLogic() + ges.GetWillpower()
}

// func (ges *GameEntityStats) GetArmorValue() int {
// 	return ges.GetBody()
// }

func (ges *GameEntityStats) GetAcidResistance() int {
	return ges.GetBody()
}

func (ges *GameEntityStats) GetColdResistance() int {
	return ges.GetBody()
}

func (ges *GameEntityStats) GetFallingResistance() int {
	return ges.GetBody()
}

func (ges *GameEntityStats) GetElectricityResistance() int {
	return ges.GetBody()
}

func (ges *GameEntityStats) GetFireResistance() int {
	return ges.GetBody()
}

func (ges *GameEntityStats) GetFatigueResistance() int {
	return ges.GetBody()
}

func (ges *GameEntityStats) GetLiftCarry() int {
	return ges.GetBody() * 10
}

func (ges *GameEntityStats) GetInitative() int {
	return ges.GetReaction() + ges.GetIntuition()
}

func (ges *GameEntityStats) GetInitativeDice() int {
	return 1
}
