package mud

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// var eventBus = EventBus.New()

// const (
// 	EventPlayerEnter = "player.enter"
// 	EventPlayerExit  = "player.exit"
// )

func NewProdLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewDevLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
}

type GameServer struct {
	// EventBus      EventBus.Bus
	Log zerolog.Logger
	// Accounts      map[string]*Account
	CommandManager *CommandManager
	// CommandParser *CommandParser
}

func NewGameServer() *GameServer {
	return &GameServer{
		Log: NewDevLogger(),
		// EventBus:      EventBus.New(),
		// Accounts:      make(map[string]*Account),
		CommandManager: NewCommandManager(),
		// CommandParser: NewCommandParser(),
	}
}

// func (gs *GameServer) loadAccounts(dirPath string) error {
// 	files, errReadDir := os.ReadDir(dirPath)
// 	if errReadDir != nil {
// 		return errReadDir
// 	}

// 	for _, file := range files {
// 		if file.IsDir() {
// 			continue
// 		}

// 		filePath := dirPath + "/" + file.Name()
// 		fileContent, errReadFile := os.ReadFile(filePath)
// 		if errReadFile != nil {
// 			return errReadFile
// 		}

// 		var account *Account
// 		errUnmarshal := json.Unmarshal(fileContent, &account)
// 		if errUnmarshal != nil {
// 			return errUnmarshal
// 		}

// 		gs.Accounts[strings.ToLower(account.Username)] = account
// 	}

// 	return nil
// }

func (gs *GameServer) Start() {
	gs.Log.Info().Msg("Starting server")

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", ":4000") // Port 4000
	if err != nil {
		gs.Log.Fatal().Err(err).Msg("Error starting server")
		return
	}
	defer listener.Close()

	// Register commands
	registerCommands(gs.CommandManager)
	setupRooms()

	gs.Log.Info().Msg("Server started")

	// Accept connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			gs.Log.Error().Err(err).Msg("Error accepting connection")
			continue
		}

		gs.Log.Debug().
			Str("remote_addr", conn.RemoteAddr().String()).
			Msg("New client connected")

		// Handle the client connection in a separate goroutine
		go gs.handleConnection(conn)
	}
}

func (gs *GameServer) handleConnection(conn net.Conn) {
	gs.Log.Debug().
		Str("remote_addr", conn.RemoteAddr().String()).
		Msg("Handling connection")
	defer conn.Close()

	// Create a new player for this connection
	player := &Player{
		Name:   "Player", // Default name; could be replaced later
		Conn:   conn,
		RoomID: "room1", // Starting room
	}

	// Send a welcome message to the client
	io.WriteString(player.Conn, "Welcome to the Go MUD server!\n")
	io.WriteString(player.Conn, "Type 'quit' to exit.\n")

	// Create a buffered reader for reading input from the client
	reader := bufio.NewReader(player.Conn)
	for {
		// Read input from the client
		io.WriteString(player.Conn, "> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			gs.Log.Error().
				Str("remote_addr", player.Conn.RemoteAddr().String()).
				Err(err).Msg("Client disconnected")

			return
		}

		// Clean up the input
		input = strings.TrimSpace(input)
		gs.Log.Debug().
			Str("remote_addr", conn.RemoteAddr().String()).
			Str("input", input).
			Msg("Received text")

		// fmt.Printf("Received: %s from %s\n", input, conn.RemoteAddr())

		gs.CommandManager.ParseAndExecute(input, player)

		// Parse the command and arguments
		// parts := strings.SplitN(input, " ", 2)
		// command := strings.ToLower(parts[0])
		// args := ""
		// if len(parts) > 1 {
		// 	args = parts[1]
		// }

		// // Find the command handler
		// if handler, exists := CommandMap[command]; exists {
		// 	// Execute the command handler
		// 	handler(player, args)
		// } else {
		// 	// Unknown command
		// 	io.WriteString(player.Conn, "Unknown command. Try 'look', 'move', or 'quit'.\n")
		// }
	}
}

// func (gs *GameServer) Start() {
// 	network := "tcp"
// 	address := ":4000"
// 	gs.Log.Info().Msg("Starting server")

// 	// Load accounts
// 	if err := gs.loadAccounts("_data/accounts"); err != nil {
// 		gs.Log.Fatal().
// 			Err(err).
// 			Msg("Error loading accounts")

// 		os.Exit(1)
// 	}

// 	gs.Log.Debug().
// 		Int("num_accounts", len(gs.Accounts)).
// 		Msg("Loaded accounts")

// 	// Register commands
// 	gs.CommandParser.RegisterCommand(lookCommand)
// 	gs.CommandParser.RegisterCommand(moveCommand)

// 	listener, err := net.Listen(network, address)
// 	if err != nil {
// 		gs.Log.Fatal().
// 			Err(err).
// 			Str("network", network).
// 			Str("address", address).
// 			Msg("Error starting telnet server")
// 	}
// 	defer listener.Close()

// 	gs.Log.Info().
// 		Str("network", network).
// 		Str("address", address).
// 		Msg("Telnet server started")

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			gs.Log.Error().
// 				Err(err).
// 				Str("network", network).
// 				Str("address", address).
// 				Msg("Error accepting connection")

// 			continue
// 		}
// 		go gs.handleConnection(conn)
// 	}
// }

// var i = 1 // TODO: Remove this

// func (gs *GameServer) handleConnection(conn net.Conn) {
// 	gs.Log.Debug().Msg("Handling connection")

// 	// START: Banner display
// 	defer conn.Close()
// 	banner := `
// Welcome to the MUD server!
// ===========================
// Press return to continue...
// `
// 	io.WriteString(conn,conn, cfmt.Sprint(banner))

// 	// Read input from the player
// 	scanner := bufio.NewScanner(conn)
// 	if scanner.Scan() {
// 		io.WriteString(conn,conn, cfmt.Sprint("Welcome to the game!\n"))
// 	}
// 	// END: Banner display

// 	// Create a new player
// 	startingRoom := setupWorld()
// 	player := NewPlayer(fmt.Sprintf("Hero%d", i), conn)
// 	i++
// 	player.Room = startingRoom

// 	// Start listening for player output
// 	go func() {
// 		for msg := range player.Out {
// 			io.WriteString(conn,conn, msg)
// 		}
// 	}()

// 	// Load player into the room and render the room
// 	eventBus.Publish(EventPlayerEnter, player, startingRoom.ID)

// 	// START: Game loop
// 	gs.Log.Debug().Msg("Entering game loop")
// 	for {
// 		io.WriteString(conn,conn, "> ")
// 		if !scanner.Scan() {
// 			break
// 		}
// 		input := scanner.Text()

// 		gs.Log.Debug().
// 			Str("input", input).
// 			Msg("Received text")

// 		gs.CommandParser.ParseAndExecute(input, player)
// 	}
// 	// END: Game loop
// }
