package persistence

import (
	"github.com/pkg/errors"
	"go-ddd/domain/entity"
	"go-ddd/domain/repository"
	"gorm.io/gorm"
)

type user struct{}

func NewUser() repository.IUser {
	return &user{}
}

func (u user) Create(db *gorm.DB, user *entity.User) (uint, error) {
	if err := db.Create(user).Error; err != nil {
		return 0, errors.WithStack(err)
	}
	return user.ID, nil
}

func (u user) EmailExists(db *gorm.DB, email string) (bool, error) {
	var count int64
	if err := db.Model(&entity.User{}).Where(&entity.User{Email: email}).Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}
	return count > 0, nil
}

func (u user) GetByEmail(db *gorm.DB, email string) (*entity.User, error) {
	var dest entity.User
	err := db.Where(&entity.User{Email: email}).First(&dest).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &dest, nil
}
