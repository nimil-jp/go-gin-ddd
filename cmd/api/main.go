package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ken109/gin-jwt"

	"go-gin-ddd/config"
	"go-gin-ddd/constant"
	"go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/log"
	"go-gin-ddd/infrastructure/persistence"
	"go-gin-ddd/interface/handler"
	"go-gin-ddd/interface/middleware"
	"go-gin-ddd/usecase"
)

func main() {
	logger := log.Logger()

	err := jwt.SetUp(
		jwt.Option{
			Realm:            constant.DefaultRealm,
			SigningAlgorithm: jwt.HS256,
			SecretKey:        []byte(config.Env.App.Secret),
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Info("Succeeded in setting up JWT.")

	engine := gin.New()

	engine.Use(middleware.Log(logger, time.RFC3339, false))
	engine.Use(middleware.RecoveryWithLog(logger, true))

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// cors
	engine.Use(
		cors.New(
			cors.Config{
				AllowOriginFunc: func(origin string) bool {
					return true
				},
				AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// persistence
	userPersistence := persistence.NewUser()

	// ----- use case -----
	userUseCase := usecase.NewUser(emailDriver, userPersistence)

	// ----- handler -----
	userHandler := handler.NewUser(userUseCase)

	// routes
	{
		user := engine.Group("user")
		post(user, "", userHandler.Create)
		post(user, "login", userHandler.Login)
		get(user, "refresh-token", userHandler.RefreshToken)
		patch(user, "reset-password-request", userHandler.ResetPasswordRequest)
		patch(user, "reset-password", userHandler.ResetPassword)
	}

	logger.Info("Succeeded in setting up routes.")

	// serve
	var port = ":8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Info("Succeeded in listen and serve.")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %+v", err)
	}

	logger.Info("Server exiting")
}
