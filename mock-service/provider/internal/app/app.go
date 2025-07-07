package app

import (
	"fmt"
	"os"
	"os/signal"
	"provider-api/config"
	controller "provider-api/internal/controller/http"
	"provider-api/internal/service"
	"provider-api/pkg/httpserver"
	"provider-api/pkg/logger"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// Logger
	logger := logger.New(cfg.Log.Level)

	// Middleware
	services := service.NewContainer(logger)
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
	err := httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
