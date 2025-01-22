package main

// import (
// 	"log/slog"
// 	"math"
// 	"sync"

// 	"github.com/google/uuid"
// 	ee "github.com/vansante/go-event-emitter"
// )

// type GameEntity struct {
// 	sync.RWMutex `yaml:"-"`
// 	Listeners    []ee.Listener `yaml:"-"`

// 	ID              string           `yaml:"id"`
// 	Name            string           `yaml:"name"`
// 	Metatype        string           `yaml:"metatype"`
// 	Age             int              `yaml:"age"`
// 	Sex             string           `yaml:"sex"`
// 	Height          int              `yaml:"height"`
// 	Weight          int              `yaml:"weight"`
// 	Ethnicity       string           `yaml:"ethnicity"`
// 	StreetCred      int              `yaml:"street_cred"`
// 	Notoriety       int              `yaml:"notoriety"`
// 	PublicAwareness int              `yaml:"public_awareness"`
// 	Karma           int              `yaml:"karma"`
// 	TotalKarma      int              `yaml:"total_karma"`
// 	Description     string           `yaml:"description"`
// 	Attributes      Attributes       `yaml:"attributes"`
// 	PhysicalDamage  PhysicalDamage   `yaml:"physical_damage"`
// 	StunDamage      StunDamage       `yaml:"stun_damage"`
// 	Edge            Edge             `yaml:"edge"`
// 	Room            *Room            `yaml:"-"`
// 	RoomID          string           `yaml:"room_id"`
// 	Area            *Area            `yaml:"-"`
// 	AreaID          string           `yaml:"area_id"`
// 	Inventory       Inventory        `yaml:"inventory"`
// 	Equipment       map[string]*Item `yaml:"equipment"`
// }

// func NewGameEntity() GameEntity {
// 	return GameEntity{
// 		ID:        uuid.New().String(),
// 		Equipment: make(map[string]*Item),
// 		Listeners: make([]ee.Listener, 0),
// 	}
// }

// func (e *GameEntity) GetName() string {
// 	return e.Name
// }

// func (e *GameEntity) GetID() string {
// 	return e.ID
// }

// func (e *GameEntity) SetRoom(room *Room) {
// 	e.Room = room
// 	e.RoomID = room.ReferenceID
// }

// // UseEdge - Decreases the available Edge by 1.
// func (e *GameEntity) UseEdge() bool {
// 	if e.Edge.Available <= 0 || e.Edge.Max <= 0 {
// 		slog.Warn("Cannot use Edge. No available edge to use.",
// 			slog.String("character_id", e.ID),
// 			slog.Int("max_edge_points", e.Edge.Max),
// 			slog.Int("remaining_edge_points", e.Edge.Available))
// 		return false
// 	}

// 	e.Edge.Available -= 1
// 	slog.Debug("Edge point used",
// 		slog.String("character_id", e.ID),
// 		slog.Int("max_edge_points", e.Edge.Max),
// 		slog.Int("remaining_edge_points", e.Edge.Available))
// 	return true
// }

// // Burn Edge - Decreases the maximum Edge by 1 and the available Edge by 1.
// func (e *GameEntity) BurnEdge() bool {
// 	if e.Edge.Max <= 0 || e.Edge.Available <= 0 {
// 		slog.Warn("Cannot burn Edge. Edge attribute is already zero.",
// 			slog.String("character_id", e.ID),
// 			slog.Int("max_edge_points", e.Edge.Max),
// 			slog.Int("remaining_edge_points", e.Edge.Available))
// 		return false
// 	}

// 	e.Edge.Max -= 1
// 	e.Edge.Available -= 1

// 	slog.Debug("Edge point burned",
// 		slog.String("character_id", e.ID),
// 		slog.Int("max_edge_points", e.Edge.Max),
// 		slog.Int("remaining_edge_points", e.Edge.Available))
// 	return true
// }

// // Regain Edge - Increases the available Edge by 1.
// func (e *GameEntity) RegainEdge() bool {
// 	if e.Edge.Available >= e.Edge.Max {
// 		slog.Warn("Cannot regain Edge. Edge attribute is already at max.",
// 			slog.String("character_id", e.ID),
// 			slog.Int("max_edge_points", e.Edge.Max),
// 			slog.Int("remaining_edge_points", e.Edge.Available))
// 		return false
// 	}

// 	e.Edge.Available += 1
// 	slog.Debug("Edge point regained",
// 		slog.String("character_id", e.ID),
// 		slog.Int("max_edge_points", e.Edge.Max),
// 		slog.Int("remaining_edge_points", e.Edge.Available))

// 	return true
// }

// // Damage
// // Physical Damage Level		(Body + Armor) / 2

