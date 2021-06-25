package persistence

import (
	"go-ddd/domain/model"
	"go-ddd/domain/repository"
	"gorm.io/gorm"
)

type user struct{}

func NewUser() repository.IUser {
	return &user{}
}

func (u user) Create(db *gorm.DB, user *model.User) (uint, error) {
	if err := db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (u user) EmailExists(db *gorm.DB, email string) bool {
	var count int64
	db.Model(&model.User{}).Where(&model.User{Email: email}).Count(&count)
	return count > 0
}

func (u user) GetByEmail(db *gorm.DB, email string) (*model.User, error) {
	var dest model.User
	err := db.Where(&model.User{Email: email}).First(&dest).Error
	if err != nil {
		return nil, err
	}
	return &dest, nil
}
