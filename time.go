package main

// import (
// 	"fmt"
// 	"log/slog"
// 	"time"
// )

// const (
// 	GameDayLength      = 1440 // Minutes in a game day (24 hours * 60 minutes)
// 	GameTicksPerMinute = 10   // Number of ticks for one in-game minute
// )

// type GameTime struct {
// 	Minutes         int // Total minutes since the start of the game day
// 	TickAccumulator int // Tracks the number of ticks since the last minute increment
// }

// func (t *GameTime) CurrentHour() int {
// 	return (t.Minutes / 60) % 24
// }

// func (t *GameTime) CurrentMinute() int {
// 	return t.Minutes % 60
// }

// // Advance processes ticks and increments minutes when enough ticks accumulate
// func (t *GameTime) Advance(ticks int) {
// 	t.TickAccumulator += ticks

// 	// Convert ticks to minutes when accumulator exceeds GameTicksPerMinute
// 	if t.TickAccumulator >= GameTicksPerMinute {
// 		minutesToAdd := t.TickAccumulator / GameTicksPerMinute
// 		t.Minutes = (t.Minutes + minutesToAdd) % GameDayLength
// 		t.TickAccumulator %= GameTicksPerMinute
// 	}
// }

// func (t *GameTime) String() string {
// 	hour := t.CurrentHour()
// 	minute := t.CurrentMinute()
// 	return fmt.Sprintf("%02d:%02d", hour, minute)
// }

// func startTicker() {
// 	slog.Debug("Starting game ticker",
// 		slog.Duration("tick_duration", tickDuration))

// 	ticker := time.NewTicker(tickDuration)
// 	defer ticker.Stop()

// 	slog.Info("Game ticker started", slog.Duration("tick_duration", tickDuration))

// 	for range ticker.C {
// 		handleGameTick()
// 	}
// }

// // func adjustTickDuration(ticker *time.Ticker, newDuration time.Duration) {
// // 	ticker.Reset(newDuration)
// // 	slog.Info("Game tick duration adjusted", slog.Duration("new_duration", newDuration))
// // }

// var gameTime = &GameTime{Minutes: 0}

// func handleGameTick() {
// 	gameTime.Advance(1) // Advance by one tick

// 	if gameTime.TickAccumulator == 0 {
// 		slog.Debug("Game time updated", slog.String("time", gameTime.String()))
// 	}

// 	triggerTimeBasedEvents()
// }

// func triggerTimeBasedEvents() {
// 	hour := gameTime.CurrentHour()

// 	switch hour {
// 	case 6: // Sunrise
// 		for _, c := range CharacterMgr.GetOnlineCharacters() {
// 			c.Send("The sun rises over the horizon, bathing the land in light.")
// 		}
// 	case 18: // Sunset
// 		for _, c := range CharacterMgr.GetOnlineCharacters() {
// 			c.Send("The sun sets, and darkness envelops the world.")
// 		}
// 	}
// }

// // Calculate the number of in-game minutes until the next occurrence of a given hour
// func calculateTimeUntil(targetHour int) int {
// 	currentHour := gameTime.CurrentHour()
// 	currentMinute := gameTime.CurrentMinute()

// 	if currentHour < targetHour {
// 		return (targetHour-currentHour)*60 - currentMinute
// 	} else if currentHour > targetHour {
// 		return (24-currentHour+targetHour)*60 - currentMinute
// 	}
// 	return -currentMinute
// }

// // Converts total minutes into a human-readable format like "1 hour 33 minutes"
// func formatMinutesAsTime(minutes int) string {
// 	hours := minutes / 60
// 	minutes = minutes % 60

// 	if hours > 0 && minutes > 0 {
// 		return fmt.Sprintf("%d hour%s %d minute%s", hours, pluralize(hours), minutes, pluralize(minutes))
// 	} else if hours > 0 {
// 		return fmt.Sprintf("%d hour%s", hours, pluralize(hours))
// 	} else {
// 		return fmt.Sprintf("%d minute%s", minutes, pluralize(minutes))
// 	}
// }

// func pluralize(value int) string {
// 	if value == 1 {
// 		return ""
// 	}
// 	return "s"
// }
