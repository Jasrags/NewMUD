package game

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
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

func RenderMenu(s ssh.Session, title string, options []string) {
	// Create and render the menu list
	l := list.New(options).
		Enumerator(list.Arabic).
		EnumeratorStyle(boldWhiteText).
		ItemStyle(boldGreenText)

	// Render the menu
	io.WriteString(s, lipgloss.JoinVertical(lipgloss.Left,
		greenText.Render(title),
		l.String()+"\n",
		"",
	))
}

func ReadChoice(s ssh.Session, options []string) (int, error) {
	// Create terminal for input
	t := term.NewTerminal(s, "")

	for {
		// Write the prompt
		io.WriteString(s, "\r"+boldWhiteText.Render("Enter choice: "))

		// Read user input
		input, err := t.ReadLine()
		if err != nil {
			return 0, fmt.Errorf("error reading input: %w", err)
		}

		// Parse and validate choice
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(options) {
			io.WriteString(s, redText.Render("Invalid choice, please try again.\n"))
			continue
		}

		return choice, nil
	}
}

// func WriteToTerminal(writer termenv.Output, content string) {
// 	// Styled content
// 	style := termenv.String(content).Bold().Foreground(termenv.Color("34")) // Blue bold text
// 	writer.WriteString(style.String())
// }

func PromptForMenu(s ssh.Session, title string, options []string) (string, error) {
	// Render the menu using lipgloss
	menu := lipgloss.JoinVertical(
		lipgloss.Left,
		cfmt.Sprintf("{{%s}}::green|bold\n", title),
		list.New(options).
			Enumerator(list.Arabic).
			EnumeratorStyle(boldWhiteText).
			ItemStyle(boldGreenText).
			String(),
		"",
	)

	// Display the menu
	io.WriteString(s, menu+"\n")

	// Create a terminal for input
	t := term.NewTerminal(s, "")
	for {
		// Prompt for user input
		io.WriteString(s, cfmt.Sprint("{{Enter choice:}}::white|bold\n"))
		input, err := t.ReadLine()
		if err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}

		// Validate choice
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(options) {
			io.WriteString(s, cfmt.Sprint("{{Invalid choice, please try again.}}::red\n"))
			continue
		}

		// Return the selected option
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
	io.WriteString(s, whiteText.Render("%s\n", message))
}

// void send_to_all(char *messg)

// void send_to_room(char *messg, int room)
func SendToRoom(s ssh.Session, message string,
	room *Room) {
	// for _, c := range room.Characters {
	// io.WriteString(s, cfmt.Sprintf("{{%s}}::white\n", message))
	// }
}
