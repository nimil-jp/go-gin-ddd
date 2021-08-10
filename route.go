package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
	"github.com/pkg/errors"
	"go-ddd/constant"
	"go-ddd/infrastructure/persistence"
	"go-ddd/interface/handler"
	xerrors2 "go-ddd/pkg/xerrors"
	"go-ddd/usecase"
)

func registerRoute(engine *gin.Engine) {
	// dependencies injection
	// persistence
	userPersistence := persistence.NewUser()

	// use case
	userUseCase := usecase.NewUser(userPersistence)

	// handler
	userHandler := handler.NewUser(userUseCase)

	// define routes
	{
		user := engine.Group("user")
		post(user, "", userHandler.Create)
		post(user, "login", userHandler.Login)
		get(user, "refresh-token", userHandler.RefreshToken)
		patch(user, "reset-password-request", userHandler.ResetPasswordRequest)
		patch(user, "reset-password", userHandler.ResetPassword)

		userA := user.Group("")
		userA.Use(jwt.Verify(constant.DefaultRealm))
		userA.GET(
			"auth", func(c *gin.Context) {
				c.Status(200)
			},
		)
	}
}

type handlerFunc func(c *gin.Context) error

func get(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.GET(relativePath, hf(handlerFunc))
}

func post(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.POST(relativePath, hf(handlerFunc))
}

func put(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.PUT(relativePath, hf(handlerFunc))
}

func patch(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.PATCH(relativePath, hf(handlerFunc))
}

func delete(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.DELETE(relativePath, hf(handlerFunc))
}

func options(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.OPTIONS(relativePath, hf(handlerFunc))
}

func head(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.HEAD(relativePath, hf(handlerFunc))
}

func hf(handlerFunc handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handlerFunc(c)

		if err != nil {
			switch v := err.(type) {
			case *xerrors2.Expected:
				if v.StatusOk() {
					return
				} else {
					c.JSON(v.StatusCode(), v.Message())
				}
			case *xerrors2.Validation:
				c.JSON(http.StatusBadRequest, v)
			default:
				if gin.Mode() == gin.DebugMode {
					c.String(http.StatusInternalServerError, "%+v", v)
				} else {
					c.Status(http.StatusInternalServerError)
				}
			}

			_ = c.Error(errors.Errorf("%+v", err))
		}
	}
}
