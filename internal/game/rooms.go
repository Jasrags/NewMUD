package game

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/muesli/reflow/wordwrap"
	ee "github.com/vansante/go-event-emitter"
)

// TODO: do we want to persist the room state between resets (mobs, items, etc)?

type Exit struct {
	Room      *Room  `yaml:"-"`
	RoomID    string `yaml:"room_id"`
	Direction string `yaml:"direction"`
	Door      *Door  `yaml:"door"`
}

type Door struct {
	IsClosed       bool     `yaml:"is_closed"`
	IsLocked       bool     `yaml:"is_locked"`
	KeyIDs         []string `yaml:"key_ids"`
	PickDifficulty int      `yaml:"pick_difficulty"`
}

type Corrdinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type DefaultItem struct {
	ID string `yaml:"id"`
	// RespawnChance    int    `yaml:"respawn_chance"`
	MaxLoad int `yaml:"max_load"`
	// ReplaceOnRespawn bool   `yaml:"replace_on_respawn"`
	Quantity int `yaml:"quantity"`
}

type DefaultMob struct {
	ID string `yaml:"id"`
	// RespawnChance    int    `yaml:"respawn_chance"`
}

// TODO: Add Doors and Locks
// TODO: Keep track of items in the room between resets
// TODO: Keep track of mobs in the room between resets
// TODO: Check respawn chance of items and mobs on update
type Room struct {
	sync.RWMutex `yaml:"-"`
	Listeners    []ee.Listener `yaml:"-"`

	ID           string           `yaml:"id"`
	ReferenceID  string           `yaml:"reference_id"`
	UUID         string           `yaml:"uuid"`
	AreaID       string           `yaml:"area_id"`
	Area         *Area            `yaml:"-"`
	Title        string           `yaml:"title"`
	Description  string           `yaml:"description"`
	Exits        map[string]*Exit `yaml:"exits"`
	Corrdinates  *Corrdinates     `yaml:"corrdinates"`
	Inventory    Inventory        `yaml:"inventory"`
	Characters   []*Character     `yaml:"-"`
	Mobs         []*Mob           `yaml:"-"`
	DefaultItems []DefaultItem    `yaml:"default_items"` // IDs of items to load into the room
	DefaultMobs  []DefaultMob     `yaml:"default_mobs"`  // IDs of mobs to load into the room
	SpawnedMobs  []*Mob           `yaml:"-"`             // Mobs that have been spawned into the room
}

func NewRoom() *Room {
	return &Room{
		UUID:  uuid.New().String(),
		Exits: make(map[string]*Exit),
	}
}

func (r *Room) FindInteractableByName(name string) Interactable {
	for _, c := range r.Characters {
		if strings.EqualFold(c.Name, name) {
			return c
		}
	}
	for _, m := range r.Mobs {
		if strings.EqualFold(m.Name, name) {
			return m
		}
	}
	return nil
}

// func (r *Room) Init() {
// 	slog.Debug("Initializing room",
// 		slog.String("room_id", r.ID))

// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomCharacterEnter, r.onRoomCharacterEnter))
// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomCharacterLeave, r.onRoomCharacterLeave))
// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomMobEnter, r.onRoomMobEnter))
// 	r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomMobLeave, r.onRoomMobLeave))
// 	// r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomSpawn, r.onRoomSpawn))
// 	// r.Listeners = append(r.Listeners, *EventMgr.Subscribe(EventRoomUpdate, r.onRoomUpdate))
// }

// FindMobByName searches for a mob in the room by name and returns the first match or nil if not found
func (r *Room) FindMobByName(name string) *Mob {
	r.RLock()
	defer r.RUnlock()

	for _, mob := range r.Mobs {
		if strings.EqualFold(mob.Name, name) {
			return mob
		}
	}
	return nil
}

func (r *Room) HasExit(dir string) bool {
	return r.Exits[dir] != nil
}

// FindCharacterByName searches for a character in the room by name and returns the first match or nil if not found
func (r *Room) FindCharacterByName(name string) *Character {
	r.RLock()
	defer r.RUnlock()

	for _, char := range r.Characters {
		if strings.EqualFold(char.Name, name) {
			return char
		}
	}
	return nil
}

