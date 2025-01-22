package main

import (
	"math"
	"sync"

	"github.com/google/uuid"
	ee "github.com/vansante/go-event-emitter"
)

type GameEntity struct {
	sync.RWMutex `yaml:"-"`
	Listeners    []ee.Listener `yaml:"-"`

	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Attributes  Attributes       `yaml:"attributes"`
	Room        *Room            `yaml:"-"`
	RoomID      string           `yaml:"room_id"`
	Area        *Area            `yaml:"-"`
	AreaID      string           `yaml:"area_id"`
	Inventory   Inventory        `yaml:"inventory"`
	Equipment   map[string]*Item `yaml:"equipment"`
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

// Attributes

// GetInitiative (Reaction + Intuition) + 1D6 + Attribute/Initiative Dice bonus
func (e *GameEntity) GetInitative() int {
	e.Attributes.Recalculate()
	poolSize := 1
	_, _, results := RollDice(poolSize)
	return (e.Attributes.Reaction.TotalValue + e.Attributes.Intuition.TotalValue) + RollResultsTotal(results)
}

// Astral Initiative					(Intuition x 2) + 2D6					—
// Matrix AR Initiative				(Reaction + Intuition) + 1D6				—
// Matrix VR Initiative (Hot Sim)		(Data Processing + Intuition) + 4D6		—
// Matrix VR Initiative (Cold Sim)		(Data Processing + Intuition) + 3D6		—

// GetLiftCarry calculates and returns the Lift Carry of the character.
func (e *GameEntity) GetLiftCarry() int {
	e.Attributes.Recalculate()
	baseCarryWeight := 10
	// Formula: (STR + BOD) * 10
	return (e.Attributes.Strength.TotalValue + e.Attributes.Body.TotalValue) * baseCarryWeight
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
func (e *GameEntity) IsEncumbered() bool {
	return e.GetEncumbrancePenalty() > 0
}

// GetComposure calculates and returns the Composure of the character.
// (WIL + CHA)
func (e *GameEntity) GetComposure() int {
	e.Attributes.Recalculate()
	// Formula: (WIL + CHA)
	return e.Attributes.Willpower.TotalValue + e.Attributes.Charisma.TotalValue
}

// GetJudgeIntentions calculates and returns the Judge Intentions of the character.
func (e *GameEntity) GetJudgeIntentions() int {
	e.Attributes.Recalculate()
	// Formula: (INT + CHA)
	return e.Attributes.Intuition.TotalValue + e.Attributes.Charisma.TotalValue
}

// GetMemory calculates and returns the Memory of the character.
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

	// Formula: (Strength * 2 + Body + Reaction) / 3
	limit := (strength*2 + body + reaction) / 3.0
	return int(math.Ceil(limit)) // Round up
}

// GetSocialLimit calculates and returns the Social Limit of the character.
func (e *GameEntity) GetSocialLimit() int {
	e.Attributes.Recalculate()
	charisma := float64(e.Attributes.Charisma.TotalValue)
	willpower := float64(e.Attributes.Willpower.TotalValue)
	essence := e.Attributes.Essence.TotalValue // Already a float64

	// Formula: (Charisma * 2 + Willpower + Essence) / 3
	limit := (charisma*2 + willpower + essence) / 3.0
	return int(math.Ceil(limit)) // Round up
}

// GetMentalLimit calculates and returns the Mental Limit of the character.
func (e *GameEntity) GetMentalLimit() int {
	e.Attributes.Recalculate()
	logic := float64(e.Attributes.Logic.TotalValue)
	intuition := float64(e.Attributes.Intuition.TotalValue)
	willpower := float64(e.Attributes.Willpower.TotalValue)

	// Formula: (Logic * 2 + Intuition + Willpower) / 3
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
