package repository

import (
	"go-ddd/domain/model"
	"go-ddd/interface/request"
	"gorm.io/gorm"
)

type IUser interface {
	Create(db *gorm.DB, task *model.User) (uint, error)
	GetAll(db *gorm.DB, paging *request.Paging) ([]*model.User, uint, error)
	Update(db *gorm.DB, task *model.User) error
}
