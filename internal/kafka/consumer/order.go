package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	kfk "top-up-api/pkg/kafka"
	"top-up-api/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type OrderConsumer struct {
	logger   logger.Interface
	service  service.OrderService
	consumer kfk.Consumer
}

func NewOrderConsumer(l logger.Interface, s service.OrderService, c kfk.Consumer) *OrderConsumer {
	return &OrderConsumer{logger: l, service: s, consumer: c}
}

func (c *OrderConsumer) StartOrderConfirmConsumer(ctx context.Context, topic, groupID string) error {
	if err := c.consumer.Consume(ctx, topic, groupID, func(msg *kafka.Message) error {
		var orderConfirmRequest schema.OrderConfirmRequest
		if err := json.Unmarshal(msg.Value, &orderConfirmRequest); err != nil {
			c.logger.Warn("failed to unmarshal order confirm event: ", zap.Error(err))
			return nil
		}
		if err := c.service.ConfirmOrder(ctx, orderConfirmRequest); err != nil {
			c.logger.Error(errors.New("failed to process confirm event: "), zap.Error(err))
		}
		return nil
	}, func(err error) {
		c.logger.Warn("consume: ", zap.Error(err))
	}); err != nil {
		return fmt.Errorf("failed to start order confirm consumer: %w", err)
	}

	return nil
}

func (c *OrderConsumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return err
	}
	return nil
}