func (r *Room) GetExits() {
	r.RLock()
	defer r.RUnlock()

	// adjacents := map[string]Corrdinates{
	// 	"north":     {X: r.Corrdinates.X, Y: r.Corrdinates.Y + 1, Z: r.Corrdinates.Z},
	// 	"south":     {X: r.Corrdinates.X, Y: r.Corrdinates.Y - 1, Z: r.Corrdinates.Z},
	// 	"east":      {X: r.Corrdinates.X + 1, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z},
	// 	"west":      {X: r.Corrdinates.X - 1, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z},
	// 	"up":        {X: r.Corrdinates.X, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z + 1},
	// 	"down":      {X: r.Corrdinates.X, Y: r.Corrdinates.Y, Z: r.Corrdinates.Z - 1},
	// 	"northeast": {X: r.Corrdinates.X + 1, Y: r.Corrdinates.Y + 1, Z: r.Corrdinates.Z},
	// 	"northwest": {X: r.Corrdinates.X - 1, Y: r.Corrdinates.Y + 1, Z: r.Corrdinates.Z},
	// 	"southeast": {X: r.Corrdinates.X + 1, Y: r.Corrdinates.Y - 1, Z: r.Corrdinates.Z},
	// 	"southwest": {X: r.Corrdinates.X - 1, Y: r.Corrdinates.Y - 1, Z: r.Corrdinates.Z},
	// }

	// var exits []string
	// for direction := range r.Exits {
	// 	exits = append(exits, direction)
	// }
}

func (r *Room) AddCharacter(c *Character) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Adding character to room",
		slog.String("room_id", r.ID),
		slog.String("character_id", c.ID))

	r.Characters = append(r.Characters, c)
}

func (r *Room) RemoveCharacter(c *Character) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Removing character from room",
		slog.String("room_id", r.ID),
		slog.String("character_id", c.ID))

	for i, char := range r.Characters {
		if char.ID == c.ID {
			r.Characters = append(r.Characters[:i], r.Characters[i+1:]...)
			break
		}
	}
}

func (r *Room) AddMob(m *Mob) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Adding mob to room",
		slog.String("room_id", r.ID),
		slog.String("mob_id", m.ID))

	r.Mobs = append(r.Mobs, m)
}

func (r *Room) RemoveMob(m *Mob) {
	r.Lock()
	defer r.Unlock()

	slog.Debug("Removing mob from room",
		slog.String("room_id", r.ID),
		slog.String("mob_id", m.ID))

	for i, mob := range r.Mobs {
		if mob.ID == m.ID {
			r.Mobs = append(r.Mobs[:i], r.Mobs[i+1:]...)
			break
		}
	}
}

func (r *Room) Broadcast(msg string, excludeIDs []string) {
	excludes := make(map[string]bool)

	for _, id := range excludeIDs {
		excludes[id] = true
	}

	for _, char := range r.Characters {
		if _, ok := excludes[char.ID]; !ok {
			char.Send(msg)
		}
	}
}

// // Event functions
// func (r *Room) onRoomCharacterEnter(arguments ...interface{}) {
// 	slog.Debug("Room character enter event",
// 		slog.String("room_id", r.ID),
// 		slog.Any("args", arguments))

// 	arg := arguments[0].(*RoomCharacterEnter)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A character has entered the room", []string{arg.Character.ID})
// }

// func (r *Room) onRoomCharacterLeave(arguments ...interface{}) {
// 	slog.Debug("Room character leave event",
// 		slog.String("room_id", r.ID))

// 	arg := arguments[0].(*RoomCharacterLeave)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A character has left the room", []string{arg.Character.ID})
// }

// func (r *Room) onRoomMobEnter(arguments ...interface{}) {
// 	slog.Debug("Room mob enter event",
// 		slog.String("room_id", r.ID))

// 	arg := arguments[0].(*RoomMobEnter)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A mob has entered the room", []string{arg.Mob.ID})

// }

// func (r *Room) onRoomMobLeave(arguments ...interface{}) {
// 	slog.Debug("Room mob leave event",
// 		slog.String("room_id", r.ID))

// 	arg := arguments[0].(*RoomMobLeave)

// 	if arg.Room.ID != r.ID {
// 		return
// 	}

// 	r.Broadcast("A mob has left the room", []string{arg.Mob.ID})
// }

// RenderRoom renders the room to a string for the player.
func RenderRoom(user *Account, char *Character, room *Room) string {
	// var builder strings.Builder

	// Create the room title text
	roomTitle := fmt.Sprintf("%-10s", boldCyanText.Render(room.Title))
	if char.Role == CharacterRoleAdmin {
		roomTitle = fmt.Sprintf("%s [%s]", roomTitle, whiteText.Render(room.ID))
	}

	roomText := []string{
		roomTitle,
		whiteText.Render(wordwrap.String(room.Description, 80)),
		"",
		RenderEntitiesInRoom(char),
		"",
	}
	if len(char.Room.Inventory.Items) > 0 {
		roomText = append(roomText, RenderItemsInRoom(char))
		roomText = append(roomText, "")
	}

	// if len(char.Room.Characters) > 0 {
	// roomText = append(roomText, RenderCharactersInRoom(char))
	// }
	// if len(char.Room.Mobs) > 0 {
	// roomText = append(roomText, RenderMobsInRoom(char))
	// }
	roomText = append(roomText, RenderRoomExits(char))
	// if len(char.Room.Inventory.Items) > 0 {
	// roomText = append(roomText, RenderItemsInRoom(char))
	// }
	roomText = append(roomText, "")

	return lipgloss.JoinVertical(lipgloss.Left,
		roomText...,
	)
}

