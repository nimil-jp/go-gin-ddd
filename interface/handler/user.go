package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-ddd/resource/request"
	"go-ddd/usecase"
	"go-ddd/util"
)

type User struct {
	userUseCase usecase.IUser
}

func NewUser(uuc usecase.IUser) User {
	return User{
		userUseCase: uuc,
	}
}

func (u User) Create(c *gin.Context) {
	var req request.UserCreate

	if !bind(c, &req) {
		return
	}

	id, err := u.userUseCase.Create(&req)

	if err != nil {
		util.ErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusCreated, id)
}

func (u User) Login(c *gin.Context) {
	var req request.UserLogin

	if !bind(c, &req) {
		return
	}

	res, err := u.userUseCase.Login(&req)

	if err != nil {
		util.ErrorJSON(c, err)
		return
	}

	if res == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u User) RefreshToken(c *gin.Context) {
	res, err := u.userUseCase.RefreshToken(c.Query("refresh_token"))

	if err != nil {
		util.ErrorJSON(c, err)
		return
	}

	if res == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, res)
}
