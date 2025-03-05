package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddLabel(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		expected []string
	}{
		{
			name:     "Single label",
			labels:   []string{"TestLabel"},
			expected: []string{"testlabel"},
		},
		{
			name:     "Multiple labels",
			labels:   []string{"TestLabel", "AnotherLabel"},
			expected: []string{"testlabel", "anotherlabel"},
		},
		{
			name:     "Duplicate label",
			labels:   []string{"TestLabel", "TestLabel", "AnotherLabel"},
			expected: []string{"testlabel", "anotherlabel"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ged := NewGameEntityDynamic()
			for _, label := range tt.labels {
				ged.AddLabel(label)
			}
			for _, expectedLabel := range tt.expected {
				assert.Contains(t, ged.Labels, expectedLabel)
			}
			assert.Len(t, ged.Labels, len(tt.expected))
		})
	}
}
func TestHasLabels(t *testing.T) {
	tests := []struct {
		name     string
		labels   []string
		check    []string
		expected bool
	}{
		{
			name:     "Single label exists",
			labels:   []string{"TestLabel"},
			check:    []string{"TestLabel"},
			expected: true,
		},
		{
			name:     "Single label does not exist",
			labels:   []string{"TestLabel"},
			check:    []string{"NonExistentLabel"},
			expected: false,
		},
		{
			name:     "Multiple labels exist",
			labels:   []string{"TestLabel", "AnotherLabel"},
			check:    []string{"TestLabel", "AnotherLabel"},
			expected: true,
		},
		{
			name:     "One of multiple labels does not exist",
			labels:   []string{"TestLabel", "AnotherLabel"},
			check:    []string{"TestLabel", "NonExistentLabel"},
			expected: true,
		},
		{
			name:     "No labels exist",
			labels:   []string{"TestLabel", "AnotherLabel"},
			check:    []string{"NonExistentLabel1", "NonExistentLabel2"},
			expected: false,
		},
		{
			name:     "Empty check labels",
			labels:   []string{"TestLabel", "AnotherLabel"},
			check:    []string{},
			expected: false,
		},
		{
			name:     "Empty entity labels",
			labels:   []string{},
			check:    []string{"TestLabel"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ged := NewGameEntityDynamic()
			for _, label := range tt.labels {
				ged.AddLabel(label)
			}
			result := ged.HasLabels(tt.check...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestRemoveLabel(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		remove   string
		expected []string
	}{
		{
			name:     "Remove existing label",
			initial:  []string{"TestLabel"},
			remove:   "TestLabel",
			expected: []string{},
		},
		{
			name:     "Remove one of multiple labels",
			initial:  []string{"TestLabel", "AnotherLabel"},
			remove:   "TestLabel",
			expected: []string{"anotherlabel"},
		},
		{
			name:     "Remove non-existent label",
			initial:  []string{"TestLabel"},
			remove:   "NonExistentLabel",
			expected: []string{"testlabel"},
		},
		{
			name:     "Remove label with different case",
			initial:  []string{"TestLabel"},
			remove:   "testlabel",
			expected: []string{},
		},
		{
			name:     "Remove label from empty list",
			initial:  []string{},
			remove:   "TestLabel",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ged := NewGameEntityDynamic()
			for _, label := range tt.initial {
				ged.AddLabel(label)
			}
			ged.RemoveLabel(tt.remove)
			assert.Equal(t, tt.expected, ged.Labels)
		})
	}
}
func TestGetCharacterDisposition(t *testing.T) {
	tests := []struct {
		name           string
		dispositions   map[string]string
		characterID    string
		expectedResult string
	}{
		{
			name:           "Character disposition exists",
			dispositions:   map[string]string{"char1": "Friendly"},
			characterID:    "char1",
			expectedResult: "Friendly",
		},
		{
			name:           "Character disposition does not exist",
			dispositions:   map[string]string{"char1": "Friendly"},
			characterID:    "char2",
			expectedResult: DispositionNeutral,
		},
		{
			name:           "Empty dispositions",
			dispositions:   map[string]string{},
			characterID:    "char1",
			expectedResult: DispositionNeutral,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ged := NewGameEntityDynamic()
			ged.CharacterDispositions = tt.dispositions
			result := ged.GetCharacterDisposition(tt.characterID)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
