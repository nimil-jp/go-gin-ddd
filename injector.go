package main

import (
	"github.com/gin-gonic/gin"
	"go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/persistence"
	"go-gin-ddd/interface/handler"
	"go-gin-ddd/usecase"
)

func inject(engine *gin.Engine) {
	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// persistence
	userPersistence := persistence.NewUser()

	// ----- use case -----
	userUseCase := usecase.NewUser(emailDriver, userPersistence)

	// ----- handler -----
	handler.NewUser(engine.Group("user"), userUseCase)
}
