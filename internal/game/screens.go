package game

import (
	"log/slog"
	"strings"

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
	for {
		// Prompt for username or registration.
		WriteString(s, "{{Enter your username to continue or type}}::white {{new}}::green|bold {{to register:}}::white"+CRLF)
		WriteString(s, "{{Username:}}::white|bold ")

		username, err := InputPrompt(s, "")
		if err != nil {
			return StateError, nil
		}
		username = strings.TrimSpace(username)

		// Handle "new" user registration.
		if strings.EqualFold(username, "new") {
			return StateRegistration, nil
		}

		// Prompt for password.
		WriteString(s, "{{Password:}}::white|bold ")
		password, err := PasswordPrompt(s, "")
		if err != nil {
			return StateError, nil
		}

		// Check user credentials.
		u := AccountMgr.GetByUsername(username)
		if u == nil || !u.CheckPassword(password) {
			slog.Warn("Invalid login attempt", slog.String("username", username))
			WriteString(s, "{{Invalid username or password.}}::red"+CRLF)
			continue // Retry the login prompt.
		}

		// TODO: Check if the user is already logged in.
		// TODO: Check if the user is banned.

		// Login successful.
		WriteStringF(s, "{{Welcome back, %s!}}::green|bold"+CRLF, username)
		return StateMainMenu, u
	}
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

	// Prompt for a valid username.
	var username string
	for {
		var err error
		username, err = InputPrompt(s, cfmt.Sprint("{{Enter your username: }}::white|bold"))
		if err != nil {
			return StateError, nil
		}

		// Validate the username.
		if err := ValidateCharacterName(username); err != nil {
			slog.Error("Invalid username", slog.Any("error", err))
			WriteString(s, cfmt.Sprintf("{{Invalid username: %s}}::red"+CRLF, err.Error()))
			continue
		}

		break // Valid username entered.
	}

	// Prompt for a valid password and confirmation.
	var password string
	for {
		var err error
		password, err = PasswordPrompt(s, cfmt.Sprint("{{Enter your password:}}::white|bold "))
		if err != nil {
			return StateError, nil
		}

		if err := ValidatePassword(password); err != nil {
			slog.Error("Invalid password", slog.Any("error", err))
			WriteString(s, cfmt.Sprintf("{{Invalid password: %s}}::red"+CRLF, err.Error()))
			continue
		}

		confirmPassword, err := PasswordPrompt(s, cfmt.Sprint("{{Confirm your password:}}::white|bold "))
		if err != nil {
			return StateError, nil
		}

		if confirmPassword == "" {
			WriteString(s, "{{Password cannot be empty.}}::red"+CRLF)
			continue
		}

		if password != confirmPassword {
			WriteString(s, "{{Passwords do not match.}}::red"+CRLF)
			continue
		}

		break // Passwords are valid and match.
	}

	// Create a new user account.
	u := NewAccount()
	u.Username = username
	u.SetPassword(password)
	u.Save()
	AccountMgr.AddAccount(u)

	return StateMainMenu, u
}

func PromptMainMenu(s ssh.Session, a *Account) string {
	options := []MenuOption{
		{"Enter Game", "enter_game", "Enter Game"},
		{"Create Character", "create_character", "Create Character"},
		{"Change Password", "change_password", "Change your password"},
		{"Quit", "quit", "Exit the game"},
	}

	for {
		option, err := PromptForMenu(s, "Main Menu", options)
		if err != nil {
			slog.Error("Error prompting for menu", slog.Any("error", err))
			return StateError
		}

		slog.Debug("Main menu state",
			slog.String("option", option))

		// Handle menu selection
		switch option {
		case "enter_game":
			return StateEnterGame
		case "create_character":
			return StateCharacterCreate
		case "change_password":
			return StateChangePassword
		case "quit":
			return StateQuit
		}
	}
}

