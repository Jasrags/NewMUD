package game

import (
	"log/slog"
	"math"
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type GameEntity struct {
	sync.RWMutex `yaml:"-"`
	Listeners    []ee.Listener `yaml:"-"`

	ID              string           `yaml:"id"`
	Name            string           `yaml:"name"`
	Title           string           `yaml:"title"`
	Metatype        string           `yaml:"metatype"`
	Age             int              `yaml:"age"`
	Sex             string           `yaml:"sex"`
	Height          int              `yaml:"height"`
	Weight          int              `yaml:"weight"`
	Ethnicity       string           `yaml:"ethnicity"`
	StreetCred      int              `yaml:"street_cred"`
	Notoriety       int              `yaml:"notoriety"`
	PublicAwareness int              `yaml:"public_awareness"`
	Karma           int              `yaml:"karma"`
	TotalKarma      int              `yaml:"total_karma"`
	Description     string           `yaml:"description"`
	Attributes      Attributes       `yaml:"attributes"`
	PhysicalDamage  PhysicalDamage   `yaml:"physical_damage"`
	StunDamage      StunDamage       `yaml:"stun_damage"`
	Edge            Edge             `yaml:"edge"`
	Room            *Room            `yaml:"-"`
	RoomID          string           `yaml:"room_id"`
	Area            *Area            `yaml:"-"`
	AreaID          string           `yaml:"area_id"`
	Inventory       Inventory        `yaml:"inventory"`
	Equipment       map[string]*Item `yaml:"equipment"`
}

func NewGameEntity() GameEntity {
	return GameEntity{
		ID:        uuid.New().String(),
		Equipment: make(map[string]*Item),
		Listeners: make([]ee.Listener, 0),
	}
}

func (e *GameEntity) GetName() string {
	return e.Name
}

func (e *GameEntity) GetID() string {
	return e.ID
}

func (e *GameEntity) SetRoom(room *Room) {
	e.Room = room
	e.RoomID = room.ReferenceID
}

// Recalculate triggers the recalculation of all attributes and derivied values.
func (e *GameEntity) Recalculate() {
	// Start with attributes
	e.Attributes.Recalculate()

	// Then, recalculate derived values
	// Composure
	e.Attributes.Composure.Base = (e.Attributes.Charisma.TotalValue + e.Attributes.Willpower.TotalValue)
	// Judge Intentions
	e.Attributes.JudgeIntentions.Base = (e.Attributes.Intuition.TotalValue + e.Attributes.Charisma.TotalValue)
	// Lift
	e.Attributes.Lift.Base = (e.Attributes.Body.TotalValue + e.Attributes.Strength.TotalValue) * 15
	// Carry
	e.Attributes.Carry.Base = (e.Attributes.Body.TotalValue + e.Attributes.Strength.TotalValue) * 10
	// Memory
	e.Attributes.Memory.Base = (e.Attributes.Logic.TotalValue + e.Attributes.Willpower.TotalValue)
	// Initiative
	e.Attributes.Initiative.Base = (e.Attributes.Reaction.TotalValue + e.Attributes.Intuition.TotalValue)
	// Walk Rate=Base Walk Rate (by metatype)+Agility
	// Run Rate=Base Run Rate (by metatype)+(Agility×2)
	// Swim
	e.Attributes.Swim.Base = e.Attributes.Agility.TotalValue
	e.Attributes.Recalculate()
}

// UseEdge - Decreases the available Edge by 1.
func (e *GameEntity) UseEdge() bool {
	if e.Edge.Available <= 0 || e.Edge.Max <= 0 {
		slog.Warn("Cannot use Edge. No available edge to use.",
			slog.String("character_id", e.ID),
			slog.Int("max_edge_points", e.Edge.Max),
			slog.Int("remaining_edge_points", e.Edge.Available))
		return false
	}

	e.Edge.Available -= 1
	slog.Debug("Edge point used",
		slog.String("character_id", e.ID),
		slog.Int("max_edge_points", e.Edge.Max),
		slog.Int("remaining_edge_points", e.Edge.Available))
	return true
}

// Burn Edge - Decreases the maximum Edge by 1 and the available Edge by 1.
func (e *GameEntity) BurnEdge() bool {
	if e.Edge.Max <= 0 || e.Edge.Available <= 0 {
		slog.Warn("Cannot burn Edge. Edge attribute is already zero.",
			slog.String("character_id", e.ID),
			slog.Int("max_edge_points", e.Edge.Max),
			slog.Int("remaining_edge_points", e.Edge.Available))
		return false
	}

	e.Edge.Max -= 1
	e.Edge.Available -= 1

	slog.Debug("Edge point burned",
		slog.String("character_id", e.ID),
		slog.Int("max_edge_points", e.Edge.Max),
		slog.Int("remaining_edge_points", e.Edge.Available))
	return true
}

// Regain Edge - Increases the available Edge by 1.
func (e *GameEntity) RegainEdge() bool {
	if e.Edge.Available >= e.Edge.Max {
		slog.Warn("Cannot regain Edge. Edge attribute is already at max.",
			slog.String("character_id", e.ID),
			slog.Int("max_edge_points", e.Edge.Max),
			slog.Int("remaining_edge_points", e.Edge.Available))
		return false
	}

	e.Edge.Available += 1
	slog.Debug("Edge point regained",
		slog.String("character_id", e.ID),
		slog.Int("max_edge_points", e.Edge.Max),
		slog.Int("remaining_edge_points", e.Edge.Available))

	return true
}

// Damage
// Physical Damage Level		(Body + Armor) / 2

// The Physical Condition Monitor has boxes equal to half the character’s ((BOD/2) + 8) Body (rounded up) + 8;
// the Stun Condition Monitor has boxes equaling half the character’s ((WILL/2) + 8)Willpower (rounded up) + 8.
// When a row of the Condition Monitor is filled up, the player character takes a –1 penalty to all subsequent tests. This penalty stacks for each row of the Condition Monitor that is filled in.

// GetPhysicalConditionMax calculates and returns the Physical Condition Max of the character.
func (e *GameEntity) GetPhysicalConditionMax() int {
	e.Attributes.Recalculate()
	// Formula: [Body / 2] + 8 (rounded up)
	// TODO: Add bonuses to Body before calculating; round up final results
	return int(math.Ceil(float64(e.Attributes.Body.TotalValue)/2.0) + 8)
}

// GetStunConditionMax calculates and returns the Stun Condition Max of the character.
func (e *GameEntity) GetStunConditionMax() int {
	e.Attributes.Recalculate()
	// Formula: [Willpower / 2] + 8 (rounded up)
	// TODO: Add bonuses to Willpower before calculating; round up final results
	return int(math.Ceil(float64(e.Attributes.Willpower.TotalValue)/2.0) + 8)
}

// GetOverflowConditionMax calculates and returns the Overflow Condition Max of the character.
func (e *GameEntity) GetOverflowConditionMax() int {
	e.Attributes.Recalculate()
	// Formula: Body + Augmentation bonuses (rounded up)
	// TODO: Implement Augmentation bonuses
	return e.Attributes.Body.TotalValue
}

type PhysicalDamage struct {
	Current  int // Current damage boxes filled
	Max      int // Max damage boxes in the track
	Overflow int // Overflow boxes beyond the Physical track
}

type StunDamage struct {
	Current int // Current damage boxes filled
	Max     int // Max damage boxes in the track
}

type Edge struct {
	Max       int `yaml:"max"`
	Available int `yaml:"available"`
}

// TODO: implement movement
func (e *GameEntity) GetMovement() int {
	e.Attributes.Recalculate()
	// Formula: (Reaction + Agility) / 2
	return (e.Attributes.Reaction.TotalValue + e.Attributes.Agility.TotalValue) / 2
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
	e.Attributes.Recalculate()
	// Formula: [Body / 2] + 1 (rounded up)
	return int(math.Ceil(float64(e.Attributes.Body.TotalValue)/2.0) + 1)
}

// Attributes

// GetInitiative (Reaction + Intuition) + 1D6 + Attribute/Initiative Dice bonus
func (e *GameEntity) GetInitative() int {
	e.Attributes.Recalculate()
	poolSize := 1
	_, _, results := RollDice(poolSize)
	// Formula: (Reaction + Intuition) + 1D6
	// TODO: Add appropriate attribute and Initiative
	return (e.Attributes.Reaction.TotalValue + e.Attributes.Intuition.TotalValue) + RollResultsTotal(results)
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

// GetLiftCarry calculates and returns the Lift Carry of the character.
// The baseline for lifting weight is 15 kilograms per point of Strength. Anything more than that requires a Strength + Body Test. Each hit increases the max weight lifted by 15 kilograms. Lifting weight above your head, as with a clean & jerk, is more difficult. The baseline for lifting weight above the head is 5 kilograms per point Strength. Each hit on the Lifting Test increases the maximum weight you can lift by 5 kilograms.
// Carrying weight is significantly different than lifting weight. Characters can carry Strength x 10 kilograms in gear without effort. Additional weight requires a Lifting Test. Each hit increases the maximum by 10 kilograms.
func (e *GameEntity) GetLiftCarry() float64 {
	e.Attributes.Recalculate()
	baseCarryWeight := 10
	// Lift Formula: STR * 15
	// Carry Formula: STR * 10
	return float64(e.Attributes.Strength.TotalValue+e.Attributes.Body.TotalValue) * float64(baseCarryWeight)
}

// GetCurrentCarryWeight calculates and returns the current carry weight of the character.
func (e *GameEntity) GetCurrentCarryWeight() float64 {
	e.Attributes.Recalculate()
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

// GetComposure calculates and returns the Composure of the character.
// Some situations are tough to deal with, even for hardened professionals like shadowrunners. When a character is faced with an emotionally overwhelming situation there are only two choices. Stay and fight or turn into a quivering lump of goo. To find out which one happens, make a Willpower + Charisma Test, with a threshold based on the severity of the situation.Take note that repeating similar situations over and again eventually eliminates the need to perform this test. Staring down a group of well-armed gangers will be scary at first, but after a character does it a few times the fear gives way to instinct.
// (WIL + CHA)
func (e *GameEntity) GetComposure() int {
	e.Attributes.Recalculate()
	// Formula: (WIL + CHA)
	return e.Attributes.Willpower.TotalValue + e.Attributes.Charisma.TotalValue
}

// GetJudgeIntentions calculates and returns the Judge Intentions of the character.
// Reading another person is also a matter of instinct. A character can use their instincts to guess at the intentions of another person or to gauge how much they can trust someone. Make an Opposed Intuition + Charisma Test against the target’s Willpower + Charisma. This is not an exact science. A successful test doesn’t mean the target will never betray you (intentions have been known to change), and deceptive characters can gain another’s confidence easily. This primarily serves as a benchmark or gut instinct about how much you can trust the person you are dealing with.
func (e *GameEntity) GetJudgeIntentions() int {
	e.Attributes.Recalculate()
	// Formula: (INT + CHA)
	return e.Attributes.Intuition.TotalValue + e.Attributes.Charisma.TotalValue
}

// GetMemory calculates and returns the Memory of the character.
// While there are numerous mnemonic devices, and even a few select pieces of bioware, designed for remembering information, memory is not a skill. If a character needs to recall information make a Logic + Willpower Test. Use the Knowledge Skill Table to determine the threshold. If a character actively tries to memorize information, make a Logic + Willpower Test at the time of memorization. Each hit adds a dice to the Recall Test later on.
// Glitches can have a devastating effect on memory. A glitch means the character misremembers some portion of the information, such as order of numbers in a passcode. A critical glitch means the character has completely fooled himself into believing and thus remembering something that never actually happened.
func (e *GameEntity) GetMemory() int {
	e.Attributes.Recalculate()
	// Formula: (LOG + WIL)
	return e.Attributes.Logic.TotalValue + e.Attributes.Willpower.TotalValue
}

// Judge Intentions (INT + CHA)
// Memory (LOG + WIL)

// GetPhysicalLimit calculates and returns the Physical Limit of the character.
func (e *GameEntity) GetPhysicalLimit() int {
	e.Attributes.Recalculate()
	strength := float64(e.Attributes.Strength.TotalValue)
	body := float64(e.Attributes.Body.TotalValue)
	reaction := float64(e.Attributes.Reaction.TotalValue)

	// Formula: [(Strength x 2) + Body + Reaction] / 3 (round up)
	limit := (strength*2 + body + reaction) / 3.0
	return int(math.Ceil(limit)) // Round up
}

// GetSocialLimit calculates and returns the Social Limit of the character.
func (e *GameEntity) GetSocialLimit() int {
	e.Attributes.Recalculate()
	charisma := float64(e.Attributes.Charisma.TotalValue)
	willpower := float64(e.Attributes.Willpower.TotalValue)
	essence := e.Attributes.Essence.TotalValue

	// Formula: [(Charisma x 2) + Willpower + Essence] / 3 (round up)
	limit := (charisma*2 + willpower + essence) / 3.0
	return int(math.Ceil(limit)) // Round up
}

// GetMentalLimit calculates and returns the Mental Limit of the character.
func (e *GameEntity) GetMentalLimit() int {
	e.Attributes.Recalculate()
	logic := float64(e.Attributes.Logic.TotalValue)
	intuition := float64(e.Attributes.Intuition.TotalValue)
	willpower := float64(e.Attributes.Willpower.TotalValue)

	// Formula: [(Logic x 2) + Intuition + Willpower] / 3 (round up)
	limit := (logic*2 + intuition + willpower) / 3.0
	return int(math.Ceil(limit)) // Round up
}

// GetEncumbrancePenalty calculates and returns the encumbrance penalty of the character.
func (e *GameEntity) GetEncumbrancePenalty() int {
	e.Attributes.Recalculate()
	currentWeight := e.GetCurrentCarryWeight()
	maxCarryWeight := float64(e.GetLiftCarry())
	excessWeight := currentWeight - maxCarryWeight

	if excessWeight <= 0 {
		return 0 // No penalty if within carrying capacity
	}

	// Formula: Excess weight / 15 kg (rounded up)
	// // -1 penalty for every 15 kg over capacity
	return int(math.Ceil(excessWeight / 15.0))
}

// GetAdjustedPhysicalLimit calculates and returns the adjusted physical limit of the character.
func (e *GameEntity) GetAdjustedPhysicalLimit() int {
	e.Attributes.Recalculate()
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

// Condition Monitor Boxes
// Physical 							[Body x 2] + 8												Add bonuses to Body before calculating; round up final results
// Stun								[Willpower x 2] + 8											Add bonuses to Willpower before calculating; round up final results
// Overflow							Body + Augmentation bonuses									-
