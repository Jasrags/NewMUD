package main

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Jasrags/NewMUD/characters"
	"github.com/Jasrags/NewMUD/commands"
	"github.com/Jasrags/NewMUD/connections"
	"github.com/Jasrags/NewMUD/items"
	"github.com/Jasrags/NewMUD/rooms"
	"github.com/Jasrags/NewMUD/users"
	"github.com/fsnotify/fsnotify"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
)

var (
	sigChan = make(chan os.Signal, 1)
	wg      sync.WaitGroup
)

func main() {
	// gs := mud.NewGameServer()
	// gs.Start()

	setupConfig()
	loadAllDataFiles()

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	listener := setupServer(&wg)
	defer listener.Close()

	// block until a signal comes in
	<-sigChan

	slog.Info("Shutting down")

	wg.Wait()
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
	for _, key := range viper.AllKeys() {
		slog.Debug("config",
			slog.String("key", key),
			slog.String("value", viper.GetString(key)))
	}
}

func setupLogger() {
	// Parse and set log level
	logLevel := viper.GetString("server.log_level")
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
	logHandler := viper.GetString("server.log_handler")
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

func setupServer(wg *sync.WaitGroup) net.Listener {
	address := strings.Join([]string{viper.GetString("server.host"), viper.GetString("server.port")}, ":")

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Error starting server",
			slog.String("address", address),
			slog.Any("error", err))
		return nil
	}

	slog.Info("Server started", slog.String("address", address))

	// Start a goroutine to accept incoming connections, so that we can use a signal to stop the server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				slog.Error("Error accepting connection",
					slog.Any("error", err))
				continue
			}

			// TODO: Check for max connections

			wg.Add(1)
			go handleConnection(connections.NewNetConnection(conn), wg)
		}
	}()

	return listener

}

func handleConnection(netConn *connections.NetConnection, wg *sync.WaitGroup) {
	defer wg.Done()

	slog.Debug("New client connected",
		slog.String("id", netConn.ID),
		slog.String("remote_address", netConn.Conn.RemoteAddr().String()))

	displayBanner(netConn.Conn)
	gameLoop(netConn.Conn)

}

func loadAllDataFiles() {
	rooms.LoadDataFiles()
	items.LoadDataFiles()
	users.LoadDataFiles()
	characters.LoadDataFiles()
	commands.RegisterCommands()
}

func displayBanner(conn net.Conn) {
	slog.Debug("Displaying banner")

	banner := `
{{Welcome to the MUD server!}}::green
{{==========================}}::white|bold

`
	io.WriteString(conn, cfmt.Sprint(banner))
}

func gameLoop(conn net.Conn) {
	slog.Debug("Game Loop")

	reader := bufio.NewReader(conn)
	for {
		io.WriteString(conn, "> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("Error reading input",
				slog.Any("error", err))
			return
		}

		input = strings.TrimSpace(input)
		slog.Debug("Received text",
			slog.String("input", input))
	}

	// // Create a new player for this connection
	// player := NewPlayer(gs.Log, conn)
	// if i == 0 {
	// 	player.Role = "admin"
	// 	player.Name = "Admin"
	// } else {
	// 	player.Role = "player"
	// 	player.Name = "Player"
	// }
	// i++

	// player.MoveTo(gs.GameContext.RoomManager.GetRoom("limbo:the_void"))
	// io.WriteString(player.Conn, RenderRoom(player, player.Room))

	// gs.GameLoop(conn, player)
}
