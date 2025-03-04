package game

import (
	"log/slog"
	"math"

	"github.com/google/uuid"
)

const (
	PositionStanding    = "Standing"
	PositionSitting     = "Sitting"
	PositionKneeling    = "Kneeling"
	PositionLying       = "Lying"
	PositionProne       = "Prone"
	PositionCrouching   = "Crouching"
	PositionResting     = "Resting"
	PositionSleeping    = "Sleeping"
	PositionUnconscious = "Unconscious"

	SexMale      = "Male"
	SexFemale    = "Female"
	SexNonBinary = "Non-Binary"
)

type (
	PhysicalDamage struct {
		Current  int // Current damage boxes filled
		Max      int // Max damage boxes in the track
		Overflow int // Overflow boxes beyond the Physical track
	}
	StunDamage struct {
		Current int // Current damage boxes filled
		Max     int // Max damage boxes in the track
	}
	Edge struct {
		Max       int `yaml:"max"`       // Max edge points
		Available int `yaml:"available"` // Available edge points
	}
	GameEntity struct {
		ID              string                   `yaml:"id"`
		Name            string                   `yaml:"name"`
		Title           string                   `yaml:"title"`
		Description     string                   `yaml:"description"`
		LongDescription string                   `yaml:"long_description"`
		MetatypeID      string                   `yaml:"metatype_id"`
		Age             int                      `yaml:"age"`
		Sex             string                   `yaml:"sex"`
		Height          int                      `yaml:"height"`
		Weight          int                      `yaml:"weight"`
		StreetCred      int                      `yaml:"street_cred"`
		Notoriety       int                      `yaml:"notoriety"`
		PublicAwareness int                      `yaml:"public_awareness"`
		Body            int                      `yaml:"body"`
		Agility         int                      `yaml:"agility"`
		Reaction        int                      `yaml:"reaction"`
		Strength        int                      `yaml:"strength"`
		Willpower       int                      `yaml:"willpower"`
		Logic           int                      `yaml:"logic"`
		Intuition       int                      `yaml:"intuition"`
		Charisma        int                      `yaml:"charisma"`
		Essence         float64                  `yaml:"essence"`
		Magic           int                      `yaml:"magic"`
		Resonance       int                      `yaml:"resonance"`
		Edge            int                      `yaml:"edge"`
		PhysicalDamage  int                      `yaml:"physical_damage"`
		StunDamage      int                      `yaml:"stun_damage"`
		OverflowDamage  int                      `yaml:"overflow_damage"`
		Inventory       Inventory                `yaml:"inventory"`
		Equipment       map[string]*ItemInstance `yaml:"equipment"`
		Qualtities      map[string]*Quality      `yaml:"qualities"`
		Skills          map[string]*Skill        `yaml:"skills"`
		PositionState   string                   `yaml:"position_state"`
	}
)

func NewGameEntity() GameEntity {
	return GameEntity{
		ID:            uuid.New().String(),
		Equipment:     make(map[string]*ItemInstance),
		PositionState: PositionStanding,
		Qualtities:    make(map[string]*Quality),
		Skills:        make(map[string]*Skill),
	}
}

func (e *GameEntity) GetBody() int {
	return e.Body
}

func (e *GameEntity) GetAgility() int {
	return e.Agility
}

func (e *GameEntity) GetReaction() int {
	return e.Reaction
}

func (e *GameEntity) GetStrength() int {
	return e.Strength
}

func (e *GameEntity) GetWillpower() int {
	return e.Willpower
}

func (e *GameEntity) GetLogic() int {
	return e.Logic
}

func (e *GameEntity) GetIntuition() int {
	return e.Intuition
}

func (e *GameEntity) GetCharisma() int {
	return e.Charisma
}

func (e *GameEntity) GetEssence() float64 {
	return e.Essence
}

func (e *GameEntity) GetMagic() int {
	return e.Magic
}

func (e *GameEntity) GetResonance() int {
	return e.Resonance
}

// GetInitative calculates and returns the Initiative of the character.
// Formula: (Reaction + Intuition) + initiative bonuses
func (e *GameEntity) GetInitative() int {
	return (e.GetReaction() + e.GetIntuition())
}

// GetInitativeDice returns the Initiative Dice of the character.
// Formula: (1 + initiative_dice bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetInitativeDice() int {
	return 1
}

// LIMITS

// GetPhysicalLimit calculates and returns the Physical Limit of the character.
// Formula: [(Strength x 2) + Body + Reaction] / 3 (round up) + physical_limit modifiers
func (e *GameEntity) GetPhysicalLimit() int {
	strength := float64(e.GetStrength())
	body := float64(e.GetBody())
	reaction := float64(e.GetReaction())

	limit := (strength*2 + body + reaction) / 3.0

	return int(math.Ceil(limit)) // Round up
}

