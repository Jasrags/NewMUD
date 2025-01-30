package game

import (
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
)

// TODO: support short versions of the prompt placeholders
// TODO: need a DoPrompt function that will handle the actual printing of the prompt
var (
	promptPlaceholders = map[string]func(*Character) string{
		"{{time}}": GetFormattedGameTime,
		"{{date}}": GetFormattedGameDate,
	}
)

// RenderPrompt dynamically substitutes placeholders in the prompt string
func RenderPrompt(char *Character) string {
	prompt := char.Prompt

	for pattern, function := range promptPlaceholders {
		prompt = strings.ReplaceAll(prompt, pattern, function(char))
	}

	return cfmt.Sprintf("%s ", prompt)
}

// GetFormattedGameTime returns the in-game time formatted as HH:MM AM/PM
func GetFormattedGameTime(char *Character) string {
	return GameTimeMgr.GetFormattedTime()
}

// GetFormattedGameDate returns the in-game date formatted as "Month Day, Year"
func GetFormattedGameDate(char *Character) string {
	return GameTimeMgr.GetFormattedDate(true)
}

func pluralize(value int) string {
	if value == 1 {
		return ""
	}
	return "s"
}
