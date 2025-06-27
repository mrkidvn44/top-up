package app

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"top-up-api/config"
	controller "top-up-api/internal/controller/http"
	"top-up-api/internal/db"
	"top-up-api/internal/service"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/httpserver"
	"top-up-api/pkg/logger"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/validator"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	env := os.Getenv("ENV")
	if env == "PROD" {
		err := os.MkdirAll("./.log", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile(
			"./.log/server.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		gin.DefaultWriter = io.MultiWriter(os.Stdout, file)

	}
	logger := logger.New(cfg.Log.Level)

	//connect postgres with gorm
	db, err := db.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	validator := validator.NewValidator()
	// Middleware
	redis := redis.NewRedis(cfg.Redis)
	logger.Info(fmt.Sprintf("redis connected to %s", redis))
	auth := auth.NewAuthService([]byte(cfg.JWT.Secret))
	logger.Info(fmt.Sprintf("auth service %s", auth))
	services := service.NewContainer(db.Database, logger, redis, auth, validator)
	// HTTP Server
	handler := gin.Default()
	controller.NewRouter(handler, services)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
