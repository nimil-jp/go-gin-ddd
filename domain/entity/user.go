package entity

import (
	"time"

	"go-ddd/domain"
	"go-ddd/pkg/xerrors"
	"go-ddd/resource/request"
)

type User struct {
	domain.SoftDeleteModel
	Email    string `json:"email"`
	Password string `json:"password"`

	RecoveryToken *string `json:"-"`
}

func NewUser(verr *xerrors.Validation, dto *request.UserCreate) (*User, error) {
	var user = User{
		Email: dto.Email,
	}

	ok, err := user.setPassword(verr, dto.Password, dto.PasswordConfirm)
	if err != nil || !ok {
		return nil, err
	}

	return &user, nil
}

func (u *User) setPassword(verr *xerrors.Validation, password, passwordConfirm string) (ok bool, err error) {
	if password != passwordConfirm {
		verr.Add("PasswordConfirm", "パスワードと一致しません")
		return false, nil
	}

	password, err = genHashedPassword(password)
	if err != nil {
		return false, err
	}

	u.Password = password
	return true, nil
}

func (u User) PasswordIsValid(password string) bool {
	return passwordIsValid(u.Password, password)
}

func (u *User) ResetPasswordRequest() (token string, duration time.Duration, expire time.Time, err error) {
	token, duration, expire, err = genRecoveryToken()
	if err != nil {
		return
	}
	u.RecoveryToken = &token
	return
}

func (u *User) ResetPassword(verr *xerrors.Validation, dto *request.UserResetPassword) error {
	if !recoveryTokenIsValid(dto.RecoveryToken) {
		verr.Add("RecoveryToken", "リカバリートークンが無効です")
		return nil
	}

	ok, err := u.setPassword(verr, dto.Password, dto.PasswordConfirm)
	if err != nil || !ok {
		return err
	}
	u.RecoveryToken = emptyPointer()
	return nil
}
