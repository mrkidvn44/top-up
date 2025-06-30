package kafka

import (
	"fmt"
	"top-up-api/config"
)

const ServiceOrder = "order"

type ConsumerFactory struct {
	config *config.Kafka
}

type ProducerFactory struct {
	config *config.Kafka
}

func NewConsumerFactory(cfg *config.Kafka) *ConsumerFactory {
	return &ConsumerFactory{
		config: cfg,
	}
}

func NewProducerFactory(cfg *config.Kafka) *ProducerFactory {
	return &ProducerFactory{
		config: cfg,
	}
}

func (f *ConsumerFactory) CreateConsumer(serviceName string) (Consumer, error) {
	consumer, err := NewKafkaConsumer(f.config.Brokers, f.config.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s Kafka consumer: %w", serviceName, err)
	}
	return consumer, nil
}

func (f *ProducerFactory) CreateProducer() (Producer, error) {
	producer, err := NewKafkaProducer(f.config.Brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	return producer, nil
}
