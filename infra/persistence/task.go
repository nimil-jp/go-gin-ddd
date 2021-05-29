package persistence

import (
	"go-ddd/domain/model"
	"go-ddd/domain/repository"
	"go-ddd/interface/request"
	"gorm.io/gorm"
)

type task struct{}

func NewTask() repository.ITask {
	return &task{}
}

func (t task) Create(db *gorm.DB, task *model.Task) (uint, error) {
	panic("implement me")
}

func (t task) GetAll(db *gorm.DB, paging *request.Paging) ([]*model.Task, uint, error) {
	panic("implement me")
}

func (t task) Update(db *gorm.DB, task *model.Task) error {
	panic("implement me")
}
