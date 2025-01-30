package game

import (
	"fmt"
	"log/slog"
	"time"
)

// TODO: Support short versions of game time
// TODO: Support saving/loading game time

// Game Constants
const (
	GameDayLength      = 1440 // Minutes in a game day (24 hours * 60 minutes)
	GameTicksPerMinute = 10   // Number of ticks for one in-game minute
	DaysInWeek         = 7
)

// Gregorian calendar month lengths
var (
	GameTimeMgr  = NewGameTime()
	MonthLengths = []int{
		31, // January
		28, // February (default, leap year handled separately)
		31, // March
		30, // April
		31, // May
		30, // June
		31, // July
		31, // August
		30, // September
		31, // October
		30, // November
		31, // December
	}
)

// GameTime struct with date tracking
type GameTime struct {
	Minutes         int // Total minutes since the start of the game day
	TickAccumulator int // Tracks the number of ticks since the last minute increment
	Day             int // Current day in the month
	Month           int // Current month (1-12)
	Year            int // Current year
}

// Initializes the game time
func NewGameTime() *GameTime {
	return &GameTime{
		Minutes:         0,
		Day:             1,    // Start at January 1st
		Month:           1,    // January
		Year:            1000, // Default game start year
		TickAccumulator: 0,
	}
}

// Returns the current hour in 24-hour format
func (t *GameTime) CurrentHour() int {
	return (t.Minutes / 60) % 24
}

// Returns the current minute
func (t *GameTime) CurrentMinute() int {
	return t.Minutes % 60
}

// Returns whether the current year is a leap year
func (t *GameTime) IsLeapYear() bool {
	year := t.Year
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// Gets the number of days in the current month
func (t *GameTime) DaysInMonth() int {
	if t.Month == 2 && t.IsLeapYear() {
		return 29
	}
	return MonthLengths[t.Month-1]
}

// Advances game time based on ticks and increments the date accordingly
func (t *GameTime) Advance(ticks int) {
	t.TickAccumulator += ticks

	// Convert ticks to minutes
	if t.TickAccumulator >= GameTicksPerMinute {
		minutesToAdd := t.TickAccumulator / GameTicksPerMinute
		t.Minutes += minutesToAdd
		t.TickAccumulator %= GameTicksPerMinute
	}

	// Handle day rollover
	for t.Minutes >= GameDayLength {
		t.Minutes -= GameDayLength
		t.Day++

		if t.Day > t.DaysInMonth() { // If we exceed month length
			t.Day = 1
			t.Month++
			if t.Month > 12 { // If we exceed December, go to next year
				t.Month = 1
				t.Year++
			}
		}
	}
}

// Returns a formatted string of the current in-game time (HH:MM AM/PM)
func (t *GameTime) GetFormattedTime() string {
	return fmt.Sprintf("{{%02d}}::yellow:{{%02d}}::yellow %s",
		t.CurrentHour()%12,
		t.CurrentMinute(),
		t.GetAmPm(),
	)
}

// Returns the in-game date formatted as "Month Day, Year"
func (t *GameTime) GetFormattedDate(short bool) string {
	months := []string{"January", "February", "March", "April", "May", "June", "July",
		"August", "September", "October", "November", "December"}
	if short {
		months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul",
			"Aug", "Sep", "Oct", "Nov", "Dec"}
	}

	return fmt.Sprintf("{{%s %d, %d}}::cyan", months[t.Month-1], t.Day, t.Year)
}

// Returns whether the current time is AM or PM
func (t *GameTime) GetAmPm() string {
	if t.CurrentHour() >= 12 {
		return "PM"
	}
	return "AM"
}

func (t *GameTime) StartTicker(tickDuration time.Duration) {
	slog.Debug("Starting game ticker",
		slog.Duration("tick_duration", tickDuration))

	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	slog.Info("Game ticker started", slog.Duration("tick_duration", tickDuration))

	for range ticker.C {
		handleGameTick()
	}
}

func handleGameTick() {
	GameTimeMgr.Advance(1) // Advance by one tick

	if GameTimeMgr.TickAccumulator == 0 {
		slog.Debug("Game time updated",
			slog.String("time", GameTimeMgr.GetFormattedTime()),
			slog.String("date", GameTimeMgr.GetFormattedDate(false)))
	}

	triggerTimeBasedEvents()
}

func triggerTimeBasedEvents() {
	hour := GameTimeMgr.CurrentHour()

	switch hour {
	case 6: // Sunrise
		for _, c := range CharacterMgr.GetOnlineCharacters() {
			c.Send("The sun rises over the horizon, bathing the land in light.")
		}
	case 18: // Sunset
		for _, c := range CharacterMgr.GetOnlineCharacters() {
			c.Send("The sun sets, and darkness envelops the world.")
		}
	}
}

// Calculate the number of in-game minutes until the next occurrence of a given hour
func calculateTimeUntil(targetHour int) int {
	currentHour := GameTimeMgr.CurrentHour()
	currentMinute := GameTimeMgr.CurrentMinute()

	if currentHour < targetHour {
		return (targetHour-currentHour)*60 - currentMinute
	} else if currentHour > targetHour {
		return (24-currentHour+targetHour)*60 - currentMinute
	}
	return -currentMinute
}

// Converts total minutes into a human-readable format like "1 hour 33 minutes"
func formatMinutesAsTime(minutes int) string {
	hours := minutes / 60
	minutes = minutes % 60

	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%d hour%s %d minute%s", hours, pluralize(hours), minutes, pluralize(minutes))
	} else if hours > 0 {
		return fmt.Sprintf("%d hour%s", hours, pluralize(hours))
	} else {
		return fmt.Sprintf("%d minute%s", minutes, pluralize(minutes))
	}
}
