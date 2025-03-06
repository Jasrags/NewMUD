package game

// Street Samurai
// Covert Ops Specialist
// OCCULT INVESTIGATOR
// STREET SHAMAN
// COMBAT MAGE
// BRAWLING ADEPT
// WEAPONS SPECIALIST
// FACE
// TANK
// DECKER
// TECHNOMANCER
// GUNSLINGER ADEPT
// DRONE RIGGER
// SMUGGLER
// SPRAWL GANGER
// BOUNTY HUNTER

// func TestGetSocialLimit(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		charisma  Attribute[int]
// 		willpower Attribute[int]
// 		essence   Attribute[float64]
// 		expected  int
// 	}{
// 		{
// 			name:      "Elf Adept",
// 			charisma:  Attribute[int]{Base: 3},
// 			willpower: Attribute[int]{Base: 2},
// 			essence:   Attribute[float64]{Base: 6.0},
// 			expected:  5,
// 		},
// 		{
// 			name:      "Troll Tank",
// 			charisma:  Attribute[int]{Base: 2},
// 			willpower: Attribute[int]{Base: 3},
// 			essence:   Attribute[float64]{Base: 1.56},
// 			expected:  3,
// 		},
// 		{
// 			name:      "Elf Face",
// 			charisma:  Attribute[int]{Base: 7},
// 			willpower: Attribute[int]{Base: 4},
// 			essence:   Attribute[float64]{Base: 6},
// 			expected:  8,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Charisma = tt.charisma
// 			entity.Willpower = tt.willpower
// 			entity.Essence = tt.essence

// 			result := entity.GetSocialLimit()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetPhysicalLimit(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		strength Attribute[int]
// 		body     Attribute[int]
// 		reaction Attribute[int]
// 		expected int
// 	}{
// 		{
// 			name:     "Elf Adept",
// 			strength: Attribute[int]{Base: 2},
// 			body:     Attribute[int]{Base: 3},
// 			reaction: Attribute[int]{Base: 3},
// 			expected: 4,
// 		},
// 		{
// 			name:     "Troll Tank",
// 			strength: Attribute[int]{Base: 7},
// 			body:     Attribute[int]{Base: 10},
// 			reaction: Attribute[int]{Base: 3},
// 			expected: 9,
// 		},
// 		{
// 			name:     "Elf Face",
// 			strength: Attribute[int]{Base: 2},
// 			body:     Attribute[int]{Base: 3},
// 			reaction: Attribute[int]{Base: 3},
// 			expected: 4,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Strength = tt.strength
// 			entity.Body = tt.body
// 			entity.Reaction = tt.reaction

// 			result := entity.GetPhysicalLimit()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetMentalLimit(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		logic     Attribute[int]
// 		intuition Attribute[int]
// 		willpower Attribute[int]
// 		expected  int
// 	}{
// 		{
// 			name:      "Elf Adept",
// 			logic:     Attribute[int]{Base: 5},
// 			intuition: Attribute[int]{Base: 6},
// 			willpower: Attribute[int]{Base: 5},
// 			expected:  7,
// 		},
// 		{
// 			name:      "Troll Tank",
// 			logic:     Attribute[int]{Base: 2},
// 			intuition: Attribute[int]{Base: 3},
// 			willpower: Attribute[int]{Base: 3},
// 			expected:  4,
// 		},
// 		{
// 			name:      "Elf Face",
// 			logic:     Attribute[int]{Base: 4},
// 			intuition: Attribute[int]{Base: 4},
// 			willpower: Attribute[int]{Base: 4},
// 			expected:  6,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Logic = tt.logic
// 			entity.Intuition = tt.intuition
// 			entity.Willpower = tt.willpower

// 			result := entity.GetMentalLimit()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetComposure(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		charisma  Attribute[int]
// 		willpower Attribute[int]
// 		expected  int
// 	}{
// 		{
// 			name:      "Elf Adept",
// 			charisma:  Attribute[int]{Base: 3},
// 			willpower: Attribute[int]{Base: 2},
// 			expected:  5,
// 		},
// 		{
// 			name:      "Troll Tank",
// 			charisma:  Attribute[int]{Base: 2},
// 			willpower: Attribute[int]{Base: 3},
// 			expected:  5,
// 		},
// 		{
// 			name:      "Elf Face",
// 			charisma:  Attribute[int]{Base: 7},
// 			willpower: Attribute[int]{Base: 4},
// 			expected:  11,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Charisma = tt.charisma
// 			entity.Willpower = tt.willpower

