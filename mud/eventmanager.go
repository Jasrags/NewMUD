package mud

import (
	"github.com/asaskevich/EventBus"
	"github.com/rs/zerolog"
)

type EventManager struct {
	Log zerolog.Logger
	bus EventBus.Bus
}

func NewEventManager() *EventManager {
	return &EventManager{
		Log: NewDevLogger(),
		bus: EventBus.New(),
	}
}

func (em *EventManager) Subscribe(event string, fn interface{}) {
	em.Log.Debug().Str("event", event).Msg("Subscribing to event")

	em.bus.Subscribe(event, fn)
}

func (em *EventManager) Publish(event string, args ...interface{}) {
	em.Log.Debug().Str("event", event).Msg("Publishing event")
	em.bus.Publish(event, args...)
}

func (em *EventManager) Unsubscribe(event string, fn interface{}) {
	em.Log.Debug().Str("event", event).Msg("Unsubscribing from event")
	em.bus.Unsubscribe(event, fn)
}