// // The Physical Condition Monitor has boxes equal to half the character’s ((BOD/2) + 8) Body (rounded up) + 8;
// // the Stun Condition Monitor has boxes equaling half the character’s ((WILL/2) + 8)Willpower (rounded up) + 8.
// // When a row of the Condition Monitor is filled up, the player character takes a –1 penalty to all subsequent tests. This penalty stacks for each row of the Condition Monitor that is filled in.

// // GetPhysicalConditionMax calculates and returns the Physical Condition Max of the character.
// func (e *GameEntity) GetPhysicalConditionMax() int {
// 	e.Attributes.Recalculate()
// 	// Formula: [Body / 2] + 8 (rounded up)
// 	// TODO: Add bonuses to Body before calculating; round up final results
// 	return int(math.Ceil(float64(e.Attributes.Body.TotalValue)/2.0) + 8)
// }

// // GetStunConditionMax calculates and returns the Stun Condition Max of the character.
// func (e *GameEntity) GetStunConditionMax() int {
// 	e.Attributes.Recalculate()
// 	// Formula: [Willpower / 2] + 8 (rounded up)
// 	// TODO: Add bonuses to Willpower before calculating; round up final results
// 	return int(math.Ceil(float64(e.Attributes.Willpower.TotalValue)/2.0) + 8)
// }

// // GetOverflowConditionMax calculates and returns the Overflow Condition Max of the character.
// func (e *GameEntity) GetOverflowConditionMax() int {
// 	e.Attributes.Recalculate()
// 	// Formula: Body + Augmentation bonuses (rounded up)
// 	// TODO: Implement Augmentation bonuses
// 	return e.Attributes.Body.TotalValue
// }

// type PhysicalDamage struct {
// 	Current  int // Current damage boxes filled
// 	Max      int // Max damage boxes in the track
// 	Overflow int // Overflow boxes beyond the Physical track
// }

// type StunDamage struct {
// 	Current int // Current damage boxes filled
// 	Max     int // Max damage boxes in the track
// }

// type Edge struct {
// 	Max       int `yaml:"max"`
// 	Available int `yaml:"available"`
// }

// // TODO: implement movement
// func (e *GameEntity) GetMovement() int {
// 	e.Attributes.Recalculate()
// 	// Formula: (Reaction + Agility) / 2
// 	return (e.Attributes.Reaction.TotalValue + e.Attributes.Agility.TotalValue) / 2
// }

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

// /*
// Injuries cause pain, bleeding, and other distractions that interfere with doing all sorts of actions. Wound modifiers are accumulated with every third box of damage and are cumulative between damage tracks and with other negative modifiers such as spells or adverse conditions.

// Wound modifiers are applied to all tests not about reducing the number of boxes you’re about to take on your Condition Monitor (such as damage resistance, resisting direct combat spells, toxin resistance, and so on). The Wound Modifier penalty is also applied to the character’s Initiative attribute and therefore their Initiative Score during combat.
// */
// func (e *GameEntity) GetWoundModifiers() int {
// 	e.Attributes.Recalculate()
// 	// Formula: [Body / 2] + 1 (rounded up)
// 	return int(math.Ceil(float64(e.Attributes.Body.TotalValue)/2.0) + 1)
// }

// // Attributes

// // GetInitiative (Reaction + Intuition) + 1D6 + Attribute/Initiative Dice bonus
// func (e *GameEntity) GetInitative() int {
// 	e.Attributes.Recalculate()
// 	poolSize := 1
// 	_, _, results := RollDice(poolSize)
// 	// Formula: (Reaction + Intuition) + 1D6
// 	// TODO: Add appropriate attribute and Initiative
// 	return (e.Attributes.Reaction.TotalValue + e.Attributes.Intuition.TotalValue) + RollResultsTotal(results)
// }

// func (e *GameEntity) GetAdjustedInitiative() int {
// 	baseInitiative := e.GetInitative()
// 	return baseInitiative + e.GetWoundModifiers()
// }

// func (e *GameEntity) ApplyWoundModifiers(baseDicePool int) int {
// 	return baseDicePool + e.GetWoundModifiers()
// }

// // Astral Initiative					(Intuition x 2) + 2D6					—
// // Matrix AR Initiative				(Reaction + Intuition) + 1D6				—
// // Matrix VR Initiative (Hot Sim)		(Data Processing + Intuition) + 4D6		—
// // Matrix VR Initiative (Cold Sim)		(Data Processing + Intuition) + 3D6		—

// // GetLiftCarry calculates and returns the Lift Carry of the character.
// func (e *GameEntity) GetLiftCarry() float64 {
// 	e.Attributes.Recalculate()
// 	baseCarryWeight := 10
// 	// Formula: (STR + BOD) * 10
// 	return float64(e.Attributes.Strength.TotalValue+e.Attributes.Body.TotalValue) * float64(baseCarryWeight)
// }

