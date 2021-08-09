package repository

import (
	"go-ddd/domain/entity"
	"gorm.io/gorm"
)

type IUser interface {
	Create(db *gorm.DB, user *entity.User) (uint, error)
	EmailExists(db *gorm.DB, email string) (bool, error)
	GetByEmail(db *gorm.DB, email string) (*entity.User, error)
}