func PromptChangePassword(s ssh.Session, a *Account) string {
	for {
		// Prompt for the current password.
		currentPassword, err := PasswordPrompt(s, cfmt.Sprint("{{Enter your current password:}}::white|bold "))
		if err != nil {
			return StateError
		}

		if !a.CheckPassword(currentPassword) {
			WriteString(s, "{{Invalid password.}}::red"+CRLF)
			// Re-prompt for the current password.
			continue
		}

		// Prompt for the new password until it is confirmed correctly.
		var newPassword string
		for {
			newPassword, err = PasswordPrompt(s, cfmt.Sprint("{{Enter your new password:}}::white|bold "))
			if err != nil {
				return StateError
			}

			confirmNewPassword, err := PasswordPrompt(s, cfmt.Sprint("{{Confirm your new password:}}::white|bold "))
			if err != nil {
				return StateError
			}

			if newPassword != confirmNewPassword {
				WriteString(s, "{{Passwords do not match.}}::red"+CRLF)
				// Re-prompt for the new password.
				continue
			}

			// The new password is valid and confirmed.
			break
		}

		// Update and save the new password.
		a.SetPassword(newPassword)
		a.Save()

		WriteString(s, "{{Password changed successfully.}}::green"+CRLF)
		return StateMainMenu
	}
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

// func PromptCharacterCreate(s ssh.Session, a *Account) string {
// promptEnterCharacterName:

// 	// Check if they already have the max number of characters per account
// 	maxCharacterCount := viper.GetInt("server.max_character_count")
// 	if len(a.Characters) >= maxCharacterCount {
// 		noun := pluralizer.PluralizeNounPhrase("character", len(a.Characters))
// 		PressEnterPrompt(s, cfmt.Sprintf("{{You already have %s. You cannot create any more.}}::red", noun))
// 		PressEnterPrompt(s, "{{Press enter to continue...}}::white|bold")

// 		return StateMainMenu
// 	}

// 	// Step 1: Prompt for character name
// 	WriteString(s, "{{Enter your character's name:}}::cyan ")
// 	name, err := InputPrompt(s, "")
// 	if err != nil {
// 		slog.Error("Error reading character name", slog.Any("error", err))
// 		WriteString(s, "{{Error reading input. Returning to main menu.}}::red"+CRLF)

// 		return StateMainMenu
// 	}
// 	name = strings.TrimSpace(name)

// 	if err := ValidateCharacterName(name); err != nil {
// 		slog.Error("Invalid character name", slog.Any("error", err))
// 		WriteString(s, cfmt.Sprintf("{{Invalid name: %s}}::red"+CRLF, err.Error()))

// 		goto promptEnterCharacterName
// 	}

// 	// promptEnterCharacterDescription:
// 	// Step 2: Prompt for character description
// 	// TODO: maybe move this after metatype and archtype selection
// 	// TODO: Once we have a archtype, metatype and other personal information we can generate a "default" short and long description that can be changed later.

// 	// Step 2: Prompt for metatype
// 	// Display metatype options
// 	// Allow showing details for the metatype including suggested archtypes
// promptSelectMetatype:
// 	metatypeChoice, err := MenuPrompt(s, "Select a Metatype:", EntityMgr.GetMetatypeMenuOptions())
// 	if err != nil {
// 		slog.Error("Error selecting metatype", slog.Any("error", err))

// 		goto promptSelectMetatype
// 	}
// 	slog.Info("Metatype selected",
// 		slog.String("metatype", metatypeChoice))

// 	// Step 3: Prompt for archtype
// 	// Display archtype options
// 	// Allow showing details for the archtype (Highlight good/neutral/bad metatype choices for the selected archtype)
// promptSelectArchetype:
// 	archtypeChoice, err := MenuPrompt(s, "Select a archtype:", EntityMgr.GetPregenMenuOptions())
// 	if err != nil {
// 		slog.Error("Error selecting metatype", slog.Any("error", err))

// 		goto promptSelectArchetype
// 	}
// 	slog.Info("Archtype selected",
// 		slog.String("archtype", archtypeChoice))

// 	// promptSelectItemPack:
// 	// Step 4: Prompt for item pack purchase (Optional)
// 	// Set a base nuyen level for the character
// 	// Display item pack options
// 	// Allow showing details for the item pack
// 	// Select item pack and adjust nuyen

// 	// Step 5: Build the character
// 	// Apply base metatype attributes (min/max)
// 	// Apply base archtype attributes adjust within min/max if needed
// 	// Apply any metatype qualties
// 	// Add any item pack items to the inventory

// 	// Step 4: Create the character
// 	metatype := EntityMgr.GetMetatype(metatypeChoice)
// 	archtype := EntityMgr.GetPregen(archtypeChoice)

// 	// slog.Debug("Creating character",
// 	// 	slog.Any("metatype", metatype.Attributes),
// 	// 	slog.Any("archtype", archtype.Attributes))
// 	char := NewCharacter()
// 	char.AccountID = a.ID
// 	char.Name = name
// 	char.Title = "The Brave"         // TODO: Set this from our character creation
// 	char.Description = "Description" // TODO: Generate this from descriptive character data
// 	char.Metatype = "Human"          // TODO: Set this from our character creation
// 	char.Age = 25                    // TODO: Set this from our character creation
// 	char.Sex = "Male"                // TODO: Set this from our character creation
// 	char.Height = 180                // TODO: Set this from our character creation
// 	char.Weight = 75                 // TODO: Set this from our character creation
// 	char.Ethnicity = "Caucasian"     // TODO: Set this from our character creation
// 	// StreetCred      int              `yaml:"street_cred"`
// 	// Notoriety       int              `yaml:"notoriety"`
// 	// PublicAwareness int              `yaml:"public_awareness"`
// 	// Karma           int              `yaml:"karma"`
// 	// TotalKarma      int              `yaml:"total_karma"`
// 	// Attributes      Attributes       `yaml:"attributes"`
// 	char.Attributes.Body.Base = archtype.Attributes.Body.Base
// 	char.Attributes.Body.Min = metatype.Attributes.Body.Min
// 	char.Attributes.Body.Max = metatype.Attributes.Body.Max
// 	char.Attributes.Body.AugMax = metatype.Attributes.Body.AugMax

// 	char.Attributes.Agility.Base = archtype.Attributes.Agility.Base

// 	char.Attributes.Reaction.Base = archtype.Attributes.Reaction.Base

// 	char.Attributes.Strength.Base = archtype.Attributes.Strength.Base

// 	char.Attributes.Willpower.Base = archtype.Attributes.Willpower.Base

// 	char.Attributes.Logic.Base = archtype.Attributes.Logic.Base

// 	char.Attributes.Intuition.Base = archtype.Attributes.Intuition.Base

// 	char.Attributes.Charisma.Base = archtype.Attributes.Charisma.Base

// 	// Edge
// 	// Initiative

// 	char.Attributes.Essence.Base = archtype.Attributes.Essence.Base

// 	char.Attributes.Magic.Base = archtype.Attributes.Magic.Base

// 	char.Attributes.Resonance.Base = archtype.Attributes.Resonance.Base

// 	// char.Attributes.Body = Attribute[int]{
// 	// 	Name:   "Body",
// 	// 	Base:   archtype.Attributes.Body.Base,
// 	// 	Min:    metatype.Attributes.Body.Min,
// 	// 	Max:    metatype.Attributes.Body.Max,
// 	// 	AugMax: metatype.Attributes.Body.AugMax,
// 	// }
// 	// char.Attributes.Agility = Attribute[int]{
// 	// 	Name:   "Agility",
// 	// 	Base:   archtype.Attributes.Agility.Base,
// 	// 	Min:    metatype.Attributes.Agility.Min,
// 	// 	Max:    metatype.Attributes.Agility.Max,
// 	// 	AugMax: metatype.Attributes.Agility.AugMax,
// 	// }
// 	// char.Attributes.Reaction = Attribute[int]{
// 	// 	Name:   "Reaction",
// 	// 	Base:   archtype.Attributes.Reaction.Base,
// 	// 	Min:    metatype.Attributes.Reaction.Min,
// 	// 	Max:    metatype.Attributes.Reaction.Max,
// 	// 	AugMax: metatype.Attributes.Reaction.AugMax,
// 	// }
// 	// char.Attributes.Strength = Attribute[int]{
// 	// 	Name:   "Strength",
// 	// 	Base:   archtype.Attributes.Strength.Base,
// 	// 	Min:    metatype.Attributes.Strength.Min,
// 	// 	Max:    metatype.Attributes.Strength.Max,
// 	// 	AugMax: metatype.Attributes.Strength.AugMax,
// 	// }
// 	// char.Attributes.Willpower = Attribute[int]{
// 	// 	Name:   "Willpower",
// 	// 	Base:   archtype.Attributes.Willpower.Base,
// 	// 	Min:    metatype.Attributes.Willpower.Min,
// 	// 	Max:    metatype.Attributes.Willpower.Max,
// 	// 	AugMax: metatype.Attributes.Willpower.AugMax,
// 	// }
// 	// char.Attributes.Logic = Attribute[int]{
// 	// 	Name:   "Logic",
// 	// 	Base:   archtype.Attributes.Logic.Base,
// 	// 	Min:    metatype.Attributes.Logic.Min,
// 	// 	Max:    metatype.Attributes.Logic.Max,
// 	// 	AugMax: metatype.Attributes.Logic.AugMax,
// 	// }
// 	// char.Attributes.Intuition = Attribute[int]{
// 	// 	Name:   "Intuition",
// 	// 	Base:   archtype.Attributes.Intuition.Base,
// 	// 	Min:    metatype.Attributes.Intuition.Min,
// 	// 	Max:    metatype.Attributes.Intuition.Max,
// 	// 	AugMax: metatype.Attributes.Intuition.AugMax,
// 	// }
// 	// char.Attributes.Charisma = Attribute[int]{
// 	// 	Name:   "Charisma",
// 	// 	Base:   archtype.Attributes.Charisma.Base,
// 	// 	Min:    metatype.Attributes.Charisma.Min,
// 	// 	Max:    metatype.Attributes.Charisma.Max,
// 	// 	AugMax: metatype.Attributes.Charisma.AugMax,
// 	// }
// 	// char.Attributes.Essence = Attribute[float64]{
// 	// 	Name:   "Essence",
// 	// 	Base:   archtype.Attributes.Essence.Base,
// 	// 	Min:    metatype.Attributes.Essence.Min,
// 	// 	Max:    metatype.Attributes.Essence.Max,
// 	// 	AugMax: metatype.Attributes.Essence.AugMax,
// 	// }
// 	// char.Attributes.Magic = Attribute[int]{
// 	// 	Name:   "Magic",
// 	// 	Base:   archtype.Attributes.Magic.Base,
// 	// 	Min:    metatype.Attributes.Magic.Min,
// 	// 	Max:    metatype.Attributes.Magic.Max,
// 	// 	AugMax: metatype.Attributes.Magic.AugMax,
// 	// }
// 	// char.Attributes.Resonance = Attribute[int]{
// 	// 	Name:   "Resonance",
// 	// 	Base:   archtype.Attributes.Resonance.Base,
// 	// 	Min:    metatype.Attributes.Resonance.Min,
// 	// 	Max:    metatype.Attributes.Resonance.Max,
// 	// 	AugMax: metatype.Attributes.Resonance.AugMax,
// 	// }

// 	// PhysicalDamage  PhysicalDamage   `yaml:"physical_damage"`
// 	// StunDamage      StunDamage       `yaml:"stun_damage"`
// 	// Edge            Edge             `yaml:"edge"`
// 	// Room            *Room            `yaml:"-"`
// 	// RoomID          string           `yaml:"room_id"`
// 	// Area            *Area            `yaml:"-"`
// 	// AreaID          string           `yaml:"area_id"`
// 	// Inventory       Inventory        `yaml:"inventory"`
// 	// Equipment       map[string]*Item `yaml:"equipment"`
// 	// Qualtities      []Quality        `yaml:"qualities"`
// 	// Skills          []Skill          `yaml:"skills"`

// 	char.Save()

// 	// Step 5: Add character to user
// 	a.Characters = append(a.Characters, char.Name)
// 	// a.Save()

// 	// Step 6: Save user
// 	// err = UserMgr.SaveUser(u)
// 	// if err != nil {
// 	// 	slog.Error("Error saving user after character creation", slog.Any("error", err))
// 	// 	WriteString(s, "{{Error saving character. Returning to main menu.}}::red"+CRLF)
// 	// 	return StateMainMenu
// 	// }

// 	// Step 7: Confirmation and return to main menu
// 	WriteStringF(s, "{{Character '%s' created successfully! Returning to main menu.}}::green"+CRLF, name)

//		return StateMainMenu
//	}
func PromptCharacterCreate(s ssh.Session, a *Account) string {
	options := []MenuOption{
		{"Choose a Pre-Generated Character", "pregen", "Select from predefined character archetypes"},
		{"Create a Custom Character", "custom", "Build a character from scratch"},
		{"Back to Main Menu", "back", "Return to the main menu"},
	}

	for {
		choice, err := PromptForMenu(s, "Character Creation", options)
		if err != nil {
			return StateError
		}

		switch choice {
		case "pregen":
			return PromptSelectPregenCharacter(s, a)
		case "custom":
			WriteString(s, "{{Custom character creation is not yet implemented. Returning to main menu.}}::yellow"+CRLF)
			return StateMainMenu
		case "back":
			return StateMainMenu
		default:
			// This branch should rarely be reached if PromptForMenu
			// validates the input, but it's good to be defensive.
			WriteString(s, "{{Invalid option. Please try again.}}::red"+CRLF)
		}
	}
}

func PromptSelectPregenCharacter(s ssh.Session, a *Account) string {
	// Build the menu options from the pre-generated characters.
	pregens := EntityMgr.GetPregens()
	options := make([]MenuOption, 0, len(pregens)+1)
	for _, pregen := range pregens {
		options = append(options, MenuOption{
			DisplayText: pregen.Title,
			Value:       pregen.ID,
			Description: pregen.GetSelectionInfo(),
		})
	}
	options = append(options, MenuOption{"Back", "back", "Return to character creation menu"})

	// Loop until a valid selection and confirmation is made.
	for {
		choice, err := PromptForMenu(s, "Select a Pre-Generated Character", options)
		if err != nil {
			return StateError
		}

		if choice == "back" {
			return StateCharacterCreate
		}

		pregen := EntityMgr.GetPregen(choice)
		if pregen == nil {
			WriteString(s, "{{Invalid selection. Try again.}}::red"+CRLF)
			continue
		}

		// Display selection details.
		WriteString(s, cfmt.Sprintf("{{%s}}::cyan"+CRLF, pregen.GetSelectionInfo()))

		// Confirm selection.
		if !YesNoPrompt(s, true) {
			WriteString(s, "{{Returning to character selection.}}::yellow"+CRLF)
			return StateCharacterCreate
		}

		// Proceed with character naming.
		return PromptSetCharacterName(s, a, pregen)
	}
}

func PromptSetCharacterName(s ssh.Session, a *Account, pregen *Pregen) string {
	for {
		// Prompt for the character's name.
		WriteString(s, "{{Enter your character's name:}}::cyan ")
		name, err := InputPrompt(s, "")
		if err != nil {
			return StateError
		}
		name = strings.TrimSpace(name)

		// Validate the character name.
		if err := ValidateCharacterName(name); err != nil {
			WriteString(s, cfmt.Sprintf("{{Invalid name: %s}}::red"+CRLF, err.Error()))
			continue
		}

		// Instantiate and save the new character.
		// This code is commented out because the actual instantiation logic may change.
		// newChar := pregen.Instantiate(name, a.ID)
		// newChar.Save()
		// a.Characters = append(a.Characters, newChar.Name)

		// Inform the user of the successful character creation.
		WriteString(s, cfmt.Sprintf("{{Character '%s' created successfully! Returning to main menu.}}::green"+CRLF, name))
		return StateMainMenu
	}
}

// func PromptEnterGame(s ssh.Session, a *Account) (string, *Character) {
// 	// Check if user has characters
// 	if len(a.Characters) == 0 {
// 		WriteString(s, "{{You have no characters. Create one to start playing.}}::red"+CRLF)

// 		return StateEnterGame, nil
// 	}

// 	// Collect available characters
// 	var characters []string
// 	for _, name := range a.Characters {
// 		char := CharacterMgr.GetCharacterByName(name)
// 		if char == nil {
// 			continue
// 		}
// 		characters = append(characters, char.Name) // No need to style names here; handled by PromptForMenu
// 	}

// 	// Use PromptForMenu to render the character selection menu
// 	option, err := PromptForMenu(s, "Select a character:", characters)
// 	if err != nil {
// 		WriteString(s, "{{An error occurred while selecting a character.}}::red"+CRLF)

// 		return StateError, nil
// 	}

// 	// Load the selected character
// 	c := CharacterMgr.GetCharacterByName(option)
// 	if c == nil {
// 		WriteString(s, "{{Character not found. Please try again.}}::red"+CRLF)

// 		return StateEnterGame, nil
// 	}

// 	c.Conn = s

// 	// Handle the character's room
// 	if c.RoomID == "" {
// 		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
// 	}
// 	c.Room = EntityMgr.GetRoom(c.RoomID)
// 	if c.Room == nil {
// 		slog.Error("Room not found",
// 			slog.String("room_id", c.RoomID))

// 		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
// 	}

// 	// Save account and character states
// 	a.Save()
// 	c.Save()

// 	// Set the character as online
// 	CharacterMgr.SetCharacterOnline(c)

// 	// Notify the player and enter the game loop
// 	WriteString(s, lipgloss.JoinVertical(lipgloss.Left,
// 		cfmt.Sprintf("{{Entering the game as %s...}}::green|bold"+CRLF, c.Name),
// 	))

// 	return StateGameLoop, c
// }

// func PromptEnterGame(s ssh.Session, a *Account) (string, *Character) {
// 	if len(a.Characters) == 0 {
// 		WriteString(s, "{{You have no characters. Create one to start playing.}}::red"+CRLF)
// 		return StateEnterGame, nil
// 	}

// 	options := make([]MenuOption, len(a.Characters)+1)
// 	for i, name := range a.Characters {
// 		options[i] = MenuOption{name, name, "Select this character"}
// 	}
// 	options[len(a.Characters)] = MenuOption{"Back to Main Menu", "back", "Return to main menu"}

// 	choice, err := PromptForMenu(s, "Select a Character", options)
// 	if err != nil {
// 		return StateError, nil
// 	}

// 	if choice == "back" {
// 		return StateMainMenu, nil
// 	}

// 	c := CharacterMgr.GetCharacterByName(choice)
// 	if c == nil {
// 		WriteString(s, "{{Character not found. Please try again.}}::red"+CRLF)
// 		return StateEnterGame, nil
// 	}

// 	c.Conn = s
// 	if c.RoomID == "" {
// 		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
// 	}
// 	c.Room = EntityMgr.GetRoom(c.RoomID)
// 	if c.Room == nil {
// 		c.SetRoom(EntityMgr.GetRoom(viper.GetString("server.starting_room")))
// 	}

// 	a.Save()
// 	c.Save()
// 	CharacterMgr.SetCharacterOnline(c)

// 	WriteString(s, cfmt.Sprintf("{{Entering the game as %s...}}::green|bold"+CRLF, c.Name))
// 	return StateGameLoop, c
// }

func PromptEnterGame(s ssh.Session, a *Account) (string, *Character) {
	// Ensure the account has at least one character.
	if len(a.Characters) == 0 {
		WriteString(s, "{{You have no characters. Create one to start playing.}}::red"+CRLF)
		return StateEnterGame, nil
	}

	// Build menu options for each character and an option to return to the main menu.
	options := make([]MenuOption, len(a.Characters)+1)
	for i, name := range a.Characters {
		options[i] = MenuOption{
			DisplayText: name,
			Value:       name,
			Description: "Select this character",
		}
	}
	options[len(a.Characters)] = MenuOption{
		DisplayText: "Back to Main Menu",
		Value:       "back",
		Description: "Return to main menu",
	}

	// Loop until a valid character is selected.
	for {
		choice, err := PromptForMenu(s, "Select a Character", options)
		if err != nil {
			return StateError, nil
		}

		if choice == "back" {
			return StateMainMenu, nil
		}

		c := CharacterMgr.GetCharacterByName(choice)
		if c == nil {
			WriteString(s, "{{Character not found. Please try again.}}::red"+CRLF)
			continue
		}

		// Set the session connection for the character.
		c.Conn = s

		// Retrieve the starting room ID from configuration.
		startingRoomID := viper.GetString("server.starting_room")

		// Check if the character's current room is valid; otherwise, assign the starting room.
		if c.RoomID == "" || EntityMgr.GetRoom(c.RoomID) == nil {
			startingRoom := EntityMgr.GetRoom(startingRoomID)
			c.SetRoom(startingRoom)
			c.Room = startingRoom
		} else {
			c.Room = EntityMgr.GetRoom(c.RoomID)
		}

		// Save any changes to the account and character, and mark the character as online.
		a.Save()
		c.Save()
		CharacterMgr.SetCharacterOnline(c)

		// Notify the user and proceed into the game.
		WriteString(s, cfmt.Sprintf("{{Entering the game as %s...}}::green|bold"+CRLF, c.Name))
		return StateGameLoop, c
	}
}

func PromptGameLoop(s ssh.Session, a *Account, c *Character) string {
	// Add the character to their current room.
	c.Room.AddCharacter(c)

	// Render the room on initial entry.
	WriteString(s, RenderRoom(a, c, c.Room))
	WriteString(s, CRLF)

	// Game loop: repeatedly prompt the user for input.
	for {
		// Render and display the prompt.
		prompt := RenderPrompt(c)
		WriteStringF(s, "{{%s}}::white|bold ", prompt)

		// Read user input.
		input, err := InputPrompt(s, "")
		if err != nil {
			slog.Error("Error reading input", slog.Any("error", err))
			return StateExitGame
		}

		// Trim extra whitespace.
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Check for the "quit" command.
		if strings.EqualFold(input, "quit") {
			return StateExitGame
		}

		// Parse and execute the entered command.
		CommandMgr.ParseAndExecute(s, input, a, c, c.Room)
	}
}

func PromptExitGame(s ssh.Session, a *Account, c *Character) string {
	// Broadcast that the character is leaving the game.
	exitMessage := cfmt.Sprintf("%s leaves the game."+CRLF, c.Name)
	c.Room.Broadcast(exitMessage, []string{c.ID})

	// Send a goodbye message to the user.
	WriteStringF(s, "{{Goodbye, %s!}}::green"+CRLF, a.Username)

	// Mark the character as offline.
	CharacterMgr.SetCharacterOffline(c)

	return StateMainMenu
}
