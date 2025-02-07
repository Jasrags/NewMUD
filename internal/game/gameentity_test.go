package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// var (
// 	testStreetSamurai = &Character{
// 		AccountID: "test_user",
// 		Role:      CharacterRolePlayer,
// 		GameEntity: GameEntity{
// 			Name:            "Street Samurai",
// 			Title:           "Street Samurai",
// 			ID:              "ID",
// 			Metatype:        "Ork",
// 			Age:             25,
// 			Sex:             "Male",
// 			Height:          180,
// 			Weight:          80,
// 			Ethnicity:       "White",
// 			StreetCred:      2,
// 			Notoriety:       2,
// 			PublicAwareness: 2,
// 			Karma:           2,
// 			TotalKarma:      5,
// 			Description:     "A street samurai character",
// 			Attributes: &Attributes{
// 				Body:      Attribute[int]{Base: 7},
// 				Agility:   Attribute[int]{Base: 6},
// 				Reaction:  Attribute[int]{Base: 5}, // (7)
// 				Strength:  Attribute[int]{Base: 5},
// 				Willpower: Attribute[int]{Base: 3},
// 				Logic:     Attribute[int]{Base: 2},
// 				Intuition: Attribute[int]{Base: 3},
// 				Charisma:  Attribute[int]{Base: 2},
// 				Essence:   Attribute[float64]{Base: 0.88},
// 				Magic:     Attribute[int]{Base: 0},
// 				Resonance: Attribute[int]{Base: 0},
// 			},
// 			PhysicalDamage: PhysicalDamage{
// 				Current:  0,
// 				Max:      10,
// 				Overflow: 0,
// 			},
// 			StunDamage: StunDamage{
// 				Current: 0,
// 				Max:     10,
// 			},
// 			Edge: Edge{
// 				Max:       5,
// 				Available: 5,
// 			},
// 			Equipment: map[string]*Item{},
// 		},
// 	}
// 	testCovertOpsSpecialist = &Character{
// 		AccountID: "test_user",
// 		Role:      CharacterRolePlayer,
// 		GameEntity: GameEntity{
// 			Name:            "Covert Ops Specialist",
// 			Title:           "Covert Ops Specialist",
// 			ID:              "ID",
// 			Metatype:        "Dwarf",
// 			Age:             25,
// 			Sex:             "Male",
// 			Height:          180,
// 			Weight:          80,
// 			Ethnicity:       "White",
// 			StreetCred:      0,
// 			Notoriety:       0,
// 			PublicAwareness: 0,
// 			Karma:           0,
// 			TotalKarma:      0,
// 			Description:     "A covert ops specialist character",
// 			Attributes: &Attributes{
// 				Body:      Attribute[int]{Base: 5},
// 				Agility:   Attribute[int]{Base: 6},
// 				Reaction:  Attribute[int]{Base: 4},
// 				Strength:  Attribute[int]{Base: 5},
// 				Willpower: Attribute[int]{Base: 4},
// 				Logic:     Attribute[int]{Base: 4},
// 				Intuition: Attribute[int]{Base: 5},
// 				Charisma:  Attribute[int]{Base: 4},
// 				Essence:   Attribute[float64]{Base: 5.6},
// 				Magic:     Attribute[int]{Base: 0},
// 				Resonance: Attribute[int]{Base: 0},
// 			},
// 			Equipment: map[string]*Item{},
// 		},
// 	}
// 	// OCCULT INVESTIGATOR
// 	// STREET SHAMAN
// 	// COMBAT MAGE
// 	// BRAWLING ADEPT
// 	// WEAPONS SPECIALIST
// 	// FACE
// 	// TANK
// 	// DECKER
// 	// TECHNOMANCER
// 	// GUNSLINGER ADEPT
// 	// DRONE RIGGER
// 	// SMUGGLER
// 	// SPRAWL GANGER
// 	// BOUNTY HUNTER
// )

