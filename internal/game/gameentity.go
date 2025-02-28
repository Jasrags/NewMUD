package game

import (
	"math"
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
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
		sync.RWMutex `yaml:"-"`
		Listeners    []ee.Listener `yaml:"-"`

		// Information
		ID              string `yaml:"id"`
		Name            string `yaml:"name"`
		Title           string `yaml:"title"`
		Description     string `yaml:"description"`
		LongDescription string `yaml:"long_description"`
		MetatypeID      string `yaml:"metatype_id"`
		Age             int    `yaml:"age"`
		Sex             string `yaml:"sex"`
		Height          int    `yaml:"height"`
		Weight          int    `yaml:"weight"`
		StreetCred      int    `yaml:"street_cred"`
		Notoriety       int    `yaml:"notoriety"`
		PublicAwareness int    `yaml:"public_awareness"`
		Edge            Edge   `yaml:"edge"`
		// Attributes
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
		// Damage
		PhysicalDamage PhysicalDamage `yaml:"physical_damage"`
		StunDamage     StunDamage     `yaml:"stun_damage"`
		// Other
		Room   *Room  `yaml:"-"` // TODO: move away from using the Room pointer and using the RoomID to reference the room via EntityMgr
		RoomID string `yaml:"room_id"`
		// Area            *Area            `yaml:"-"`
		// AreaID     string           `yaml:"area_id"`
		Inventory     Inventory          `yaml:"inventory"`
		Equipment     map[string]*Item   `yaml:"equipment"`
		Qualtities    map[string]Quality `yaml:"qualities"`
		Skills        map[string]Skill   `yaml:"skills"`
		PositionState string             `yaml:"position_state"`
	}
)

func NewGameEntity() GameEntity {
	return GameEntity{
		ID:            uuid.New().String(),
		Equipment:     make(map[string]*Item),
		Listeners:     make([]ee.Listener, 0),
		PositionState: PositionStanding,
		Qualtities:    make(map[string]Quality),
		Skills:        make(map[string]Skill),
	}
}

func (e *GameEntity) SetRoom(room *Room) {
	e.Room = room
	e.RoomID = room.ReferenceID
}

// Recalculate triggers the recalculation of all attributes and derivied values.
func (e *GameEntity) Recalculate() {
	// Start with base attributes
	e.Body.Recalculate()
	e.Agility.Recalculate()
	e.Reaction.Recalculate()
	e.Strength.Recalculate()
	e.Willpower.Recalculate()
	e.Logic.Recalculate()
	e.Intuition.Recalculate()
	e.Charisma.Recalculate()
	e.Essence.Recalculate()
	e.Magic.Recalculate()
	e.Resonance.Recalculate()
}

// GetInitative calculates and returns the Initiative of the character.
// Formula: (Reaction + Intuition) + initiative bonuses
func (e *GameEntity) GetInitative() int {
	e.Recalculate()

	return (e.Reaction.TotalValue + e.Intuition.TotalValue)
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
	e.Recalculate()

	strength := float64(e.Strength.TotalValue)
	body := float64(e.Body.TotalValue)
	reaction := float64(e.Reaction.TotalValue)

	limit := (strength*2 + body + reaction) / 3.0

	return int(math.Ceil(limit)) // Round up
}

// GetAdjustedPhysicalLimit calculates and returns the adjusted physical limit of the character.
// func (e *GameEntity) GetAdjustedPhysicalLimit() int {
// 	// e.Attributes.Recalculate()
// 	basePhysicalLimit := e.GetPhysicalLimit()
// 	penalty := e.GetEncumbrancePenalty()

// 	// Physical Limit cannot go below 1
// 	adjustedLimit := basePhysicalLimit - penalty
// 	if adjustedLimit < 1 {
// 		adjustedLimit = 1
// 	}

// 	// Formula: (STR + BOD) / 2
// 	return adjustedLimit
// }

// GetSocialLimit calculates and returns the Social Limit of the character.
// Formula: [(Charisma x 2) + Willpower + Essence] / 3 (round up) + social_limit modifiers
func (e *GameEntity) GetSocialLimit() int {
	e.Recalculate()

	charisma := float64(e.Charisma.TotalValue)
	willpower := float64(e.Willpower.TotalValue)
	essence := e.Essence.TotalValue
	limit := (charisma*2 + willpower + essence) / 3.0

	return int(math.Ceil(limit)) // Round up
}

