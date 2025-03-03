package game

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/google/uuid"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/reflow/wordwrap"
	ee "github.com/vansante/go-event-emitter"
)

const (
	BiasNone Bias = "None"
	BiasGood Bias = "Metahuman"

	RoomTagPeaceful = "Peaceful"
	RoomTagElevator = "Elevator"

	// Corporate & High-Security
	RoomTagCorporate  = "Corporate"
	RoomTagExecutive  = "Executive"
	RoomTagDataCenter = "DataCenter"
	RoomTagSecurity   = "Security"

	// Street-Level & Underground
	RoomTagStreet      = "Street"
	RoomTagSlum        = "Slum"
	RoomTagGangHideout = "GangHideout"
	RoomTagUnderground = "Underground"

	// Commerce & Services
	RoomTagBlackMarket = "BlackMarket"
	RoomTagCyberClinic = "CyberClinic"
	RoomTagBar         = "Bar"
	RoomTagNightclub   = "Nightclub"
	RoomTagArcade      = "Arcade"
	RoomTagHotel       = "Hotel"
	RoomTagRestaurant  = "Restaurant"

	// Safehouses & Hideouts
	RoomTagSafehouse         = "Safehouse"
	RoomTagWarehouse         = "Warehouse"
	RoomTagUndergroundLab    = "UndergroundLab"
	RoomTagAbandonedBuilding = "AbandonedBuilding"

	// Matrix & Digital Spaces
	RoomTagMatrixNode = "MatrixNode"
	RoomTagHost       = "Host"
	RoomTagVRClub     = "VRClub"

	// Combat Zones & Hostile Environments
	RoomTagCombatZone = "CombatZone"
	RoomTagMilitary   = "Military"
	RoomTagDroneBay   = "DroneBay"
	RoomTagToxicZone  = "ToxicZone"

	// Magic-Related
	RoomTagShamanicLodge = "ShamanicLodge"
	RoomTagMagicShop     = "MagicShop"
	RoomTagAstralPlane   = "AstralPlane"
	RoomTagRitualSite    = "RitualSite"

	ExitTypeStairs    = "stairs"
	ExitTypeEscalator = "escalator"
	ExitTypeLadder    = "ladder"
	ExitTypeRope      = "rope"
	ExitTypeRamp      = "ramp"
	ExitTypeSlide     = "slide"
	ExitTypeJumpPad   = "jump_pad"
	ExitTypePassage   = "passage" // Default
)

// TODO: do we want to persist the room state between resets (mobs, items, etc)?

type (
	ExitType string
	// RoomTag string
	Bias string
	Exit struct {
		Room      *Room  `yaml:"-"`
		RoomID    string `yaml:"room_id"`
		Direction string `yaml:"direction"`
		Door      *Door  `yaml:"door"`
		Type      string `yaml:"type,omitempty"`
	}
	Door struct {
		IsClosed       bool     `yaml:"is_closed"`
		IsLocked       bool     `yaml:"is_locked"`
		KeyIDs         []string `yaml:"key_ids"`
		PickDifficulty int      `yaml:"pick_difficulty"`
	}
	Corrdinates struct {
		X int `yaml:"x"`
		Y int `yaml:"y"`
		Z int `yaml:"z"`
	}
	Spawn struct {
		ItemID   string `yaml:"item_id"`
		MobID    string `yaml:"mob_id"`
		Chance   int    `yaml:"chance"`
		Quantity int    `yaml:"quantity"`
	}
	// TODO: Add Doors and Locks
	// TODO: Keep track of items in the room between resets
	// TODO: Keep track of mobs in the room between resets
	// TODO: Check respawn chance of items and mobs on update
	Room struct {
		sync.RWMutex `yaml:"-"`
		Listeners    []ee.Listener `yaml:"-"`

		ID          string           `yaml:"id"`
		ReferenceID string           `yaml:"reference_id"`
		UUID        string           `yaml:"uuid"`
		AreaID      string           `yaml:"area_id"`
		Area        *Area            `yaml:"-"`
		Title       string           `yaml:"title"`
		Description string           `yaml:"description"`
		Tags        []string         `yaml:"tags"`
		Bias        Bias             `yaml:"bias"`
		Exits       map[string]*Exit `yaml:"exits"`
		Corrdinates *Corrdinates     `yaml:"corrdinates"`
		Inventory   Inventory        `yaml:"inventory"`
		Characters  []*Character     `yaml:"-"`
		Mobs        []*Mob           `yaml:"-"`
		Spawns      []Spawn          `yaml:"spawns,omitempty"`
		SpawnedMobs []*Mob           `yaml:"-"` // Mobs that have been spawned into the room
	}
)

