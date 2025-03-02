package game

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/exp/rand"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func RollChance(chance int) bool {
	rand.Seed(uint64(time.Now().UnixNano()))
	randomNumber := rand.Intn(101)

	return randomNumber <= chance

	// r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	// rint := r.Int()
	// slog.Debug("Rolling chance",
	// 	slog.Int("chance", chance),
	// 	slog.Int("roll", rint))

	// return rint >= chance
	// return r.Int() >= chance
}

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
				result.WriteString(line + CRLF)
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

func RemoveFile(filePath string) error {
	return os.Remove(filePath)
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
	return cfmt.Sprintf("{{%s}}::green\n{{Description: %s}}::white"+CRLF, bp.Name, bp.Description)
}

func RenderMobDescription(mob *Mob) string {
	return cfmt.Sprintf("{{%s}}::red\n{{Description: %s}}::white"+CRLF, mob.Name, mob.Description)
}

func RenderCharacterDescription(char *Character) string {
	return cfmt.Sprintf("{{%s}}::blue\n{{Description: %s}}::white"+CRLF, char.Name, char.Description)
}

func RenderExitDescription(direction string) string {
	return cfmt.Sprintf("{{To the %s, you see an exit.}}::cyan"+CRLF, direction)
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

func WriteString(w io.Writer, s string) (int, error) {
	s = strings.ReplaceAll(s, ""+CRLF, CRLF)
	return io.WriteString(w, cfmt.Sprint(s))
}

func WriteStringF(w io.Writer, s string, a ...interface{}) (int, error) {
	s = strings.ReplaceAll(s, ""+CRLF, CRLF)
	return io.WriteString(w, cfmt.Sprintf(s, a...))
}

func PressEnterPrompt(s ssh.Session, label string) {
	WriteString(s, label)
	term := term.NewTerminal(s, "")
	if _, err := term.ReadLine(); err != nil {
		slog.Error("Error reading input", slog.Any("error", err))
		s.Close()
	}
}

func YesNoPrompt(s ssh.Session, def bool) bool {
	choices := "{{Y}}::green|bold{{/}}::white|bold{{n}}::red"
	if !def {
		choices = "{{y}}::green{{/}}::white|bold{{N}}::red|bold"
	}

	term := term.NewTerminal(s, "")
	for {
		WriteStringF(s, "{{Do you want to continue?}}::white|bold {{(}}::white|bold%s{{):}}::white|bold ", choices)
		input, err := term.ReadLine()
		if err != nil {
			slog.Error("Error reading input", slog.Any("error", err))
			s.Close()
		}

		input = strings.ToLower(strings.TrimSpace(input))

		if input == "" {
			return def
		}

		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			WriteStringF(s, "{{Invalid choice, %q please try again.}}::red"+CRLF, input)
		}
	}
}

func InputPrompt(s ssh.Session, prompt string) (string, error) {
	t := term.NewTerminal(s, prompt)
	input, err := t.ReadLine()
	if err != nil {
		slog.Error("Error reading input", slog.Any("error", err))
		s.Close()

		return "", err
	}

	return strings.TrimSpace(input), nil
}

func PasswordPrompt(s ssh.Session, prompt string) (string, error) {
	t := term.NewTerminal(s, prompt)
	input, err := t.ReadPassword(prompt)
	if err != nil {
		slog.Error("Error reading password", slog.Any("error", err))
		s.Close()

		return "", err
	}

	return strings.TrimSpace(input), nil
}

func SendToChar(s ssh.Session, message string) {
	WriteStringF(s, "%s", message)
}

// void send_to_all(char *messg)

// void send_to_room(char *messg, int room)
func SendToRoom(s ssh.Session, message string,
	room *Room) {
	// for _, c := range room.Characters {
	// 	io.WriteString(s, cfmt.Sprintf("{{%s}}::white"+CRLF, message))
	// }
}

type MenuOption struct {
	DisplayText string
	Value       string
	Description string
}

func PromptForMenu(s ssh.Session, title string, options []MenuOption) (string, error) {
	for {
		var menuBuilder strings.Builder
		menuBuilder.WriteString(cfmt.Sprintf("\n{{%s}}::white|bold|underline\n\n", title))

		for i, option := range options {
			menuBuilder.WriteString(cfmt.Sprintf("{{%d}}::green|bold. {{%s}}::white|bold\n", i+1, option.DisplayText))
		}
		menuBuilder.WriteString(cfmt.Sprint("\n{{Enter choice or info <choice> for details:}}::white|bold "))

		WriteString(s, menuBuilder.String())

		input, err := InputPrompt(s, "")
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)

		// Handle info request
		if strings.HasPrefix(strings.ToLower(input), "info ") {
			detailChoice := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(input), "info "))
			numChoice, err := strconv.Atoi(detailChoice)
			if err == nil && numChoice > 0 && numChoice <= len(options) {
				// WriteString(s, borderStyle.Render(cfmt.Sprintf("{{%s}}::cyan", options[numChoice-1].Description)))
				WriteStringF(s, "\n{{%s}}::cyan\n", options[numChoice-1].Description)
				continue
			}
			WriteString(s, "{{Invalid selection. Please try again.}}::red\n")
			continue
		}

		// Check if input is a number
		numChoice, err := strconv.Atoi(input)
		if err == nil && numChoice > 0 && numChoice <= len(options) {
			return options[numChoice-1].Value, nil
		}

		// Check if input matches an option value
		for _, option := range options {
			if strings.EqualFold(input, option.Value) {
				return option.Value, nil
			}
		}

		WriteString(s, "{{Invalid selection. Please try again.}}::red\n")
	}
}

// HasAnyTag returns true if filterTags is empty or if any tag in filterTags is found (case-insensitive)
// in the object's tags.
func HasAnyTag(objectTags []string, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, ft := range filterTags {
		for _, ot := range objectTags {
			if strings.EqualFold(ot, ft) {
				return true
			}
		}
	}
	return false
}

// FindMobsByName searches the current room's mobs and returns all instances
// that match the provided name (case-insensitive).
func FindMobsByName(room *Room, name string) []*Mob {
	var matches []*Mob
	room.RLock()
	defer room.RUnlock()

	for _, mob := range room.Mobs {
		if strings.EqualFold(mob.Name, name) {
			matches = append(matches, mob)
		}
	}
	return matches
}

// RenderAttribute renders a single attribute for display.
func RenderAttribute[T int | float64](name string, attr Attribute[T]) string {
	var output strings.Builder

	if attr.Base == 0 {
		return ""
	}

	output.WriteString(attrNameStyle.Render(fmt.Sprintf("%-10s", name)))
	output.WriteString(attrValueStyle.Render(fmt.Sprintf(" %-2v", renderValue(attr.Base))))
	if attr.TotalValue != attr.Base {
		style := attrPosModStyle
		if attr.TotalValue < attr.Base {
			style = attrNegModStyle
		}
		output.WriteString(style.Render(fmt.Sprintf(" (%v)", renderValue(attr.TotalValue))))
	}

	return output.String()
}

// renderValue formats the value of an attribute for display.
func renderValue[T int | float64](value T) string {
	switch v := any(value).(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
func RenderKeyValue(key, value string) string {
	return fmt.Sprintf("%s: %s", attrNameStyle.Render(key), attrTextValueStyle.Render(value))
}
