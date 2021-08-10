package entity

import (
	"go-ddd/domain"
	"go-ddd/resource/request"
	"go-ddd/util/xerrors"
)

type User struct {
	domain.SoftDeleteModel
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(verr *xerrors.Validation, dto *request.UserCreate) (*User, error) {
	if dto.Password != dto.PasswordConfirm {
		verr.Add("PasswordConfirm", "パスワードと一致しません")
		return nil, nil
	}

	password, err := genHashedPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    dto.Email,
		Password: password,
	}, nil
}

func (u User) PasswordIsValid(password string) bool {
	return passwordIsValid(u.Password, password)
}