// GetMentalLimit calculates and returns the Mental Limit of the character.
// Formula: [(Logic x 2) + Intuition + Willpower] / 3 (round up) + mental_limit modifiers
func (e *GameEntity) GetMentalLimit() int {
	e.Recalculate()

	logic := float64(e.Logic.TotalValue)
	intuition := float64(e.Intuition.TotalValue)
	willpower := float64(e.Willpower.TotalValue)
	limit := (logic*2 + intuition + willpower) / 3.0

	return int(math.Ceil(limit)) // Round up
}

// ATTRIBUTE-ONLY TESTS

// GetComposure calculates and returns the Composure of the character.
// Formula: (WIL + CHA)
func (e *GameEntity) GetComposure() int {
	e.Recalculate()

	return e.Willpower.TotalValue + e.Charisma.TotalValue
}

// GetJudgeIntentions calculates and returns the Judge Intentions of the character.
// Formula: (INT + CHA)
func (e *GameEntity) GetJudgeIntentions() int {
	e.Recalculate()

	return e.Intuition.TotalValue + e.Charisma.TotalValue
}

// TODO: Implement these rules when doing lift/carry tests
// The baseline for lifting weight is 15 kilograms per point of Strength. Anything more than that requires a Strength + Body Test. Each hit increases the max weight lifted by 15 kilograms. Lifting weight above your head, as with a clean & jerk, is more difficult. The baseline for lifting weight above the head is 5 kilograms per point Strength. Each hit on the Lifting Test increases the maximum weight you can lift by 5 kilograms.
// Carrying weight is significantly different than lifting weight. Characters can carry Strength x 10 kilograms in gear without effort. Additional weight requires a Lifting Test. Each hit increases the maximum by 10 kilograms.

// GetLiftCarry calculates and returns the Lift Carry of the character.
func (e *GameEntity) GetLiftCarry() float64 {
	e.Recalculate()

	baseCarryWeight := 10
	return float64(e.Strength.TotalValue+e.Body.TotalValue) * float64(baseCarryWeight)
}

// GetLiftWeight calculates and returns the Lift Weight of the character.
// Lift Formula: STR * 15
func (e *GameEntity) GetLiftWeight() float64 {
	e.Recalculate()

	baseWeight := 15

	// Carry Formula: STR * 10
	return float64(e.Strength.TotalValue * baseWeight)
}

// GetCarryWeight calculates and returns the Carry Weight of the character.
// Carry Formula: STR * 10
func (e *GameEntity) GetCarryWeight() float64 {
	e.Recalculate()

	baseWeight := 10

	return float64(e.Strength.TotalValue * baseWeight)
}

// GetCurrentCarryWeight calculates and returns the current carry weight of the character.
func (e *GameEntity) GetCurrentCarryWeight() float64 {
	totalWeight := 0.0

	for _, item := range e.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp != nil {
			totalWeight += bp.Weight
		}
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
	e.Recalculate()

	return e.Logic.TotalValue + e.Willpower.TotalValue
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
	e.Recalculate()

	totalArmorValueBonus := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if armor, ok := bp.BaseStats["armor"]; ok {
			totalArmorValueBonus += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.Body.TotalValue + totalArmorValueBonus
}

// Formula: (Body + acid_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetAcidResistance() int {
	e.Recalculate()

	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if armor, ok := bp.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.Body.TotalValue + totalResistance
}

// GetColdResistance calculates and returns the Cold Resistance of the character.
// Formula: (Body + cold_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetColdResistance() int {
	e.Recalculate()

	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if armor, ok := bp.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.Body.TotalValue + totalResistance
}

// GetFallingResistance calculates and returns the Falling Resistance of the character.
// Formula: (Body + armor bonuses)
func (e *GameEntity) GetFallingResistance() int {
	e.Recalculate()

	return e.Body.TotalValue + e.GetArmorValue()
}

// GetElectricityResistance calculates and returns the Electricity Resistance of the character.
// Formula: (Body + electrical_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetElectricityResistance() int {
	e.Recalculate()

	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if armor, ok := bp.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.Body.TotalValue + totalResistance
}

// GetFireResistance calculates and returns the Fire Resistance of the character.
// Formula: (Body + fire_resistance bonuses)
// TODO: Implement bonuses
func (e *GameEntity) GetFireResistance() int {
	e.Recalculate()

	totalResistance := 0

	// Check equiped items for armor
	for _, item := range e.Equipment {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if armor, ok := bp.BaseStats["acid_resistance"]; ok {
			totalResistance += armor
		}
	}
	// TODO: Check for other bonuses (qualtites, racial traits, spells, etc)

	return e.Body.TotalValue + totalResistance
}

