package mud

import "github.com/rs/zerolog"

type Account struct {
	Log      zerolog.Logger `json:"-"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}

func NewAccount(l zerolog.Logger) *Account {
	return &Account{
		Log: l,
	}
}
