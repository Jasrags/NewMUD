package main

import (
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Jasrags/NewMUD/characters"
	"github.com/Jasrags/NewMUD/commands"
	"github.com/Jasrags/NewMUD/connections"
	"github.com/Jasrags/NewMUD/items"
	"github.com/Jasrags/NewMUD/mobs"
	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"

	"github.com/fsnotify/fsnotify"
	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
	"golang.org/x/term"
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

func setupServer() {
	address := net.JoinHostPort(viper.GetString("server.host"), viper.GetString("server.port"))

	slog.Info("Starting server",
		slog.String("address", address))

	server := &ssh.Server{
		Addr: address,
		// IdleTimeout:              viper.GetDuration("server.idle_timeout"), // TODO: reenable timeout later when we fix the connection close issues
		Handler: handleConnection,
		// ConnCallback:             ConnCallback,
		// ConnectionFailedCallback: ConnectionFailedCallback,
	}
	defer server.Close()

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server",
			slog.String("address", address),
			slog.Any("error", err))
		return
	}

	// slog.Info("Server started", slog.String("address", address))

	// return server
}

// func ConnCallback(ctx ssh.Context, conn net.Conn) net.Conn {
// 	slog.Info("New connection",
// 		slog.String("remote_address", conn.RemoteAddr().String()))

// 	return conn
// }

// func ConnectionFailedCallback(conn net.Conn, err error) {
// 	slog.Error("Connection failed",
// 		slog.Any("error", err))
// 	conn.Close()
// }

const (
	// This will skip straight to the game loop
	StateDebug           = "debug"
	StateWelcome         = "welcome"
	StateLogin           = "login"
	StateRegistration    = "registration"
	StateMainMenu        = "main_menu"
	StateCharacterSelect = "character_select"
	StateEnterGame       = "enter_game"
	StateExitGame        = "exit_game"
)

// func setWinsize(f *os.File, w, h int) {
// 	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
// 		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
// }

func PromptForInput(s ssh.Session, prompt string) (string, error) {
	t := term.NewTerminal(s, prompt)
	input, err := t.ReadLine()
	if err != nil {
		slog.Error("Error reading input", slog.Any("error", err))
		s.Close()

		return "", err
	}

	return strings.TrimSpace(input), nil
}

func PromptForPassword(s ssh.Session, prompt string) string {
	t := term.NewTerminal(s, prompt)
	input, err := t.ReadPassword(prompt)
	if err != nil {
		slog.Error("Error reading password", slog.Any("error", err))
		s.Close()

		return ""
	}

	return strings.TrimSpace(input)
}

func handleConnection(s ssh.Session) {
	defer s.Close()
	// _, winCh, _ := s.Pty()

	// // Set the window size
	// go func() {
	// 	for win := range winCh {
	// 		slog.Debug("Window size changed",
	// 			slog.Int("width", win.Width),
	// 			slog.Int("height", win.Height),
	// 		)
	// 	}
	// }()

	var user *users.User
	// user := &users.User{
	// 	NetConn: connections.NewNetConnection(s),
	// 	// State:   StateWelcome,
	// }

	// connections.Mgr.Add(user.NetConn)

	state := viper.GetString("server.initial_state")
	for {
		switch state {
		// Skip straight to the game loop
		case StateDebug:
			slog.Debug("Debug state")
			user = users.Mgr.GetByUsername("admin")
			user.NetConn = connections.NewNetConnection(s)
			if user.RoomID == "" {
				user.RoomID = viper.GetString("server.starting_room")
			}

			user.Room = rooms.Mgr.GetRoom(user.RoomID)

			if user.Room == nil {
				slog.Error("Starting room not found", slog.String("room_id", user.RoomID))
				return
			}

			t := time.Now()
			user.LastLoginAt = &t

			state = StateEnterGame
		case StateWelcome:
			slog.Debug("Welcome state")
			banner := `
{{Welcome to the MUD server!}}::green
{{==========================}}::white|bold
`

			io.WriteString(s, cfmt.Sprint(banner))
			input, _ := PromptForInput(s, cfmt.Sprint("Type {{login}}::green or {{register}}::green to continue: "))

			switch input {
			case "login":
				state = StateLogin
			case "register":
				state = StateRegistration
			default:
				io.WriteString(s, cfmt.Sprint("{{Invalid option}}::red\n"))
				state = StateWelcome
			}
		case StateLogin:
			slog.Debug("Login state")
			// Collect username
			username, _ := PromptForInput(s, cfmt.Sprint("{{Enter your username:}}::green "))

			// Collect password
			password := PromptForPassword(s, cfmt.Sprint("{{Enter your password:}}::green "))

			slog.Debug("Received username and password",
				slog.String("username", username),
				slog.String("password", password))

			// Check if user exists
			user = users.Mgr.GetByUsername(username)

			// If user does not exist, we need to go to the registration process
			if user == nil {
				io.WriteString(s, "User does not exist\n")
				state = StateRegistration
				continue
			}

			// Validate password against user's hashed password
			if !user.CheckPassword(password) {
				io.WriteString(s, cfmt.Sprint("{{Invalid username or password}}::red\n"))

				slog.Debug("Invalid username or password")
				state = StateLogin
				return
			}

			// TODO: Check if user is already logged in

			// TODO: Check if user is banned

			t := time.Now()
			user.LastLoginAt = &t

			// state = StateMainMenu
			state = StateMainMenu
			continue
		case StateRegistration:
			slog.Debug("Registration state")
			io.WriteString(s, cfmt.Sprint("{{Registration}}::green\n"))
			PromptForInput(s, cfmt.Sprint("{{Press enter to continue...}}::green"))
			state = StateMainMenu
		case StateMainMenu:
			slog.Debug("Main menu state")
			io.WriteString(s, cfmt.Sprint("{{Main Menu}}::green\n"))
			PromptForInput(s, cfmt.Sprint("{{Press enter to continue...}}::green"))
			state = StateCharacterSelect
		case StateCharacterSelect:
			slog.Debug("Character select state")
			io.WriteString(s, cfmt.Sprint("{{Character Select}}::green\n"))
			PromptForInput(s, cfmt.Sprint("{{Press enter to continue...}}::green"))
			state = StateEnterGame
		case StateEnterGame:
			slog.Debug("Enter game state")
			input, err := PromptForInput(s, cfmt.Sprint("{{>}}::white|bold "))
			if err != nil {
				slog.Error("Error reading input", slog.Any("error", err))
				state = StateExitGame
				break
			}

			if input == "" {
				continue
			}
			commands.Mgr.ParseAndExecute(s, input, nil, nil)
		case StateExitGame:
			slog.Debug("Exit game state")
			// user.NetConn.Close()
			s.Close()
			return
		default:
			slog.Error("Invalid state", slog.String("user_state", user.State))
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
	rooms.Mgr.LoadDataFiles()
	items.Mgr.LoadDataFiles()
	mobs.Mgr.LoadDataFiles()
	users.Mgr.LoadDataFiles()
	characters.Mgr.LoadDataFiles()
	commands.Mgr.RegisterCommands()
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
