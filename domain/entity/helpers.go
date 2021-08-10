package entity

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

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