// 			result := entity.GetComposure()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetJudgeIntentions(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		intuition Attribute[int]
// 		charisma  Attribute[int]
// 		expected  int
// 	}{
// 		{
// 			name:      "Elf Adept",
// 			intuition: Attribute[int]{Base: 6},
// 			charisma:  Attribute[int]{Base: 3},
// 			expected:  9,
// 		},
// 		{
// 			name:      "Troll Tank",
// 			intuition: Attribute[int]{Base: 3},
// 			charisma:  Attribute[int]{Base: 2},
// 			expected:  5,
// 		},
// 		{
// 			name:      "Elf Face",
// 			intuition: Attribute[int]{Base: 4},
// 			charisma:  Attribute[int]{Base: 7},
// 			expected:  11,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Intuition = tt.intuition
// 			entity.Charisma = tt.charisma

// 			result := entity.GetJudgeIntentions()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetMemory(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		logic     Attribute[int]
// 		willpower Attribute[int]
// 		expected  int
// 	}{
// 		{
// 			name:      "Elf Adept",
// 			logic:     Attribute[int]{Base: 5},
// 			willpower: Attribute[int]{Base: 6},
// 			expected:  11,
// 		},
// 		{
// 			name:      "Troll Tank",
// 			logic:     Attribute[int]{Base: 2},
// 			willpower: Attribute[int]{Base: 3},
// 			expected:  5,
// 		},
// 		{
// 			name:      "Elf Face",
// 			logic:     Attribute[int]{Base: 4},
// 			willpower: Attribute[int]{Base: 4},
// 			expected:  8,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Logic = tt.logic
// 			entity.Willpower = tt.willpower

// 			result := entity.GetMemory()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetLiftCarry(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		strength Attribute[int]
// 		body     Attribute[int]
// 		expected float64
// 	}{
// 		{
// 			name:     "Elf Adept",
// 			strength: Attribute[int]{Base: 2},
// 			body:     Attribute[int]{Base: 3},
// 			expected: 50.0,
// 		},
// 		{
// 			name:     "Troll Tank",
// 			strength: Attribute[int]{Base: 7},
// 			body:     Attribute[int]{Base: 10},
// 			expected: 170.0,
// 		},
// 		{
// 			name:     "Elf Face",
// 			strength: Attribute[int]{Base: 2},
// 			body:     Attribute[int]{Base: 3},
// 			expected: 50.0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Strength = tt.strength
// 			entity.Body = tt.body

//				result := entity.GetLiftCarry()
//				assert.Equal(t, tt.expected, result)
//			})
//		}
//	}
// func TestUseEdge(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		edge         Edge
// 		expected     bool
// 		expectedEdge int
// 	}{
// 		{
// 			name:         "Use edge with available points",
// 			edge:         Edge{Max: 5, Available: 5},
// 			expectedEdge: 4,
// 			expected:     true,
// 		},
// 		{
// 			name:         "Use edge with no points",
// 			edge:         Edge{Max: 5, Available: 0},
// 			expectedEdge: 0,
// 			expected:     false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Edge = tt.edge

//				result := entity.UseEdge()
//				assert.Equal(t, tt.expected, result)
//				assert.Equal(t, tt.expectedEdge, entity.Edge.Available)
//			})
//		}
//	}
// func TestRegainEdge(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		edge         Edge
// 		expected     bool
// 		expectedMax  int
// 		expectedEdge int
// 	}{
// 		{
// 			name:         "Regain edge with available capacity",
// 			edge:         Edge{Max: 5, Available: 4},
// 			expectedEdge: 5,
// 			expected:     true,
// 		},
// 		{
// 			name:         "Regain edge at maximum capacity",
// 			edge:         Edge{Max: 5, Available: 5},
// 			expectedEdge: 5,
// 			expected:     false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Edge = tt.edge

