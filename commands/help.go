package commands

import (
	"io"
	"strings"

	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

func Help(s ssh.Session, args []string, user *users.User, room *rooms.Room) {
	uniqueCommands := make(map[string]*Command)
	for _, cmd := range Mgr.GetCommands() {
		uniqueCommands[cmd.Name] = cmd
	}

	var builder strings.Builder
	builder.WriteString(cfmt.Sprintf("{{Available commands:}}::white|bold\n"))
	for _, cmd := range uniqueCommands {
		builder.WriteString(cfmt.Sprintf("{{%s}}::cyan - %s (aliases: %s)\n", cmd.Name, cmd.Description, strings.Join(cmd.Aliases, ", ")))
	}

	io.WriteString(s, builder.String())
}
