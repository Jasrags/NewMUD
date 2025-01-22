package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Jasrags/NewMUD/internal/game"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Add this near the top of your main.go file
var tickDuration time.Duration

func main() {
	setupConfig()
	loadAllDataFiles()
	go game.StartTicker(tickDuration)
	game.SetupServer()

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

// func handleConnection(s ssh.Session) {
// 	defer s.Close()

// 	_, winCh, _ := s.Pty()

// 	// Set the window size
// 	go func() {
// 		for win := range winCh {
// 			slog.Debug("Window size changed",
// 				slog.Int("width", win.Width),
// 				slog.Int("height", win.Height),
// 			)
// 		}
// 	}()

// 	var account *game.Account
// 	var char *game.Character
// 	// var room *Room
// 	var state = game.StateWelcome

// 	for {
// 		switch state {
// 		case game.StateWelcome:
// 			state = game.PromptWelcome(s)
// 		case game.StateLogin:
// 			state, account = game.PromptLogin(s)
// 		case game.StateRegistration:
// 			state, account = game.PromptRegistration(s)
// 		case game.StateMainMenu:
// 			state = game.PromptMainMenu(s, account)
// 		case game.StateChangePassword:
// 			state = game.PromptChangePassword(s, account)
// 			// case StateCharacterSelect:
// 			// state, char = game.PromptCharacterSelect(s, user)
// 		case game.StateCharacterCreate:
// 			state = game.PromptCharacterCreate(s, account)
// 		case game.StateEnterGame:
// 			state, char = game.PromptEnterGame(s, account)
// 		case game.StateGameLoop:
// 			state = game.PromptGameLoop(s, account, char)
// 		case game.StateExitGame:
// 			state = game.PromptExitGame(s, account, char)
// 		case game.StateQuit:
// 			fallthrough
// 		case game.StateError:
// 			s.Close()
// 			char = nil
// 			account = nil
// 			return
// 		default:
// 			slog.Error("Invalid state", slog.String("user_state", state))
// 			s.Close()
// 			char = nil
// 			account = nil
// 		}
// 	}
// }

func loadAllDataFiles() {
	slog.Info("Loading data files")

	game.EntityMgr.LoadDataFiles()
	game.AccountMgr.LoadDataFiles()
	game.CharacterMgr.LoadDataFiles()
	RegisterCommands()
}

func RegisterCommands() {
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "history",
		Description: "Show the list of commands executed in this session.",
		Usage:       []string{"history"},
		Func:        game.DoHistory,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "stats",
		Description: "Display your current attributes and stats.",
		Usage:       []string{"stats"},
		Func:        game.DoStats,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "time",
		Description: "Display the current in-game time.",
		Usage:       []string{"time", "time details"},
		Func:        game.DoTime,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "pick",
		Description: "Pick a lock",
		Usage:       []string{"pick [direction]"},
		// Aliases:     []string{"p"},
		Func: game.DoPick,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "lock",
		Description: "Lock a door",
		Usage:       []string{"lock [direction]"},
		// Aliases:     []string{"l"},
		Func: game.DoLock,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "unlock",
		Description: "Unlock a door",
		Usage:       []string{"unlock [direction]"},
		// Aliases:     []string{"u"},
		Func: game.DoUnlock,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "open",
		Description: "Open a door",
		Usage:       []string{"open [direction]"},
		// Aliases:     []string{"o"},
		Func: game.DoOpen,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "close",
		Description: "Close a door",
		Usage:       []string{"close [direction]"},
		// Aliases:     []string{"c"},
		Func: game.DoClose,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "who",
		Description: "List players currently in the game",
		Usage:       []string{"who"},
		Aliases:     []string{"w"},
		Func:        game.DoWho,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "look",
		Description: "Look around the room",
		Usage: []string{
			"look [item|character|mob|direction]",
		},
		Aliases:     []string{"l"},
		Func:        game.DoLook,
		SuggestFunc: game.SuggestLook,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "get",
		Description: "Get an item from the room.",
		Usage: []string{
			"get [<quantity>] <item>",
			"get all <item>",
			"get all",
		},
		Func:        game.DoGet,
		SuggestFunc: game.SuggestGet,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "give",
		Description: "Give an item",
		Usage:       []string{"give <character> [<quantity>] <item>"},
		Func:        game.DoGive,
		SuggestFunc: game.SuggestGive,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "drop",
		Description: "Drop items in the room.",
		Usage: []string{
			"drop [<quantity>] <item>",
			"drop all <item>",
			"drop all",
		},
		Func:        game.DoDrop,
		SuggestFunc: game.SuggestDrop,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "help",
		Description: "List available commands",
		Usage: []string{
			"help",
			"help <command>",
		},
		Aliases: []string{"h"},
		Func:    game.DoHelp,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "move",
		Description: "Move to a different room",
		Usage:       []string{"move [direction]"},
		Aliases:     []string{"m", "n", "s", "e", "w", "u", "d", "north", "south", "east", "west", "up", "down"},
		Func:        game.DoMove,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "inventory",
		Description: "List your inventory",
		Usage:       []string{"inventory"},
		Aliases:     []string{"i"},
		Func:        game.DoInventory,
	})
	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "say",
		Description: "Say something to everyone in the room.",
		Usage:       []string{"say <message>"},
		Func:        game.DoSay,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "tell",
		Description: "Send a private message to a specific character.",
		Usage:       []string{"tell <username> <message>"},
		Func:        game.DoTell,
		SuggestFunc: game.SuggestTell,
	})

	game.CommandMgr.RegisterCommand(game.Command{
		Name:        "spawn",
		Description: "Spawn an item or mob into the room",
		Usage: []string{
			"spawn item <item>",
			"spawn mob <mob>",
		},
		RequiredRoles: []game.CharacterRole{game.CharacterRoleAdmin},
		Func:          game.DoSpawn,
	})
}
