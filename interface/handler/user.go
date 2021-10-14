package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-gin-ddd/pkg/context"
	"go-gin-ddd/resource/request"
	"go-gin-ddd/usecase"
)

type User struct {
	userUseCase usecase.IUser
}

func NewUser(uuc usecase.IUser) *User {
	return &User{
		userUseCase: uuc,
	}
}

func (u User) Create(ctx context.Context, c *gin.Context) error {
	var req request.UserCreate

	if !bind(c, &req) {
		return nil
	}

	id, err := u.userUseCase.Create(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, id)
	return nil
}

func (u User) ResetPasswordRequest(ctx context.Context, c *gin.Context) error {
	var req request.UserResetPasswordRequest

	if !bind(c, &req) {
		return nil
	}

	res, err := u.userUseCase.ResetPasswordRequest(ctx, &req)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, res)
	return nil
}

func (u User) ResetPassword(ctx context.Context, c *gin.Context) error {
	var req request.UserResetPassword

	if !bind(c, &req) {
		return nil
	}

	err := u.userUseCase.ResetPassword(ctx, &req)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (u User) Login(ctx context.Context, c *gin.Context) error {
	var req request.UserLogin

	if !bind(c, &req) {
		return nil
	}

	res, err := u.userUseCase.Login(ctx, &req)
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

func (u User) RefreshToken(_ context.Context, c *gin.Context) error {
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
