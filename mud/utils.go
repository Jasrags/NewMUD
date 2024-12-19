package mud

import (
	"fmt"
	"strings"

	"github.com/i582/cfmt/cmd/cfmt"
)

func CreateEntityRef(area, id string) string {
	return fmt.Sprintf("%s:%s", area, id)
}

func RenderRoom(player *Player, room *Room) string {
	var builder strings.Builder
	var args []interface{}

	builder.WriteString("{{%s}}::green|bold")
	args = append(args, room.Title)
	if player.Role == "admin" {
		builder.WriteString("{{ (ID: %s)}}::white")
		args = append(args, room.ID)
	}
	// builder.WriteString("\n")
	builder.WriteString("\n{{%s}}::white\n\n")
	args = append(args, WrapText(room.Description, 80))

	// Display players in the room
	if len(room.Players) == 1 {
		builder.WriteString("\n{{You are alone in the room.}}::cyan\n")
	} else if len(room.Players) >= 2 {
		builder.WriteString("\n{{Players in the room:}}::cyan|bold\n")
		for _, p := range room.Players {
			if p.Name != player.Name {
				builder.WriteString("{{ - %s}}::cyan\n")
				args = append(args, p.Name)
			}
		}
	}

	// Display exits
	if len(room.Exits) == 0 {
		builder.WriteString("{{There are no exits.}}::red\n")
	} else {
		builder.WriteString("{{Exits:}}::yellow|bold\n")
		for direction, _ := range player.Room.Exits {
			builder.WriteString("{{ - %s}}::yellow\n")
			args = append(args, direction)
		}
	}

	// // Display other players in the room
	// otherPlayers := []string{}
	// for _, p := range room.Players {
	// 	if p.Name != player.Name {
	// 		otherPlayers = append(otherPlayers, p.Name)
	// 	}
	// }
	// if len(otherPlayers) > 0 {
	// 	builder.WriteString("\n{{Players here:}}::cyan|bold\n")
	// 	for _, p := range otherPlayers {
	// 		builder.WriteString("{{ - %s}}::cyan\n")
	// 		args = append(args, p)
	// 	}
	// }

	return cfmt.Sprintf(builder.String(), args...)
}

// WrapText splits text into lines of the specified width without breaking words.
func WrapText(text string, width int) string {
	var result strings.Builder
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		// Check if adding the next word would exceed the width
		if len(line)+len(word)+1 > width {
			result.WriteString(line + "\n") // Write the current line and start a new one
			line = word
		} else {
			if line != "" {
				line += " " // Add a space before appending the word
			}
			line += word
		}
	}
	if line != "" {
		result.WriteString(line) // Append the last line
	}

	return result.String()
}
