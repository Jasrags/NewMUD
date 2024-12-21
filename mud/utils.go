package mud

import (
	"fmt"
	"os"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
	"gopkg.in/yaml.v3"
)

// CreateEntityRef creates an entity reference from an area and ID.
func CreateEntityRef(area, id string) string {
	return fmt.Sprintf("%s:%s", area, id)
}

// ParseEntityRef parses an entity reference into its area and ID parts.
func ParseEntityRef(entityRef string) (area, id string) {
	parts := strings.Split(entityRef, ":")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// RenderRoom renders the room to a string for the player.
func RenderRoom(player *Player, room *Room) string {
	var builder strings.Builder
	var args []interface{}

	// Optionally display the room ID for admins
	if player.Role == "admin" {
		builder.WriteString(cfmt.Sprintf("{{[%s] }}::green", room.ID))
	}

	// Display the room title
	builder.WriteString(cfmt.Sprintf("{{%s}}::#4287f5\n", room.Title))

	// Display the room description
	builder.WriteString(cfmt.Sprintf("{{%s}}::white\n", WrapText(room.Description, 80)))

	// Display players in the room
	playerCount := len(room.Players)
	var playerNames []string
	for _, p := range room.Players {
		if p.Name != player.Name {
			color := "cyan"
			if p.Role == "admin" {
				color = "yellow"
			}

			playerNames = append(playerNames, cfmt.Sprintf("{{%s}}::%s", p.Name, color))
		}
	}
	if playerCount == 1 {
		builder.WriteString(cfmt.Sprint("{{You are alone in the room.}}::cyan\n"))
	} else if playerCount >= 2 {
		builder.WriteString(cfmt.Sprintf("{{There is %d other person in the room: }}::cyan|bold", playerCount-1))
	} else {
		builder.WriteString(cfmt.Sprintf("{{There are %d other people in the room: }}::cyan|bold", playerCount-1))
	}
	if len(playerNames) > 0 {
		builder.WriteString(cfmt.Sprintf("{{%s}}::cyan", WrapText(strings.Join(playerNames, ", "), 80)))
	}
	builder.WriteString("\n")

	// Display exits
	if len(room.Exits) == 0 {
		builder.WriteString(cfmt.Sprint("{{There are no exits.}}::red\n"))
	} else {
		builder.WriteString(cfmt.Sprint("{{Exits:}}::#2359b0\n"))
		for _, exit := range player.Room.Exits {
			builder.WriteString(cfmt.Sprintf("{{ %-5s - %s}}::#2359b0\n", exit.Direction, exit.Room.Title))
		}
	}

	return cfmt.Sprintf(builder.String(), args...)
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
