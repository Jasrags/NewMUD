package main

import (
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/spf13/viper"
)

const (
	// This will skip straight to the game loop
	StateDebug           = "debug"
	StateWelcome         = "welcome"
	StateLogin           = "login"
	StateRegistration    = "registration"
	StateMainMenu        = "main_menu"
	StateChangePassword  = "change_password"
	StateCharacterSelect = "character_select"
	StateCharacterCreate = "character_create"
	StateEnterGame       = "enter_game"
	StateGameLoop        = "game_loop"
	StateExitGame        = "exit_game"
	StateQuit            = "quit"
	StateError           = "error"
)

func promptWelcome(s ssh.Session) string {
	slog.Debug("Welcome state",
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))
	var builder strings.Builder
	builder.WriteString("{{     ::::::::  :::    :::     :::     :::::::::   ::::::::  :::       ::: ::::    ::::  :::    ::: :::::::::  }}::#ff8700\n")
	builder.WriteString("{{    :+:    :+: :+:    :+:   :+: :+:   :+:    :+: :+:    :+: :+:       :+: +:+:+: :+:+:+ :+:    :+: :+:    :+: }}::#ff5f00\n")
	builder.WriteString("{{    +:+        +:+    +:+  +:+   +:+  +:+    +:+ +:+    +:+ +:+       +:+ +:+ +:+:+ +:+ +:+    +:+ +:+    +:+ }}::#ff0000\n")
	builder.WriteString("{{    +#++:++#++ +#++:++#++ +#++:++#++: +#+    +:+ +#+    +:+ +#+  +:+  +#+ +#+  +:+  +#+ +#+    +:+ +#+    +:+ }}::#d70000\n")
	builder.WriteString("{{           +#+ +#+    +#+ +#+     +#+ +#+    +#+ +#+    +#+ +#+ +#+#+ +#+ +#+       +#+ +#+    +#+ +#+    +#+ }}::#af0000\n")
	builder.WriteString("{{    #+#    #+# #+#    #+# #+#     #+# #+#    #+# #+#    #+#  #+#+# #+#+#  #+#       #+# #+#    #+# #+#    #+# }}::#870000\n")
	builder.WriteString("{{     ########  ###    ### ###     ### #########   ########    ###   ###   ###       ###  ########  #########  }}::#5f0000\n")

	// Check if login is enabled
	if !viper.GetBool("server.login_enabled") {
		builder.WriteString(cfmt.Sprint("\n{{Login is disabled.}}::red\n"))
	}

	io.WriteString(s, cfmt.Sprint(builder.String()))

	if _, err := PromptForInput(s, cfmt.Sprint("{{Press enter to continue...}}::white|bold\n")); err != nil {
		return StateError
	}

	return StateLogin
}

