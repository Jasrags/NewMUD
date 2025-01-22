package main

// import (
// 	"github.com/spf13/viper"
// 	ee "github.com/vansante/go-event-emitter"
// )

// var (
// 	EventMgr = NewEventManager()
// )

// type EventManager struct {
// 	em ee.EventEmitter
// 	ob ee.Observable
// }

// // NewEventManager creates a new Manager.
// func NewEventManager() *EventManager {
// 	e := ee.NewEmitter(viper.GetBool("server.async_events"))

// 	return &EventManager{
// 		em: e,
// 		ob: e,
// 	}
// }

// func (mgr *EventManager) Publish(event ee.EventType, args ...interface{}) {
// 	mgr.em.EmitEvent(event, args...)
// }
// func (mgr *EventManager) Subscribe(event ee.EventType, fn ee.HandleFunc) *ee.Listener {
// 	return mgr.ob.AddListener(event, fn)
// }
