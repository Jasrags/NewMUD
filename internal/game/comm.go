package game

import (
	"log/slog"
	"net"

	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
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
		Addr:    address,
		Handler: handleConnection,
	}
	defer server.Close()

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Error starting server",
			slog.String("address", address),
			slog.Any("error", err))
		return
	}
}
