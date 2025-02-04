package game

import (
	"log/slog"
	"net"

	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
)

// GameContext holds shared state during the connection.
type GameContext struct {
	Account   *Account
	Character *Character
}

// stateHandler is a function that processes a state and returns the next state.
type stateHandler func(s ssh.Session, ctx *GameContext) string

func welcomeState(s ssh.Session, ctx *GameContext) string {
	return PromptWelcome(s)
}

func loginState(s ssh.Session, ctx *GameContext) string {
	state, acc := PromptLogin(s)
	ctx.Account = acc
	return state
}

func registrationState(s ssh.Session, ctx *GameContext) string {
	state, acc := PromptRegistration(s)
	ctx.Account = acc
	return state
}

func mainMenuState(s ssh.Session, ctx *GameContext) string {
	return PromptMainMenu(s, ctx.Account)
}

func changePasswordState(s ssh.Session, ctx *GameContext) string {
	return PromptChangePassword(s, ctx.Account)
}

func characterCreateState(s ssh.Session, ctx *GameContext) string {
	state, char := PromptCharacterCreate(s, ctx.Account)
	ctx.Character = char
	return state
}

func enterGameState(s ssh.Session, ctx *GameContext) string {
	state, char := PromptEnterGame(s, ctx.Account)
	ctx.Character = char
	return state
}

func gameLoopState(s ssh.Session, ctx *GameContext) string {
	return PromptGameLoop(s, ctx.Account, ctx.Character)
}

func exitGameState(s ssh.Session, ctx *GameContext) string {
	return PromptExitGame(s, ctx.Account, ctx.Character)
}

var stateHandlers = map[string]stateHandler{
	StateWelcome:         welcomeState,
	StateLogin:           loginState,
	StateRegistration:    registrationState,
	StateMainMenu:        mainMenuState,
	StateChangePassword:  changePasswordState,
	StateCharacterCreate: characterCreateState,
	StateEnterGame:       enterGameState,
	StateGameLoop:        gameLoopState,
	StateExitGame:        exitGameState,
}

func handleConnection(s ssh.Session) {
	defer s.Close()

	_, winCh, _ := s.Pty()
	go func() {
		for win := range winCh {
			slog.Debug("Window size changed",
				slog.Int("width", win.Width),
				slog.Int("height", win.Height),
			)
		}
	}()

	ctx := &GameContext{}
	state := StateWelcome

	for {
		// Look up the handler for the current state.
		handler, ok := stateHandlers[state]
		if !ok {
			slog.Error("Invalid state", slog.String("state", state))
			break
		}

		// Process the state and get the next one.
		state = handler(s, ctx)

		// If the state indicates exit, break out of the loop.
		if state == StateQuit || state == StateError {
			break
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
