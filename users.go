package main

import (
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	ee "github.com/vansante/go-event-emitter"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	sync.RWMutex
	Listeners []ee.Listener `json:"-"`

	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Password    []byte     `json:"password"`
	Characters  []string   `json:"characters"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	State       string     `json:"-"`
}

func NewUser() *User {
	return &User{
		CreatedAt: time.Now(),
	}
}

func (u *User) Init() {
	slog.Debug("Initializing user",
		slog.String("user_id", u.ID))
}

func (u *User) SetPassword(password string) {
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

func (u *User) CheckPassword(password string) bool {
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

func (u *User) Save() error {
	u.Lock()
	defer u.Unlock()

	dataFilePath := viper.GetString("data.users_path")

	slog.Info("Saving user data",
		slog.String("datafile_path", dataFilePath),
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	filePath := filepath.Join(dataFilePath, strings.ToLower(u.Username)+".json")

	t := time.Now()
	u.UpdatedAt = &t

	if err := SaveJSON(filePath, u); err != nil {
		slog.Error("failed to save user data",
			slog.Any("error", err))
		return err
		// file, err := os.Create(filePath)
		// if err != nil {
		// slog.Error("failed creating file",
		// slog.String("file", filePath),
		// slog.Any("error", err))
	}
	// defer file.Close()

	// encoder := json.NewEncoder(file)
	// encoder.SetIndent("", "  ")
	// encoder.SetEscapeHTML(false)
	// if err := encoder.Encode(u); err != nil {
	// 	slog.Error("failed to encode user data",
	// 		slog.Any("error", err))
	// }

	return nil
}
