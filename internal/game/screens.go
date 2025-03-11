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

func PromptCharacterCreate(s ssh.Session, a *Account) (string, *Character) {
	options := []MenuOption{
		{"Create a Pre-Generated Character", "pregen", "Select from predefined character archetypes"},
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
			c := NewCharacter()
			return PromptPregenCharacterMenu(s, a, c)
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

// TODO: Set a base nuyen level for the character
// TODO: We should use a item pack to set the starting gear when we have that implemented
// TODO: Should we allow changing of the metatype for a pregen?
func PromptPregenCharacterMenu(s ssh.Session, a *Account, c *Character) (string, *Character) {

	// Build the menu options
	options := []MenuOption{
		{"Set Character Template", "template", "Select a pre-generated character template"},
		{"Set Character Name", "name", "Enter the character's name"},
		{"Set Character Sex", "sex", "Select the character's sex"},
		{"Set Character Age", "age", "Enter the character's age"},
		{"Set Character Height", "height", "Enter the character's height in cm"},
		{"Set Character Weight", "weight", "Enter the character's weight in kg"},
		{"Set Character Short Description", "short_desc", "Enter a short description"},
		{"Set Character Long Description", "long_desc", "Enter a detailed description"},
		{"Save Character", "save", "Save the character (all fields must be completed)"},
		{"Back", "back", "Return to the previous menu"},
	}

	for {
		// Display current creation progress.
		displayCharacterProgress(s, c)

		choice, err := PromptForMenu(s, "Create a Pre-Generated Character", options)
		if err != nil {
			return StateError, nil
		}

		switch choice {
		case "template":
			// Optionally, you can reuse your existing pregen selection
			state, newChar := PromptSelectPregenTemplate(s, a, c)
			if state == StateError {
				return state, nil
			}
			c = newChar
		case "name":
			state, newChar := PromptSetCharacterName(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "sex":
			state, newChar := PromptSetCharacterSex(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "age":
			state, newChar := PromptSetCharacterAge(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "height":
			state, newChar := PromptSetCharacterHeight(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "weight":
			state, newChar := PromptSetCharacterWeight(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "short_desc":
			state, newChar := PromptSetCharacterShortDescription(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "long_desc":
			state, newChar := PromptSetCharacterLongDescription(s, a, c)
			if state == StateError {
				continue
			}
			c = newChar
		case "save":
			// Validate that all required fields are present.
			missing := validateCharacterFields(c)
			if len(missing) > 0 {
				WriteString(s, fmt.Sprintf("The following fields are missing: %s. Please complete them before saving."+CRLF, strings.Join(missing, ", ")))
				continue
			}

			// Save the character
			a.Characters = append(a.Characters, c.Name)
			CharacterMgr.AddCharacter(c)
			c.Save()
			WriteStringF(s, "Character '%s' created successfully! Returning to main menu."+CRLF, c.Name)
			return StateMainMenu, c
		case "back":
			return StateCharacterCreate, nil
		default:
			WriteString(s, "Invalid selection. Please try again."+CRLF)
		}
	}
}

// displayCharacterProgress prints the current state of the character creation.
func displayCharacterProgress(s ssh.Session, char *Character) {
	var progress strings.Builder

	titleString := "{{%s:}}::white|bold|underline "
	unsetString := "{{<unset>}}::red" + CRLF
	setString := "{{%v}}::green" + CRLF

	// Template: If MetatypeID is set, try to get the corresponding template title.
	progress.WriteString(fmt.Sprintf(titleString, "Template"))
	if char.MetatypeID == "" {
		progress.WriteString(unsetString)
	} else {
		p := EntityMgr.GetPregen(char.PregenID)
		progress.WriteString(fmt.Sprintf(setString, p.Title))
	}

	// Name
	progress.WriteString(fmt.Sprintf(titleString, "Name"))
	if strings.TrimSpace(char.Name) == "" {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.Name))
	}

	// Sex
	progress.WriteString(fmt.Sprintf(titleString, "Sex"))
	if strings.TrimSpace(char.Sex) == "" {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.Sex))
	}

	// Age
	progress.WriteString(fmt.Sprintf(titleString, "Age"))
	if char.Age <= 0 {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.Age))
	}

	// Height
	progress.WriteString(fmt.Sprintf(titleString, "Height"))
	if char.Height <= 0 {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.Height))
	}

	// Weight
	progress.WriteString(fmt.Sprintf(titleString, "Weight"))
	if char.Weight <= 0 {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.Weight))
	}

	// Short Description
	progress.WriteString(fmt.Sprintf(titleString, "Short Description"))
	if strings.TrimSpace(char.Description) == "" {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.Description))
	}

	// Long Description
	progress.WriteString(fmt.Sprintf(titleString, "Long Description"))
	if strings.TrimSpace(char.LongDescription) == "" {
		progress.WriteString(unsetString)
	} else {
		progress.WriteString(fmt.Sprintf(setString, char.LongDescription))
	}

	// WriteString(s, borderStyle.Render(progress.String()))
	// Write the progress block with a newline after.
	WriteString(s, CRLF+progress.String()+CRLF)
}

