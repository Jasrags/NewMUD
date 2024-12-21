package users

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser() *User {
	return &User{}
}

func GetByUserID(id string) *User {
	return &User{}
}

func Save(u *User) {
	dataFilePath := viper.GetString("data.users_path")

	slog.Info("Saving user data",
		slog.String("datafile_path", dataFilePath),
		slog.String("id", u.ID),
		slog.String("username", u.Username))

	filePath := filepath.Join(dataFilePath, u.ID+".json")
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("failed creating file",
			slog.String("file", filePath),
			slog.Any("error", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(u); err != nil {
		slog.Error("failed to encode user data",
			slog.Any("error", err))
	}
}
