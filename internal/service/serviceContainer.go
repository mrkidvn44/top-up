package service

import (
	"context"
	"fmt"
	"top-up-api/config"
	"top-up-api/internal/repository"
	"top-up-api/pkg/auth"
	kfk "top-up-api/pkg/kafka"
	"top-up-api/pkg/logger"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/validator"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	// Config
	config *config.Config

	// Core dependencies
	DB        *gorm.DB
	Redis     redis.Interface
	Logger    logger.Interface
	Auth      auth.Interface
	Validator validator.Interface

	// Services
	UserService            *UserService
	ProviderService        *ProviderService
	CardDetailService      *CardDetailService
	PurchaseHistoryService *PurchaseHistoryService
	OrderService           *OrderService
}

// NewContainer creates and initializes all dependencies
func NewContainer(
	database *gorm.DB,
	logger logger.Interface,
	redis redis.Interface,
	auth auth.Interface,
	validator validator.Interface,
	config *config.Config,
) *Container {
	// Initialize Kafka factories
	kafkaConsumerFactory := kfk.NewConsumerFactory(&config.Kafka)
	orderKafkaConsumer, err := kafkaConsumerFactory.CreateConsumer(kfk.ServiceOrder)
	if err != nil {
		logger.Error(err)
	}

	// Initialize repositories
	userRepository := repository.NewUserRepository(database)
	providerRepository := repository.NewProviderRepository(database)
	cardDetailRepository := repository.NewCardDetailRepository(database)
	purchaseHistoryRepository := repository.NewPurchaseHistoryRepository(database)

	// Initialize services
	userService := NewUserService(*userRepository)
	providerService := NewProviderService(*providerRepository)
	cardDetailService := NewCardDetailService(*cardDetailRepository)
	purchaseHistoryService := NewPurchaseHistoryService(*purchaseHistoryRepository)
	orderService := NewOrderService(*cardDetailRepository, *purchaseHistoryRepository, redis, orderKafkaConsumer)

	return &Container{
		// Config
		config: config,
		// Core dependencies
		DB:        database,
		Redis:     redis,
		Logger:    logger,
		Auth:      auth,
		Validator: validator,

		// Services
		UserService:            userService,
		ProviderService:        providerService,
		CardDetailService:      cardDetailService,
		PurchaseHistoryService: purchaseHistoryService,
		OrderService:           orderService,
	}
}

func (c *Container) StartKafkaConsumers(ctx context.Context) {
	c.Logger.Info("Starting Kafka consumers for all services...")
	baseGroupID := c.config.Kafka.GroupID
	// Start OrderService Kafka consumers
	go func() {
		if err := c.OrderService.StartOrderConfirmConsumer(ctx, c.config.OrderGroup.ConfirmTopic, baseGroupID); err != nil {
			c.Logger.Error(fmt.Errorf("service container: %w", err))
		}
	}()

	c.Logger.Info("All service Kafka consumers started successfully")
}
