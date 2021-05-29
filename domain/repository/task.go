package repository

import (
	"go-ddd/domain/model"
	"go-ddd/interface/request"
	"gorm.io/gorm"
)

type ITask interface {
	Create(db *gorm.DB, task *model.Task) (uint, error)
	GetAll(db *gorm.DB, paging *request.Paging) ([]*model.Task, uint, error)
	Update(db *gorm.DB, task *model.Task) error
}
