package commands

import (
	"io"
	"strings"

	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/gliderlabs/ssh"
)

func Look(s ssh.Session, args []string, user *users.User, room *rooms.Room) {
	var builder strings.Builder

	io.WriteString(s, builder.String())
}
