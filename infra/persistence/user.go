package persistence

import (
	"go-ddd/domain/model"
	"go-ddd/domain/repository"
	"go-ddd/interface/request"
	"gorm.io/gorm"
)

type user struct{}

func NewUser() repository.IUser {
	return &user{}
}

func (t user) Create(db *gorm.DB, task *model.User) (uint, error) {
	panic("implement me")
}

func (t user) GetAll(db *gorm.DB, paging *request.Paging) ([]*model.User, uint, error) {
	panic("implement me")
}

func (t user) Update(db *gorm.DB, task *model.User) error {
	panic("implement me")
}
