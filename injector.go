package main

import (
	"github.com/gin-gonic/gin"
	"go-ddd/infrastructure/email"
	"go-ddd/infrastructure/persistence"
	"go-ddd/interface/handler"
	"go-ddd/usecase"
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