func RenderEntitiesInRoom(char *Character) string {
	var builder strings.Builder

	// Render Characters in the Room
	charCount := len(char.Room.Characters)
	charDescriptions := []string{}
	for _, c := range char.Room.Characters {
		if c.Name != char.Name {
			charDescriptions = append(charDescriptions, fmt.Sprintf(
				"%s, a %s",
				boldCyanText.Render(c.Name), c.Metatype))
		}
	}

	if charCount > 1 && len(charDescriptions) > 0 {
		if len(charDescriptions) == 1 {
			builder.WriteString(whiteText.Render("You notice one figure: "))
			builder.WriteString(charDescriptions[0])
			builder.WriteString(". ")
		} else {
			builder.WriteString(whiteText.Render(fmt.Sprintf("You notice %d figures: ", len(charDescriptions))))
			builder.WriteString(strings.Join(charDescriptions[:len(charDescriptions)-1], ", "))
			builder.WriteString(whiteText.Render(", and "))
			builder.WriteString(charDescriptions[len(charDescriptions)-1])
			builder.WriteString(". ")
		}
	}

	// Render Mobs in the Room
	mobCount := len(char.Room.Mobs)
	mobNameCounts := make(map[string]int)
	for _, m := range char.Room.Mobs {
		mobNameCounts[m.Name]++
	}

	if mobCount > 0 {
		mobDescriptions := []string{}
		for name, count := range mobNameCounts {
			mobDescriptions = append(mobDescriptions, pluralizer.PluralizeNounPhrase(name, count))
		}

		if len(mobDescriptions) == 1 {
			builder.WriteString(whiteText.Render("Nearby, "))
			builder.WriteString(fmt.Sprintf("%s stands, their features twisted with a mix of curiosity and malice.",
				boldMagentaText.Render(mobDescriptions[0])))
		} else {
			builder.WriteString(whiteText.Render("Nearby, "))
			builder.WriteString(strings.Join(mobDescriptions[:len(mobDescriptions)-1], ", "))
			builder.WriteString(whiteText.Render(", and "))
			builder.WriteString(boldMagentaText.Render(mobDescriptions[len(mobDescriptions)-1]))
			builder.WriteString(" stand, their features twisted with a mix of curiosity and malice.")
		}
	}

	return wordwrap.String(builder.String(), 80)
}

func RenderItemsInRoom(char *Character) string {
	var builder strings.Builder

	// Count and map item names
	itemCount := len(char.Room.Inventory.Items)
	itemNameCounts := make(map[string]int)
	for _, i := range char.Room.Inventory.Items {
		bp := EntityMgr.GetItemBlueprintByInstance(i)
		itemNameCounts[bp.Name]++
	}

	// Build the item description
	if itemCount > 0 {
		// Introductory text (white)
		builder.WriteString(whiteText.Render("Scattered throughout the room, you spot "))

		// Collect item names with counts
		itemNames := []string{}
		for name, count := range itemNameCounts {
			itemNames = append(itemNames, boldWhiteText.Render(pluralizer.PluralizeNounPhrase(name, count)))
		}

		// Join item names with commas (white) and append to the builder
		builder.WriteString(strings.Join(itemNames, whiteText.Render(", ")))
	}

	return builder.String()
}

func RenderRoomExits(char *Character) string {
	// Handle no exits case early
	if len(char.Room.Exits) == 0 {
		return redText.Render("There are no exits")
	}

	// Build exits descriptions
	exitStrings := make([]string, 0, len(char.Room.Exits)) // Preallocate capacity
	for dir, exit := range char.Room.Exits {
		// Determine the door description
		doorDescription := "a passage"
		if exit.Door != nil {
			if exit.Door.IsClosed {
				doorDescription = "a closed door"
			} else {
				doorDescription = "an open doorway"
			}
		}

		// Format the exit description
		exitStrings = append(exitStrings, fmt.Sprintf(
			"%s %s, %s %s %s %s.",
			whiteText.Render("To the"),
			boldWhiteText.Underline(true).Render(dir),
			whiteText.Render("you see"),
			whiteText.Render(doorDescription),
			whiteText.Render("leading to"),
			boldWhiteText.Italic(true).Render(exit.Room.Title),
		))
	}

	// Join and word-wrap exit descriptions
	return wordwrap.String(strings.Join(exitStrings, " "), 80)
}
