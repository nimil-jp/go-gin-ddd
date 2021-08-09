package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
	"github.com/pkg/errors"
	"go-ddd/constant"
	"go-ddd/infrastructure/persistence"
	"go-ddd/interface/handler"
	"go-ddd/usecase"
	"go-ddd/util/xerrors"
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
			var (
				eerr xerrors.Expected
				verr xerrors.Validation
			)

			if errors.As(err, &eerr) {
				c.JSON(eerr.StatusCode(), eerr.Message())
			} else if errors.As(err, &verr) {
				c.JSON(http.StatusBadRequest, verr)
			} else {
				if gin.Mode() == gin.DebugMode {
					c.JSON(http.StatusInternalServerError, err)
				} else {
					c.Status(http.StatusInternalServerError)
				}
			}
			log.Printf("%+v\n", err)
			return
		}
	}
}
