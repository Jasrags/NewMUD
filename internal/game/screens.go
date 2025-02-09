package game

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
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
	StateCharacterDelete = "character_delete"
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
		{"Delete Character", "delete_character", "Delete a character"},
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
		case "delete_character":
			return StateCharacterDelete
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
func PromptCharacterCreate(s ssh.Session, a *Account) (string, *Character) {
	options := []MenuOption{
		{"Choose a Pre-Generated Character", "pregen", "Select from predefined character archetypes"},
		{"Create a Custom Character", "custom", "Build a character from scratch"},
		{"Back to Main Menu", "back", "Return to the main menu"},
	}

	for {
		choice, err := PromptForMenu(s, "Character Creation", options)
		if err != nil {
			return StateError, nil
		}

		switch choice {
		case "pregen":
			return PromptSelectPregenCharacter(s, a)
		case "custom":
			WriteString(s, "{{Custom character creation is not yet implemented. Returning to main menu.}}::yellow"+CRLF)
			return StateMainMenu, nil
		case "back":
			return StateMainMenu, nil
		default:
			// This branch should rarely be reached if PromptForMenu
			// validates the input, but it's good to be defensive.
			WriteString(s, "{{Invalid option. Please try again.}}::red"+CRLF)
		}
	}
}

func PromptSelectPregenCharacter(s ssh.Session, a *Account) (string, *Character) {
	// Build menu options
	pregens := EntityMgr.GetPregens()
	pregenMap := make(map[string]*Pregen)
	options := make([]MenuOption, 0, len(pregens)+1)

	for _, pregen := range pregens {
		pregenMap[pregen.ID] = pregen
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
			slog.Error("Error prompting for pregen selection", slog.Any("error", err))
			return StateError, nil
		}

		if choice == "back" {
			return StateCharacterCreate, nil
		}

		pregen, exists := pregenMap[choice]
		if !exists {
			WriteString(s, "{{Invalid selection. Try again.}}::red"+CRLF)
			continue
		}

		// Display selection details.
		WriteString(s, cfmt.Sprintf("{{%s}}::cyan"+CRLF, pregen.GetSelectionInfo()))

		// Confirm selection.
		if !YesNoPrompt(s, true) {
			WriteString(s, "{{Returning to character selection.}}::yellow"+CRLF)
			continue // Instead of returning, continue prompting
		}

		metatype := EntityMgr.GetMetatype(pregen.MetatypeID)

		char := NewCharacter()
		char.AccountID = a.ID
		char.MetatypeID = pregen.MetatypeID
		char.Body = Attribute[int]{Name: "Body", Base: pregen.Body.Base, Min: metatype.Body.Min, Max: metatype.Body.Max, AugMax: metatype.Body.AugMax}
		char.Agility = Attribute[int]{Name: "Agility", Base: pregen.Agility.Base, Min: metatype.Agility.Min, Max: metatype.Agility.Max, AugMax: metatype.Agility.AugMax}
		char.Reaction = Attribute[int]{Name: "Reaction", Base: pregen.Reaction.Base, Min: metatype.Reaction.Min, Max: metatype.Reaction.Max, AugMax: metatype.Reaction.AugMax}
		char.Strength = Attribute[int]{Name: "Strength", Base: pregen.Strength.Base, Min: metatype.Strength.Min, Max: metatype.Strength.Max, AugMax: metatype.Strength.AugMax}
		char.Willpower = Attribute[int]{Name: "Willpower", Base: pregen.Willpower.Base, Min: metatype.Willpower.Min, Max: metatype.Willpower.Max, AugMax: metatype.Willpower.AugMax}
		char.Logic = Attribute[int]{Name: "Logic", Base: pregen.Logic.Base, Min: metatype.Logic.Min, Max: metatype.Logic.Max, AugMax: metatype.Logic.AugMax}
		char.Intuition = Attribute[int]{Name: "Intuition", Base: pregen.Intuition.Base, Min: metatype.Intuition.Min, Max: metatype.Intuition.Max, AugMax: metatype.Intuition.AugMax}
		char.Charisma = Attribute[int]{Name: "Charisma", Base: pregen.Charisma.Base, Min: metatype.Charisma.Min, Max: metatype.Charisma.Max, AugMax: metatype.Charisma.AugMax}
		char.Essence = Attribute[float64]{Name: "Essence", Base: pregen.Essence.Base, Min: metatype.Essence.Min, Max: metatype.Essence.Max, AugMax: metatype.Essence.AugMax}
		char.Magic = Attribute[int]{Name: "Magic", Base: pregen.Magic.Base, Min: metatype.Magic.Min, Max: metatype.Magic.Max, AugMax: metatype.Magic.AugMax}
		char.Resonance = Attribute[int]{Name: "Resonance", Base: pregen.Resonance.Base, Min: metatype.Resonance.Min, Max: metatype.Resonance.Max, AugMax: metatype.Resonance.AugMax}

		char.Skills = pregen.Skills
		char.Qualtities = pregen.Qualtities

		// Proceed with character details.
		return PromptSetCharacterDetails(s, a, char)
	}
}

