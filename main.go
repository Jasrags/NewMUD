package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
)

// Add this near the top of your main.go file
var tickDuration time.Duration

func main() {
	setupConfig()
	loadAllDataFiles()
	go startTicker()
	setupServer()

	slog.Info("Shutting down")
}

func setupConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	setupLogger()

	// Parse tick duration
	var err error
	tickDuration, err = time.ParseDuration(viper.GetString("server.tick_duration"))
	if err != nil {
		panic(fmt.Sprintf("Invalid tick_duration format: %s", viper.GetString("server.tick_duration")))
	}

	// Update configuration on change
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Config updated",
			slog.String("file", e.Name))
		outputConfig()
	})
	viper.WatchConfig()

	outputConfig()
}

func outputConfig() {
	slog.Info("Current configuration")

	for _, key := range viper.AllKeys() {
		slog.Debug("config",
			slog.String("key", key),
			slog.String("value", viper.GetString(key)))
	}
}

func setupLogger() {
	logLevel := viper.GetString("server.log_level")
	logHandler := viper.GetString("server.log_handler")

	slog.Info("Setting up logger",
		slog.String("log_level", logLevel),
		slog.String("log_handler", logHandler))

	// Parse and set log level
	var programLevel slog.Level
	switch logLevel {
	case "debug":
		programLevel = slog.LevelDebug
	case "info":
		programLevel = slog.LevelInfo
	case "warn":
		programLevel = slog.LevelWarn
	case "error":
		programLevel = slog.LevelError
	default:
		panic("Invalid log level")
	}

	// Setup log handler
	var logger *slog.Logger
	handlerOptions := &slog.HandlerOptions{Level: programLevel}

	switch logHandler {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	default:
		panic("Invalid log handler")
	}

	// Set default logger
	slog.SetDefault(logger)
}

// func setWinsize(f *os.File, w, h int) {
// 	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
// 		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
// }

func handleConnection(s ssh.Session) {
	defer s.Close()

	_, winCh, _ := s.Pty()

	// Set the window size
	go func() {
		for win := range winCh {
			slog.Debug("Window size changed",
				slog.Int("width", win.Width),
				slog.Int("height", win.Height),
			)
		}
	}()

	var account *Account
	var char *Character
	// var room *Room
	var state = StateWelcome

	for {
		switch state {
		case StateWelcome:
			state = promptWelcome(s)
		case StateLogin:
			state, account = promptLogin(s)
		case StateRegistration:
			state, account = promptRegistration(s)
		case StateMainMenu:
			state = promptMainMenu(s, account)
		case StateChangePassword:
			state = promptChangePassword(s, account)
			// case StateCharacterSelect:
			// state, char = promptCharacterSelect(s, user)
		case StateCharacterCreate:
			state = promptCharacterCreate(s, account)
		case StateEnterGame:
			state, char = promptEnterGame(s, account)
		case StateGameLoop:
			state = promptGameLoop(s, account, char)
		case StateExitGame:
			state = promptExitGame(s, account, char)
		case StateQuit:
			fallthrough
		case StateError:
			s.Close()
			char = nil
			account = nil
			return
		default:
			slog.Error("Invalid state", slog.String("user_state", state))
			s.Close()
			char = nil
			account = nil
		}
	}
}

func loadAllDataFiles() {
	slog.Info("Loading data files")

	EntityMgr.LoadDataFiles()
	AccountMgr.LoadDataFiles()
	CharacterMgr.LoadDataFiles()
	RegisterCommands()
}

func RegisterCommands() {
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
