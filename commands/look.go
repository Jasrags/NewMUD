package commands

import (
	"io"
	"strings"

	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

func Look(s ssh.Session, args []string, user *users.User, room *rooms.Room) {
	var builder strings.Builder
	builder.WriteString(cfmt.Sprintf("{{%s}}::green\n", room.Title))
	builder.WriteString(cfmt.Sprintf("{{%s}}::white\n", room.Description))

	io.WriteString(s, builder.String())
}