// GetAdjustedPhysicalLimit calculates and returns the adjusted physical limit of the character.
func (e *GameEntity) GetAdjustedPhysicalLimit() int {
	basePhysicalLimit := e.GetPhysicalLimit()
	penalty := e.GetEncumbrancePenalty()

	// Physical Limit cannot go below 1
	adjustedLimit := basePhysicalLimit - penalty
	if adjustedLimit < 1 {
		adjustedLimit = 1
	}

	// Formula: (STR + BOD) / 2
	return adjustedLimit
}

// GetSocialLimit calculates and returns the Social Limit of the character.
// Formula: [(Charisma x 2) + Willpower + Essence] / 3 (round up) + social_limit modifiers
func (e *GameEntity) GetSocialLimit() int {
	charisma := float64(e.GetCharisma())
	willpower := float64(e.GetWillpower())
	essence := e.GetEssence()
	limit := (charisma*2 + willpower + essence) / 3.0

	return int(math.Ceil(limit)) // Round up
}

// GetMentalLimit calculates and returns the Mental Limit of the character.
// Formula: [(Logic x 2) + Intuition + Willpower] / 3 (round up) + mental_limit modifiers
func (e *GameEntity) GetMentalLimit() int {
	logic := float64(e.GetLogic())
	intuition := float64(e.GetIntuition())
	willpower := float64(e.GetWillpower())
	limit := (logic*2 + intuition + willpower) / 3.0

	return int(math.Ceil(limit)) // Round up
}

// ATTRIBUTE-ONLY TESTS

// GetComposure calculates and returns the Composure of the character.
// Formula: (WIL + CHA)
func (e *GameEntity) GetComposure() int {
	return e.GetWillpower() + e.GetCharisma()
}

// GetJudgeIntentions calculates and returns the Judge Intentions of the character.
// Formula: (INT + CHA)
func (e *GameEntity) GetJudgeIntentions() int {
	return e.GetIntuition() + e.GetCharisma()
}

// TODO: Implement these rules when doing lift/carry tests
// The baseline for lifting weight is 15 kilograms per point of Strength. Anything more than that requires a Strength + Body Test. Each hit increases the max weight lifted by 15 kilograms. Lifting weight above your head, as with a clean & jerk, is more difficult. The baseline for lifting weight above the head is 5 kilograms per point Strength. Each hit on the Lifting Test increases the maximum weight you can lift by 5 kilograms.
// Carrying weight is significantly different than lifting weight. Characters can carry Strength x 10 kilograms in gear without effort. Additional weight requires a Lifting Test. Each hit increases the maximum by 10 kilograms.

// GetLiftCarry calculates and returns the Lift Carry of the character.
func (e *GameEntity) GetLiftCarry() float64 {
	baseCarryWeight := 10
	return float64(e.GetStrength()+e.GetBody()) * float64(baseCarryWeight)
}

// GetLiftWeight calculates and returns the Lift Weight of the character.
// Lift Formula: STR * 15
func (e *GameEntity) GetLiftWeight() float64 {
	baseWeight := 15

	// Carry Formula: STR * 10
	return float64(e.GetStrength() * baseWeight)
}

// GetCarryWeight calculates and returns the Carry Weight of the character.
// Carry Formula: STR * 10
func (e *GameEntity) GetCarryWeight() float64 {
	baseWeight := 10

	return float64(e.GetStrength() * baseWeight)
}

// GetCurrentCarryWeight calculates and returns the current carry weight of the character.
func (e *GameEntity) GetCurrentCarryWeight() float64 {
	totalWeight := 0.0

	for _, item := range e.Inventory.Items {
		if item.Blueprint == nil {
			slog.Warn("GetCurrentCarryWeight: item blueprint is nil",
				slog.String("item_id", item.InstanceID),
				slog.String("item_blueprint_id", item.BlueprintID))

			continue
		}

		totalWeight += item.Blueprint.Weight
	}

	return totalWeight
}

// GetEncumberancePenalty calculates and returns if the character is encumbered.
// TODO: Implement encumbered penatlies for combat
// If a character overburdens himself with gear, he suffers encumbrance modifiers. For every 15 kilograms (or part thereof) by which you exceed your carrying capacity, you suffer a –1 modifier to your Physical Limit (minimum limit of 1). This means that a character with Strength 3 (Carrying Capacity 30) that is trudging along with 50 kilograms of equipment suffers a –2 penalty to his Physical Limit.
func (e *GameEntity) IsEncumbered() bool {
	return e.GetEncumbrancePenalty() > 0
}

