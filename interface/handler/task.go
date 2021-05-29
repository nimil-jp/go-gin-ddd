package handler

import (
	"github.com/gin-gonic/gin"
	"go-ddd/usecase"
)

type ITask interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
}

type task struct {
	taskUseCase usecase.ITask
}

func NewTask(tu usecase.ITask) ITask {
	return &task{
		taskUseCase: tu,
	}
}

func (t task) Create(c *gin.Context) {
	panic("implement me")
}

func (t task) GetAll(c *gin.Context) {
	panic("implement me")
}

func (t task) Update(c *gin.Context) {
	panic("implement me")
}
