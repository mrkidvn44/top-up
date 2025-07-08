package service

import (
	"top-up-api/config"
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
	SkuService             *SkuService
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

	// Initialize repositories
	userRepository := repository.NewUserRepository(database)
	providerRepository := repository.NewProviderRepository(database)
	skuRepository := repository.NewSkuRepository(database)
	purchaseHistoryRepository := repository.NewPurchaseHistoryRepository(database)

	// Initialize services
	userService := NewUserService(userRepository)
	providerService := NewProviderService(providerRepository)
	skuService := NewSkuService(skuRepository)
	purchaseHistoryService := NewPurchaseHistoryService(purchaseHistoryRepository)
	orderService := NewOrderService(skuRepository, purchaseHistoryRepository, redis)

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
		SkuService:             skuService,
		PurchaseHistoryService: purchaseHistoryService,
		OrderService:           orderService,
	}
}
