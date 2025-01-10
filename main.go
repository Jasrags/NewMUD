package main

import (
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
)

func main() {
	// gs := mud.NewGameServer()
	// gs.Start()

	setupConfig()
	loadAllDataFiles()
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

	var user *User
	var char *Character
	// var room *Room
	var state = StateWelcome

	for {
		switch state {
		case StateWelcome:
			state = promptWelcome(s)
		case StateLogin:
			state, user = promptLogin(s)
		case StateRegistration:
			state, user = promptRegistration(s)
		case StateMainMenu:
			state = promptMainMenu(s, user)
		case StateChangePassword:
			state = promptChangePassword(s, user)
		// case StateCharacterSelect:
		// state, char = promptCharacterSelect(s, user)
		case StateEnterGame:
			state, char = promptEnterGame(s, user)
		case StateGameLoop:
			state = promptGameLoop(s, user, char)
		case StateExitGame:
			state = promptExitGame(s, user, char)
		case StateQuit:
			fallthrough
		case StateError:
			s.Close()
			char = nil
			user = nil
			return
		default:
			slog.Error("Invalid state", slog.String("user_state", state))
			s.Close()
			char = nil
			user = nil
		}
	}
}

// func handleConnection(netConn *connections.NetConnection, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	slog.Debug("New client connected",
// 		slog.String("id", netConn.ID),
// 		slog.String("remote_address", netConn.Conn.RemoteAddr().String()))

// 	user := users.NewUser()
// 	user.Conn = netConn

// 	// // Initialize state machine
// 	// sm := NewStateMachine(user)
// 	// sm.RegisterState(StateWelcome, func(input string) { handleWelcome(sm) })
// 	// sm.RegisterState(StateLogin, func(input string) { handleLogin(sm) })
// 	// sm.RegisterState(StateMainMenu, func(input string) { handleMainMenu(sm) })

// 	// sm.TransitionTo(StateWelcome)

// 	// // Main input loop
// 	// reader := bufio.NewReader(netConn.Conn)
// 	// for {
// 	// 	// Read player input
// 	// 	input, err := reader.ReadString('\n')
// 	// 	if err != nil {
// 	// 		if errors.Is(err, io.EOF) {
// 	// 			slog.Info("Client disconnected", slog.String("id", user.ID))
// 	// 		} else {
// 	// 			slog.Error("Error reading input", slog.String("id", user.ID), slog.Any("error", err))
// 	// 		}
// 	// 		return
// 	// 	}

// 	// 	// Trim input and route to state machine
// 	// 	input = strings.TrimSpace(input)
// 	// 	sm.HandleInput(input)
// 	// }

// 	displayBanner(netConn.Conn)
// 	promptLogin(netConn.Conn)
// 	gameLoop(netConn.Conn)

// }

func loadAllDataFiles() {
	slog.Info("Loading data files")

	EntityMgr.LoadDataFiles()
	UserMgr.LoadDataFiles()
	CharacterMgr.LoadDataFiles()
	RegisterCommands()
}

func RegisterCommands() {
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
		Aliases: []string{"l"},
		Func:    DoLook,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "get",
		Description: "Get an item",
		Usage: []string{
			"get all",
			"get <item>",
			"get <number> <items>",
			"get all <items>",
		},
		Aliases: []string{"g"},
		Func:    DoGet,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "give",
		Description: "Give an item",
		Usage: []string{
			"give <item> [to] <character>",
			"give 2 <items> [to] <character>",
			"give all [to] <character>",
		},
		Aliases: []string{"gi"},
		Func:    DoGive,
	})
	CommandMgr.RegisterCommand(Command{
		Name:        "drop",
		Description: "Drop an item",
		Usage: []string{
			"drop all",
			"drop <item>",
			"drop <number> <items>",
			"drop all <items>",
		},
		Aliases: []string{"d"},
		Func:    DoDrop,
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
		Description: "Say something to the room or to a character or mob",
		Usage: []string{
			"say <message>",
			"say @<name> <message>",
		},
		Func: DoSay,
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
	// }
}

// func displayBanner(conn net.Conn) {
// 	slog.Debug("Displaying banner")

// 	banner := `
// {{Welcome to the MUD server!}}::green
// {{==========================}}::white|bold
// {{Press enter to continue...}}::green

// `

// 	io.WriteString(conn, cfmt.Sprint(banner))

// 	bufio.NewReader(conn).ReadString('\n')
// }

// func promptLogin(conn net.Conn) {
// 	slog.Debug("Prompting for login")

// 	// Prompt for username
// 	io.WriteString(conn, "Please enter your username: ")
// 	reader := bufio.NewReader(conn)
// 	username, err := reader.ReadString('\n')
// 	if err != nil {
// 		slog.Error("Error reading username",
// 			slog.Any("error", err))
// 		return
// 	}

// 	username = strings.TrimSpace(username)
// 	slog.Debug("Received username",
// 		slog.String("username", username))

// 	// Prompt for password
// 	io.WriteString(conn, "Please enter your password: ")
// 	// io.WriteString(conn, "\xff\xfb\x01") // IAC WILL ECHO

// 	password, err := reader.ReadString('\n')
// 	if err != nil {
// 		slog.Error("Error reading password",
// 			slog.Any("error", err))
// 		return
// 	}

// 	// io.WriteString(conn, "\xff\xfc\x01") // IAC WONT ECHO
// 	io.WriteString(conn, "\n")

// 	password = strings.TrimSpace(password)
// 	slog.Debug("Received password",
// 		slog.String("password", password))

// 	// Check if user exists
// 	u := users.Mgr.GetByUsername(username)
// 	if u == nil {
// 		// TODO: User does not exist, we need to go to the registration process
// 		io.WriteString(conn, "User does not exist\n")
// 		conn.Close()
// 		return
// 	}

// 	// Check if password matches the user's password
// 	if ok := u.CheckPassword(password); !ok {
// 		io.WriteString(conn, "Invalid username or password \n")
// 		conn.Close()
// 	}

// 	t := time.Now()
// 	u.CreatedAt = t
// 	u.LastLoginAt = &t

// 	u.Save()

// 	// TODO: Check if user is already logged in

// 	// TODO: Check if user is banned

// 	users.Mgr.SetOnline(u)
// }

// func gameLoop(conn net.Conn) {
// 	slog.Debug("Game Loop")

// 	reader := bufio.NewReader(conn)
// 	for {
// 		io.WriteString(conn, "> ")
// 		input, err := reader.ReadString('\n')
// 		if err != nil {
// 			slog.Error("Error reading input",
// 				slog.Any("error", err))

// 			return
// 		}

// 		input = strings.TrimSpace(input)
// 		slog.Debug("Received text",
// 			slog.String("input", input))
// 	}

// 	// // Create a new player for this connection
// 	// player := NewPlayer(gs.Log, conn)
// 	// if i == 0 {
// 	// 	player.Role = "admin"
// 	// 	player.Name = "Admin"
// 	// } else {
// 	// 	player.Role = "player"
// 	// 	player.Name = "Player"
// 	// }
// 	// i++

// 	// player.MoveTo(gs.GameContext.RoomManager.GetRoom("limbo:the_void"))
// 	// io.WriteString(player.Conn, RenderRoom(player, player.Room))

// 	// gs.GameLoop(conn, player)
// }
