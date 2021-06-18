package model

import (
	"go-ddd/domain"
)

type User struct {
	domain.SoftDeleteModel
	UserID   string `json:"user_id"`
	Password bool   `json:"password"`
}