// GetEncumbrancePenalty calculates and returns the encumbrance penalty of the character.
// Formula: Excess weight / 15 kg (rounded up), -1 penalty for every 15 kg over capacity
func (e *GameEntity) GetEncumbrancePenalty() int {
	currentWeight := e.GetCurrentCarryWeight()
	maxCarryWeight := float64(e.GetLiftCarry())
	excessWeight := currentWeight - maxCarryWeight

	if excessWeight <= 0 {
		return 0 // No penalty if within carrying capacity
	}

	return int(math.Ceil(excessWeight / 15.0))
}

// GetMemory calculates and returns the Memory of the character.
// Formula: (LOG + WIL)
func (e *GameEntity) GetMemory() int {
	return e.GetLogic() + e.GetWillpower()
}

// TOXIN RESISTANCES
// Contact
// Ingestion
// Inhaliation
// Injection

// ADDICITION RESISTANCE
// Resist Physical Addiction
// Resist Psychological Addiction

// DAMAGE RESISTANCES
// GetArmorValue calculates and returns the Armor value of the character.
// Formula: (Body + armor bonuses)
func (e *GameEntity) GetArmorValue() int {
	totalArmorValueBonus := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		if item.Blueprint == nil {
			slog.Warn("GetArmorValue: item blueprint is nil",
				slog.String("item_id", item.InstanceID),
				slog.String("item_blueprint_id", item.BlueprintID))

			continue
		}

		if armor, ok := item.Blueprint.BaseStats["armor"]; ok {
			totalArmorValueBonus += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.GetBody() + totalArmorValueBonus
}

// Formula: (Body + acid_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetAcidResistance() int {
	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		if item.Blueprint == nil {
			slog.Warn("GetAcidResistance: item blueprint is nil",
				slog.String("item_id", item.InstanceID),
				slog.String("item_blueprint_id", item.BlueprintID))

			continue
		}

		if armor, ok := item.Blueprint.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.GetBody() + totalResistance
}

// GetColdResistance calculates and returns the Cold Resistance of the character.
// Formula: (Body + cold_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetColdResistance() int {
	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		if item.Blueprint == nil {
			slog.Warn("GetColdResistance: item blueprint is nil",
				slog.String("item_id", item.InstanceID),
				slog.String("item_blueprint_id", item.BlueprintID))

			continue
		}

		if armor, ok := item.Blueprint.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.GetBody() + totalResistance
}

// GetFallingResistance calculates and returns the Falling Resistance of the character.
// Formula: (Body + armor bonuses)
func (e *GameEntity) GetFallingResistance() int {
	return e.GetBody() + e.GetArmorValue()
}

// GetElectricityResistance calculates and returns the Electricity Resistance of the character.
// Formula: (Body + electrical_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetElectricityResistance() int {
	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		if item.Blueprint == nil {
			slog.Warn("GetElectricityResistance: item blueprint is nil",
				slog.String("item_id", item.InstanceID),
				slog.String("item_blueprint_id", item.BlueprintID))

			continue
		}

		if armor, ok := item.Blueprint.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.GetBody() + totalResistance
}

// GetFireResistance calculates and returns the Fire Resistance of the character.
// Formula: (Body + fire_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetFireResistance() int {
	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		if item.Blueprint == nil {
			slog.Warn("GetFireResistance: item blueprint is nil",
				slog.String("item_id", item.InstanceID),
				slog.String("item_blueprint_id", item.BlueprintID))

			continue
		}

		if armor, ok := item.Blueprint.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.GetBody() + totalResistance
}

// GetFatigueResistance calculates and returns the Fatigue Resistance of the character.
// Formula: (Body + Willpower) + fatigue_resistance bonuses
// TODO: Implement bonuses
func (e *GameEntity) GetFatigueResistance() int {
	return e.GetBody() + e.GetWillpower()
}

// DAMAGE APPLICATION (pg. 172)
// Physical Damage
// Stun Damage
// Cold Damage
// Fire Damage
// Electric Damage
// Fire Damage
// Falling Damage
// Fatigue Damage

// DAMAGE MONITORS

// GetPhysicalConditionMax calculates and returns the Physical Condition Max of the character.
// Formula: [Body / 2] + 8 (rounded up)
// TODO: Implement bonuses
func (e *GameEntity) GetPhysicalConditionMax() int {
	return int(math.Ceil(float64(e.GetBody())/2.0) + 8)
}

// GetStunConditionMax calculates and returns the Stun Condition Max of the character.
// Formula: [Willpower / 2] + 8 (rounded up)
// TODO: Implement bonuses
func (e *GameEntity) GetStunConditionMax() int {
	return int(math.Ceil(float64(e.GetWillpower())/2.0) + 8)
}

// GetOverflowConditionMax calculates and returns the Overflow Condition Max of the character.
// Formula: Body + Augmentation bonuses (rounded up)
// TODO: Implement bonuses
func (e *GameEntity) GetOverflowConditionMax() int {
	return e.GetBody()
}