func PromptSetCharacterDetails(s ssh.Session, a *Account, char *Character) (string, *Character) {
	if state, updatedChar := PromptSetCharacterName(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	// Prompt for additional details
	if state, updatedChar := PromptSetCharacterSex(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	if state, updatedChar := PromptSetCharacterAge(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	if state, updatedChar := PromptSetCharacterHeight(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	if state, updatedChar := PromptSetCharacterWeight(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	if state, updatedChar := PromptSetCharacterShortDescription(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	if state, updatedChar := PromptSetCharacterLongDescription(s, a, char); state == StateError {
		return state, nil
	} else {
		char = updatedChar
	}

	// Save and return
	a.Characters = append(a.Characters, char.Name)
	CharacterMgr.AddCharacter(char)
	char.Save()

	WriteString(s, cfmt.Sprintf("{{Character '%s' created successfully! Returning to main menu.}}::green"+CRLF, char.Name))

	return StateMainMenu, char
}

func PromptSetCharacterName(s ssh.Session, a *Account, char *Character) (string, *Character) {
	WriteString(s, "{{Enter your character's name:}}::white|bold ")
	name, err := InputPrompt(s, "")
	if err != nil {
		return StateError, nil
	}

	if err := ValidateCharacterName(name); err != nil {
		WriteString(s, cfmt.Sprintf("{{Invalid name: %s}}::red"+CRLF, err.Error()))
		return StateError, nil
	}

	char.Name = strings.TrimSpace(name)

	return "", char
}

func PromptSetCharacterSex(s ssh.Session, a *Account, char *Character) (string, *Character) {
	options := []MenuOption{
		{"Male", "male", "Male character"},
		{"Female", "female", "Female character"},
		{"Non-Binary", "non_binary", "Non-binary character"},
		// {"Custom", "custom", "Enter a custom gender identity"},
	}

	for {
		choice, err := PromptForMenu(s, "Select Your Character's Sex", options)
		if err != nil {
			return StateError, nil
		}

		// if choice == "custom" {
		// 	WriteString(s, "{{Enter your character's gender identity:}}::cyan ")
		// 	customGender, err := InputPrompt(s, "")
		// 	if err != nil {
		// 		return StateError, nil
		// 	}
		// 	char.Sex = strings.TrimSpace(customGender)
		// } else {
		char.Sex = choice
		// }

		return "", char
	}
}

func PromptSetCharacterAge(s ssh.Session, a *Account, char *Character) (string, *Character) {
	for {
		WriteString(s, "{{Enter your character's age (numeric):}}::white|bold ")
		input, err := InputPrompt(s, "")
		if err != nil {
			return StateError, nil
		}

		age, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || age < 0 {
			WriteString(s, "{{Invalid age. Please enter a positive number.}}::red"+CRLF)
			continue
		}

		char.Age = age
		return "", char
	}
}

func PromptSetCharacterHeight(s ssh.Session, a *Account, char *Character) (string, *Character) {
	for {
		WriteString(s, "{{Enter your character's height in cm:}}::white|bold ")
		input, err := InputPrompt(s, "")
		if err != nil {
			return StateError, nil
		}

		height, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || height < 50 || height > 300 {
			WriteString(s, "{{Invalid height. Enter a value between 50cm and 300cm.}}::red"+CRLF)
			continue
		}

		char.Height = height
		return "", char
	}
}

func PromptSetCharacterWeight(s ssh.Session, a *Account, char *Character) (string, *Character) {
	for {
		WriteString(s, "{{Enter your character's weight in kg:}}::white|bold ")
		input, err := InputPrompt(s, "")
		if err != nil {
			return StateError, nil
		}

		weight, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || weight < 20 || weight > 500 {
			WriteString(s, "{{Invalid weight. Enter a value between 20kg and 500kg.}}::red"+CRLF)
			continue
		}

		char.Weight = weight
		return "", char
	}
}

// TODO: Using details about this character, generate a random description
func PromptSetCharacterShortDescription(s ssh.Session, a *Account, char *Character) (string, *Character) {
	WriteString(s, "{{Enter a short description of your character:}}::white|bold ")
	description, err := InputPrompt(s, "")
	if err != nil {
		return StateError, nil
	}

	char.Description = strings.TrimSpace(description)
	return "", char
}

// TODO: Using details about this character, generate a random description
func PromptSetCharacterLongDescription(s ssh.Session, a *Account, char *Character) (string, *Character) {
	WriteString(s, "{{Enter a long description of your character (background, details, etc.):}}::white|bold ")
	description, err := InputPrompt(s, "")
	if err != nil {
		return StateError, nil
	}

	char.LongDescription = strings.TrimSpace(description)
	return "", char
}

func PromptEnterGame(s ssh.Session, a *Account) (string, *Character) {
	// Ensure the account has at least one character; otherwise, redirect to character creation.
	if len(a.Characters) == 0 {
		WriteString(s, "{{You have no characters. Create one to start playing.}}::red"+CRLF)
		return PromptCharacterCreate(s, a)
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

		c.Recalculate()
		c.Save()
		CharacterMgr.SetCharacterOnline(c)

		// Notify the user and proceed into the game.
		WriteString(s, cfmt.Sprintf("{{Entering the game as %s...}}::green|bold"+CRLF, c.Name))
		return StateGameLoop, c
	}
}

func PromptCharacterDelete(s ssh.Session, a *Account) string {
	// Ensure the account has at least one character.
	if len(a.Characters) == 0 {
		WriteString(s, "{{You have no characters to delete.}}::red"+CRLF)
		return StateMainMenu
	}

	// Build menu options for each character and an option to return to the main menu.
	options := make([]MenuOption, len(a.Characters)+1)
	for i, name := range a.Characters {
		options[i] = MenuOption{
			DisplayText: name,
			Value:       name,
			Description: "Delete this character",
		}
	}
	options[len(a.Characters)] = MenuOption{
		DisplayText: "Back to Main Menu",
		Value:       "back",
		Description: "Return to main menu",
	}

	for {
		choice, err := PromptForMenu(s, "Select a Character to Delete", options)
		if err != nil {
			slog.Error("Error prompting for character deletion", slog.Any("error", err))
			return StateError
		}

		if choice == "back" {
			return StateMainMenu
		}

		c := CharacterMgr.GetCharacterByName(choice)
		if c == nil {
			WriteString(s, "{{Character not found. Please try again.}}::red"+CRLF)
			continue
		}

		// Confirm deletion
		WriteString(s, cfmt.Sprintf("\n{{Are you sure you want to delete %s? This action cannot be undone.}}::red|bold\n", c.Name))
		if !YesNoPrompt(s, false) {
			WriteString(s, "{{Character deletion canceled.}}::yellow"+CRLF)
			continue
		}

		// Remove the character
		CharacterMgr.RemoveCharacter(c)
		a.Characters = removeCharacterFromAccount(a.Characters, choice)
		a.Save()

		if err := RemoveFile(filepath.Join(CharactersFilepath, fmt.Sprintf("%s.yml", strings.ToLower(c.Name)))); err != nil {
			slog.Error("Error removing character file", slog.Any("error", err))
		}

		WriteString(s, cfmt.Sprintf("\n{{Character %s has been deleted.}}::green\n", c.Name))
		return StateMainMenu
	}
}

func removeCharacterFromAccount(characters []string, name string) []string {
	updated := []string{}
	for _, char := range characters {
		if char != name {
			updated = append(updated, char)
		}
	}
	return updated
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
