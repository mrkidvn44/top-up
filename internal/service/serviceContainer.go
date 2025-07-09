package service

import (
	"top-up-api/config"
	grpcClient "top-up-api/internal/grpc/client"
	"top-up-api/internal/repository"
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
	Validator validator.Interface

	// Services
	ProviderService        ProviderService
	SkuService             SkuService
	PurchaseHistoryService PurchaseHistoryService
	OrderService           OrderService
}

// NewContainer creates and initializes all dependencies
func NewContainer(
	database *gorm.DB,
	logger logger.Interface,
	redis redis.Interface,
	validator validator.Interface,
	config *config.Config,
	grpcClients grpcClient.GRPCServiceClient,
) *Container {

	// Initialize repositories
	providerRepository := repository.NewProviderRepository(database)
	skuRepository := repository.NewSkuRepository(database)
	purchaseHistoryRepository := repository.NewPurchaseHistoryRepository(database)

	// Initialize services
	providerService := NewProviderService(providerRepository)
	skuService := NewSkuService(skuRepository)
	purchaseHistoryService := NewPurchaseHistoryService(purchaseHistoryRepository)
	orderService := NewOrderService(skuRepository, purchaseHistoryRepository, redis, grpcClients.ProviderGRPCClient)

	return &Container{
		// Core dependencies
		DB:        database,
		Redis:     redis,
		Logger:    logger,
		Validator: validator,

		// Services
		ProviderService:        providerService,
		SkuService:             skuService,
		PurchaseHistoryService: purchaseHistoryService,
		OrderService:           orderService,
	}
}
