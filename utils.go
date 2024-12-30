package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"gopkg.in/yaml.v3"
)

// CreateEntityRef creates an entity reference from an area and ID.
func CreateEntityRef(area, id string) string {
	return strings.ToLower(fmt.Sprintf("%s:%s", area, id))
}

// ParseEntityRef parses an entity reference into its area and ID parts.
func ParseEntityRef(entityRef string) (area, id string) {
	parts := strings.Split(strings.ToLower(entityRef), ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// WrapText splits text into lines of the specified width without breaking words.
func WrapText(text string, width int) string {
	var result strings.Builder
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > width {
			result.WriteString(line + "\n")
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}
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

// RenderRoom renders the room to a string for the player.
func RenderRoom(user *User, room *Room) string {
	var builder strings.Builder

	char := user.ActiveCharacter

	// Optionally display the room ID for admins
	if user.ActiveCharacter.Role == CharacterRoleAdmin {
		builder.WriteString(cfmt.Sprintf("{{[%s] }}::green", room.ID))
	}

	// Display the room title
	builder.WriteString(cfmt.Sprintf("{{%s}}::#4287f5\n", room.Title))

	// Display the room description
	builder.WriteString(cfmt.Sprintf("{{%s}}::white\n", WrapText(room.Description, 80)))

	// Display players in the room
	charCount := len(room.Characters)
	var charNames []string
	for _, c := range room.Characters {
		if c.Name != char.Name {
			color := "cyan"
			if c.Role == CharacterRoleAdmin {
				color = "yellow"
			}

			charNames = append(charNames, cfmt.Sprintf("{{%s}}::%s", c.Name, color))
		}
	}
	if charCount == 1 {
		builder.WriteString(cfmt.Sprint("{{You are alone in the room.}}::cyan\n"))
	} else if charCount >= 2 {
		builder.WriteString(cfmt.Sprintf("{{There is %d other person in the room: }}::cyan|bold", charCount-1))
	} else {
		builder.WriteString(cfmt.Sprintf("{{There are %d other people in the room: }}::cyan|bold", charCount-1))
	}
	if len(charNames) > 0 {
		builder.WriteString(cfmt.Sprintf("{{%s}}::cyan", WrapText(strings.Join(charNames, ", "), 80)))
	}
	builder.WriteString("\n")

	// Display exits
	if len(room.Exits) == 0 {
		builder.WriteString(cfmt.Sprint("{{There are no exits.}}::red\n"))
	} else {
		builder.WriteString(cfmt.Sprint("{{Exits:}}::#2359b0\n"))
		for _, exit := range char.Room.Exits {
			builder.WriteString(cfmt.Sprintf("{{ %-5s - %s}}::#2359b0\n", exit.Direction, exit.Room.Title))
		}
	}

	return cfmt.Sprint(builder.String())
}
