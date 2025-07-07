package service

import (
	"provider-api/pkg/logger"
)

// Container holds all application dependencies
type Container struct {
	// Core dependencies
	Logger logger.Interface

	// Services
	OrderService *OrderService
}

// NewContainer creates and initializes all dependencies
func NewContainer(
	logger logger.Interface,
) *Container {

	// Initialize services
	orderService := NewOrderService(logger)

	return &Container{
		// Core dependencies
		Logger: logger,

		// Services
		OrderService: orderService,
	}
}