func TestGetSocialLimit(t *testing.T) {
	tests := []struct {
		name string
		// attributes *Attributes
		charisma  Attribute[int]
		willpower Attribute[int]
		essence   Attribute[float64]
		expected  int
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			charisma:  Attribute[int]{Base: 3},
			willpower: Attribute[int]{Base: 2},
			essence:   Attribute[float64]{Base: 6.0},
			// },
			expected: 5,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			charisma:  Attribute[int]{Base: 2},
			willpower: Attribute[int]{Base: 3},
			essence:   Attribute[float64]{Base: 1.56},
			// },
			expected: 3,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			charisma:  Attribute[int]{Base: 7},
			willpower: Attribute[int]{Base: 4},
			essence:   Attribute[float64]{Base: 6},
			// },
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Charisma = tt.charisma
			entity.Willpower = tt.willpower
			entity.Essence = tt.essence
			// entity.Attributes = tt.attributes

			result := entity.GetSocialLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetPhysicalLimit(t *testing.T) {
	tests := []struct {
		name string
		// attributes *Attributes
		strength Attribute[int]
		body     Attribute[int]
		reaction Attribute[int]
		expected int
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			strength: Attribute[int]{Base: 2},
			body:     Attribute[int]{Base: 3},
			reaction: Attribute[int]{Base: 3},
			// },
			expected: 4,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			strength: Attribute[int]{Base: 7},
			body:     Attribute[int]{Base: 10},
			reaction: Attribute[int]{Base: 3},
			// },
			expected: 9,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			strength: Attribute[int]{Base: 2},
			body:     Attribute[int]{Base: 3},
			reaction: Attribute[int]{Base: 3},
			// },
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Strength = tt.strength
			entity.Body = tt.body
			entity.Reaction = tt.reaction
			// entity.Attributes = tt.attributes

			result := entity.GetPhysicalLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetMentalLimit(t *testing.T) {
	tests := []struct {
		name string
		// attributes *Attributes
		logic     Attribute[int]
		intuition Attribute[int]
		willpower Attribute[int]
		expected  int
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			logic:     Attribute[int]{Base: 5},
			intuition: Attribute[int]{Base: 6},
			willpower: Attribute[int]{Base: 5},
			// },
			expected: 7,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			logic:     Attribute[int]{Base: 2},
			intuition: Attribute[int]{Base: 3},
			willpower: Attribute[int]{Base: 3},
			// },
			expected: 4,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			logic:     Attribute[int]{Base: 4},
			intuition: Attribute[int]{Base: 4},
			willpower: Attribute[int]{Base: 4},
			// },
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Logic = tt.logic
			entity.Intuition = tt.intuition
			entity.Willpower = tt.willpower
			// entity.Attributes = tt.attributes

			result := entity.GetMentalLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetComposure(t *testing.T) {
	tests := []struct {
		name string
		// attributes *Attributes
		charisma  Attribute[int]
		willpower Attribute[int]
		expected  int
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			charisma:  Attribute[int]{Base: 3},
			willpower: Attribute[int]{Base: 2},
			// },
			expected: 5,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			charisma:  Attribute[int]{Base: 2},
			willpower: Attribute[int]{Base: 3},
			// },
			expected: 5,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			charisma:  Attribute[int]{Base: 7},
			willpower: Attribute[int]{Base: 4},
			// },
			expected: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Charisma = tt.charisma
			entity.Willpower = tt.willpower
			// entity.Attributes = tt.attributes

			result := entity.GetComposure()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetJudgeIntentions(t *testing.T) {
	tests := []struct {
		name       string
		attributes *Attributes
		intuition  Attribute[int]
		charisma   Attribute[int]
		expected   int
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			intuition: Attribute[int]{Base: 6},
			charisma:  Attribute[int]{Base: 3},
			// },
			expected: 9,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			intuition: Attribute[int]{Base: 3},
			charisma:  Attribute[int]{Base: 2},
			// },
			expected: 5,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			intuition: Attribute[int]{Base: 4},
			charisma:  Attribute[int]{Base: 7},
			// },
			expected: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Intuition = tt.intuition
			entity.Charisma = tt.charisma
			// entity.Attributes = tt.attributes

			result := entity.GetJudgeIntentions()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetMemory(t *testing.T) {
	tests := []struct {
		name string
		// attributes *Attributes
		logic     Attribute[int]
		willpower Attribute[int]
		expected  int
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			logic:     Attribute[int]{Base: 5},
			willpower: Attribute[int]{Base: 6},
			// },
			expected: 11,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			logic:     Attribute[int]{Base: 2},
			willpower: Attribute[int]{Base: 3},
			// },
			expected: 5,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			logic:     Attribute[int]{Base: 4},
			willpower: Attribute[int]{Base: 4},
			// },
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Logic = tt.logic
			entity.Willpower = tt.willpower
			// entity.Attributes = tt.attributes

			result := entity.GetMemory()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetLiftCarry(t *testing.T) {
	tests := []struct {
		name string
		// attributes *Attributes
		strength Attribute[int]
		body     Attribute[int]
		expected float64
	}{
		{
			name: "Elf Adept",
			// attributes: &Attributes{
			strength: Attribute[int]{Base: 2},
			body:     Attribute[int]{Base: 3},
			// },
			expected: 50.0,
		},
		{
			name: "Troll Tank",
			// attributes: &Attributes{
			strength: Attribute[int]{Base: 7},
			body:     Attribute[int]{Base: 10},
			// },
			expected: 170.0,
		},
		{
			name: "Elf Face",
			// attributes: &Attributes{
			strength: Attribute[int]{Base: 2},
			body:     Attribute[int]{Base: 3},
			// },
			expected: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Strength = tt.strength
			entity.Body = tt.body
			// entity.Attributes = tt.attributes

			result := entity.GetLiftCarry()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestUseEdge(t *testing.T) {
	tests := []struct {
		name         string
		edge         Edge
		expected     bool
		expectedEdge int
	}{
		{
			name:         "Use edge with available points",
			edge:         Edge{Max: 5, Available: 5},
			expectedEdge: 4,
			expected:     true,
		},
		{
			name:         "Use edge with no points",
			edge:         Edge{Max: 5, Available: 0},
			expectedEdge: 0,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Edge = tt.edge

			result := entity.UseEdge()
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedEdge, entity.Edge.Available)
		})
	}
}
func TestRegainEdge(t *testing.T) {
	tests := []struct {
		name         string
		edge         Edge
		expected     bool
		expectedMax  int
		expectedEdge int
	}{
		{
			name:         "Regain edge with available capacity",
			edge:         Edge{Max: 5, Available: 4},
			expectedEdge: 5,
			expected:     true,
		},
		{
			name:         "Regain edge at maximum capacity",
			edge:         Edge{Max: 5, Available: 5},
			expectedEdge: 5,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Edge = tt.edge

			result := entity.RegainEdge()
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedEdge, entity.Edge.Available)
		})
	}
}
func TestBurnEdge(t *testing.T) {
	tests := []struct {
		name         string
		edge         Edge
		expected     bool
		expectedMax  int
		expectedEdge int
	}{
		{
			name:         "Burn edge with available points",
			edge:         Edge{Max: 5, Available: 5},
			expectedMax:  4,
			expectedEdge: 4,
			expected:     true,
		},
		{
			name:         "Burn edge with no available points",
			edge:         Edge{Max: 5, Available: 0},
			expectedMax:  5,
			expectedEdge: 0,
			expected:     false,
		},
		{
			name:         "Burn edge with max already zero",
			edge:         Edge{Max: 0, Available: 0},
			expectedMax:  0,
			expectedEdge: 0,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Edge = tt.edge

			result := entity.BurnEdge()
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expectedMax, entity.Edge.Max)
			assert.Equal(t, tt.expectedEdge, entity.Edge.Available)
		})
	}
}
