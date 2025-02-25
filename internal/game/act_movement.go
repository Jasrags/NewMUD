package game

import (
	"github.com/gliderlabs/ssh"
)

/*
Usage:
  - move <north,n,south,s,east,e,west,w,up,u,down,d>
  - <north,n,south,s,east,e,west,w,up,u,down,d>
*/
func DoMove(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if cmd == "move" && len(args) == 0 {
		WriteString(s, "{{Move where?}}::red"+CRLF)
		return
	}

	dir := ParseDirection(cmd)

	// Check if the exit exists
	if exit, ok := char.Room.Exits[dir]; ok {
		if exit.Door != nil && exit.Door.IsClosed {
			WriteStringF(s, "{{The door to the %s is closed.}}::red"+CRLF, dir)
			return
		}

		char.MoveToRoom(exit.Room)
		char.Save()

		WriteStringF(s, "You move %s."+CRLF, dir)
		WriteString(s, RenderRoom(user, char, nil))
	} else {
		WriteString(s, "{{You can't go that way.}}::red"+CRLF)
		return
	}
}
