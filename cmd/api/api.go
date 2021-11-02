package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	middle "go-gin-ddd/interface/middleware"
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
	engine.Use(
		cors.New(
			cors.Config{
				AllowOriginFunc: func(origin string) bool {
					return true
				},
				AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Request-Id"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)

	// cookie
	var (
		corsSecure   bool
		corsSameSite http.SameSite
	)

	switch gin.Mode() {
	case gin.ReleaseMode:
		corsSecure = true
		corsSameSite = http.SameSiteNoneMode
	case gin.DebugMode:
		corsSecure = false
		corsSameSite = http.SameSiteLaxMode
	}

	store := cookie.NewStore([]byte(config.Env.App.Secret))
	store.Options(
		sessions.Options{
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 365,
			Secure:   corsSecure,
			HttpOnly: true,
			SameSite: corsSameSite,
		},
	)
	engine.Use(sessions.Sessions(config.UserSession, store))

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

	r.Group("", []gin.HandlerFunc{middle.Authentication}, func(r *router.Router) {
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
