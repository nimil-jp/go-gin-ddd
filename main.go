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
	"github.com/gin-gonic/gin"
	"github.com/ken109/gin-jwt"
	"go-ddd/config"
	"go-ddd/constant"
)

func main() {
	err := jwt.SetUp(
		jwt.Option{
			Realm:            constant.DefaultRealm,
			SigningAlgorithm: jwt.HS256,
			SecretKey:        []byte(config.Env.App.Secret),
			Timeout:          time.Second * 30,
		},
	)
	if err != nil {
		panic(err)
	}

	engine := gin.New()

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

	registerRoute(engine)

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
