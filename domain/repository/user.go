package repository

import (
	"go-ddd/domain/model"
	"gorm.io/gorm"
)

type IUser interface {
	Create(db *gorm.DB, user *model.User) (uint, error)
	EmailExists(db *gorm.DB, email string) bool
	GetByEmail(db *gorm.DB, email string) (*model.User, error)
}
