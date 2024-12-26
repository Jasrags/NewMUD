package connections

var (
	Mgr = NewManager()
)

type Manager struct {
	connections map[string]*NetConnection
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*NetConnection),
	}
}

func (mgr *Manager) Add(nc *NetConnection) {
	mgr.connections[nc.ID] = nc
}

func (mgr *Manager) Get(id string) *NetConnection {
	return mgr.connections[id]
}

func (mgr *Manager) Remove(id string) {
	delete(mgr.connections, id)
}
