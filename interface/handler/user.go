package handler

import (
	"github.com/gin-gonic/gin"
	"go-ddd/usecase"
)

type IUser interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
}

type user struct {
	userUseCase usecase.IUser
}

func NewUser(tu usecase.IUser) IUser {
	return &user{
		userUseCase: tu,
	}
}

func (t user) Create(c *gin.Context) {
	panic("implement me")
}

func (t user) GetAll(c *gin.Context) {
	panic("implement me")
}

func (t user) Update(c *gin.Context) {
	panic("implement me")
}