// DEFENSES
// Ranged attacks against you
// Raged Defense
//  - Full Defense
// Melee attacks against you
// Melee Defense
//  - Full Defense
//  - Club parry
//  - Knife parry
//  - Unarmed Strike Block
// Sensor-aided attacks against you
// Sensro Defense

// EDGE
// UseEdge - Decreases the available Edge by 1.
// func (e *GameEntity) UseEdge() bool {
// 	if e.Edge.Available <= 0 || e.Edge.Max <= 0 {
// 		return false
// 	}

// 	e.Edge.Available -= 1
// 	return true
// }

// // Burn Edge - Decreases the maximum Edge by 1 and the available Edge by 1.
// func (e *GameEntity) BurnEdge() bool {
// 	if e.Edge.Max <= 0 || e.Edge.Available <= 0 {
// 		return false
// 	}

// 	e.Edge.Max -= 1
// 	e.Edge.Available -= 1

// 	return true
// }

// // Regain Edge - Increases the available Edge by 1.
// func (e *GameEntity) RegainEdge() bool {
// 	if e.Edge.Available >= e.Edge.Max {
// 		return false
// 	}

// 	e.Edge.Available += 1

// 	return true
// }

// Damage
// Physical Damage Level		(Body + Armor) / 2

// The Physical Condition Monitor has boxes equal to half the character’s ((BOD/2) + 8) Body (rounded up) + 8;
// the Stun Condition Monitor has boxes equaling half the character’s ((WILL/2) + 8)Willpower (rounded up) + 8.
// When a row of the Condition Monitor is filled up, the player character takes a –1 penalty to all subsequent tests. This penalty stacks for each row of the Condition Monitor that is filled in.

// TODO: implement movement
func (e *GameEntity) GetMovement() int {
	// e.Attributes.Recalculate()
	// Formula: (Reaction + Agility) / 2
	return (e.GetReaction() + e.GetAgility()) / 2
}

// func (e *GameEntity) ApplyStunDamage(damage int) {
// 	e.StunDamage.Current += damage

// 	// Check for Stun overflow
// 	if e.StunDamage.Current > e.StunDamage.Max {
// 		excessStun := e.StunDamage.Current - e.StunDamage.Max
// 		e.StunDamage.Current = e.StunDamage.Max

// 		// Convert Stun overflow to Physical damage
// 		physicalOverflow := excessStun / 2 // Every 2 Stun = 1 Physical
// 		e.ApplyPhysicalDamage(physicalOverflow)
// 	}
// }

// func (e *GameEntity) ApplyPhysicalDamage(damage int) {
// 	remainingCapacity := e.PhysicalDamage.Max - e.PhysicalDamage.Current

// 	if damage > remainingCapacity {
// 		// Overflow logic
// 		e.PhysicalDamage.Overflow += damage - remainingCapacity
// 		e.PhysicalDamage.Current = e.PhysicalDamage.Max
// 	} else {
// 		e.PhysicalDamage.Current += damage
// 	}
// }

/*
Injuries cause pain, bleeding, and other distractions that interfere with doing all sorts of actions. Wound modifiers are accumulated with every third box of damage and are cumulative between damage tracks and with other negative modifiers such as spells or adverse conditions.

Wound modifiers are applied to all tests not about reducing the number of boxes you’re about to take on your Condition Monitor (such as damage resistance, resisting direct combat spells, toxin resistance, and so on). The Wound Modifier penalty is also applied to the character’s Initiative attribute and therefore their Initiative Score during combat.
*/
func (e *GameEntity) GetWoundModifiers() int {
	// Formula: [Body / 2] + 1 (rounded up)
	return int(math.Ceil(float64(e.GetBody())/2.0) + 1)
}

// Attributes

func (e *GameEntity) RollInitative() int {
	poolSize := 1
	_, _, results := RollDice(poolSize)
	// Formula: (Reaction + Intuition) + 1D6
	return e.GetInitative() + RollResultsTotal(results)
}

func (e *GameEntity) GetAdjustedInitiative() int {
	baseInitiative := e.GetInitative()
	return baseInitiative + e.GetWoundModifiers()
}

func (e *GameEntity) ApplyWoundModifiers(baseDicePool int) int {
	return baseDicePool + e.GetWoundModifiers()
}

// Astral Initiative					(Intuition x 2) + 2D6					—
// Matrix AR Initiative				(Reaction + Intuition) + 1D6				—
// Matrix VR Initiative (Hot Sim)		(Data Processing + Intuition) + 4D6		—
// Matrix VR Initiative (Cold Sim)		(Data Processing + Intuition) + 3D6		—

// Condition Monitor Boxes
// Physical 							[Body x 2] + 8												Add bonuses to Body before calculating; round up final results
// Stun								[Willpower x 2] + 8											Add bonuses to Willpower before calculating; round up final results
// Overflow							Body + Augmentation bonuses									-
