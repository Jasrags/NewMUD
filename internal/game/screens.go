package game

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gliderlabs/ssh"
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

func PromptWelcome(s ssh.Session) string {
	var output strings.Builder

	output.WriteString("{{  ::::::::  :::    :::     :::     :::::::::   ::::::::  :::       ::: ::::    ::::  :::    ::: :::::::::  }}::#ff8700" + CRLF)
	output.WriteString("{{ :+:    :+: :+:    :+:   :+: :+:   :+:    :+: :+:    :+: :+:       :+: +:+:+: :+:+:+ :+:    :+: :+:    :+: }}::#ff5f00" + CRLF)
	output.WriteString("{{ +:+        +:+    +:+  +:+   +:+  +:+    +:+ +:+    +:+ +:+       +:+ +:+ +:+:+ +:+ +:+    +:+ +:+    +:+ }}::#ff0000" + CRLF)
	output.WriteString("{{ +#++:++#++ +#++:++#++ +#++:++#++: +#+    +:+ +#+    +:+ +#+  +:+  +#+ +#+  +:+  +#+ +#+    +:+ +#+    +:+ }}::#d70000" + CRLF)
	output.WriteString("{{        +#+ +#+    +#+ +#+     +#+ +#+    +#+ +#+    +#+ +#+ +#+#+ +#+ +#+       +#+ +#+    +#+ +#+    +#+ }}::#af0000" + CRLF)
	output.WriteString("{{ #+#    #+# #+#    #+# #+#     #+# #+#    #+# #+#    #+#  #+#+# #+#+#  #+#       #+# #+#    #+# #+#    #+# }}::#870000" + CRLF)
	output.WriteString("{{  ########  ###    ### ###     ### #########   ########    ###   ###   ###       ###  ########  #########  }}::#5f0000" + CRLF)

	if !viper.GetBool("server.login_enabled") {
		output.WriteString(cfmt.Sprint("{{Login is disabled.}}::red" + CRLF))
	}
	WriteString(s, output.String())

	PressEnterPrompt(s, "{{Press enter to continue...}}::white|bold")

	return StateLogin
}

func PromptLogin(s ssh.Session) (string, *Account) {
promptUsername:
	WriteString(s, "{{Enter your username to continue or type}}::white {{new}}::green|bold {{to register:}}::white"+CRLF)

	WriteString(s, "{{Username:}}::white|bold ")
	username, err := InputPrompt(s, "")
	if err != nil {
		return StateError, nil
	}

	// Handle "new" user registration
	if strings.EqualFold(username, "new") {
		return StateRegistration, nil
	}

	WriteString(s, "{{Password:}}::white|bold ")
	password, err := PasswordPrompt(s, "")
	if err != nil {
		return StateError, nil
	}

	u := AccountMgr.GetByUsername(username)
	if u == nil || !u.CheckPassword(password) {
		// Log and display error
		slog.Warn("Invalid login attempt",
			slog.String("username", username))

		WriteString(s, "{{Invalid username or password.}}::red"+CRLF)
		goto promptUsername
	}

	// TODO: Check if user is already logged in
	// TODO: Check if user is banned

	// Login successful
	WriteStringF(s, "{{Welcome back, %s!}}::green|bold"+CRLF, username)

	return StateMainMenu, u
}