// 			result := entity.RegainEdge()
// 			assert.Equal(t, tt.expected, result)
// 			assert.Equal(t, tt.expectedEdge, entity.Edge.Available)
// 		})
// 	}
// }
// func TestBurnEdge(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		edge         Edge
// 		expected     bool
// 		expectedMax  int
// 		expectedEdge int
// 	}{
// 		{
// 			name:         "Burn edge with available points",
// 			edge:         Edge{Max: 5, Available: 5},
// 			expectedMax:  4,
// 			expectedEdge: 4,
// 			expected:     true,
// 		},
// 		{
// 			name:         "Burn edge with no available points",
// 			edge:         Edge{Max: 5, Available: 0},
// 			expectedMax:  5,
// 			expectedEdge: 0,
// 			expected:     false,
// 		},
// 		{
// 			name:         "Burn edge with max already zero",
// 			edge:         Edge{Max: 0, Available: 0},
// 			expectedMax:  0,
// 			expectedEdge: 0,
// 			expected:     false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Edge = tt.edge

// 			result := entity.BurnEdge()
// 			assert.Equal(t, tt.expected, result)
// 			assert.Equal(t, tt.expectedMax, entity.Edge.Max)
// 			assert.Equal(t, tt.expectedEdge, entity.Edge.Available)
// 		})
// 	}
// }

// func TestGetPhysicalConditionMax(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		body     Attribute[int]
// 		expected int
// 	}{
// 		{
// 			name:     "Body 3",
// 			body:     Attribute[int]{Base: 3},
// 			expected: 10,
// 		},
// 		{
// 			name:     "Body 4",
// 			body:     Attribute[int]{Base: 4},
// 			expected: 10,
// 		},
// 		{
// 			name:     "Body 5",
// 			body:     Attribute[int]{Base: 5},
// 			expected: 11,
// 		},
// 		{
// 			name:     "Body 7",
// 			body:     Attribute[int]{Base: 7},
// 			expected: 12,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Body = tt.body

// 			result := entity.GetPhysicalConditionMax()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetStunConditionMax(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		willpower Attribute[int]
// 		expected  int
// 	}{
// 		{
// 			name:      "Willpower 3",
// 			willpower: Attribute[int]{Base: 3},
// 			expected:  10,
// 		},
// 		{
// 			name:      "Willpower 4",
// 			willpower: Attribute[int]{Base: 4},
// 			expected:  10,
// 		},
// 		{
// 			name:      "Willpower 5",
// 			willpower: Attribute[int]{Base: 5},
// 			expected:  11,
// 		},
// 		{
// 			name:      "Willpower 7",
// 			willpower: Attribute[int]{Base: 7},
// 			expected:  12,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Willpower = tt.willpower

// 			result := entity.GetStunConditionMax()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetOverflowConditionMax(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		body     Attribute[int]
// 		expected int
// 	}{
// 		{
// 			name:     "Body 3",
// 			body:     Attribute[int]{Base: 3},
// 			expected: 3,
// 		},
// 		{
// 			name:     "Body 4",
// 			body:     Attribute[int]{Base: 4},
// 			expected: 4,
// 		},
// 		{
// 			name:     "Body 5",
// 			body:     Attribute[int]{Base: 5},
// 			expected: 5,
// 		},
// 		{
// 			name:     "Body 7",
// 			body:     Attribute[int]{Base: 7},
// 			expected: 7,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Body = tt.body

// 			result := entity.GetOverflowConditionMax()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
// func TestGetInitative(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		reaction  Attribute[int]
// 		intuition Attribute[int]
// 		expected  int
// 	}{
// 		{
// 			name:      "Elf Adept",
// 			reaction:  Attribute[int]{Base: 3},
// 			intuition: Attribute[int]{Base: 6},
// 			expected:  9,
// 		},
// 		{
// 			name:      "Troll Tank",
// 			reaction:  Attribute[int]{Base: 3},
// 			intuition: Attribute[int]{Base: 3},
// 			expected:  6,
// 		},
// 		{
// 			name:      "Elf Face",
// 			reaction:  Attribute[int]{Base: 4},
// 			intuition: Attribute[int]{Base: 4},
// 			expected:  8,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			entity := NewGameEntity()
// 			entity.Reaction = tt.reaction
// 			entity.Intuition = tt.intuition

// 			result := entity.GetInitative()
// 			assert.Equal(t, tt.expected, result)
// 		})
// 	}
// }
