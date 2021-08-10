package entity

import (
	"time"

	"github.com/noknow-hub/go_crypto"
	"github.com/pkg/errors"
	"go-ddd/config"
	"golang.org/x/crypto/bcrypt"
)

func emptyPointer() *string {
	v := ""
	return &v
}

func genHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(hashedPassword), nil
}

func passwordIsValid(hashed string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func genRecoveryToken() (string, time.Duration, time.Time, error) {
	duration := time.Hour * 2
	expire := time.Now().Add(duration)
	token, err := crypto.EncryptCTR(
		expire.Format(time.RFC3339),
		config.Env.App.Secret,
	)
	return token, duration, expire, errors.WithStack(err)
}

func recoveryTokenIsValid(token string) bool {
	decrypted, err := crypto.DecryptCTR(token, config.Env.App.Secret)
	expire, err := time.Parse(time.RFC3339, decrypted)
	return !(err != nil || time.Now().After(expire))
}