func NewRoom() *Room {
	return &Room{
		UUID:  uuid.New().String(),
		Exits: make(map[string]*Exit),
	}
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
	var builder strings.Builder
	if char.Role == CharacterRoleAdmin {
		builder.WriteString(cfmt.Sprintf("{{[}}::white|bold{{%s}}::yellow{{]}}::white|bold", char.Room.ID))
		builder.WriteString(CRLF)
	}

	builder.WriteString(cfmt.Sprintf("{{%s}}::cyan|bold", char.Room.Title))
	if len(char.Room.Tags) > 0 {
		builder.WriteString(cfmt.Sprintf(" {{[}}::white|bold{{%s}}::green{{]}}::white|bold", strings.Join(char.Room.Tags, ", ")))
	}
	builder.WriteString(CRLF + HT)
	builder.WriteString(wordwrap.String(cfmt.Sprint(char.Room.Description), 80) + "" + CRLF)
	builder.WriteString("" + CRLF)
	builder.WriteString(RenderEntitiesInRoom(char) + "" + CRLF)
	builder.WriteString("" + CRLF)

	if len(char.Room.Inventory.Items) > 0 {
		builder.WriteString(RenderItemsInRoom(char) + "" + CRLF)
		builder.WriteString("" + CRLF)
	}

	builder.WriteString(RenderRoomExits(char) + "" + CRLF)
	builder.WriteString("" + CRLF)

	return builder.String()
}

func RenderEntitiesInRoom(char *Character) string {
	var builder strings.Builder

	// Total entity count minus the character itself
	entityCount := len(char.Room.Characters) - 1 + len(char.Room.Mobs)
	entityDescriptions := []string{}
	for _, c := range char.Room.Characters {
		if c.Name != char.Name {
			metatype := EntityMgr.GetMetatype(c.MetatypeID)
			entityDescriptions = append(entityDescriptions, cfmt.Sprintf(
				"{{%s (%s)}}::cyan|bold", c.Name, metatype.Name))
		}
	}
	// Count and map mob names for pluralization
	mobNameCounts := make(map[string]int)
	for _, m := range char.Room.Mobs {
		mobNameCounts[m.Name]++
	}

	figureText := pluralizer.PluralizeNounPhrase("figure", entityCount)
	switch entityCount {
	case 0:
		builder.WriteString(cfmt.Sprintf("You are the only one here."))
	default:
		builder.WriteString(cfmt.Sprintf("You notice {{%s}}::bold: ", figureText))
	}

	for name, count := range mobNameCounts {
		entityDescriptions = append(entityDescriptions, cfmt.Sprintf("{{%s}}::green", pluralizer.PluralizeNounPhrase(name, count)))
	}
	builder.WriteString(strings.Join(entityDescriptions, ", "))

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
		builder.WriteString(cfmt.Sprint("Scattered throughout the room, you spot "))

		// Collect item names with counts
		itemNames := []string{}
		for name, count := range itemNameCounts {
			itemNames = append(itemNames, cfmt.Sprintf("{{%s}}::magenta", pluralizer.PluralizeNounPhrase(name, count)))
		}

		// Join item names with commas (white) and append to the builder
		builder.WriteString(strings.Join(itemNames, cfmt.Sprint(", ")))
	}

	return wordwrap.String(builder.String(), 80)
}

func RenderRoomExits(char *Character) string {
	var builder strings.Builder
	if len(char.Room.Exits) == 0 {
		return cfmt.Sprintf("{{There are no exits}}::red")
	}
	exitStrings := make([]string, 0, len(char.Room.Exits))
	for dir, exit := range char.Room.Exits {
		// Determine exit description based on exit type
		var exitDescription string
		switch exit.Type {
		case ExitTypeEscalator:
			exitDescription = "an escalator"
		case ExitTypeStairs:
			exitDescription = "a set of stairs"
		case ExitTypeLadder:
			exitDescription = "a ladder"
		case ExitTypeRope:
			exitDescription = "a rope"
		case ExitTypeRamp:
			exitDescription = "a ramp"
		case ExitTypeSlide:
			exitDescription = "a slide"
		case ExitTypeJumpPad:
			exitDescription = "a jump pad"
		default:
			exitDescription = "a passage"
		}

		// Adjust for door state if there is one
		if exit.Door != nil {
			if exit.Door.IsClosed {
				exitDescription = fmt.Sprintf("a closed %s", exit.Type)
			} else {
				exitDescription = fmt.Sprintf("an open %s", exit.Type)
			}
		}

		// Special phrasing for vertical movement
		if dir == "up" {
			// For a rope, you might want different phrasing
			if exit.Type == ExitTypeRope {
				exitStrings = append(exitStrings,
					cfmt.Sprintf("There is %s from the floor above.", exitDescription))
			} else {
				exitStrings = append(exitStrings,
					cfmt.Sprintf("There is %s leading {{up}}::yellow to {{%s}}::yellow.", exitDescription, exit.Room.Title))
			}
		} else if dir == "down" {
			exitStrings = append(exitStrings,
				cfmt.Sprintf("There is %s leading {{down}}::yellow to {{%s}}::yellow.", exitDescription, exit.Room.Title))
		} else {
			// Generic phrasing for non-vertical exits
			exitStrings = append(exitStrings,
				cfmt.Sprintf("To the {{%s}}::yellow, there is %s leading to {{%s}}::yellow.", dir, exitDescription, exit.Room.Title))
		}
	}

	builder.WriteString(strings.Join(exitStrings, " "))

	return wordwrap.String(builder.String(), 80)

}
