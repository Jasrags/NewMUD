package events

import (
	ebus "github.com/asaskevich/EventBus"
)

var (
	bus = ebus.New()
)

type Event interface {
	Type() string
}

func QueueEvent(e Event) {
}
