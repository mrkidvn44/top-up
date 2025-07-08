package consumer

import (
	"context"
	"errors"
	"fmt"
	"top-up-api/config"
	"top-up-api/internal/service"
	kfk "top-up-api/pkg/kafka"
	"top-up-api/pkg/logger"
)

// Container holds all application dependencies
type Consumers struct {
	// Config
	config *config.Kafka
	// Dependency
	logger logger.Interface
	// Consumer
	orderConsumer *OrderConsumer
}

// NewContainer creates and initializes all dependencies
func NewConsumers(
	config *config.Kafka,
	services *service.Container,
) *Consumers {
	// Initialize Kafka factories
	kafkaConsumerFactory := kfk.NewConsumerFactory(config)
	orderKafkaConsumer, err := kafkaConsumerFactory.CreateConsumer(kfk.ServiceOrder)
	if err != nil {
		services.Logger.Error(err)
	}
	orderConsumer := NewOrderConsumer(services.Logger, services.OrderService, orderKafkaConsumer)

	return &Consumers{
		// Config
		config: config,

		// Dependency
		logger: services.Logger,

		// Consumer
		orderConsumer: orderConsumer,
	}
}

func (c *Consumers) StartKafkaConsumers(ctx context.Context) {
	c.logger.Info("Starting Kafka consumers for all services...")
	baseGroupID := c.config.GroupID
	// Start OrderService Kafka consumers
	go func() {
		if err := c.orderConsumer.StartOrderConfirmConsumer(ctx, c.config.OrderGroup.ConfirmTopic, baseGroupID); err != nil {
			c.logger.Error(fmt.Errorf("consumers: %w", err))
		}
	}()

	c.logger.Info("All service Kafka consumers started successfully")
}

func (c *Consumers) CloseKafkaConsumers() error {
	var errs []error

	if c.orderConsumer != nil {
		if err := c.orderConsumer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