func promptLogin(s ssh.Session) (string, *Account) {
	slog.Debug("Login state",
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

promptUsername:
	io.WriteString(s, cfmt.Sprint("{{Enter your username to continue or type}}::white|bold {{new}}::green|bold {{to register:}}::white|bold\n"))

	username, err := PromptForInput(s, cfmt.Sprint("{{Username:}}::white|bold "))
	if err != nil {
		return StateError, nil
	}

	// New user registration
	if strings.EqualFold(username, "new") {
		return StateRegistration, nil
	}

	password, err := PromptForPassword(s, cfmt.Sprint("{{Password:}}::white|bold "))
	if err != nil {
		return StateError, nil
	}

	// Check if user exists
	u := AccountMgr.GetByUsername(username)

	if u == nil {
		slog.Warn("User does not exist")
		io.WriteString(s, cfmt.Sprint("{{Invalid username or password}}::red\n"))

		goto promptUsername
	}

	// Validate password against user's hashed password
	if !u.CheckPassword(password) {
		slog.Warn("Invalid password")

		io.WriteString(s, cfmt.Sprint("{{Invalid username or password}}::red\n"))

		goto promptUsername
	}

	// TODO: Check if user is already logged in

	// TODO: Check if user is banned

	return StateMainMenu, u
}

func promptRegistration(s ssh.Session) (string, *Account) {
	slog.Debug("Registration state",
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	if !viper.GetBool("server.registration_enabled") {
		io.WriteString(s, cfmt.Sprint("\n{{Registration is disabled.}}::red\n"))

		return StateLogin, nil
	}

promptUsername:
	io.WriteString(s, cfmt.Sprint("{{User registration}}::green\n"))
	username, err := PromptForInput(s, cfmt.Sprint("{{Enter your username: }}::white|bold"))
	if err != nil {
		return StateError, nil
	}

	// Check if username is empty
	if username == "" {
		io.WriteString(s, cfmt.Sprint("{{Username cannot be empty.}}::red\n"))
		goto promptUsername
	}

	// Check if username is within the allowed length
	if len(username) < viper.GetInt("server.username_min_length") || len(username) > viper.GetInt("server.username_max_length") {
		io.WriteString(s, cfmt.Sprintf("{{Username must be between %d and %d characters.}}::red\n", viper.GetInt("server.username_min_length"), viper.GetInt("server.username_max_length")))
		goto promptUsername
	}

	// Check if username already exists
	if AccountMgr.Exists(username) {
		io.WriteString(s, cfmt.Sprint("{{Username already exists.}}::red\n"))
		goto promptUsername
	}

	// Check if username is banned
	if AccountMgr.IsBannedName(username) {
		io.WriteString(s, cfmt.Sprint("{{Username is not allowed.}}::red\n"))
		goto promptUsername
	}

promptPassword:
	password, err := PromptForPassword(s, cfmt.Sprint("{{Enter your password:}}::white|bold "))
	if err != nil {
		return StateError, nil
	}

	// Check if password is empty
	if password == "" {
		io.WriteString(s, cfmt.Sprint("{{Password cannot be empty.}}::red\n"))
		goto promptPassword
	}

	// Check if password is within the allowed length
	if len(password) < viper.GetInt("server.password_min_length") || len(password) > viper.GetInt("server.password_max_length") {
		io.WriteString(s, cfmt.Sprintf("{{Password must be between %d and %d characters.}}::red\n", viper.GetInt("server.password_min_length"), viper.GetInt("server.password_max_length")))
		goto promptPassword
	}

	confirmPassword, err := PromptForPassword(s, cfmt.Sprint("{{Confirm your password:}}::white|bold "))
	if err != nil {
		return StateError, nil
	}

	// Check if confirm password is empty
	if confirmPassword == "" {
		io.WriteString(s, cfmt.Sprint("{{Password cannot be empty.}}::red\n"))
		goto promptPassword
	}

	// Check if passwords match
	if password != confirmPassword {
		io.WriteString(s, cfmt.Sprint("{{Passwords do not match.}}::red\n"))
		goto promptPassword
	}

	// Create a new user
	u := NewAccount()
	u.Username = username
	u.SetPassword(password)
	u.Save()
	AccountMgr.AddAccount(u)

	return StateMainMenu, u
}

func promptMainMenu(s ssh.Session, u *Account) string {
	slog.Debug("Main menu state",
		slog.String("username", u.Username),
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	option, err := PromptForMenu(s, cfmt.Sprint("{{Main Menu}}::green\n"),
		[]string{"Enter Game", "Create Character", "Change Password", "Quit"})
	if err != nil {
		return StateError
	}

	switch option {
	case "Enter Game":
		return StateEnterGame
	case "Create Character":
		return StateCharacterCreate
	case "Change Password":
		return StateChangePassword
	case "Quit":
		return StateQuit
	}

	slog.Debug("Selected option",
		slog.String("option", option))

	if _, err := PromptForInput(s, cfmt.Sprint("{{Press enter to continue...}}::white|bold\n")); err != nil {
		return StateError
	}

	return StateMainMenu
}

func promptChangePassword(s ssh.Session, u *Account) string {
	slog.Debug("Change password state",
		slog.String("username", u.Username),
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	password, err := PromptForPassword(s, cfmt.Sprint("{{Enter your current password:}}::white|bold "))
	if err != nil {
		return StateError
	}

	if !u.CheckPassword(password) {
		io.WriteString(s, cfmt.Sprint("{{Invalid password.}}::red\n"))
		return StateChangePassword
	}

	newPassword, err := PromptForPassword(s, cfmt.Sprint("{{Enter your new password:}}::white|bold "))
	if err != nil {
		return StateError
	}

	confirmNewPassword, err := PromptForPassword(s, cfmt.Sprint("{{Confirm your new password:}}::white|bold "))
	if err != nil {
		return StateError
	}

	if newPassword != confirmNewPassword {
		io.WriteString(s, cfmt.Sprint("{{Passwords do not match.}}::red\n"))
		return StateChangePassword
	}

	u.SetPassword(newPassword)
	u.Save()

	io.WriteString(s, cfmt.Sprint("{{Password changed successfully.}}::green\n"))

	return StateMainMenu
}

func promptCharacterCreate(s ssh.Session, u *Account) string {
	slog.Debug("Character create state",
		slog.String("username", u.Username),
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	// Step 1: Prompt for character name
	io.WriteString(s, cfmt.Sprintf("{{Enter your character's name:}}::cyan\n"))
	name, err := PromptForInput(s, "> ")
	if err != nil {
		slog.Error("Error reading character name", slog.Any("error", err))
		io.WriteString(s, cfmt.Sprintf("{{Error reading input. Returning to main menu.}}::red\n"))
		return StateMainMenu
	}
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		io.WriteString(s, cfmt.Sprintf("{{Name cannot be empty. Returning to main menu.}}::red\n"))
		return StateMainMenu
	}

	// Step 2: Prompt for character description
	io.WriteString(s, cfmt.Sprintf("{{Enter a short description for your character:}}::cyan\n"))
	description, err := PromptForInput(s, "> ")
	if err != nil {
		slog.Error("Error reading character description", slog.Any("error", err))
		io.WriteString(s, cfmt.Sprintf("{{Error reading input. Returning to main menu.}}::red\n"))
		return StateMainMenu
	}
	description = strings.TrimSpace(description)

	// Step 3: Set base attributes
	io.WriteString(s, cfmt.Sprintf("{{Setting base attributes...}}::green\n"))
	baseAttributes := &Attributes{
		Body:      Attribute[int]{Name: "Body", Base: 5},
		Agility:   Attribute[int]{Name: "Agility", Base: 6},
		Reaction:  Attribute[int]{Name: "Reaction", Base: 4},
		Strength:  Attribute[int]{Name: "Strength", Base: 5},
		Willpower: Attribute[int]{Name: "Willpower", Base: 4},
		Logic:     Attribute[int]{Name: "Logic", Base: 4},
		Intuition: Attribute[int]{Name: "Intuition", Base: 5},
		Charisma:  Attribute[int]{Name: "Charisma", Base: 4},
		Edge:      Attribute[int]{Name: "Edge", Base: 5},
		Essence:   Attribute[float64]{Name: "Essence", Base: 5.6},
		Magic:     Attribute[int]{Name: "Magic", Base: 0},
		Resonance: Attribute[int]{Name: "Resonance", Base: 0},
	}

	// Step 4: Create the character
	char := &Character{
		GameEntity: GameEntity{
			ID:          uuid.New().String(),
			Name:        name,
			Description: description,
			Attributes:  baseAttributes,
			Equipment:   make(map[string]*Item),
		},
		UserID:    u.ID,
		Role:      CharacterRolePlayer,
		CreatedAt: time.Now(),
	}
	char.Save()

	// Step 5: Add character to user
	u.Characters = append(u.Characters, char.Name)
	u.Save()

	// Step 6: Save user
	// err = UserMgr.SaveUser(u)
	// if err != nil {
	// 	slog.Error("Error saving user after character creation", slog.Any("error", err))
	// 	io.WriteString(s, cfmt.Sprintf("{{Error saving character. Returning to main menu.}}::red\n"))
	// 	return StateMainMenu
	// }

	// Step 7: Confirmation and return to main menu
	io.WriteString(s, cfmt.Sprintf("{{Character '%s' created successfully! Returning to main menu.}}::green\n", name))
	return StateMainMenu
}

// func promptCharacterSelect(s ssh.Session, u *User) (string, *Character) {
// 	slog.Debug("Character select state",
// 		slog.String("username", u.Username),
// 		slog.String("remote_address", s.RemoteAddr().String()),
// 		slog.String("session_id", s.Context().SessionID()))

// 	if len(u.Characters) == 0 {
// 		io.WriteString(s, cfmt.Sprintf("{{You have no characters.}}::red\n"))
// 		return StateCharacterSelect, nil
// 	}

// 	io.WriteString(s, cfmt.Sprintf("{{Select a character:}}::green\n"))
// 	for i, c := range u.Characters {
// 		io.WriteString(s, cfmt.Sprintf("{{%d. %s}}::green\n", i+1, c))
// 	}

// 	if _, err := PromptForInput(s, cfmt.Sprint("{{Press enter to continue...}}::white|bold\n")); err != nil {
// 		return StateError, nil
// 	}

// 	return StateCharacterSelect, nil
// }

func promptEnterGame(s ssh.Session, u *Account) (string, *Character) {
	slog.Debug("Enter game state",
		slog.String("username", u.Username),
		// slog.String("character_name", c.Name),
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	io.WriteString(s, cfmt.Sprintf("{{Welcome to the game, %s!}}::green\n", u.Username))

	// // Check if user has characters
	// if len(u.Characters) == 0 {
	// 	io.WriteString(s, cfmt.Sprintf("{{You have no characters.}}::red\n"))

	// 	// TODO: Remove this when we have character creation
	// 	c := NewCharacter()
	// 	c.Name = u.Username
	// 	c.Save()
	// 	CharacterMgr.AddCharacter(c)
	// 	u.AddCharacter(c)
	// 	u.Save()

	// 	return StateMainMenu, nil
	// }

	var characters []string
	for _, name := range u.Characters {
		char := CharacterMgr.GetCharacterByName(name)
		if char == nil {
			continue
		}

		characters = append(characters, char.Name)
	}

	// Prompt to select a character
	option, err := PromptForMenu(s, cfmt.Sprint("{{Select a character:}}::green\n"), characters)
	if err != nil {
		return StateError, nil
	}

	option = strings.ToLower(option)

	slog.Debug("Selected character",
		slog.String("character", option))

	// Load the character
	c := CharacterMgr.GetCharacterByName(option)
	if c == nil {
		io.WriteString(s, cfmt.Sprintf("{{Character not found.}}::red\n"))

		return StateEnterGame, nil
	}

	// // TODO: Remove me when we have attributes set
	// // c.Attributes = NewAttributes()
	// c.Attributes = &Attributes{
	// 	Body:      Attribute[int]{Name: "Body", Base: 5},
	// 	Agility:   Attribute[int]{Name: "Agility", Base: 6},
	// 	Reaction:  Attribute[int]{Name: "Reaction", Base: 4},
	// 	Strength:  Attribute[int]{Name: "Strength", Base: 5},
	// 	Willpower: Attribute[int]{Name: "Willpower", Base: 4},
	// 	Logic:     Attribute[int]{Name: "Logic", Base: 4},
	// 	Intuition: Attribute[int]{Name: "Intuition", Base: 5},
	// 	Charisma:  Attribute[int]{Name: "Charisma", Base: 4},
	// 	Edge:      Attribute[int]{Name: "Edge", Base: 5},
	// 	Essence:   Attribute[float64]{Name: "Essence", Base: 5.6},
	// 	Magic:     Attribute[int]{Name: "Magic", Base: 0},
	// 	Resonance: Attribute[int]{Name: "Resonance", Base: 0},
	// }
	// c.Save()

	c.Conn = s

	// If the character has no room, set the starting room
	if c.RoomID == "" {
		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
	}

	// If the room is not found, set the starting room
	c.Room = EntityMgr.GetRoom(c.RoomID)
	if c.Room == nil {
		slog.Error("Room not found",
			slog.String("room_id", c.RoomID))
		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
	}

	u.Save()
	c.Save()

	CharacterMgr.SetCharacterOnline(c)

	return StateGameLoop, c
}

func promptGameLoop(s ssh.Session, u *Account, c *Character) string {
	slog.Debug("Game loop state",
		slog.String("username", u.Username),
		slog.String("character_name", c.Name),
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	// Add our character to the room
	c.Room.AddCharacter(c)
	// Render the room on initial entry to the game loop
	io.WriteString(s, RenderRoom(u, c, c.Room))

	for {
		input, err := PromptForInput(s, cfmt.Sprint("{{>}}::white|bold "))
		if err != nil {
			slog.Error("Error reading input", slog.Any("error", err))
			return StateExitGame
		}

		if input == "" {
			continue
		}

		if strings.EqualFold(input, "quit") {
			return StateExitGame
		}

		CommandMgr.ParseAndExecute(s, input, u, c, c.Room)
	}
}

func promptExitGame(s ssh.Session, u *Account, c *Character) string {
	slog.Debug("Exit game state",
		slog.String("username", u.Username),
		slog.String("character_name", c.Name),
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	c.Room.Broadcast(cfmt.Sprintf("\n%s leaves the game.\n", c.Name), []string{c.ID})
	io.WriteString(s, cfmt.Sprintf("{{Goodbye, %s!}}::green\n", u.Username))

	CharacterMgr.SetCharacterOffline(c)

	c = nil

	return StateMainMenu
}
