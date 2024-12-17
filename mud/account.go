package mud

import "github.com/rs/zerolog"

type Account struct {
	Log      zerolog.Logger `json:"-"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}

func NewAccount(username, password string) *Account {
	return &Account{
		Log:      NewDevLogger(),
		Username: username,
		Password: password,
	}
}
