package vobj

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go-gin-ddd/config"
	"go-gin-ddd/pkg/context"
)

type Password string

func NewPassword(ctx context.Context, password, passwordConfirm string) (*Password, error) {
	if password != passwordConfirm {
		ctx.FieldError("PasswordConfirm", "パスワードと一致しません")
		return nil, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.BcryptHashCost)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	value := Password(hashedPassword)
	return &value, nil
}

func (p Password) IsValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(password)) == nil
}

// sql

func (p *Password) Scan(value interface{}) error {
	nullString := &sql.NullString{}
	err := nullString.Scan(value)
	*p = Password(nullString.String)

	return errors.WithStack(err)
}

func (p Password) Value() (driver.Value, error) {
	return string(p), nil
}

// GormDataType gorm common data type
func (p Password) GormDataType() string {
	return "password"
}

// GormDBDataType gorm db data type
func (p Password) GormDBDataType(_ *gorm.DB, _ *schema.Field) string {
	return "longtext"
}

// json

func (p Password) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p + "\""), nil
}

func (p *Password) UnmarshalJSON(b []byte) error {
	*p = Password(b)
	return nil
}
