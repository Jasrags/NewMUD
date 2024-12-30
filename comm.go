package main

import (
	"log/slog"
	"net"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

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

func RunTheGame() {
	slog.Info("Starting game server")
}

func GameLoop() {
	slog.Info("Game loop started")

}

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
