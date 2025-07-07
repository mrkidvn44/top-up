package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	_workerCount = 50
)

type Consumer interface {
	Consume(ctx context.Context, topic string, groupID string, handler func(msg *kafka.Message) error, errHandler func(err error)) error
	Close() error
}

type Producer interface {
	Produce(ctx context.Context, topic string, key string, value interface{}) error
	Close() error
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
}

type KafkaProducer struct {
	producer *kafka.Producer
}

var _ Consumer = (*KafkaConsumer)(nil)
var _ Producer = (*KafkaProducer)(nil)

func NewKafkaConsumer(brokers string, groupID string) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer}, nil
}

func NewKafkaProducer(brokers string) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

func (k *KafkaConsumer) Consume(ctx context.Context, topic string, groupID string, handler func(msg *kafka.Message) error, errHandler func(err error)) error {
	if err := k.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}

	msgCh := make(chan *kafka.Message)
	defer close(msgCh)

	for range _workerCount {
		go worker(msgCh, handler, errHandler)
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		msg, err := k.consumer.ReadMessage(-1)
		if err != nil {
			if errHandler != nil {
				errHandler(err)
				continue
			}
			return err
		}

		select {
		case msgCh <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (k *KafkaConsumer) Close() error {
	if k.consumer != nil {
		return k.consumer.Close()
	}
	return nil
}

func (k *KafkaProducer) Produce(ctx context.Context, topic string, key string, value interface{}) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value.(string)),
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err := k.producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func (k *KafkaProducer) Close() error {
	if k.producer != nil {
		k.producer.Flush(15 * 1000) // Wait for all messages to be delivered
		k.producer.Close()
	}
	return nil
}

func worker(msgCh <-chan *kafka.Message, handler func(msg *kafka.Message) error, errHandler func(err error)) {
	for msg := range msgCh {
		if err := handler(msg); err != nil && errHandler != nil {
			errHandler(err)
		}
	}
}
