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
func (e *GameEntity) GetInitative() int {
	poolSize := 1
	_, _, results := RollDice(poolSize)
	return (e.Attributes.Reaction.TotalValue + e.Attributes.Intuition.TotalValue) + RollResultsTotal(results)
}

// Initiative							(Reaction + Intuition) + 1D6			Add appropriate attribute and Initiative Dice bonu
// Astral Initiative					(Intuition x 2) + 2D6					—
// Matrix AR Initiative				(Reaction + Intuition) + 1D6				—
// Matrix VR Initiative (Hot Sim)		(Data Processing + Intuition) + 4D6		—
// Matrix VR Initiative (Cold Sim)		(Data Processing + Intuition) + 3D6		—

// Lifting and Carrying (STR + BOD)
func (e *GameEntity) GetLiftCarry() int {
	baseCarryWeight := 10
	return (e.Attributes.Strength.TotalValue + e.Attributes.Body.TotalValue) * baseCarryWeight
}

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

// TODO: Implement encumbered penatlies for combat
func (e *GameEntity) IsEncumbered() bool {
	return e.GetEncumbrancePenalty() > 0
}

// Composure (WIL + CHA)
// Judge Intentions (INT + CHA)
// Memory (LOG + WIL)

// Inherent Limits
func (e *GameEntity) GetMentalLimit() int {
	e.Attributes.Recalculate()
	logic := e.Attributes.Logic.TotalValue
	intuition := e.Attributes.Intuition.TotalValue
	willpower := e.Attributes.Willpower.TotalValue
	return (logic*2 + intuition + willpower) / 3
}

func (e *GameEntity) GetSocialLimit() int {
	e.Attributes.Recalculate()
	charisma := e.Attributes.Charisma.TotalValue
	willpower := e.Attributes.Willpower.TotalValue
	essence := e.Attributes.Essence.TotalValue
	return (charisma*2 + willpower + int(essence)) / 3
}

func (e *GameEntity) GetPhysicalLimit() int {
	e.Attributes.Recalculate()
	strength := e.Attributes.Strength.TotalValue
	body := e.Attributes.Body.TotalValue
	reaction := e.Attributes.Reaction.TotalValue
	return int(math.Ceil(float64(strength*2+body+reaction) / 3.0))
}

func (e *GameEntity) GetEncumbrancePenalty() int {
	e.Attributes.Recalculate()
	currentWeight := e.GetCurrentCarryWeight()
	maxCarryWeight := float64(e.GetLiftCarry()) // Carrying capacity from Strength
	excessWeight := currentWeight - maxCarryWeight

	if excessWeight <= 0 {
		return 0 // No penalty if within carrying capacity
	}

	return int(math.Ceil(excessWeight / 15.0)) // -1 penalty for every 15 kg over capacity
}

func (e *GameEntity) GetAdjustedPhysicalLimit() int {
	e.Attributes.Recalculate()
	basePhysicalLimit := e.GetPhysicalLimit()
	penalty := e.GetEncumbrancePenalty()

	// Physical Limit cannot go below 1
	adjustedLimit := basePhysicalLimit - penalty
	if adjustedLimit < 1 {
		adjustedLimit = 1
	}

	return adjustedLimit
}

// Condition Monitor Boxes
// Physical 							[Body x 2] + 8												Add bonuses to Body before calculating; round up final results
// Stun								[Willpower x 2] + 8											Add bonuses to Willpower before calculating; round up final results
// Overflow							Body + Augmentation bonuses									-
