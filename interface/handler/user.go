package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-ddd/resource/request"
	"go-ddd/usecase"
)

type User struct {
	userUseCase usecase.IUser
}

func NewUser(route *gin.RouterGroup, uuc usecase.IUser) {
	handler := User{
		userUseCase: uuc,
	}

	post(route, "", handler.Create)
	post(route, "login", handler.Login)
	get(route, "refresh-token", handler.RefreshToken)
	patch(route, "reset-password-request", handler.ResetPasswordRequest)
	patch(route, "reset-password", handler.ResetPassword)
}

func (u User) Create(c *gin.Context) error {
	var req request.UserCreate

	if !bind(c, &req) {
		return nil
	}

	id, err := u.userUseCase.Create(newCtx(), &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}

func (u User) ResetPasswordRequest(c *gin.Context) error {
	var req request.UserResetPasswordRequest

	if !bind(c, &req) {
		return nil
	}

	res, err := u.userUseCase.ResetPasswordRequest(newCtx(), &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, res)
	return nil
}

func (u User) ResetPassword(c *gin.Context) error {
	var req request.UserResetPassword

	if !bind(c, &req) {
		return nil
	}

	err := u.userUseCase.ResetPassword(newCtx(), &req)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u User) Login(c *gin.Context) error {
	var req request.UserLogin

	if !bind(c, &req) {
		return nil
	}

	res, err := u.userUseCase.Login(newCtx(), &req)
	if err != nil {
		return err
	}

	if res == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	c.JSON(http.StatusOK, res)
	return nil
}

func (u User) RefreshToken(c *gin.Context) error {
	res, err := u.userUseCase.RefreshToken(c.Query("refresh_token"))
	if err != nil {
		return err
	}

	if res == nil {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	c.JSON(http.StatusOK, res)
	return nil
}
