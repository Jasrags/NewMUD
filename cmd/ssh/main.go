package main

import (
	"log/slog"

	"github.com/Jasrags/NewMUD/internal/game"
)

func main() {
	gs := game.NewGameServer()
	gs.Init()
	gs.Start()
	slog.Info("Shutting down")
}
