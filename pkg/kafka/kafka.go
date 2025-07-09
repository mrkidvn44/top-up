package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

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

type kafkaConsumer struct {
	consumer *kafka.Consumer
	wg       sync.WaitGroup
}

type kafkaProducer struct {
	producer *kafka.Producer
}

var _ Consumer = (*kafkaConsumer)(nil)
var _ Producer = (*kafkaProducer)(nil)

func NewKafkaConsumer(brokers string, groupID string) (*kafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{consumer: consumer, wg: sync.WaitGroup{}}, nil
}

func NewKafkaProducer(brokers string) (*kafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{producer: producer}, nil
}

func (k *kafkaConsumer) Consume(ctx context.Context, topic string, groupID string, handler func(msg *kafka.Message) error, errHandler func(err error)) error {
	if err := k.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}

	msgCh := make(chan *kafka.Message)
	defer close(msgCh)
	for range _workerCount {
		k.wg.Add(1)
		go worker(msgCh, handler, errHandler, &k.wg)
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		msg, err := k.consumer.ReadMessage(500 * time.Millisecond)
		if err != nil && errHandler != nil && !err.(kafka.Error).IsTimeout() {
			errHandler(err)
			continue
		}
		if msg == nil {
			continue
		}
		select {
		case msgCh <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func (k *kafkaConsumer) Close() error {
	if k.consumer != nil {
		k.wg.Wait()
		return k.consumer.Close()
	}
	return nil
}

func (k *kafkaProducer) Produce(ctx context.Context, topic string, key string, value interface{}) error {
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

func (k *kafkaProducer) Close() error {
	if k.producer != nil {
		k.producer.Flush(15 * 1000) // Wait for all messages to be delivered
		k.producer.Close()
	}
	return nil
}

func worker(msgCh <-chan *kafka.Message, handler func(msg *kafka.Message) error, errHandler func(err error), wg *sync.WaitGroup) {
	defer handlePanic(errHandler)
	defer wg.Done()

	for msg := range msgCh {
		if err := handler(msg); err != nil && errHandler != nil {
			errHandler(err)
		}
	}
}

func handlePanic(errHandler func(err error)) {
	if r := recover(); r != nil {
		if errHandler != nil {
			errHandler(fmt.Errorf("worker panic recovered: %v", r))
		}
	}
}
