package game

import (
	"github.com/gliderlabs/ssh"
)

// TODO: We need a RP consistent way to communicate directly with other individuals or groups of individuals I.E. for shadowrun it could be via comlinks and some group or party system

const (
	CommandCategoryAdministration CommandCategory = "Administration"
	CommandCategoryCommunication  CommandCategory = "Communication"
	CommandCategoryInformative    CommandCategory = "Informative"
	CommandCategoryMovement       CommandCategory = "Movement"
	CommandCategoryInteraction    CommandCategory = "Interaction"
)

type (
	CommandCategory string
	CommandFunc     func(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room)

	SuggestFunc func(line string, args []string, char *Character, room *Room) []string

	Command struct {
		Name            string
		Description     string
		CommandCategory CommandCategory
		Usage           []string
		Aliases         []string
		RequiredRoles   []CharacterRole
		Func            CommandFunc
		SuggestFunc     SuggestFunc // Optional suggestion logic
	}
)

func RegisterCommands() {
	CommandMgr.RegisterCommand(Command{
		Name:            "equipment",
		Aliases:         []string{"eq"},
		Description:     "List your equipment",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"equipment"},
		Func:            DoEquipment,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "equip",
		Description:     "Equip an item",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"equip <item>"},
		Func:            DoEquip,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "unequip",
		Description:     "Unequip an item",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"unequip <item>"},
		Func:            DoUnequip,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "list",
		Description:     "List game entities",
		CommandCategory: CommandCategoryAdministration,
		Usage:           []string{"list <mobs|m> [tags]", "list <rooms|r> [tags]", "list <items|i> [tags]"},
		RequiredRoles:   []CharacterRole{CharacterRoleAdmin},
		Func:            DoList,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "goto",
		Description:     "Teleport to a room or character",
		CommandCategory: CommandCategoryAdministration,
		Usage:           []string{"goto <room_id>", "goto <character_name>"},
		RequiredRoles:   []CharacterRole{CharacterRoleAdmin},
		Func:            DoGoto,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "mobstats",
		Description:     "Display the stats of a mob",
		CommandCategory: CommandCategoryAdministration,
		Usage:           []string{"mobstats <mob>"},
		Func:            DoMobStats,
		RequiredRoles:   []CharacterRole{CharacterRoleAdmin},
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "prompt",
		Description:     "Get and set your prompt",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"prompt [prompt]"},
		Func:            DoPrompt,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "history",
		Description:     "Show the list of commands executed in this session.",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"history"},
		Func:            DoHistory,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "stats",
		Description:     "Display your current attributes and stats.",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"stats"},
		Func:            DoStats,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "time",
		Description:     "Display the current in-game time.",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"time", "time details"},
		Func:            DoTime,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "pick",
		Description:     "Pick a lock",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"pick [direction]"},
		Func:            DoPick,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "lock",
		Description:     "Lock a door",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"lock [direction]"},
		Func:            DoLock,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "unlock",
		Description:     "Unlock a door",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"unlock [direction]"},
		Func:            DoUnlock,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "open",
		Description:     "Open a door",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"open [direction]"},
		Func:            DoOpen,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "close",
		Description:     "Close a door",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"close [direction]"},
		Func:            DoClose,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "who",
		Description:     "List players currently in the game",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"who"},
		Aliases:         []string{"w"},
		Func:            DoWho,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "look",
		Description:     "Look around the room",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"look [item|character|mob|direction]"},
		Aliases:         []string{"l"},
		Func:            DoLook,
		SuggestFunc:     SuggestLook,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "get",
		Description:     "Get an item from the room.",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"get [<quantity>] <item>", "get all <item>", "get all"},
		Func:            DoGet,
		SuggestFunc:     SuggestGet,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "give",
		Description:     "Give an item",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"give <character> [<quantity>] <item>"},
		Func:            DoGive,
		SuggestFunc:     SuggestGive,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "drop",
		Description:     "Drop items in the room.",
		CommandCategory: CommandCategoryInteraction,
		Usage:           []string{"drop [<quantity>] <item>", "drop all <item>", "drop all"},
		Func:            DoDrop,
		SuggestFunc:     SuggestDrop,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "help",
		Description:     "List available commands",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"help", "help <command>"},
		Aliases:         []string{"h"},
		Func:            DoHelp,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "move",
		Description:     "Move to a different room",
		CommandCategory: CommandCategoryMovement,
		Usage:           []string{"move [direction]"},
		Aliases:         []string{"m", "n", "s", "e", "w", "u", "d", "north", "south", "east", "west", "up", "down"},
		Func:            DoMove,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "inventory",
		Description:     "List your inventory",
		CommandCategory: CommandCategoryInformative,
		Usage:           []string{"inventory"},
		Aliases:         []string{"i"},
		Func:            DoInventory,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "say",
		Description:     "Say something to everyone in the room.",
		CommandCategory: CommandCategoryCommunication,
		Usage:           []string{"say <message>"},
		Func:            DoSay,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "tell",
		Description:     "Send a private message to a specific character.",
		CommandCategory: CommandCategoryCommunication,
		Usage:           []string{"tell <username> <message>"},
		Func:            DoTell,
		SuggestFunc:     SuggestTell,
	})
	CommandMgr.RegisterCommand(Command{
		Name:            "spawn",
		Description:     "Spawn an item or mob into the room",
		CommandCategory: CommandCategoryAdministration,
		Usage:           []string{"spawn item <item>", "spawn mob <mob>"},
		RequiredRoles:   []CharacterRole{CharacterRoleAdmin},
		Func:            DoSpawn,
	})
}
