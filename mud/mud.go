package mud

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/rs/zerolog"
)

var eventBus = EventBus.New()

const (
	EventPlayerEnter = "player.enter"
	EventPlayerExit  = "player.exit"
)

func NewProdLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func NewDevLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
}

type GameServer struct {
	// EventBus      EventBus.Bus
	Log           zerolog.Logger
	Accounts      map[string]*Account
	CommandParser *CommandParser
}

func NewGameServer() *GameServer {
	return &GameServer{
		Log: NewDevLogger(),
		// EventBus:      EventBus.New(),
		Accounts:      make(map[string]*Account),
		CommandParser: NewCommandParser(),
	}
}

func (gs *GameServer) loadAccounts(dirPath string) error {
	files, errReadDir := os.ReadDir(dirPath)
	if errReadDir != nil {
		return errReadDir
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := dirPath + "/" + file.Name()
		fileContent, errReadFile := os.ReadFile(filePath)
		if errReadFile != nil {
			return errReadFile
		}

		var account *Account
		errUnmarshal := json.Unmarshal(fileContent, &account)
		if errUnmarshal != nil {
			return errUnmarshal
		}

		gs.Accounts[strings.ToLower(account.Username)] = account
	}

	return nil
}

func (gs *GameServer) Start() {
	network := "tcp"
	address := ":4000"
	gs.Log.Info().Msg("Starting server")

	// Load accounts
	if err := gs.loadAccounts("_data/accounts"); err != nil {
		gs.Log.Fatal().
			Err(err).
			Msg("Error loading accounts")

		os.Exit(1)
	}

	gs.Log.Debug().
		Int("num_accounts", len(gs.Accounts)).
		Msg("Loaded accounts")

	// Register commands
	gs.CommandParser.RegisterCommand(lookCommand)
	gs.CommandParser.RegisterCommand(moveCommand)

	listener, err := net.Listen(network, address)
	if err != nil {
		gs.Log.Fatal().
			Err(err).
			Str("network", network).
			Str("address", address).
			Msg("Error starting telnet server")
	}
	defer listener.Close()

	gs.Log.Info().
		Str("network", network).
		Str("address", address).
		Msg("Telnet server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			gs.Log.Error().
				Err(err).
				Str("network", network).
				Str("address", address).
				Msg("Error accepting connection")

			continue
		}
		go gs.handleConnection(conn)
	}
}

var i = 1 // TODO: Remove this

func (gs *GameServer) handleConnection(conn net.Conn) {
	gs.Log.Debug().Msg("Handling connection")

	// START: Banner display
	defer conn.Close()
	banner := `
Welcome to the MUD server!
===========================
Press return to continue...
`
	io.WriteString(conn, cfmt.Sprint(banner))

	// Read input from the player
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		io.WriteString(conn, cfmt.Sprint("Welcome to the game!\n"))
	}
	// END: Banner display

	// Create a new player
	startingRoom := setupWorld()
	player := NewPlayer(fmt.Sprintf("Hero%d", i), conn)
	i++
	player.Room = startingRoom

	// Start listening for player output
	go func() {
		for msg := range player.Out {
			io.WriteString(conn, msg)
		}
	}()

	// Load player into the room and render the room
	eventBus.Publish(EventPlayerEnter, player, startingRoom.ID)

	// START: Game loop
	gs.Log.Debug().Msg("Entering game loop")
	for {
		io.WriteString(conn, "> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		gs.Log.Debug().
			Str("input", input).
			Msg("Received text")

		gs.CommandParser.ParseAndExecute(input, player)
	}
	// END: Game loop
}
