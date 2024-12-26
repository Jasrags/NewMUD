package connections

import (
	"log/slog"
	"sync"

	"github.com/gliderlabs/ssh"
)

var (
	lock           sync.RWMutex
	netConnections = make(map[string]ssh.Session)
)

type NetConnection struct {
	Session ssh.Session
	ID      string
}

func NewNetConnection(s ssh.Session) *NetConnection {
	return &NetConnection{
		Session: s,
		ID:      s.Context().SessionID(),
	}
}

// func (nc *NetConnection) Write(data string) {
// 	io.WriteString(nc.Session., data)
// }

func (nc *NetConnection) Close() {
	nc.Session.Close()
}

func Add(s ssh.Session) {
	slog.Debug("Adding connection",
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))
	lock.Lock()
	defer lock.Unlock()

	netConn := NewNetConnection(s)
	netConnections[netConn.ID] = netConn.Session

	slog.Debug("Added connection",
		slog.String("connection_id", netConn.ID),
		slog.String("remote_address", s.RemoteAddr().String()))
}