// GetFatigueResistance calculates and returns the Fatigue Resistance of the character.
// Formula: (Body + Willpower) + fatigue_resistance bonuses
// TODO: Implement bonuses
func (e *GameEntity) GetFatigueResistance() int {
	e.Recalculate()

	return e.Body.TotalValue + e.Willpower.TotalValue
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
	e.Recalculate()

	return int(math.Ceil(float64(e.Body.TotalValue)/2.0) + 8)
}

// GetStunConditionMax calculates and returns the Stun Condition Max of the character.
// Formula: [Willpower / 2] + 8 (rounded up)
// TODO: Implement bonuses
func (e *GameEntity) GetStunConditionMax() int {
	e.Recalculate()

	return int(math.Ceil(float64(e.Willpower.TotalValue)/2.0) + 8)
}

// GetOverflowConditionMax calculates and returns the Overflow Condition Max of the character.
// Formula: Body + Augmentation bonuses (rounded up)
// TODO: Implement bonuses
func (e *GameEntity) GetOverflowConditionMax() int {
	e.Recalculate()

	return e.Body.TotalValue
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
func (e *GameEntity) UseEdge() bool {
	if e.Edge.Available <= 0 || e.Edge.Max <= 0 {
		return false
	}

	e.Edge.Available -= 1
	return true
}

// Burn Edge - Decreases the maximum Edge by 1 and the available Edge by 1.
func (e *GameEntity) BurnEdge() bool {
	if e.Edge.Max <= 0 || e.Edge.Available <= 0 {
		return false
	}

	e.Edge.Max -= 1
	e.Edge.Available -= 1

	return true
}

// Regain Edge - Increases the available Edge by 1.
func (e *GameEntity) RegainEdge() bool {
	if e.Edge.Available >= e.Edge.Max {
		return false
	}

	e.Edge.Available += 1

	return true
}

// Damage
// Physical Damage Level		(Body + Armor) / 2

// The Physical Condition Monitor has boxes equal to half the character’s ((BOD/2) + 8) Body (rounded up) + 8;
// the Stun Condition Monitor has boxes equaling half the character’s ((WILL/2) + 8)Willpower (rounded up) + 8.
// When a row of the Condition Monitor is filled up, the player character takes a –1 penalty to all subsequent tests. This penalty stacks for each row of the Condition Monitor that is filled in.

// TODO: implement movement
func (e *GameEntity) GetMovement() int {
	// e.Attributes.Recalculate()
	// Formula: (Reaction + Agility) / 2
	return (e.Reaction.TotalValue + e.Agility.TotalValue) / 2
}

func (e *GameEntity) ApplyStunDamage(damage int) {
	e.StunDamage.Current += damage

	// Check for Stun overflow
	if e.StunDamage.Current > e.StunDamage.Max {
		excessStun := e.StunDamage.Current - e.StunDamage.Max
		e.StunDamage.Current = e.StunDamage.Max

		// Convert Stun overflow to Physical damage
		physicalOverflow := excessStun / 2 // Every 2 Stun = 1 Physical
		e.ApplyPhysicalDamage(physicalOverflow)
	}
}

func (e *GameEntity) ApplyPhysicalDamage(damage int) {
	remainingCapacity := e.PhysicalDamage.Max - e.PhysicalDamage.Current

	if damage > remainingCapacity {
		// Overflow logic
		e.PhysicalDamage.Overflow += damage - remainingCapacity
		e.PhysicalDamage.Current = e.PhysicalDamage.Max
	} else {
		e.PhysicalDamage.Current += damage
	}
}

/*
Injuries cause pain, bleeding, and other distractions that interfere with doing all sorts of actions. Wound modifiers are accumulated with every third box of damage and are cumulative between damage tracks and with other negative modifiers such as spells or adverse conditions.

Wound modifiers are applied to all tests not about reducing the number of boxes you’re about to take on your Condition Monitor (such as damage resistance, resisting direct combat spells, toxin resistance, and so on). The Wound Modifier penalty is also applied to the character’s Initiative attribute and therefore their Initiative Score during combat.
*/
func (e *GameEntity) GetWoundModifiers() int {
	// e.Attributes.Recalculate()
	// Formula: [Body / 2] + 1 (rounded up)
	return int(math.Ceil(float64(e.Body.TotalValue)/2.0) + 1)
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
