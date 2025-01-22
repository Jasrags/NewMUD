package game

import (
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	ee "github.com/vansante/go-event-emitter"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	sync.RWMutex `yaml:"-"`
	Listeners    []ee.Listener `yaml:"-"`

	ID          string     `yaml:"id"`
	Username    string     `yaml:"username"`
	Password    string     `yaml:"password"`
	Characters  []string   `yaml:"characters"`
	CreatedAt   time.Time  `yaml:"created_at"`
	UpdatedAt   *time.Time `yaml:"updated_at"`
	LastLoginAt *time.Time `yaml:"last_login_at"`
	DeletedAt   *time.Time `yaml:"deleted_at"`
	State       string     `yaml:"-"`
}

func NewAccount() *Account {
	return &Account{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
	}
}

func (u *Account) Init() {
	slog.Debug("Initializing user",
		slog.String("user_id", u.ID))
}

func (u *Account) AddCharacter(char *Character) {
	u.Lock()
	defer u.Unlock()

	slog.Debug("Adding character to user",
		slog.String("user_id", u.ID),
		slog.String("character_id", char.Name))

	u.Characters = append(u.Characters, strings.ToLower(char.Name))
}

func (u *Account) RemoveCharacter(char *Character) {
	u.Lock()
	defer u.Unlock()

	slog.Debug("Removing character from user",

		slog.String("user_id", u.ID),
		slog.String("character_id", char.Name))

	for i, c := range u.Characters {
		if c == char.Name {
			u.Characters = append(u.Characters[:i], u.Characters[i+1:]...)
			break
		}
	}
}

func (u *Account) SetPassword(password string) {
	slog.Debug("Setting password",
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password",
			slog.Any("error", err))
	}

	u.Password = string(hashedPassword)
}

func (u *Account) CheckPassword(password string) bool {
	slog.Debug("Checking password",
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			slog.Error("Invalid password for user",
				slog.String("id", u.ID),
				slog.String("username", u.Username))
		case bcrypt.ErrHashTooShort:
		case bcrypt.ErrMismatchedHashAndPassword:
		case bcrypt.ErrPasswordTooLong:
		default:
			slog.Error("Password error",
				slog.Any("error", err))
		}

		return false
	}

	return true
}

func (u *Account) Save() error {
	u.Lock()
	defer u.Unlock()

	dataFilePath := viper.GetString("data.accounts_path")

	slog.Info("Saving user data",
		slog.String("datafile_path", dataFilePath),
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	filePath := filepath.Join(dataFilePath, strings.ToLower(u.Username)+".yml")

	t := time.Now()
	u.UpdatedAt = &t

	if err := SaveYAML(filePath, u); err != nil {
		slog.Error("failed to save user data",
			slog.Any("error", err))
		return err
	}

	return nil
}
