package game

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/spf13/viper"
)

// Validator function type for composability
type Validator func(string) error

// RunValidators applies multiple validators and aggregates errors.
func RunValidators(input string, validators ...Validator) error {
	var errorMessages []string

	for _, validator := range validators {
		if err := validator(input); err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, "; ")) // Return all errors
	}
	return nil
}

// CheckLength validates if input length is within range.
func CheckLength(minLength, maxLength int) Validator {
	return func(input string) error {
		length := utf8.RuneCountInString(input)
		if length < minLength || length > maxLength {
			return fmt.Errorf("must be between %d and %d characters", minLength, maxLength)
		}
		return nil
	}
}

// CheckInputNotInList ensures input is not in a restricted list.
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

// CheckValidCharacters ensures input matches a regex pattern.
func CheckValidCharacters(pattern string) Validator {
	return func(input string) error {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("contains invalid characters")
		}
		return nil
	}
}

// CheckRequiresNumbers ensures input contains at least one numeric character.
func CheckRequiresNumbers() Validator {
	return func(input string) error {
		matched, _ := regexp.MatchString(`[0-9]`, input)
		if !matched {
			return errors.New("must contain at least one number")
		}
		return nil
	}
}

// CheckRequiresSpecialCharacters ensures input contains at least one special character.
func CheckRequiresSpecialCharacters(pattern string) Validator {
	return func(input string) error {
		matched, _ := regexp.MatchString(pattern, input)
		if !matched {
			return errors.New("must contain at least one special character")
		}
		return nil
	}
}

// ValidateCharacterName ensures names follow rules: length, uniqueness, valid characters.
func ValidateCharacterName(name string) error {
	return RunValidators(name,
		CheckLength(viper.GetInt("server.name_min_length"), viper.GetInt("server.name_max_length")),
		CheckInputNotInList(CharacterMgr.GetCharacterNames()),
		CheckInputNotInList(viper.GetStringSlice("banned_names")),
		CheckValidCharacters(viper.GetString("server.name_regex")),
	)
}

// ValidateAccountName ensures names follow rules: length, uniqueness, valid characters.
func ValidateAccountName(name string) error {
	return RunValidators(name,
		CheckLength(viper.GetInt("server.name_min_length"), viper.GetInt("server.name_max_length")),
		CheckInputNotInList(AccountMgr.GetAccountNames()),
		CheckInputNotInList(viper.GetStringSlice("banned_names")),
		CheckValidCharacters(viper.GetString("server.name_regex")),
	)
}

// ValidateDescription ensures descriptions follow rules: length.
func ValidateDescription(description string) error {
	return RunValidators(description,
		CheckLength(viper.GetInt("server.description_min_length"), viper.GetInt("server.description_max_length")),
	)
}

// ValidatePassword ensures passwords follow rules: length, numbers, special characters.
func ValidatePassword(password string) error {
	return RunValidators(password,
		CheckLength(viper.GetInt("server.password_min_length"), viper.GetInt("server.password_max_length")),
		CheckRequiresNumbers(),
		CheckRequiresSpecialCharacters(viper.GetString("server.special_character_regex")),
	)
}
