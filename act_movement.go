package main

import (
	"io"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"golang.org/x/exp/rand"
)

// This command is for opening closed entities
func DoOpen(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Open what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if !exit.Door.IsClosed {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is already open.}}::yellow\n", direction))
		return
	}

	if exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is locked.}}::red\n", direction))
		return
	}

	exit.Door.IsClosed = false
	io.WriteString(s, cfmt.Sprintf("{{You open the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s opens the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s opens from the other side.}}::green\n", ReverseDirection(direction)), []string{})
	}
}

func DoClose(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Close what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])

	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if exit.Door.IsClosed {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is already closed.}}::yellow\n", direction))
		return
	}

	exit.Door.IsClosed = true
	io.WriteString(s, cfmt.Sprintf("{{You close the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s closes the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})

	// Notify the adjacent room
	if exit.Room != nil {
		exit.Room.Broadcast(cfmt.Sprintf("{{The door to the %s closes from the other side.}}::green\n", ReverseDirection(direction)), []string{})
	}
}

/*
Usage:
  - move <north,n,south,s,east,e,west,w,up,u,down,d>
  - <north,n,south,s,east,e,west,w,up,u,down,d>
*/
func DoMove(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if room == nil {
		io.WriteString(s, cfmt.Sprintf("{{You are not in a room.}}::red\n"))
		return
	}

	if cmd == "move" && len(args) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Move where?}}::red\n"))
		return
	}

	dir := ParseDirection(cmd)

	// Check if the exit exists
	if exit, ok := char.Room.Exits[dir]; ok {
		if exit.Door != nil && exit.Door.IsClosed {
			io.WriteString(s, cfmt.Sprintf("{{The door to the %s is closed.}}::red\n", dir))
			return
		}

		char.MoveToRoom(exit.Room)
		char.Save()

		io.WriteString(s, cfmt.Sprintf("You move %s.\n\n", dir))
		io.WriteString(s, RenderRoom(user, char, nil))
	} else {
		io.WriteString(s, cfmt.Sprintf("{{You can't go that way.}}::red\n"))
		return
	}
}

func DoLock(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Lock what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])
	if direction == "" {
		io.WriteString(s, cfmt.Sprintf("{{Invalid direction.}}::red\n"))
		return
	}

	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is already locked.}}::yellow\n", direction))
		return
	}

	if !exit.Door.IsClosed {
		io.WriteString(s, cfmt.Sprintf("{{You must close the door to the %s before locking it.}}::yellow\n", direction))
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
		io.WriteString(s, cfmt.Sprintf("{{You don't have the key to lock the door to the %s.}}::red\n", direction))
		return
	}

	exit.Door.IsLocked = true
	io.WriteString(s, cfmt.Sprintf("{{You lock the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s locks the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
}

func DoUnlock(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Unlock what?}}::yellow\n"))
		return
	}

	direction := ParseDirection(args[0])
	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if !exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is not locked.}}::yellow\n", direction))
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
		io.WriteString(s, cfmt.Sprintf("{{You don't have the key to unlock the door to the %s.}}::red\n", direction))
		return
	}

	exit.Door.IsLocked = false
	io.WriteString(s, cfmt.Sprintf("{{You unlock the door to the %s.}}::green\n", direction))
	room.Broadcast(cfmt.Sprintf("{{%s unlocks the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
}

func DoPick(s ssh.Session, cmd string, args []string, user *User, char *Character, room *Room) {
	if len(args) < 1 {
		io.WriteString(s, cfmt.Sprintf("{{Pick what?}}::yellow\n"))
		return
	}

	direction := args[0]
	exit, exists := room.Exits[direction]
	if !exists {
		io.WriteString(s, cfmt.Sprintf("{{There is no exit to the %s.}}::red\n", direction))
		return
	}

	if exit.Door == nil {
		io.WriteString(s, cfmt.Sprintf("{{There is no door to the %s.}}::red\n", direction))
		return
	}

	if !exit.Door.IsLocked {
		io.WriteString(s, cfmt.Sprintf("{{The door to the %s is not locked.}}::yellow\n", direction))
		return
	}

	// if !hasKey {
	//     if exit.Door.PickDifficulty > 0 {
	//         success := AttemptLockPick(char, exit.Door.PickDifficulty)
	//         if success {
	//             exit.Door.IsLocked = false
	//             io.WriteString(s, cfmt.Sprintf("{{You successfully pick the lock on the door to the %s.}}::green\n", direction))
	//             room.Broadcast(cfmt.Sprintf("{{%s picks the lock on the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
	//             return
	//         } else {
	//             io.WriteString(s, cfmt.Sprintf("{{You fail to pick the lock on the door to the %s.}}::red\n", direction))
	//             return
	//         }
	//     }

	//     io.WriteString(s, cfmt.Sprintf("{{You don't have the key to unlock the door to the %s.}}::red\n", direction))
	//     return
	// }

	pickRoll := rand.Intn(100) + 1 // Random roll between 1 and 100
	if pickRoll > exit.Door.PickDifficulty {
		exit.Door.IsLocked = false
		io.WriteString(s, cfmt.Sprintf("{{You successfully pick the lock on the door to the %s.}}::green\n", direction))
		room.Broadcast(cfmt.Sprintf("{{%s picks the lock on the door to the %s.}}::green\n", char.Name, direction), []string{char.ID})
	} else {
		io.WriteString(s, cfmt.Sprintf("{{You fail to pick the lock on the door to the %s.}}::red\n", direction))
	}
}
