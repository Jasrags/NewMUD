package commands

import (
	"strings"

	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/i582/cfmt/cmd/cfmt"
)

func Help(args []string, user *users.User, room *rooms.Room) {
	uniqueCommands := make(map[string]*Command)
	for _, cmd := range commandList {
		uniqueCommands[cmd.Name] = cmd
	}

	var builder strings.Builder
	builder.WriteString(cfmt.Sprintf("{{Available commands:}}::white|bold\n"))
	// for _, cmd := range uniqueCommands {
	// 	builder.WriteString(cfmt.Sprintf("{{%s}}::cyan - %s (aliases: %s)\n", cmd.Name, cmd.Description, strings.Join(cmd.Aliases, ", ")))
	// }
	// io.WriteString(player.Conn, builder.String())
}
