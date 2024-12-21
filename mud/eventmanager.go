package mud

import (
	"github.com/rs/zerolog"
	ee "github.com/vansante/go-event-emitter"
)

type EventManager struct {
	Log zerolog.Logger
	e   ee.EventEmitter
	o   ee.Observable
}

func NewEventManager(l zerolog.Logger, async bool) *EventManager {
	e := ee.NewEmitter(async)
	return &EventManager{
		Log: l,
		e:   e,
		o:   e,
		// bus: EventBus.New(),
	}
}

func (em *EventManager) Subscribe(event ee.EventType, fn ee.HandleFunc) {
	em.Log.Debug().Str("event", string(event)).Msg("Subscribing to event")
	em.o.AddListener(event, fn)
	// em.bus.Subscribe(event, fn)
}

func (em *EventManager) Publish(event ee.EventType, args ...interface{}) {
	em.Log.Debug().Str("event", string(event)).Msg("Publishing event")
	em.e.EmitEvent(event, args...)
	// em.e.Emit(event, args...)
	// em.bus.Publish(event, args...)
}

func (em *EventManager) Unsubscribe(event ee.EventType, listener *ee.Listener) {
	em.Log.Debug().Str("event", string(event)).Msg("Unsubscribing from event")
	em.o.RemoveListener(event, listener)
	// em.bus.Unsubscribe(event, fn)
}