// // GetCurrentCarryWeight calculates and returns the current carry weight of the character.
// func (e *GameEntity) GetCurrentCarryWeight() float64 {
// 	e.Attributes.Recalculate()
// 	totalWeight := 0.0

// 	for _, item := range e.Inventory.Items {
// 		bp := EntityMgr.GetItemBlueprintByInstance(item)
// 		if bp != nil {
// 			totalWeight += bp.Weight
// 		}
// 	}

// 	return totalWeight
// }

// // GetEncumberancePenalty calculates and returns if the character is encumbered.
// // TODO: Implement encumbered penatlies for combat
// func (e *GameEntity) IsEncumbered() bool {
// 	return e.GetEncumbrancePenalty() > 0
// }

// // GetComposure calculates and returns the Composure of the character.
// // (WIL + CHA)
// func (e *GameEntity) GetComposure() int {
// 	e.Attributes.Recalculate()
// 	// Formula: (WIL + CHA)
// 	return e.Attributes.Willpower.TotalValue + e.Attributes.Charisma.TotalValue
// }

// // GetJudgeIntentions calculates and returns the Judge Intentions of the character.
// func (e *GameEntity) GetJudgeIntentions() int {
// 	e.Attributes.Recalculate()
// 	// Formula: (INT + CHA)
// 	return e.Attributes.Intuition.TotalValue + e.Attributes.Charisma.TotalValue
// }

// // GetMemory calculates and returns the Memory of the character.
// func (e *GameEntity) GetMemory() int {
// 	e.Attributes.Recalculate()
// 	// Formula: (LOG + WIL)
// 	return e.Attributes.Logic.TotalValue + e.Attributes.Willpower.TotalValue
// }

// // Judge Intentions (INT + CHA)
// // Memory (LOG + WIL)

// // GetPhysicalLimit calculates and returns the Physical Limit of the character.
// func (e *GameEntity) GetPhysicalLimit() int {
// 	e.Attributes.Recalculate()
// 	strength := float64(e.Attributes.Strength.TotalValue)
// 	body := float64(e.Attributes.Body.TotalValue)
// 	reaction := float64(e.Attributes.Reaction.TotalValue)

// 	// Formula: [(Strength x 2) + Body + Reaction] / 3 (round up)
// 	limit := (strength*2 + body + reaction) / 3.0
// 	return int(math.Ceil(limit)) // Round up
// }

// // GetSocialLimit calculates and returns the Social Limit of the character.
// func (e *GameEntity) GetSocialLimit() int {
// 	e.Attributes.Recalculate()
// 	charisma := float64(e.Attributes.Charisma.TotalValue)
// 	willpower := float64(e.Attributes.Willpower.TotalValue)
// 	essence := e.Attributes.Essence.TotalValue

// 	// Formula: [(Charisma x 2) + Willpower + Essence] / 3 (round up)
// 	limit := (charisma*2 + willpower + essence) / 3.0
// 	return int(math.Ceil(limit)) // Round up
// }

// // GetMentalLimit calculates and returns the Mental Limit of the character.
// func (e *GameEntity) GetMentalLimit() int {
// 	e.Attributes.Recalculate()
// 	logic := float64(e.Attributes.Logic.TotalValue)
// 	intuition := float64(e.Attributes.Intuition.TotalValue)
// 	willpower := float64(e.Attributes.Willpower.TotalValue)

// 	// Formula: [(Logic x 2) + Intuition + Willpower] / 3 (round up)
// 	limit := (logic*2 + intuition + willpower) / 3.0
// 	return int(math.Ceil(limit)) // Round up
// }

// // GetEncumbrancePenalty calculates and returns the encumbrance penalty of the character.
// func (e *GameEntity) GetEncumbrancePenalty() int {
// 	e.Attributes.Recalculate()
// 	currentWeight := e.GetCurrentCarryWeight()
// 	maxCarryWeight := float64(e.GetLiftCarry())
// 	excessWeight := currentWeight - maxCarryWeight

// 	if excessWeight <= 0 {
// 		return 0 // No penalty if within carrying capacity
// 	}

// 	// Formula: Excess weight / 15 kg (rounded up)
// 	// // -1 penalty for every 15 kg over capacity
// 	return int(math.Ceil(excessWeight / 15.0))
// }

// // GetAdjustedPhysicalLimit calculates and returns the adjusted physical limit of the character.
// func (e *GameEntity) GetAdjustedPhysicalLimit() int {
// 	e.Attributes.Recalculate()
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

// // Condition Monitor Boxes
// // Physical 							[Body x 2] + 8												Add bonuses to Body before calculating; round up final results
// // Stun								[Willpower x 2] + 8											Add bonuses to Willpower before calculating; round up final results
// // Overflow							Body + Augmentation bonuses									-