func PromptRegistration(s ssh.Session) (string, *Account) {
	slog.Debug("Registration state",
		slog.String("remote_address", s.RemoteAddr().String()),
		slog.String("session_id", s.Context().SessionID()))

	if !viper.GetBool("server.registration_enabled") {
		WriteString(s, "\n{{Registration is disabled.}}::red"+CRLF)

		return StateLogin, nil
	}

	WriteString(s, "{{User registration}}::green"+CRLF)

promptUsername:
	username, err := InputPrompt(s, cfmt.Sprint("{{Enter your username: }}::white|bold"))
	if err != nil {
		return StateError, nil
	}

	// Validate username
	if err := ValidateCharacterName(username); err != nil {
		slog.Error("Invalid username", slog.Any("error", err))
		WriteString(s, cfmt.Sprintf("{{Invalid username: %s}}::red"+CRLF, err.Error()))

		goto promptUsername
	}

promptPassword:
	password, err := PasswordPrompt(s, cfmt.Sprint("{{Enter your password:}}::white|bold "))
	if err != nil {
		return StateError, nil
	}

	if err := ValidatePassword(password); err != nil {
		slog.Error("Invalid password", slog.Any("error", err))
		WriteString(s, cfmt.Sprintf("{{Invalid password: %s}}::red"+CRLF, err.Error()))

		goto promptPassword
	}

	confirmPassword, err := PasswordPrompt(s, cfmt.Sprint("{{Confirm your password:}}::white|bold "))
	if err != nil {

		return StateError, nil
	}

	// Check if confirm password is empty
	if confirmPassword == "" {
		WriteString(s, "{{Password cannot be empty.}}::red"+CRLF)

		goto promptPassword
	}

	// Check if passwords match
	if password != confirmPassword {
		WriteString(s, "{{Passwords do not match.}}::red"+CRLF)

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

func PromptMainMenu(s ssh.Session, a *Account) string {
	options := []string{"Enter Game", "Create Character", "Change Password", "Quit"}

	option, err := PromptForMenu(s, "Main Menu", options)
	if err != nil {
		return StateError
	}

	// Handle menu selection
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

	return StateMainMenu
}

func PromptChangePassword(s ssh.Session, a *Account) string {
	password, err := PasswordPrompt(s, cfmt.Sprint("{{Enter your current password:}}::white|bold "))
	if err != nil {
		return StateError
	}

	if !a.CheckPassword(password) {
		WriteString(s, "{{Invalid password.}}::red"+CRLF)
		return StateChangePassword
	}

	newPassword, err := PasswordPrompt(s, cfmt.Sprint("{{Enter your new password:}}::white|bold "))
	if err != nil {
		return StateError
	}

	confirmNewPassword, err := PasswordPrompt(s, cfmt.Sprint("{{Confirm your new password:}}::white|bold "))
	if err != nil {
		return StateError
	}

	if newPassword != confirmNewPassword {
		WriteString(s, "{{Passwords do not match.}}::red"+CRLF)
		return StateChangePassword
	}

	a.SetPassword(newPassword)
	a.Save()

	WriteString(s, "{{Password changed successfully.}}::green"+CRLF)

	return StateMainMenu
}

// Code prequisites: for character creation
// Hook up PromptCharacterCreate
// 1. Metatypes defined and loaded
// 2. Archetypes defined and loaded
// 3. Item packs defined and loaded (use fake items for now)

// --- Character creation steps ---
// Step 1: Prompt for character name
// Validate name (not empty, not already taken, length within limits, alphanumeric)

// Step 2: Prompt for metatype
// Display metatype options
// Allow showing details for the metatype including suggested archtypes

// Step 3: Prompt for archtype
// Display archtype options
// Allow showing details for the archtype (Highlight good/neutral/bad metatype choices for the selected archtype)

// Step 4: Prompt for item pack purchase (Optional)
// Set a base nuyen level for the character
// Display item pack options
// Allow showing details for the item pack
// Select item pack and adjust nuyen

// Step 5: Build the character
// Apply base metatype attributes (min/max)
// Apply base archtype attributes adjust within min/max if needed
// Apply any metatype qualties
// Add any item pack items to the inventory

// --- Future functions ---
// Item type support for shadowrun item types (weapons, armor, etc)
// PromptCharacterDelete
// Finish DoStats now that we have a better character definition
//

func PromptCharacterCreate(s ssh.Session, a *Account) string {
promptEnterCharacterName:
	// Step 1: Prompt for character name
	WriteString(s, "{{Enter your character's name:}}::cyan ")
	name, err := InputPrompt(s, "")
	if err != nil {
		slog.Error("Error reading character name", slog.Any("error", err))
		WriteString(s, "{{Error reading input. Returning to main menu.}}::red"+CRLF)

		return StateMainMenu
	}
	name = strings.TrimSpace(name)

	if err := ValidateCharacterName(name); err != nil {
		slog.Error("Invalid character name", slog.Any("error", err))
		WriteString(s, cfmt.Sprintf("{{Invalid name: %s}}::red"+CRLF, err.Error()))

		goto promptEnterCharacterName
	}

	// promptEnterCharacterDescription:
	// Step 2: Prompt for character description
	// TODO: maybe move this after metatype and archtype selection
	// TODO: Once we have a archtype, metatype and other personal information we can generate a "default" short and long description that can be changed later.

	// Step 2: Prompt for metatype
	// Display metatype options
	// Allow showing details for the metatype including suggested archtypes
promptSelectMetatype:
	metatypeChoice, err := MenuPrompt(s, "Select a Metatype:", EntityMgr.GetMetatypeMenuOptions())
	if err != nil {
		slog.Error("Error selecting metatype", slog.Any("error", err))

		goto promptSelectMetatype
	}
	slog.Info("Metatype selected",
		slog.String("metatype", metatypeChoice))

	// Step 3: Prompt for archtype
	// Display archtype options
	// Allow showing details for the archtype (Highlight good/neutral/bad metatype choices for the selected archtype)
promptSelectArchetype:
	archtypeChoice, err := MenuPrompt(s, "Select a archtype:", EntityMgr.GetPregenMenuOptions())
	if err != nil {
		slog.Error("Error selecting metatype", slog.Any("error", err))

		goto promptSelectArchetype
	}
	slog.Info("Archtype selected",
		slog.String("archtype", archtypeChoice))

	// promptSelectItemPack:
	// Step 4: Prompt for item pack purchase (Optional)
	// Set a base nuyen level for the character
	// Display item pack options
	// Allow showing details for the item pack
	// Select item pack and adjust nuyen

	// Step 5: Build the character
	// Apply base metatype attributes (min/max)
	// Apply base archtype attributes adjust within min/max if needed
	// Apply any metatype qualties
	// Add any item pack items to the inventory
	// WriteString(s, "{{Setting base attributes...}}::green"+CRLF)
	// baseAttributes := Attributes{
	// 	Body:      Attribute[int]{Name: "Body", Base: 5},
	// 	Agility:   Attribute[int]{Name: "Agility", Base: 6},
	// 	Reaction:  Attribute[int]{Name: "Reaction", Base: 4},
	// 	Strength:  Attribute[int]{Name: "Strength", Base: 5},
	// 	Willpower: Attribute[int]{Name: "Willpower", Base: 4},
	// 	Logic:     Attribute[int]{Name: "Logic", Base: 4},
	// 	Intuition: Attribute[int]{Name: "Intuition", Base: 5},
	// 	Charisma:  Attribute[int]{Name: "Charisma", Base: 4},
	// 	Essence:   Attribute[float64]{Name: "Essence", Base: 5.6},
	// 	Magic:     Attribute[int]{Name: "Magic", Base: 0},
	// 	Resonance: Attribute[int]{Name: "Resonance", Base: 0},
	// }

	// Step 4: Create the character
	char := NewCharacter()
	char.Prompt = DefaultPrompt
	char.Role = CharacterRolePlayer
	char.AccountID = a.ID
	char.Name = name
	char.Title = "The Brave"         // TODO: Set this from our character creation
	char.Description = "Description" // TODO: Generate this from descriptive character data
	char.Metatype = "Human"          // TODO: Set this from our character creation
	char.Age = 25                    // TODO: Set this from our character creation
	char.Sex = "Male"                // TODO: Set this from our character creation
	char.Height = 180                // TODO: Set this from our character creation
	char.Weight = 75                 // TODO: Set this from our character creation
	char.Ethnicity = "Caucasian"     // TODO: Set this from our character creation
	// StreetCred      int              `yaml:"street_cred"`
	// Notoriety       int              `yaml:"notoriety"`
	// PublicAwareness int              `yaml:"public_awareness"`
	// Karma           int              `yaml:"karma"`
	// TotalKarma      int              `yaml:"total_karma"`
	// Attributes      Attributes       `yaml:"attributes"`
	char.Attributes = NewAttributes()
	// PhysicalDamage  PhysicalDamage   `yaml:"physical_damage"`
	// StunDamage      StunDamage       `yaml:"stun_damage"`
	// Edge            Edge             `yaml:"edge"`
	// Room            *Room            `yaml:"-"`
	// RoomID          string           `yaml:"room_id"`
	// Area            *Area            `yaml:"-"`
	// AreaID          string           `yaml:"area_id"`
	// Inventory       Inventory        `yaml:"inventory"`
	// Equipment       map[string]*Item `yaml:"equipment"`
	// Qualtities      []Quality        `yaml:"qualities"`
	// Skills          []Skill          `yaml:"skills"`

	// char := &Character{
	// 	GameEntity: GameEntity{
	// 		ID:          uuid.New().String(),
	// 		Name:        name,
	// 		Description: description,
	// 		Attributes:  baseAttributes,
	// 		Equipment:   make(map[string]*Item),
	// 		Edge:        Edge{Max: 5, Available: 5},
	// 	},
	// 	UserID:    a.ID,
	// 	Role:      CharacterRolePlayer,
	// 	CreatedAt: time.Now(),
	// }
	// char.Save()

	// Step 5: Add character to user
	a.Characters = append(a.Characters, char.Name)
	// a.Save()

	// Step 6: Save user
	// err = UserMgr.SaveUser(u)
	// if err != nil {
	// 	slog.Error("Error saving user after character creation", slog.Any("error", err))
	// 	WriteString(s, "{{Error saving character. Returning to main menu.}}::red"+CRLF)
	// 	return StateMainMenu
	// }

	// Step 7: Confirmation and return to main menu
	WriteStringF(s, "{{Character '%s' created successfully! Returning to main menu.}}::green"+CRLF, name)

	return StateMainMenu
}

func PromptEnterGame(s ssh.Session, a *Account) (string, *Character) {
	// Check if user has characters
	if len(a.Characters) == 0 {
		WriteString(s, "{{You have no characters. Create one to start playing.}}::red"+CRLF)

		return StateEnterGame, nil
	}

	// Collect available characters
	var characters []string
	for _, name := range a.Characters {
		char := CharacterMgr.GetCharacterByName(name)
		if char == nil {
			continue
		}
		characters = append(characters, char.Name) // No need to style names here; handled by PromptForMenu
	}

	// Use PromptForMenu to render the character selection menu
	option, err := PromptForMenu(s, "Select a character:", characters)
	if err != nil {
		WriteString(s, "{{An error occurred while selecting a character.}}::red"+CRLF)

		return StateError, nil
	}

	// Load the selected character
	c := CharacterMgr.GetCharacterByName(option)
	if c == nil {
		WriteString(s, "{{Character not found. Please try again.}}::red"+CRLF)

		return StateEnterGame, nil
	}

	c.Conn = s

	// Handle the character's room
	if c.RoomID == "" {
		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
	}
	c.Room = EntityMgr.GetRoom(c.RoomID)
	if c.Room == nil {
		slog.Error("Room not found",
			slog.String("room_id", c.RoomID))

		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
	}

	// Save account and character states
	a.Save()
	c.Save()

	// Set the character as online
	CharacterMgr.SetCharacterOnline(c)

	// Notify the player and enter the game loop
	WriteString(s, lipgloss.JoinVertical(lipgloss.Left,
		cfmt.Sprintf("{{Entering the game as %s...}}::green|bold"+CRLF, c.Name),
	))

	return StateGameLoop, c
}

func PromptGameLoop(s ssh.Session, a *Account, c *Character) string {
	// Add our character to the room
	c.Room.AddCharacter(c)
	// Render the room on initial entry to the game loop
	WriteString(s, RenderRoom(a, c, c.Room))
	WriteString(s, ""+CRLF)

	for {
		WriteStringF(s, "{{%s}}::white|bold ", RenderPrompt(c))
		input, err := InputPrompt(s, "")
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

		CommandMgr.ParseAndExecute(s, input, a, c, c.Room)
	}
}

func PromptExitGame(s ssh.Session, a *Account, c *Character) string {
	c.Room.Broadcast(cfmt.Sprintf("%s leaves the game."+CRLF, c.Name), []string{c.ID})
	WriteStringF(s, "{{Goodbye, %s!}}::green"+CRLF, a.Username)

	CharacterMgr.SetCharacterOffline(c)

	c = nil

	return StateMainMenu
}
