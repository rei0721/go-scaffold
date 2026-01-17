package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/crypto"
)

func (a *App) initCrypto() error {
	cryptoInstance, err := crypto.NewBcrypt(
		crypto.WithBcryptCost(crypto.DefaultBcryptCost),
	)
	if err != nil {
		return fmt.Errorf("failed to create crypto: %w", err)
	}

	a.Crypto = cryptoInstance
	a.Logger.Debug("crypto initialized", "algorithm", "bcrypt", "cost", crypto.DefaultBcryptCost)
	return nil
}
