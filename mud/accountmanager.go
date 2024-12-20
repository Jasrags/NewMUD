package mud

import "github.com/rs/zerolog"

type AccountManager struct {
	Log zerolog.Logger
}

func NewAccountManager() *AccountManager {
	return &AccountManager{
		Log: NewDevLogger(),
	}
}
