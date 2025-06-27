package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"top-up-api/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Interface interface {
	PublishMessage(topic string, message interface{}) error
	SubscribeToTopic(ctx context.Context, topic string, callback func(message []byte) error) error
	Close() error
}

type KafkaClient struct {
	brokers []string
	logger  logger.Interface
}

func NewKafka(brokers []string, logger logger.Interface) *KafkaClient {
	return &KafkaClient{
		brokers: brokers,
		logger:  logger,
	}
}

func (k *KafkaClient) PublishMessage(topic string, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		k.logger.Error(errors.New("failed to marshal message: "), zap.Error(err))
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	conn, err := kafka.DialLeader(context.Background(), "tcp", k.brokers[0], topic, 0)
	if err != nil {
		k.logger.Error(errors.New("failed to dial leader: "), zap.Error(err))
		return fmt.Errorf("failed to dial leader: %w", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: jsonData},
	)
	if err != nil {
		k.logger.Error(errors.New("failed to write message: "), zap.Error(err))
		return fmt.Errorf("failed to write message: %w", err)
	}
	k.logger.Info(fmt.Sprintf("message published to topic: %s, message: %s", topic, string(jsonData)))
	return nil
}

func (k *KafkaClient) SubscribeToTopic(ctx context.Context, topic string, handler func(message []byte) error) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokers,
		Topic:   topic,
		GroupID: "top-up-api",
	})

	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			k.logger.Info("context done")
			return nil
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				k.logger.Error(errors.New("failed to read message: "), zap.Error(err))
				continue
			}

			k.logger.Info(fmt.Sprintf("message read from topic: %s, message: %s", topic, string(msg.Value)))
			if err := handler(msg.Value); err != nil {
				k.logger.Error(errors.New("failed to handle message: "), zap.Error(err))
				continue
			}
		}
	}
}

func (k *KafkaClient) Close() error {
	return nil
}
