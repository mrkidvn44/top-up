package app

import (
	"fmt"
	"os"
	"os/signal"
	"payment-api/config"
	controller "payment-api/internal/controller/http"
	"payment-api/internal/service"
	"payment-api/pkg/httpserver"
	"payment-api/pkg/logger"
	orderpb "payment-api/proto/order"
	"syscall"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(cfg *config.Config) {

	logger := logger.New(cfg.Log.Level)

	// Middleware
	services := service.NewContainer(logger, cfg)
	// gRPC Client
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - grpc.Dial: %w", err))
		return
	}
	defer conn.Close()
	grpcOrderService := orderpb.NewOrderServiceClient(conn)
	// HTTP Server
	handler := gin.Default()
	controller.NewRouter(handler, services, grpcOrderService)

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