// validateCharacterFields checks required fields and returns a slice of missing field names.
func validateCharacterFields(char *Character) []string {
	missing := []string{}
	if char.MetatypeID == "" {
		missing = append(missing, "Character Template")
	}
	if strings.TrimSpace(char.Name) == "" {
		missing = append(missing, "Character Name")
	}
	if strings.TrimSpace(char.Sex) == "" {
		missing = append(missing, "Character Sex")
	}
	if char.Age <= 0 {
		missing = append(missing, "Character Age")
	}
	if char.Height <= 0 {
		missing = append(missing, "Character Height")
	}
	if char.Weight <= 0 {
		missing = append(missing, "Character Weight")
	}
	if strings.TrimSpace(char.Description) == "" {
		missing = append(missing, "Short Description")
	}
	if strings.TrimSpace(char.LongDescription) == "" {
		missing = append(missing, "Long Description")
	}
	return missing
}

func PromptSelectPregenTemplate(s ssh.Session, a *Account, char *Character) (string, *Character) {
	// Retrieve the list of pre-generated templates.
	pregens := EntityMgr.GetPregens()
	pregenMap := make(map[string]*Pregen)
	options := make([]MenuOption, 0, len(pregens)+1)

	// Build menu options from available pregens.
	for _, pregen := range pregens {
		pregenMap[pregen.ID] = pregen
		options = append(options, MenuOption{
			DisplayText: pregen.Title,
			Value:       pregen.ID,
			Description: pregen.GetSelectionInfo(),
		})
	}
	// Add a "Back" option.
	options = append(options, MenuOption{"Back", "back", "Return to previous menu"})

	// Loop until a valid selection is made.
	for {
		choice, err := PromptForMenu(s, "Select a Pre-Generated Template", options)
		if err != nil {
			slog.Error("Error prompting for pregen template", slog.Any("error", err))
			return StateError, char
		}

		if choice == "back" {
			// Return to the calling menu without modifying the character.
			return "", char
		}

		pregen, exists := pregenMap[choice]
		if !exists {
			WriteString(s, "Invalid selection. Please try again."+CRLF)
			continue
		}

		// Show details of the selected template.
		WriteStringF(s, CRLF+"Template Selected: %s"+CRLF, pregen.Title)
		// Optionally, you can display more details by calling pregen.GetSelectionInfo()
		// and then ask for confirmation.
		if !YesNoPrompt(s, true) {
			WriteString(s, "Selection canceled. Please choose again."+CRLF)
			continue
		}

		// Update the existing character with the template data.
		char.PregenID = pregen.ID
		char.MetatypeID = pregen.MetatypeID
		char.Body = pregen.Body
		char.Agility = pregen.Agility
		char.Reaction = pregen.Reaction
		char.Strength = pregen.Strength
		char.Willpower = pregen.Willpower
		char.Logic = pregen.Logic
		char.Intuition = pregen.Intuition
		char.Charisma = pregen.Charisma
		char.Essence = pregen.Essence
		char.Magic = pregen.Magic
		char.Resonance = pregen.Resonance
		// char.Body = Attribute[int]{Base: pregen.Body.Base}
		// char.Agility = Attribute[int]{Base: pregen.Agility.Base}
		// char.Reaction = Attribute[int]{Base: pregen.Reaction.Base}
		// char.Strength = Attribute[int]{Base: pregen.Strength.Base}
		// char.Willpower = Attribute[int]{Base: pregen.Willpower.Base}
		// char.Logic = Attribute[int]{Base: pregen.Logic.Base}
		// char.Intuition = Attribute[int]{Base: pregen.Intuition.Base}
		// char.Charisma = Attribute[int]{Base: pregen.Charisma.Base}
		// char.Essence = Attribute[float64]{Base: pregen.Essence.Base}
		// char.Magic = Attribute[int]{Base: pregen.Magic.Base}
		// char.Resonance = Attribute[int]{Base: pregen.Resonance.Base}

		// Update skills and qualities from the template.
		char.Skills = pregen.Skills
		char.Qualtities = pregen.Qualtities

		// Successfully updated; return to the central menu.
		return "", char
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

	WriteStringF(s, "{{Character '%s' created successfully! Returning to main menu.}}::green"+CRLF, char.Name)

	return StateMainMenu, char
}

func PromptSetCharacterName(s ssh.Session, a *Account, char *Character) (string, *Character) {
	for {
		WriteString(s, CRLF+"{{Enter your character's name:}}::white|bold|underline ")
		name, err := InputPrompt(s, "")
		if err != nil {
			return StateError, nil
		}

		if err := ValidateCharacterName(name); err != nil {
			WriteStringF(s, "{{Invalid name: %s}}::red"+CRLF, err.Error())
			continue
		}

		name = strings.TrimSpace(name)

		WriteStringF(s, CRLF+"{{Name: %s}}::cyan"+CRLF, name)

		if !YesNoPrompt(s, true) {
			continue
		}

		char.Name = name

		return "", char
	}
}

func PromptSetCharacterSex(s ssh.Session, a *Account, char *Character) (string, *Character) {
	options := []MenuOption{
		{"Male", "male", "Male character"},
		{"Female", "female", "Female character"},
		{"Non-Binary", "non_binary", "Non-binary character"},
	}

	for {
		choice, err := PromptForMenu(s, "Select Your Character's Sex", options)
		if err != nil {
			continue
		}

		WriteStringF(s, CRLF+"{{Sex: %s}}::cyan"+CRLF, choice)

		if !YesNoPrompt(s, true) {
			continue
		}

		char.Sex = choice

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

		// Set ItemBlueprint for each item in inventory and equipment
		for _, item := range c.Inventory.Items {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp == nil {
				slog.Warn("Item blueprint not found",
					slog.String("character_id", c.ID),
					slog.String("item_blueprint_id", item.BlueprintID),
					slog.String("item_instnace_id", item.InstanceID))
				continue
			}
			item.Blueprint = bp
		}

		for slot, item := range c.Equipment.Slots {
			bp := EntityMgr.GetItemBlueprintByInstance(item)
			if bp == nil {
				slog.Warn("Item blueprint not found",
					slog.String("character_id", c.ID),
					slog.String("slot", slot),
					slog.String("item_blueprint_id", item.BlueprintID),
					slog.String("item_instnace_id", item.InstanceID))
				continue
			}
			item.Blueprint = bp
		}

		// Save any changes to the account and character, and mark the character as online.
		a.Save()

		// c.Recalculate()
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
