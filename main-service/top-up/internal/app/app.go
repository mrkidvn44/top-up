package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"top-up-api/config"
	controller "top-up-api/internal/controller/http"
	"top-up-api/internal/db"
	myGrpc "top-up-api/internal/grpc"
	"top-up-api/internal/service"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/httpserver"
	"top-up-api/pkg/logger"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/validator"

	grpc "google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

const _waitTime = 2 * time.Minute

func Run(cfg *config.Config) {
	logger := logger.New(cfg.Log.Level, cfg.Env)

	//connect postgres with gorm
	db, err := db.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.Grpc.Port)
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - net.Listen: %w", err))
		os.Exit(1)
	}
	// Validator
	validator := validator.NewValidator()

	// Middleware
	redis := redis.NewRedis(cfg.Redis)
	logger.Info(fmt.Sprintf("redis connected to %s", redis))
	auth := auth.NewAuthService([]byte(cfg.JWT.Secret))
	logger.Info(fmt.Sprintf("auth service %s", auth))

	// Services
	services := service.NewContainer(db.Database, logger, redis, auth, validator, cfg)

	// Start Kafka consumers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	services.StartKafkaConsumers(ctx)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	grpcServices := myGrpc.NewGRPCServiceContainer(services)
	grpcServices.Register(grpcServer)
	go grpcServer.Serve(lis)

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
	// httpServer
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	// Grpc Server
	grpcServer.GracefulStop()

	// Local listener
	err = lis.Close()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - lis.Close: %w", err))
	}

	// Kafka service
	err = services.CloseKafka()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - services.CloseKafka: %w", err))
	}

	// Database connection
	err = db.Close()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - db.Close: %w", err))
	}
	if cfg.Env != "dev" {
		time.Sleep(_waitTime)
	}
}
