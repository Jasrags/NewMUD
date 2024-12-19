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

type GameServer struct {
	Log            zerolog.Logger
	CommandManager *CommandManager
	AreaManager    *AreaManager
	RoomManager    *RoomManager
}

func NewGameServer() *GameServer {
	gs := &GameServer{
		Log: NewDevLogger(),
	}
	gs.CommandManager = NewCommandManager()
	gs.RoomManager = NewRoomManager()
	gs.AreaManager = NewAreaManager(gs.RoomManager)

	return gs
}

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
	gs.CommandManager.Load()

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

var i int // TODO: remove after we have player/account login

func (gs *GameServer) handleConnection(conn net.Conn) {
	gs.Log.Debug().
		Str("remote_addr", conn.RemoteAddr().String()).
		Msg("Handling connection")
	defer conn.Close()

	gs.DisplayBanner(conn)

	// Create a new player for this connection
	player := NewPlayer(conn)
	if i == 0 {
		player.Role = "admin"
		player.Name = "Admin"
	} else {
		player.Role = "player"
		player.Name = "Player"
	}
	i++

	player.MoveTo(gs.RoomManager.GetRoom("limbo:the_void"))
	io.WriteString(player.Conn, RenderRoom(player, player.Room))

	gs.GameLoop(conn, player)
}

func (gs *GameServer) DisplayBanner(conn net.Conn) {
	gs.Log.Debug().Msg("Displaying banner")

	banner := `
{{Welcome to the MUD server!}}::green
{{==========================}}::white|bold

`
	io.WriteString(conn, cfmt.Sprintf(banner))
}

func (gs *GameServer) GameLoop(conn net.Conn, player *Player) {
	gs.Log.Debug().Msg("Entering game loop")

	ctx := &GameContext{
		Log:            NewDevLogger(),
		RoomManager:    gs.RoomManager,
		AreaManager:    gs.AreaManager,
		CommandManager: gs.CommandManager,
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
