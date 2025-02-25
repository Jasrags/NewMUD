package game

import (
	"log/slog"
	"sync"

	ee "github.com/vansante/go-event-emitter"
)

const (
	AreasFilepath = "_data/areas"
	AreasFilename = "manifest.yml"
)

type (
	Area struct {
		sync.RWMutex
		Listeners []ee.Listener `yaml:"-"`

		ID          string `yaml:"id"`
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
	}
)

func NewArea() *Area {
	return &Area{}
}

func (a *Area) Init() {
	slog.Debug("Initializing area",
		slog.String("area_id", a.ID))
}
