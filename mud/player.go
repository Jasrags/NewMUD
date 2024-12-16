package mud

import "net"

type Player struct {
	Name string
	Room *Room
	Out  chan string // For sending messages to the player
	Conn net.Conn
}
