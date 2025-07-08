package service

import (
	"bytes"
	"cashier-api/internal/schema"
	kfk "cashier-api/pkg/kafka"
	"cashier-api/pkg/logger"
	orderpb "cashier-api/proto/order"
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

const (
	CacheTime = 30 * time.Minute
	CacheKey  = "order_id"
)

type IOrderService interface {
	CreateOrder(orderData string, grpcClient orderpb.OrderServiceClient) error
}

type OrderService struct {
	logger             logger.Interface
	orderKafkaProducer kfk.Producer
}

var _ IOrderService = (*OrderService)(nil)

func NewOrderService(l logger.Interface, orderProducer kfk.Producer) *OrderService {
	return &OrderService{
		logger:             l,
		orderKafkaProducer: orderProducer,
	}
}

func (s *OrderService) CreateOrder(orderData string, grpcClient orderpb.OrderServiceClient) error {
	var order schema.Order
	if err := json.Unmarshal([]byte(orderData), &order); err != nil {
		fmt.Print(err)
		return err
	}

	order.Status = "pending"
	pendingOrderDataJson, err := json.Marshal(order)
	if err != nil {
		return err
	}

	go http.Post("http://localhost:8080/v1/api/order/confirm", "application/json", bytes.NewBuffer(pendingOrderDataJson))
	// Change the status field
	if rand.IntN(100) < 98 {
		order.Status = "confirm" // Set to the desired status

		// Marshal back to JSON
		confirmOrderDataJson, err := json.Marshal(order)
		if err != nil {
			return err
		}

		ctx := context.Background()
		defer ctx.Done()

		go func(orderJson []byte) {
			time.Sleep(time.Duration(rand.IntN(1500)) * time.Millisecond)
			s.orderKafkaProducer.Produce(ctx, "my-topic", strconv.Itoa(int(order.OrderID)), string(orderJson))
		}(confirmOrderDataJson)

		order.Status = "failed"

		go func(order schema.Order) {
			time.Sleep(time.Duration(rand.IntN(1500)) * time.Millisecond)
			grpcClient.ConfirmOrder(ctx, &orderpb.OrderConfirmRequest{
				OrderId:       uint64(order.OrderID),
				UserId:        uint64(order.UserID),
				CardDetailId:  uint64(order.CardDetailID),
				TotalPrice:    int64(order.TotalPrice),
				Status:        order.Status,
				PhoneNumber:   order.PhoneNumber,
				CashBackValue: int64(order.CashBackValue),
			})
		}(order)
	}
	return nil
}
