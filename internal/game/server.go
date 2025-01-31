package game

import (
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
)

const ()

type (
	GameServer struct {
		TickDuration time.Duration
	}
)

func NewGameServer() *GameServer {
	return &GameServer{}
}

// Init initializes the game server
func (s *GameServer) Init() {
	s.SetupConfig()
	s.SetupLogger()

	// Update configuration on change
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Config updated",
			slog.String("file", e.Name))
	})
	viper.WatchConfig()

	s.OutputConfig()

	// Set game server properties
	s.TickDuration = viper.GetDuration("server.tick_duration")

	go GameTimeMgr.StartTicker(s.TickDuration)

	EntityMgr.LoadDataFiles()
	AccountMgr.LoadDataFiles()
	CharacterMgr.LoadDataFiles()

	RegisterCommands()
}

// Start starts the game server
func (s *GameServer) Start() {
	address := net.JoinHostPort(
		viper.GetString("server.host"),
		viper.GetString("server.port"))

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
}

func (s *GameServer) Stop() {
	// Stop the game server
}

func (s *GameServer) SetupConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func (s *GameServer) SetupLogger() {
	// Setup logger
	logLevel := viper.GetString("server.log_level")
	logHandler := viper.GetString("server.log_handler")

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
	case "color":
		logger = slog.New(GetColorLogHandler(os.Stderr, programLevel))
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

func (s *GameServer) OutputConfig() {
	slog.Info("Current configuration")

	for _, key := range viper.AllKeys() {
		slog.Debug("config",
			slog.String("key", key),
			slog.String("value", viper.GetString(key)))
	}
}
