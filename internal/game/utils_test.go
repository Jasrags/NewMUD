package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRollDice(t *testing.T) {
	tests := []struct {
		pool int
	}{
		{1},
		{2},
		{5},
		{10},
		{20},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("RollDice(%d)", test.pool), func(t *testing.T) {
			hits, glitches, results := RollDice(test.pool)
			assert.Len(t, results, test.pool)
			for _, result := range results {
				assert.GreaterOrEqual(t, result, 1)
				assert.LessOrEqual(t, result, 6)
			}
			assert.GreaterOrEqual(t, hits, 0)
			assert.GreaterOrEqual(t, glitches, 0)
		})
	}
}

func TestRollResultsTotal(t *testing.T) {
	tests := []struct {
		results []int
		total   int
	}{
		{[]int{1, 2, 3, 4, 5}, 15},
		{[]int{6, 6, 6, 6, 6}, 30},
		{[]int{1, 1, 1, 1, 1}, 5},
		{[]int{2, 4, 6, 8, 10}, 30},
		{[]int{}, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("RollResultsTotal(%v)", test.results), func(t *testing.T) {
			total := RollResultsTotal(test.results)
			assert.Equal(t, test.total, total)
		})
	}
}
func TestCheckGlitch(t *testing.T) {
	tests := []struct {
		pool           int
		hits           int
		glitches       int
		expectedGlitch bool
		expectedCrit   bool
	}{
		{10, 5, 0, false, false},
		{10, 0, 6, true, true},
		{10, 1, 6, true, false},
		{10, 0, 5, false, false},
		{10, 0, 0, false, false},
		{10, 3, 6, true, false},
		{10, 0, 7, true, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("CheckGlitch(pool=%d,hits=%d,glitches=%d)", test.pool, test.hits, test.glitches), func(t *testing.T) {
			glitch, criticalGlitch := CheckGlitch(test.pool, test.hits, test.glitches)
			assert.Equal(t, test.expectedGlitch, glitch)
			assert.Equal(t, test.expectedCrit, criticalGlitch)
		})
	}
}

func TestRollWithEdge(t *testing.T) {
	tests := []struct {
		pool int
	}{
		{1},
		{2},
		{5},
		{10},
		{20},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("RollWithEdge(%d)", test.pool), func(t *testing.T) {
			hits, glitches, results := RollWithEdge(test.pool)
			assert.GreaterOrEqual(t, hits, 0)
			assert.GreaterOrEqual(t, glitches, 0)
			assert.GreaterOrEqual(t, len(results), test.pool)
			for _, result := range results {
				assert.GreaterOrEqual(t, result, 1)
				assert.LessOrEqual(t, result, 6)
			}
		})
	}
}
func TestSingularize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"cats", "cat"},
		{"dogs", "dog"},
		{"buses", "buse"},
		{"", ""},
		{"a", "a"},
		{"s", "s"},
		{"class", "clas"},
		{"boss", "bos"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Singularize(%s)", test.input), func(t *testing.T) {
			result := Singularize(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestParseDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"n", "north"},
		{"north", "north"},
		{"s", "south"},
		{"south", "south"},
		{"e", "east"},
		{"east", "east"},
		{"w", "west"},
		{"west", "west"},
		{"u", "up"},
		{"up", "up"},
		{"d", "down"},
		{"down", "down"},
		{"invalid", ""},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ParseDirection(%s)", test.input), func(t *testing.T) {
			result := ParseDirection(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
func TestReverseDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"n", "south"},
		{"north", "south"},
		{"s", "north"},
		{"south", "north"},
		{"e", "west"},
		{"east", "west"},
		{"w", "east"},
		{"west", "east"},
		{"u", "down"},
		{"up", "down"},
		{"d", "up"},
		{"down", "up"},
		{"invalid", ""},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("ReverseDirection(%s)", test.input), func(t *testing.T) {
			result := ReverseDirection(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
func TestRollChance(t *testing.T) {
	tests := []struct {
		spawnChance int
	}{
		{0},
		{25},
		{50},
		{75},
		{100},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("RollChance(%d)", test.spawnChance), func(t *testing.T) {
			result := RollChance(test.spawnChance)
			assert.IsType(t, true, result)
		})
	}
}
func TestIsYAMLFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"config.yaml", true},
		{"config.yml", true},
		{"config.json", false},
		{"config.txt", false},
		{"CONFIG.YAML", true},
		{"CONFIG.YML", true},
		{"config.YaMl", true},
		{"", false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("IsYAMLFile(%s)", test.filename), func(t *testing.T) {
			result := IsYAMLFile(test.filename)
			assert.Equal(t, test.expected, result)
		})
	}
}
