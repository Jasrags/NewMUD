package game

import (
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/exp/rand"
)

// This command is for opening closed entities
func DoOpen(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Open what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if !exit.Door.IsClosed {
		WriteStringF(s, "{{The door to the %s is already open.}}::yellow"+CRLF, direction)
		return
	}

	if exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is locked.}}::red"+CRLF, direction)
		return
	}

	exit.Door.IsClosed = false
	WriteStringF(s, "{{You open the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s opens the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s opens from the other side.}}::green"+CRLF, ReverseDirection(direction)), []string{})
	}
}

func DoClose(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Close what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])

	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door.IsClosed {
		WriteStringF(s, "{{The door to the %s is already closed.}}::yellow"+CRLF, direction)
		return
	}

	exit.Door.IsClosed = true
	WriteStringF(s, "{{You close the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s closes the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s closes from the other side.}}::green"+CRLF, ReverseDirection(direction)), []string{})
	}
}

/*
Usage:
  - move <north,n,south,s,east,e,west,w,up,u,down,d>
  - <north,n,south,s,east,e,west,w,up,u,down,d>
*/
func DoMove(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if room == nil {
		WriteString(s, "{{You are not in a room.}}::red"+CRLF)
		return
	}

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

func DoLock(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Lock what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])
	if direction == "" {
		WriteString(s, "{{Invalid direction.}}::red"+CRLF)
		return
	}

	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is already locked.}}::yellow"+CRLF, direction)
		return
	}

	if !exit.Door.IsClosed {
		WriteStringF(s, "{{You must close the door to the %s before locking it.}}::yellow"+CRLF, direction)
		return
	}

	validKeys := make(map[string]bool)
	for _, key := range exit.Door.KeyIDs {
		validKeys[key] = true
	}

	hasKey := false
	for _, item := range char.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp.Type == ItemTypeKey && validKeys[bp.ID] {
			hasKey = true
			break
		}
	}

	if !hasKey {
		WriteStringF(s, "{{You don't have the key to lock the door to the %s.}}::red"+CRLF, direction)
		return
	}

	exit.Door.IsLocked = true
	WriteStringF(s, "{{You lock the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s locks the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})
}

func DoUnlock(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Unlock what?}}::yellow"+CRLF)
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if !exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is not locked.}}::yellow"+CRLF, direction)
		return
	}

	validKeys := make(map[string]bool)
	for _, key := range exit.Door.KeyIDs {
		validKeys[key] = true
	}

	// Check if character has the correct key
	hasKey := false
	for _, item := range char.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(item)
		if bp.Type == ItemTypeKey && validKeys[bp.ID] {
			hasKey = true
			break
		}
	}

	if !hasKey {
		WriteStringF(s, "{{You don't have the key to unlock the door to the %s.}}::red"+CRLF, direction)
		return
	}

	exit.Door.IsLocked = false
	WriteStringF(s, "{{You unlock the door to the %s.}}::green"+CRLF, direction)
	room.Broadcast(cfmt.Sprintf("{{%s unlocks the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})
}

func DoPick(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room) {
	if len(args) < 1 {
		WriteString(s, "{{Pick what?}}::yellow"+CRLF)
		return
	}

	direction := args[0]
	exit, exists := room.Exits[direction]
	if !exists {
		WriteStringF(s, "{{There is no exit to the %s.}}::red"+CRLF, direction)
		return
	}

	if exit.Door == nil {
		WriteStringF(s, "{{There is no door to the %s.}}::red"+CRLF, direction)
		return
	}

	if !exit.Door.IsLocked {
		WriteStringF(s, "{{The door to the %s is not locked.}}::yellow"+CRLF, direction)
		return
	}

	// if !hasKey {
	//     if exit.Door.PickDifficulty > 0 {
	//         success := AttemptLockPick(char, exit.Door.PickDifficulty)
	//         if success {
	//             exit.Door.IsLocked = false
	//             WriteString(s, "{{You successfully pick the lock on the door to the %s.}}::green"+CRLF, direction)
	//             room.Broadcast("{{%s picks the lock on the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID}
	//             return
	//         } else {
	//             WriteString(s, "{{You fail to pick the lock on the door to the %s.}}::red"+CRLF, direction)
	//             return
	//         }
	//     }

	//     WriteString(s, "{{You don't have the key to unlock the door to the %s.}}::red"+CRLF, direction)
	//     return
	// }

	pickRoll := rand.Intn(100) + 1 // Random roll between 1 and 100
	if pickRoll > exit.Door.PickDifficulty {
		exit.Door.IsLocked = false
		WriteStringF(s, "{{You successfully pick the lock on the door to the %s.}}::green"+CRLF, direction)
		room.Broadcast(cfmt.Sprintf("{{%s picks the lock on the door to the %s.}}::green"+CRLF, char.Name, direction), []string{char.ID})
	} else {
		WriteStringF(s, "{{You fail to pick the lock on the door to the %s.}}::red"+CRLF, direction)
	}
}
