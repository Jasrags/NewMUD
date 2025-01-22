package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSocialLimit(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Charisma:  Attribute[int]{Base: 3},
				Willpower: Attribute[int]{Base: 2},
				Essence:   Attribute[float64]{Base: 6.0},
			},
			expected: 5,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Charisma:  Attribute[int]{Base: 2},
				Willpower: Attribute[int]{Base: 3},
				Essence:   Attribute[float64]{Base: 1.56},
			},
			expected: 3,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Charisma:  Attribute[int]{Base: 7},
				Willpower: Attribute[int]{Base: 4},
				Essence:   Attribute[float64]{Base: 6},
			},
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetSocialLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetPhysicalLimit(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Strength: Attribute[int]{Base: 2},
				Body:     Attribute[int]{Base: 3},
				Reaction: Attribute[int]{Base: 3},
			},
			expected: 4,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Strength: Attribute[int]{Base: 7},
				Body:     Attribute[int]{Base: 10},
				Reaction: Attribute[int]{Base: 3},
			},
			expected: 9,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Strength: Attribute[int]{Base: 2},
				Body:     Attribute[int]{Base: 3},
				Reaction: Attribute[int]{Base: 3},
			},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetPhysicalLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetMentalLimit(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Logic:     Attribute[int]{Base: 5},
				Intuition: Attribute[int]{Base: 6},
				Willpower: Attribute[int]{Base: 5},
			},
			expected: 7,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Logic:     Attribute[int]{Base: 2},
				Intuition: Attribute[int]{Base: 3},
				Willpower: Attribute[int]{Base: 3},
			},
			expected: 4,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Logic:     Attribute[int]{Base: 4},
				Intuition: Attribute[int]{Base: 4},
				Willpower: Attribute[int]{Base: 4},
			},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetMentalLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetComposure(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Charisma:  Attribute[int]{Base: 3},
				Willpower: Attribute[int]{Base: 2},
			},
			expected: 5,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Charisma:  Attribute[int]{Base: 2},
				Willpower: Attribute[int]{Base: 3},
			},
			expected: 5,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Charisma:  Attribute[int]{Base: 7},
				Willpower: Attribute[int]{Base: 4},
			},
			expected: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetComposure()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetJudgeIntentions(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Intuition: Attribute[int]{Base: 6},
				Charisma:  Attribute[int]{Base: 3},
			},
			expected: 9,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Intuition: Attribute[int]{Base: 3},
				Charisma:  Attribute[int]{Base: 2},
			},
			expected: 5,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Intuition: Attribute[int]{Base: 4},
				Charisma:  Attribute[int]{Base: 7},
			},
			expected: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetJudgeIntentions()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetMemory(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Logic:     Attribute[int]{Base: 5},
				Willpower: Attribute[int]{Base: 6},
			},
			expected: 11,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Logic:     Attribute[int]{Base: 2},
				Willpower: Attribute[int]{Base: 3},
			},
			expected: 5,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Logic:     Attribute[int]{Base: 4},
				Willpower: Attribute[int]{Base: 4},
			},
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetMemory()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGetLiftCarry(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
		expected   int
	}{
		{
			name: "Elf Adept",
			attributes: Attributes{
				Strength: Attribute[int]{Base: 2},
				Body:     Attribute[int]{Base: 3},
			},
			expected: 50,
		},
		{
			name: "Troll Tank",
			attributes: Attributes{
				Strength: Attribute[int]{Base: 7},
				Body:     Attribute[int]{Base: 10},
			},
			expected: 170,
		},
		{
			name: "Elf Face",
			attributes: Attributes{
				Strength: Attribute[int]{Base: 2},
				Body:     Attribute[int]{Base: 3},
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := NewGameEntity()
			entity.Attributes = tt.attributes

			result := entity.GetLiftCarry()
			assert.Equal(t, tt.expected, result)
		})
	}
}
