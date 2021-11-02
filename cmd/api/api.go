package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ken109/gin-jwt"
	"github.com/nimil-jp/gin-utils/http/middleware"
	"github.com/nimil-jp/gin-utils/http/router"

	"go-gin-ddd/config"
	"go-gin-ddd/driver/rdb"
	"go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/log"
	"go-gin-ddd/infrastructure/persistence"
	"go-gin-ddd/interface/handler"
	"go-gin-ddd/usecase"
)

func Execute() {
	logger := log.Logger()

	err := jwt.SetUp(
		jwt.Option{
			Realm:            config.DefaultRealm,
			SigningAlgorithm: jwt.HS256,
			SecretKey:        []byte(config.Env.App.Secret),
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Info("Succeeded in setting up JWT.")

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	engine.Use(middleware.Log(log.ZapLogger(), time.RFC3339, false))
	engine.Use(middleware.RecoveryWithLog(log.ZapLogger(), true))

	// cors
	engine.Use(middleware.Cors(nil))

	// cookie
	engine.Use(middleware.Session(config.UserSession, config.Env.App.Secret, nil))

	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// persistence
	userPersistence := persistence.NewUser()

	// ----- use case -----
	userUseCase := usecase.NewUser(emailDriver, userPersistence)

	// ----- handler -----
	userHandler := handler.NewUser(userUseCase)

	r := router.New(engine, rdb.Get)

	// routes
	r.Group("user", nil, func(r *router.Router) {
		r.Post("", userHandler.Create)
		r.Post("login", userHandler.Login)
		r.Post("refresh-token", userHandler.RefreshToken)
		r.Patch("reset-password-request", userHandler.ResetPasswordRequest)
		r.Patch("reset-password", userHandler.ResetPassword)
	})

	r.Group("", []gin.HandlerFunc{middleware.Authentication(config.DefaultRealm)}, func(r *router.Router) {
		r.Group("user", nil, func(r *router.Router) {
			r.Get("me", userHandler.GetMe)
		})
	})

	logger.Info("Succeeded in setting up routes.")

	// serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
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
