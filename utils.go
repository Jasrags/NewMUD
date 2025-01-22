package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v3"
)

func Singularize(word string) string {
	if strings.HasSuffix(word, "s") && len(word) > 1 {
		return word[:len(word)-1] // Remove trailing 's'
	}
	return word
}

// WrapText splits text into lines of the specified width without breaking words.
func WrapText(text string, width int) string {
	var result strings.Builder
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > width {
			// Append the current line to the result and reset the line
			if line != "" {
				result.WriteString(line + "\n")
			}
			line = word
		} else {
			// Append the word to the current line
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	// Append the last line, if any, without an extra newline
	if line != "" {
		result.WriteString(line)
	}

	return result.String()
}

func LoadYAML(filePath string, out interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(out)
}

func SaveYAML(filePath string, in interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(in)
}

func LoadJSON(filePath string, out interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("failed reading file",
			slog.String("file", filePath),
			slog.Any("error", err))

		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(out)
}

func SaveJSON(filePath string, in interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed reading file",
			slog.String("file", filePath),
			slog.Any("error", err))

		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(in)
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func ParseDirection(dir string) string {
	switch dir {
	case "n", "north":
		dir = "north"
	case "s", "south":
		dir = "south"
	case "e", "east":
		dir = "east"
	case "w", "west":
		dir = "west"
	case "u", "up":
		dir = "up"
	case "d", "down":
		dir = "down"
	default:
		return ""
	}

	return dir
}

func ReverseDirection(dir string) string {
	switch dir {
	case "n", "north":
		dir = "south"
	case "s", "south":
		dir = "north"
	case "e", "east":
		dir = "west"
	case "w", "west":
		dir = "east"
	case "u", "up":
		dir = "down"
	case "d", "down":
		dir = "up"
	default:
		return ""
	}
	return dir
}

func RenderItemDescription(item *Item) string {
	bp := EntityMgr.GetItemBlueprintByInstance(item)
	return cfmt.Sprintf("{{%s}}::green\n{{Description: %s}}::white\n", bp.Name, bp.Description)
}

func RenderMobDescription(mob *Mob) string {
	return cfmt.Sprintf("{{%s}}::red\n{{Description: %s}}::white\n", mob.Name, mob.Description)
}

func RenderCharacterDescription(char *Character) string {
	return cfmt.Sprintf("{{%s}}::blue\n{{Description: %s}}::white\n", char.Name, char.Description)
}

func RenderExitDescription(direction string) string {
	return cfmt.Sprintf("{{To the %s, you see an exit.}}::cyan\n", direction)
}

// RollDice simulates rolling a pool of d6s. It returns the number of hits, glitches, and the results of each die.
// A hit is a roll of 5 or 6, and a glitch is a roll of 1.
func RollDice(pool int) (hits int, glitches int, results []int) {
	rand.Seed(uint64(time.Now().UnixNano()))
	results = make([]int, pool)

	for i := 0; i < pool; i++ {
		die := rand.Intn(6) + 1
		results[i] = die
		if die >= 5 {
			hits++
		} else if die == 1 {
			glitches++
		}
	}

	return hits, glitches, results
}

// RollResultsTotal calculates the total of the roll results.
func RollResultsTotal(results []int) int {
	total := 0
	for _, result := range results {
		total += result
	}
	return total
}

// CheckGlitch determines if the roll results in a glitch or critical glitch.
// A glitch is when more than half of the dice are glitches (rolls of 1).
// A critical glitch occurs if there is a glitch and no hits.
func CheckGlitch(pool int, hits int, glitches int) (bool, bool) {
	if glitches > pool/2 { // More than half are glitches
		if hits == 0 {
			return true, true // Critical glitch
		}
		return true, false // Regular glitch
	}
	return false, false // No glitch
}

// RollWithEdge adds exploding dice (re-rolling 6s) to the dice pool.
func RollWithEdge(pool int) (hits int, glitches int, results []int) {
	rand.Seed(uint64(time.Now().UnixNano()))
	results = []int{}

	for pool > 0 {
		die := rand.Intn(6) + 1
		results = append(results, die)
		if die >= 5 {
			hits++
		}
		if die == 1 {
			glitches++
		}
		if die == 6 {
			pool++ // Exploding sixes
		}
		pool--
	}

	return hits, glitches, results
}
