package events

import (
	"github.com/spf13/viper"
	ee "github.com/vansante/go-event-emitter"
)

var (
	Mgr = NewManager()
)

type Manager struct {
	em ee.EventEmitter
	ob ee.Observable
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	e := ee.NewEmitter(viper.GetBool("server.async_events"))

	return &Manager{
		em: e,
		ob: e,
	}
}

func (mgr *Manager) Publish(event ee.EventType, args ...interface{}) {
	mgr.em.EmitEvent(event, args...)
}
func (mgr *Manager) Subscribe(event ee.EventType, fn ee.HandleFunc) *ee.Listener {
	return mgr.ob.AddListener(event, fn)
}
