package main

// import (
// 	"io"
// 	"strings"

// 	"github.com/gliderlabs/ssh"
// 	"github.com/i582/cfmt/cmd/cfmt"
// )

// /*
// Usage:
//   - say <message>
// */
// // TODO: overall for communication commands we need to log messages to a database with time, to/from, and message.
// // TODO: need to implement a block/unblock function for preventing messages from certain users
// func DoSay(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
// 	if room == nil {
// 		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
// 		return
// 	}

// 	if len(args) == 0 {
// 		io.WriteString(s, cfmt.Sprintf("{{What do you want to say?}}::red\n"))
// 		return
// 	}

// 	message := strings.Join(args, " ")

// 	// Broadcast message to the room
// 	room.Broadcast(cfmt.Sprintf("{{%s says: \"%s\"}}::green\n", char.Name, message), []string{char.ID})

// 	// Message the player
// 	io.WriteString(s, cfmt.Sprintf("{{You say: \"%s\"}}::green\n", message))
// }

// func DoTell(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
// 	if room == nil {
// 		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
// 		return
// 	}

// 	if len(args) < 2 {
// 		io.WriteString(s, cfmt.Sprintf("{{Usage: tell <username> <message>.}}::red\n"))
// 		return
// 	}

// 	recipientName := args[0]
// 	message := strings.Join(args[1:], " ")

// 	var recipient *Character
// 	for _, r := range room.Characters {
// 		if strings.EqualFold(r.Name, recipientName) {
// 			recipient = r
// 			break
// 		}
// 	}

// 	if recipient == nil {
// 		io.WriteString(s, cfmt.Sprintf("{{There is no one named '%s' here.}}::yellow\n", recipientName))
// 		return
// 	}

// 	// Message the recipient
// 	io.WriteString(recipient.Conn, cfmt.Sprintf("{{%s tells you: \"%s\"}}::cyan\n", char.Name, message))

// 	// Message the sender
// 	io.WriteString(s, cfmt.Sprintf("{{You tell %s: \"%s\"}}::green\n", recipient.Name, message))

// 	// Message the room (excluding sender and recipient)
// 	room.Broadcast(cfmt.Sprintf("{{%s tells %s something privately.}}::green\n", char.Name, recipient.Name), []string{char.ID, recipient.ID})
// }

// func SuggestTell(line string, args []string, char *Character, room *Room) []string {
// 	suggestions := []string{}

// 	switch len(args) {
// 	case 0: // Suggest names of characters in the room
// 		for _, r := range room.Characters {
// 			if !strings.EqualFold(r.Name, char.Name) { // Exclude self
// 				suggestions = append(suggestions, r.Name)
// 			}
// 		}
// 	case 1: // Suggest partial names
// 		for _, r := range room.Characters {
// 			if !strings.EqualFold(r.Name, char.Name) && strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(args[0])) {
// 				suggestions = append(suggestions, r.Name)
// 			}
// 		}
// 	}

// 	return suggestions
// }
