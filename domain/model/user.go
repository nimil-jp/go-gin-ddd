package model

import (
	"go-ddd/domain"
)

type User struct {
	domain.SoftDeleteModel
	Email    string `json:"email"`
	Password string `json:"password"`
}
