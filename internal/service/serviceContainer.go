package service

import (
	"top-up-api/internal/repository"
	"top-up-api/pkg/auth"
	"top-up-api/pkg/logger"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/validator"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
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
) *Container {
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
	orderService := NewOrderService(*cardDetailRepository, *purchaseHistoryRepository, redis)

	return &Container{
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
