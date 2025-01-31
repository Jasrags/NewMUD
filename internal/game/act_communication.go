package game

import (
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
)

/*
Usage:
  - say <message>
*/
// TODO: overall for communication commands we need to log messages to a database with time, to/from, and message.
// TODO: need to implement a block/unblock function for preventing messages from certain users
func DoSay(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if room == nil {
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

	if len(args) == 0 {
		WriteString(s, "{{What do you want to say?}}::red"+CRLF)
		return
	}

	message := strings.Join(args, " ")

	// Broadcast message to the room
	room.Broadcast(cfmt.Sprintf("{{%s says: \"%s\"}}::green"+CRLF, char.Name, message), []string{char.ID})

	// Message the player
	WriteStringF(s, "{{You say: \"%s\"}}::green"+CRLF, message)
}

func DoTell(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if room == nil {
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

	if len(args) < 2 {
		WriteString(s, "{{Usage: tell <username> <message>.}}::red"+CRLF)
		return
	}

	recipientName := args[0]
	message := strings.Join(args[1:], " ")

	var recipient *Character
	for _, r := range room.Characters {
		if strings.EqualFold(r.Name, recipientName) {
			recipient = r
			break
		}
	}

	if recipient == nil {
		WriteStringF(s, "{{There is no one named '%s' here.}}::yellow"+CRLF, recipientName)
		return
	}

	// Message the recipient
	WriteStringF(recipient.Conn, "{{%s tells you: \"%s\"}}::cyan"+CRLF, char.Name, message)

	// Message the sender
	WriteStringF(s, "{{You tell %s: \"%s\"}}::green"+CRLF, recipient.Name, message)

	// Message the room (excluding sender and recipient)
	room.Broadcast(cfmt.Sprintf("{{%s tells %s something privately.}}::green"+CRLF, char.Name, recipient.Name), []string{char.ID, recipient.ID})
}

func SuggestTell(line string, args []string, char *Character, room *Room) []string {
	suggestions := []string{}

	switch len(args) {
	case 0: // Suggest names of characters in the room
		for _, r := range room.Characters {
			if !strings.EqualFold(r.Name, char.Name) { // Exclude self
				suggestions = append(suggestions, r.Name)
			}
		}
	case 1: // Suggest partial names
		for _, r := range room.Characters {
			if !strings.EqualFold(r.Name, char.Name) && strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(args[0])) {
				suggestions = append(suggestions, r.Name)
			}
		}
	}

	return suggestions
}
