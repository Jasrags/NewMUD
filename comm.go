package main

// import (
// 	"io"
// 	"log/slog"
// 	"net"
// 	"strconv"
// 	"strings"

// 	"github.com/gliderlabs/ssh"
// 	"github.com/i582/cfmt/cmd/cfmt"
// 	"github.com/spf13/viper"
// 	"golang.org/x/term"
// )

// func setupServer() {
// 	address := net.JoinHostPort(viper.GetString("server.host"), viper.GetString("server.port"))

// 	slog.Info("Starting server",
// 		slog.String("address", address))

// 	server := &ssh.Server{
// 		Addr: address,
// 		// IdleTimeout:              viper.GetDuration("server.idle_timeout"), // TODO: reenable timeout later when we fix the connection close issues
// 		Handler: handleConnection,
// 		// ConnCallback:             ConnCallback,
// 		// ConnectionFailedCallback: ConnectionFailedCallback,
// 	}
// 	defer server.Close()

// 	if err := server.ListenAndServe(); err != nil {
// 		slog.Error("Error starting server",
// 			slog.String("address", address),
// 			slog.Any("error", err))
// 		return
// 	}
// }

// // func ConnCallback(ctx ssh.Context, conn net.Conn) net.Conn {
// // 	slog.Info("New connection",
// // 		slog.String("remote_address", conn.RemoteAddr().String()))

// // 	return conn
// // }

// // func ConnectionFailedCallback(conn net.Conn, err error) {
// // 	slog.Error("Connection failed",
// // 		slog.Any("error", err))
// // 	conn.Close()
// // }

// // func RunTheGame() {
// // 	slog.Info("Starting game server")
// // }

// // func GameLoop() {
// // 	slog.Info("Game loop started")

// // }

// func PromptForMenu(s ssh.Session, title string, options []string) (string, error) {

// 	var builder strings.Builder
// 	builder.WriteString(cfmt.Sprintf("{{%s}}::green", title))
// 	t := term.NewTerminal(s, cfmt.Sprint("{{Enter choice:}}::white|bold  "))

// 	for i, option := range options {
// 		builder.WriteString(cfmt.Sprintf("{{%d:}}::green|bold %-5s\n", i+1, option))
// 	}
// 	io.WriteString(s, builder.String())

// 	for {
// 		input, err := t.ReadLine()
// 		if err != nil {
// 			slog.Error("Error reading input", slog.Any("error", err))
// 			s.Close()
// 			return "", err
// 		}

// 		choice, err := strconv.Atoi(strings.TrimSpace(input))
// 		if err != nil || choice < 1 || choice > len(options) {
// 			io.WriteString(s, cfmt.Sprintf("{{Invalid choice, please try again.}}::red\n"))
// 			continue
// 		}

// 		return options[choice-1], nil
// 	}
// }

// func PromptForInput(s ssh.Session, prompt string) (string, error) {
// 	t := term.NewTerminal(s, prompt)
// 	input, err := t.ReadLine()
// 	if err != nil {
// 		slog.Error("Error reading input", slog.Any("error", err))
// 		s.Close()

// 		return "", err
// 	}

// 	return strings.TrimSpace(input), nil
// }

// func PromptForPassword(s ssh.Session, prompt string) (string, error) {
// 	t := term.NewTerminal(s, prompt)
// 	input, err := t.ReadPassword(prompt)
// 	if err != nil {
// 		slog.Error("Error reading password", slog.Any("error", err))
// 		s.Close()

// 		return "", err
// 	}

// 	return strings.TrimSpace(input), nil
// }

// func SendToChar(s ssh.Session, message string) {
// 	io.WriteString(s, cfmt.Sprintf("{{%s}}::white\n", message))
// }

// // void send_to_all(char *messg)

// // void send_to_room(char *messg, int room)
// func SendToRoom(s ssh.Session, message string,
// 	room *Room) {
// 	// for _, c := range room.Characters {
// 	// io.WriteString(s, cfmt.Sprintf("{{%s}}::white\n", message))
// 	// }
// }
