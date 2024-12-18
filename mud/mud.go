package mud

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
)

func NewProdLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewDevLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
}

// type State struct {
// 	CommandManager *CommandManager
// 	AreaManager    *AreaManager
// 	RoomManager    *RoomManager
// }

type GameServer struct {
	Log zerolog.Logger
	// State *State
	// Accounts      map[string]*Account
	CommandManager *CommandManager
	AreaManager    *AreaManager
	RoomManager    *RoomManager
}

func NewGameServer() *GameServer {
	return &GameServer{
		Log: NewDevLogger(),
		// EventBus:      EventBus.New(),
		// Accounts:      make(map[string]*Account),
		// State: &State{
		CommandManager: NewCommandManager(),
		AreaManager:    NewAreaManager(),
		RoomManager:    NewRoomManager(),
		// },
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
	gs.RoomManager.Load()
	gs.AreaManager.Load()

	registerCommands(gs.CommandManager)

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

	gs.DisplayBanner(conn)

	// Create a new player for this connection
	player := NewPlayer(conn)
	player.Name = "Player"
	player.Room = gs.RoomManager.GetRoom("limbo:the_void")
	player.RoomID = "limbo:the_void"

	gs.GameLoop(conn, player)
}

func (gs *GameServer) DisplayBanner(conn net.Conn) {
	gs.Log.Debug().Msg("Displaying banner")

	banner := `
Welcome to the MUD server!
==========================

`
	io.WriteString(conn, cfmt.Sprintf(banner))
}

func (gs *GameServer) GameLoop(conn net.Conn, player *Player) {
	gs.Log.Debug().Msg("Entering game loop")

	ctx := &GameContext{
		Log:         NewDevLogger(),
		RoomManager: gs.RoomManager,
		AreaManager: gs.AreaManager,
	}

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

		gs.CommandManager.ParseAndExecute(ctx, input, player)
	}
}
