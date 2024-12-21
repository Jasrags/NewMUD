package mud

import (
	"github.com/rs/zerolog"
)

type AccountManager struct {
	Log zerolog.Logger
}

func NewAccountManager(l zerolog.Logger) *AccountManager {
	return &AccountManager{
		Log: l,
	}
}

// func (am *AccountManager) WithContext(ctx context.Context) context.Context {
// 	return context.WithValue(ctx, ctxKeyAccountManager{}, am)
// }
