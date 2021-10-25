package vobj

import (
	"database/sql"
	"database/sql/driver"
	"time"

	crypto "github.com/noknow-hub/go_crypto"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go-gin-ddd/config"
)

type RecoveryToken string

func NewRecoveryToken(recoveryToken string) *RecoveryToken {
	var value = RecoveryToken(recoveryToken)
	return &value
}

func (p *RecoveryToken) Generate() (time.Duration, time.Time, error) {
	duration := config.RecoveryTokenExpire
	expire := time.Now().Add(duration)
	token, err := crypto.EncryptCTR(
		expire.Format(time.RFC3339),
		config.Env.App.Secret,
	)
	*p = RecoveryToken(token)

	return duration, expire, errors.WithStack(err)
}

func (p RecoveryToken) IsValid() bool {
	decrypted, err := crypto.DecryptCTR(string(p), config.Env.App.Secret)
	expire, err := time.Parse(time.RFC3339, decrypted)
	return !(err != nil || time.Now().After(expire))
}

func (p RecoveryToken) String() string {
	return string(p)
}

func (p *RecoveryToken) Clear() {
	*p = ""
}

// sql

func (p *RecoveryToken) Scan(value interface{}) error {
	nullString := &sql.NullString{}
	err := nullString.Scan(value)
	*p = RecoveryToken(nullString.String)

	return errors.WithStack(err)
}

func (p RecoveryToken) Value() (driver.Value, error) {
	return string(p), nil
}

// GormDataType gorm common data type
func (p RecoveryToken) GormDataType() string {
	return "recovery_token"
}

// GormDBDataType gorm db data type
func (p RecoveryToken) GormDBDataType(_ *gorm.DB, _ *schema.Field) string {
	return "varchar(256)"
}

// json

func (p RecoveryToken) MarshalJSON() ([]byte, error) {
	return []byte("\"" + p + "\""), nil
}

func (p *RecoveryToken) UnmarshalJSON(b []byte) error {
	*p = RecoveryToken(b)
	return nil
}
