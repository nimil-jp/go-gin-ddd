package usecase

import (
	"go-ddd/domain/model"
	"go-ddd/domain/repository"
	"go-ddd/interface/request"
)

type ITask interface {
	Create(task *model.Task) error
	GetAll(keyword string, paging *request.Paging) ([]*model.Task, uint, error)
	Update(task *model.Task) error
}

type task struct {
	taskRepo repository.ITask
}

func NewTask(tr repository.ITask) ITask {
	return &task{
		taskRepo: tr,
	}
}

func (t task) Create(task *model.Task) error {
	panic("implement me")
}

func (t task) GetAll(keyword string, paging *request.Paging) ([]*model.Task, uint, error) {
	panic("implement me")
}

func (t task) Update(task *model.Task) error {
	panic("implement me")
}
