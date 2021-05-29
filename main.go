package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go-ddd/config"
	"go-ddd/infra/persistence"
	"go-ddd/interface/handler"
	"go-ddd/usecase"
)

func main() {
	r := gin.New()

	r.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// cors
	r.Use(
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

	// cookie
	var (
		cookieSecure   bool
		cookieSameSite http.SameSite
	)

	switch gin.Mode() {
	case gin.ReleaseMode:
		cookieSecure = true
		cookieSameSite = http.SameSiteNoneMode
		break
	case gin.DebugMode:
		cookieSecure = false
		cookieSameSite = http.SameSiteLaxMode
		break
	}

	store := cookie.NewStore([]byte(config.Env.App.Secret))
	store.Options(
		sessions.Options{
			Path:     "/",
			MaxAge:   60 * 60 * 2,
			Secure:   cookieSecure,
			HttpOnly: true,
			SameSite: cookieSameSite,
		},
	)
	r.Use(sessions.Sessions(gin.AuthUserKey, store))

	// dependencies injection
	// persistence
	taskPersistence := persistence.NewTask()

	// use case
	taskUseCase := usecase.NewTask(taskPersistence)

	// handler
	taskHandler := handler.NewTask(taskUseCase)

	// define routes
	{
		task := r.Group("task")
		task.POST("", taskHandler.Create)
		task.GET("", taskHandler.GetAll)
		task.PUT(":id", taskHandler.Update)
	}

	var port = ":8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		// serve
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
