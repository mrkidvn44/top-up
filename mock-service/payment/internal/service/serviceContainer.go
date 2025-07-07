package service

import (
	"payment-api/config"
	kfk "payment-api/pkg/kafka"
	"payment-api/pkg/logger"
)

// Container holds all application dependencies
type Container struct {
	// Config
	config *config.Config

	// Core dependencies
	Logger logger.Interface

	// Services
	OrderService *OrderService
}

// NewContainer creates and initializes all dependencies
func NewContainer(
	logger logger.Interface,
	config *config.Config,
) *Container {
	// Initialize Kafka factories
	kafkaProducerFactory := kfk.NewProducerFactory(&config.Kafka)
	orderKafkaProducer, err := kafkaProducerFactory.CreateProducer()
	if err != nil {
		logger.Error(err)
	}
	// Initialize services
	orderService := NewOrderService(logger, orderKafkaProducer)

	return &Container{
		config: config,
		// Core dependencies
		Logger: logger,

		// Services
		OrderService: orderService,
	}
}
