package usecase

import (
	"go-ddd/domain/model"
	"go-ddd/domain/repository"
	"go-ddd/interface/request"
)

type IUser interface {
	Create(task *model.User) error
	GetAll(keyword string, paging *request.Paging) ([]*model.User, uint, error)
	Update(task *model.User) error
}

type user struct {
	userRepo repository.IUser
}

func NewUser(tr repository.IUser) IUser {
	return &user{
		userRepo: tr,
	}
}

func (t user) Create(user *model.User) error {
	panic("implement me")
}

func (t user) GetAll(keyword string, paging *request.Paging) ([]*model.User, uint, error) {
	panic("implement me")
}

func (t user) Update(user *model.User) error {
	panic("implement me")
}
