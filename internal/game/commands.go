package game

import (
	"github.com/gliderlabs/ssh"
)

// TODO: We need a RP consistent way to communicate directly with other individuals or groups of individuals I.E. for shadowrun it could be via comlinks and some group or party system

const ()

type (
	CommandFunc func(s ssh.Session, cmd string, args []string, user *Account, char *Character, room *Room)

	SuggestFunc func(line string, args []string, char *Character, room *Room) []string

	Command struct {
		Name          string
		Description   string
		Usage         []string
		Aliases       []string
		RequiredRoles []CharacterRole
		Func          CommandFunc
		SuggestFunc   SuggestFunc // Optional suggestion logic
	}
)

func RegisterCommands() {
	CommandMgr.RegisterCommand(Command{
		Name:        "prompt",
		Description: "Get and set your prompt",
		Usage:       []string{"prompt [prompt]"},
		Func:        DoPrompt,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "history",
		Description: "Show the list of commands executed in this session.",
		Usage:       []string{"history"},
		Func:        DoHistory,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "stats",
		Description: "Display your current attributes and stats.",
		Usage:       []string{"stats"},
		Func:        DoStats,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "time",
		Description: "Display the current in-game time.",
		Usage:       []string{"time", "time details"},
		Func:        DoTime,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "pick",
		Description: "Pick a lock",
		Usage:       []string{"pick [direction]"},
		// Aliases:     []string{"p"},
		Func: DoPick,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "lock",
		Description: "Lock a door",
		Usage:       []string{"lock [direction]"},
		// Aliases:     []string{"l"},
		Func: DoLock,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "unlock",
		Description: "Unlock a door",
		Usage:       []string{"unlock [direction]"},
		// Aliases:     []string{"u"},
		Func: DoUnlock,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "open",
		Description: "Open a door",
		Usage:       []string{"open [direction]"},
		// Aliases:     []string{"o"},
		Func: DoOpen,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "close",
		Description: "Close a door",
		Usage:       []string{"close [direction]"},
		// Aliases:     []string{"c"},
		Func: DoClose,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "who",
		Description: "List players currently in the game",
		Usage:       []string{"who"},
		Aliases:     []string{"w"},
		Func:        DoWho,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "look",
		Description: "Look around the room",
		Usage: []string{
			"look [item|character|mob|direction]",
		},
		Aliases:     []string{"l"},
		Func:        DoLook,
		SuggestFunc: SuggestLook,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "get",
		Description: "Get an item from the room.",
		Usage: []string{
			"get [<quantity>] <item>",
			"get all <item>",
			"get all",
		},
		Func:        DoGet,
		SuggestFunc: SuggestGet,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "give",
		Description: "Give an item",
		Usage:       []string{"give <character> [<quantity>] <item>"},
		Func:        DoGive,
		SuggestFunc: SuggestGive,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "drop",
		Description: "Drop items in the room.",
		Usage: []string{
			"drop [<quantity>] <item>",
			"drop all <item>",
			"drop all",
		},
		Func:        DoDrop,
		SuggestFunc: SuggestDrop,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "help",
		Description: "List available commands",
		Usage: []string{
			"help",
			"help <command>",
		},
		Aliases: []string{"h"},
		Func:    DoHelp,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "move",
		Description: "Move to a different room",
		Usage:       []string{"move [direction]"},
		Aliases:     []string{"m", "n", "s", "e", "w", "u", "d", "north", "south", "east", "west", "up", "down"},
		Func:        DoMove,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "inventory",
		Description: "List your inventory",
		Usage:       []string{"inventory"},
		Aliases:     []string{"i"},
		Func:        DoInventory,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "say",
		Description: "Say something to everyone in the room.",
		Usage:       []string{"say <message>"},
		Func:        DoSay,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "tell",
		Description: "Send a private message to a specific character.",
		Usage:       []string{"tell <username> <message>"},
		Func:        DoTell,
		SuggestFunc: SuggestTell,
	})

	CommandMgr.RegisterCommand(Command{
		Name:        "spawn",
		Description: "Spawn an item or mob into the room",
		Usage: []string{
			"spawn item <item>",
			"spawn mob <mob>",
		},
		RequiredRoles: []CharacterRole{CharacterRoleAdmin},
		Func:          DoSpawn,
	})
}
