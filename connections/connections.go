package connections

import (
	"log/slog"
	"net"
	"sync"

	"github.com/google/uuid"
)

var (
	lock           sync.RWMutex
	netConnections = make(map[string]net.Conn)
)

type NetConnection struct {
	Conn net.Conn
	ID   string
}

func NewNetConnection(conn net.Conn) *NetConnection {
	return &NetConnection{
		Conn: conn,
		ID:   uuid.New().String(),
	}
}

func Add(conn net.Conn) {
	lock.Lock()
	defer lock.Unlock()

	netConn := NewNetConnection(conn)
	netConnections[netConn.ID] = netConn.Conn

	slog.Debug("Added connection",
		slog.String("connection_id", netConn.ID),
		slog.String("remote_address", netConn.Conn.RemoteAddr().String()))
}
