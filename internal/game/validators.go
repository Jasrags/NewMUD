package game

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/spf13/viper"
)

// ------------------------------
// Validator Types and Aggregation
// ------------------------------

// Validator is a function type for validating a string.
type Validator func(string) error

// MultiValidationError aggregates multiple validation errors.
type MultiValidationError struct {
	Errors []error
}

// Error implements the error interface for MultiValidationError.
func (m *MultiValidationError) Error() string {
	var messages []string
	for _, err := range m.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// RunValidators applies multiple validators to the input and returns a MultiValidationError if any errors occur.
func RunValidators(input string, validators ...Validator) error {
	var errs []error
	for _, validator := range validators {
		if err := validator(input); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return &MultiValidationError{Errors: errs}
	}
	return nil
}

// ------------------------------
// Validator Functions
// ------------------------------

// CheckLength validates that the input's length is between minLength and maxLength.
func CheckLength(minLength, maxLength int) Validator {
	return func(input string) error {
		length := utf8.RuneCountInString(input)
		if length < minLength || length > maxLength {
			return fmt.Errorf("must be between %d and %d characters", minLength, maxLength)
		}
		return nil
	}
}

// CheckInputNotInList ensures the input is not contained in a disallowed list.
func CheckInputNotInList(disallowedList []string) Validator {
	return func(input string) error {
		for _, banned := range disallowedList {
			if strings.EqualFold(input, banned) {
				return fmt.Errorf("%s is not allowed", input)
			}
		}
		return nil
	}
}

// CheckValidCharacters returns a validator that checks if the input matches a precompiled regex pattern.
func CheckValidCharacters(pattern string) Validator {
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Since the pattern is often supplied via configuration, panicking here can alert you to misconfiguration.
		panic(fmt.Sprintf("invalid regex pattern %q: %v", pattern, err))
	}
	return func(input string) error {
		if !re.MatchString(input) {
			return errors.New("contains invalid characters")
		}
		return nil
	}
}

// Precompile a regex for numeric characters.
var numberRegex = regexp.MustCompile(`[0-9]`)

// CheckRequiresNumbers ensures that the input contains at least one numeric character.
func CheckRequiresNumbers() Validator {
	return func(input string) error {
		if !numberRegex.MatchString(input) {
			return errors.New("must contain at least one number")
		}
		return nil
	}
}

// CheckRequiresSpecialCharacters returns a validator that ensures the input contains at least one special character.
// The special characters are determined by the provided regex pattern.
func CheckRequiresSpecialCharacters(pattern string) Validator {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid special character regex %q: %v", pattern, err))
	}
	return func(input string) error {
		if !re.MatchString(input) {
			return errors.New("must contain at least one special character")
		}
		return nil
	}
}

// ------------------------------
// High-Level Validation Functions
// ------------------------------

// Note: The following functions pull configuration values via Viper and
// assume that CharacterMgr and AccountMgr are defined elsewhere in your project.

// ValidateCharacterName ensures that a character name meets the required length,
// is unique (i.e. not already taken), is not banned, and contains only valid characters.
func ValidateCharacterName(name string) error {
	return RunValidators(name,
		CheckLength(viper.GetInt("server.name_min_length"), viper.GetInt("server.name_max_length")),
		CheckInputNotInList(CharacterMgr.GetCharacterNames()),
		CheckInputNotInList(viper.GetStringSlice("banned_names")),
		CheckValidCharacters(viper.GetString("server.name_regex")),
	)
}

// ValidateAccountName ensures that an account name meets the required criteria.
func ValidateAccountName(name string) error {
	return RunValidators(name,
		CheckLength(viper.GetInt("server.name_min_length"), viper.GetInt("server.name_max_length")),
		CheckInputNotInList(AccountMgr.GetAccountNames()),
		CheckInputNotInList(viper.GetStringSlice("banned_names")),
		CheckValidCharacters(viper.GetString("server.name_regex")),
	)
}

// ValidateDescription ensures that a description's length is within the allowed range.
func ValidateDescription(description string) error {
	return RunValidators(description,
		CheckLength(viper.GetInt("server.description_min_length"), viper.GetInt("server.description_max_length")),
	)
}

// ValidatePassword ensures that a password meets the required length,
// contains at least one number, and at least one special character.
func ValidatePassword(password string) error {
	return RunValidators(password,
		CheckLength(viper.GetInt("server.password_min_length"), viper.GetInt("server.password_max_length")),
		CheckRequiresNumbers(),
		CheckRequiresSpecialCharacters(viper.GetString("server.special_character_regex")),
	)
}
