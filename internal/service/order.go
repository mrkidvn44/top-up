package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"top-up-api/internal/mapper"
	"top-up-api/internal/model"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/util"
)

const (
	CacheTime = 30 * time.Minute
	CacheKey  = "order_id"
)

type OrderService struct {
	cardDetailRepo      repository.CardDetailRepository
	purchaseHistoryRepo repository.PurchaseHistoryRepository
	redisClient         redis.Interface
}

func NewOrderService(cardDetailRepo repository.CardDetailRepository, purchaseHistoryRepo repository.PurchaseHistoryRepository, redisClient redis.Interface) *OrderService {
	return &OrderService{cardDetailRepo: cardDetailRepo, purchaseHistoryRepo: purchaseHistoryRepo, redisClient: redisClient}
}

func (s *OrderService) CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error) {
	cardDetail, err := s.cardDetailRepo.GetCardDetailByID(ctx, order.CardDetailID)
	if err != nil {
		return nil, err
	}

	orderID := util.GenerateOrderID()
	orderResponse := mapper.OrderResponseFromOrderRequest(order, cardDetail, orderID)

	orderResponseJSON, err := json.Marshal(orderResponse)
	if err != nil {
		return nil, errors.New("failed to marshal order response: " + err.Error())
	}

	go func(payload []byte) {
		http.Post("https://mock-external-api-url", "application/json", bytes.NewBuffer(payload))
	}(orderResponseJSON)

	err = s.redisClient.Set(ctx, CacheKey+strconv.Itoa(int(orderID)), orderResponseJSON, CacheTime)
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error {
	orderID := strconv.Itoa(int(orderConfirmRequest.OrderID))
	for {
		if ok := s.redisClient.GetLock(ctx, orderID); ok {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	defer s.redisClient.ReleaseLock(ctx, orderID)

	order, err := s.redisClient.Get(ctx, CacheKey+strconv.Itoa(int(orderConfirmRequest.OrderID)))
	if err != nil {
		return errors.New("order not found or expired")
	}

	var orderResponse schema.OrderResponse
	err = json.Unmarshal([]byte(order), &orderResponse)
	if err != nil {
		return errors.New("failed to unmarshal order: " + err.Error())
	}
	if !orderResponse.CompareWithOrderConfirmRequest(orderConfirmRequest) {
		return errors.New("order mismatch")
	}

	storedOrder, err := s.purchaseHistoryRepo.GetPurchaseHistoryByOrderID(ctx, orderConfirmRequest.OrderID)
	if err == nil {
		if storedOrder.Status == model.PurchaseHistoryStatusPending {
			err = s.purchaseHistoryRepo.UpdatePurchaseHistoryStatusByOrderID(ctx, orderConfirmRequest.OrderID, orderConfirmRequest.Status)
			if err != nil {
				return err
			}
		} else {
			return errors.New("order already confirmed or failed")
		}
	}

	purchaseHistory := mapper.PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest)
	err = s.purchaseHistoryRepo.CreatePurchaseHistory(ctx, purchaseHistory)
	if err != nil {
		return err
	}

	go func(req schema.OrderProviderRequest) {
		payload, _ := json.Marshal(req)
		http.Post("https://mock-external-api-url", "application/json", bytes.NewBuffer(payload))
	}(schema.OrderProviderRequest{
		OrderID:     orderConfirmRequest.OrderID,
		PhoneNumber: orderConfirmRequest.PhoneNumber,
		TotalPrice:  orderResponse.TotalPrice,
		CardPrice:   orderResponse.CardDetail.CardPriceResponse.Value,
	})

	return nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderUpdateInfo schema.OrderUpdateRequest) error {

	err := s.purchaseHistoryRepo.UpdatePurchaseHistoryStatusWithLock(ctx, orderUpdateInfo.OrderID, orderUpdateInfo.Status)
	if err != nil {
		return err
	}

	if orderUpdateInfo.Status == model.PurchaseHistoryStatusFailed {
		go func() {
			payload, _ := json.Marshal(map[string]interface{}{
				"order_id": orderUpdateInfo.OrderID,
				"status":   orderUpdateInfo.Status,
			})
			http.Post("https://mock-api-url", "application/json", bytes.NewBuffer(payload))
		}()
	}

	return nil
}
