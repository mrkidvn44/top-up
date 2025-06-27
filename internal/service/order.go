package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"top-up-api/internal/model"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
	"top-up-api/pkg/redis"

	"github.com/google/uuid"
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

	// Capture the current timestamp in a special format
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)

	// Obtain a unique code using UUID
	uniqueID := uuid.New().ID()
	cardDetailResponse := schema.CardDetailResponseFromModel(*cardDetail)
	orderID := uint(currentTimestamp) + uint(uniqueID)
	orderResponse := schema.OrderResponse{
		OrderID:            orderID,
		UserID:             order.UserID,
		CardDetailResponse: *cardDetailResponse,
		TotalPrice:         cardDetailResponse.CardPriceResponse.Value, //MAYBE COUPON LATER
		Status:             model.PurchaseHistoryStatusPending,
		PhoneNumber:        order.PhoneNumber,
		CashBackValue:      cardDetailResponse.CashBackResponse.CalculateCashBack(cardDetail.CardPrice.Value),
	}
	orderResponseJSON, err := json.Marshal(orderResponse)
	if err != nil {
		return nil, errors.New("failed to marshal order response: " + err.Error())
	}
	err = s.redisClient.Set(ctx, CacheKey+strconv.Itoa(int(orderID)), orderResponseJSON, CacheTime)
	if err != nil {
		return nil, err
	}
	return &orderResponse, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error {
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

	s.redisClient.Del(ctx, CacheKey+strconv.Itoa(int(orderConfirmRequest.OrderID)))

	purchaseHistory := schema.PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest)
	err = s.purchaseHistoryRepo.CreatePurchaseHistory(ctx, purchaseHistory)
	if err != nil {
		return err
	}
	return nil
}
