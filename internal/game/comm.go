package game

import (
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func handleConnection(s ssh.Session) {
	defer s.Close()

	_, winCh, _ := s.Pty()

	// Set the window size
	go func() {
		for win := range winCh {
			slog.Debug("Window size changed",
				slog.Int("width", win.Width),
				slog.Int("height", win.Height),
			)
		}
	}()

	var account *Account
	var char *Character
	// var room *Room
	var state = StateWelcome

	for {
		switch state {
		case StateWelcome:
			state = PromptWelcome(s)
		case StateLogin:
			state, account = PromptLogin(s)
		case StateRegistration:
			state, account = PromptRegistration(s)
		case StateMainMenu:
			state = PromptMainMenu(s, account)
		case StateChangePassword:
			state = PromptChangePassword(s, account)
			// case StateCharacterSelect:
			// state, char = PromptCharacterSelect(s, user)
		case StateCharacterCreate:
			state = PromptCharacterCreate(s, account)
		case StateEnterGame:
			state, char = PromptEnterGame(s, account)
		case StateGameLoop:
			state = PromptGameLoop(s, account, char)
		case StateExitGame:
			state = PromptExitGame(s, account, char)
		case StateQuit:
			fallthrough
		case StateError:
			s.Close()
			char = nil
			account = nil
			return
		default:
			slog.Error("Invalid state", slog.String("user_state", state))
			s.Close()
			char = nil
			account = nil
		}
	}
}

func SetupServer() {
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

// func RunTheGame() {
// 	slog.Info("Starting game server")
// }

// func GameLoop() {
// 	slog.Info("Game loop started")

// }

func PromptForMenu(s ssh.Session, title string, options []string) (string, error) {

	var builder strings.Builder
	builder.WriteString(cfmt.Sprintf("{{%s}}::green", title))
	t := term.NewTerminal(s, cfmt.Sprint("{{Enter choice:}}::white|bold  "))

	for i, option := range options {
		builder.WriteString(cfmt.Sprintf("{{%d:}}::green|bold %-5s\n", i+1, option))
	}
	io.WriteString(s, builder.String())

	for {
		input, err := t.ReadLine()
		if err != nil {
			slog.Error("Error reading input", slog.Any("error", err))
			s.Close()
			return "", err
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(options) {
			io.WriteString(s, cfmt.Sprintf("{{Invalid choice, please try again.}}::red\n"))
			continue
		}

		return options[choice-1], nil
	}
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

func PromptForPassword(s ssh.Session, prompt string) (string, error) {
	t := term.NewTerminal(s, prompt)
	input, err := t.ReadPassword(prompt)
	if err != nil {
		slog.Error("Error reading password", slog.Any("error", err))
		s.Close()

		return "", err
	}

	return strings.TrimSpace(input), nil
}

func SendToChar(s ssh.Session, message string) {
	io.WriteString(s, cfmt.Sprintf("{{%s}}::white\n", message))
}

// void send_to_all(char *messg)

// void send_to_room(char *messg, int room)
func SendToRoom(s ssh.Session, message string,
	room *Room) {
	// for _, c := range room.Characters {
	// io.WriteString(s, cfmt.Sprintf("{{%s}}::white\n", message))
	// }
}
