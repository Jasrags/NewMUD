package mud

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
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

func (gs *GameServer) setupConfig() {
	gs.Log.Info().Msg("Setting up configuration")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		gs.Log.Fatal().
			Err(err).
			Msg("fatal error config file")
	}

	// Update configuration on change
	viper.OnConfigChange(func(e fsnotify.Event) {
		gs.Log.Info().
			Str("file", e.Name).
			Msg("Config file changed")
		gs.outputConfig()
	})
	viper.WatchConfig()

	// Attempt to parse the log_level from the config file
	serverLogLevel := viper.GetString("server.log_level")
	logLevel, err := zerolog.ParseLevel(serverLogLevel)
	if err != nil {
		gs.Log.Fatal().
			Err(err).
			Msg("unable to parse log_level")
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}

	gs.outputConfig()
}

func (gs *GameServer) outputConfig() {
	gs.Log.Info().Msg("Outputting configuration")

	for _, key := range viper.AllKeys() {
		gs.Log.Debug().
			Str("key", key).
			Str("value", viper.GetString(key)).
			Msg("Config")
	}
}

func (gs *GameServer) Start() {
	gs.Log.Info().Msg("Starting server")

	gs.setupConfig()
	address := strings.Join([]string{viper.GetString("server.host"), viper.GetString("server.port")}, ":")

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", address) // Port 4000
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
