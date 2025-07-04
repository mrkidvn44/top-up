package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"top-up-api/internal/mapper"
	"top-up-api/internal/model"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
	kfk "top-up-api/pkg/kafka"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/util"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gorm.io/gorm"
)

const (
	_lockTimeOut      = 5 * time.Minute
	_cacheTime        = 30 * time.Minute
	_cacheKey         = "order_id"
	_cashierCreateURL = "http://localhost:8081/v1/api/order/create"
	_cashierUpdateURL = "http://localhost:8081/v1/api/order/update"
	_providerURL      = "http://localhost:8082/v1/api/order/"
	_callbackURL      = "http://localhost:8080/v1/api/order/update-status"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error)
	ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error
	UpdateOrderStatus(ctx context.Context, orderUpdateInfo schema.OrderUpdateRequest) error
	StartOrderConfirmConsumer(ctx context.Context, topic, groupID string) error
}

type OrderService struct {
	cardDetailRepo       repository.ICardDetailRepository
	purchaseHistoryRepo  repository.IPurchaseHistoryRepository
	redisClient          redis.Interface
	orderConfirmConsumer kfk.Consumer
}

var _ IOrderService = (*OrderService)(nil)

func NewOrderService(
	cardDetailRepo repository.ICardDetailRepository,
	purchaseHistoryRepo repository.IPurchaseHistoryRepository,
	redisClient redis.Interface,
	orderConfirmConsumer kfk.Consumer,
) *OrderService {
	return &OrderService{
		cardDetailRepo:       cardDetailRepo,
		purchaseHistoryRepo:  purchaseHistoryRepo,
		redisClient:          redisClient,
		orderConfirmConsumer: orderConfirmConsumer,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error) {
	cardDetail, err := s.cardDetailRepo.GetCardDetailByID(ctx, order.CardDetailID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("card detail not found")
		}
		return nil, err
	}

	orderID := util.GenerateOrderID()
	orderResponse := mapper.OrderResponseFromOrderRequest(order, cardDetail, orderID)

	orderResponseJSON, err := json.Marshal(orderResponse)
	if err != nil {
		return nil, errors.New("failed to marshal order response: " + err.Error())
	}

	cacheKey := getCachKey(strconv.Itoa(int(orderID)))
	err = s.redisClient.Set(ctx, cacheKey, orderResponseJSON, _cacheTime)
	if err != nil {
		return nil, err
	}

	go util.SendPostRequest(_cashierCreateURL, orderResponseJSON)

	return orderResponse, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error {
	orderID := strconv.Itoa(int(orderConfirmRequest.OrderID))
	err := s.redisClient.TryAcquireLock(ctx, orderID, _lockTimeOut)
	if err != nil {
		return err
	}
	defer s.redisClient.ReleaseLock(ctx, orderID)

	cacheKey := getCachKey(orderID)

	orderResponse, err := s.getCachedOrder(ctx, cacheKey)
	if err != nil {
		return err
	}
	if !orderResponse.CompareWithOrderConfirmRequest(orderConfirmRequest) {
		return errors.New("order mismatch")
	}

	if orderConfirmRequest.Status == model.PurchaseHistoryStatusPending {
		return errors.New("order is pending")
	}

	if orderResponse.Status != model.PurchaseHistoryStatusPending {
		return errors.New("order already confirmed or failed")
	}

	purchaseHistory := mapper.PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest)
	err = s.purchaseHistoryRepo.CreatePurchaseHistory(ctx, purchaseHistory)
	if err != nil {
		return err
	}

	orderResponse.Status = orderConfirmRequest.Status
	err = s.updateCacheOrderStaus(ctx, cacheKey, orderResponse)
	if err != nil {
		return err
	}

	if orderConfirmRequest.Status == model.PurchaseHistoryStatusConfirm {
		go sendOrderResponse(_providerURL, orderResponse)
	}

	return nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderUpdateInfo schema.OrderUpdateRequest) error {
	orderID := strconv.Itoa(int(orderUpdateInfo.OrderID))
	err := s.redisClient.TryAcquireLock(ctx, orderID, _lockTimeOut)
	if err != nil {
		return err
	}
	defer s.redisClient.ReleaseLock(ctx, orderID)

	cacheKey := getCachKey(orderID)

	orderResponse, err := s.getCachedOrder(ctx, cacheKey)
	if err != nil {
		return err
	}

	if orderResponse.Status != model.PurchaseHistoryStatusConfirm {
		return errors.New("order is not confirmed or failed")
	}

	err = s.purchaseHistoryRepo.UpdatePurchaseHistoryStatusByOrderID(ctx, orderUpdateInfo.OrderID, orderUpdateInfo.Status)
	if err != nil {
		return err
	}

	orderResponse.Status = orderUpdateInfo.Status
	err = s.updateCacheOrderStaus(ctx, cacheKey, orderResponse)
	if err != nil {
		return err
	}

	if orderUpdateInfo.Status == model.PurchaseHistoryStatusFailed {
		go sendFailedOrder(_cashierUpdateURL, orderUpdateInfo)
	}

	return nil
}

func (s *OrderService) StartOrderConfirmConsumer(ctx context.Context, topic, groupID string) error {
	if err := s.orderConfirmConsumer.Consume(ctx, topic, groupID, func(msg *kafka.Message) error {
		var orderConfirmRequest schema.OrderConfirmRequest
		if err := json.Unmarshal(msg.Value, &orderConfirmRequest); err != nil {
			fmt.Printf("failed to unmarshal order confirm event: %v \n", err)
			return nil
		}
		if err := s.ConfirmOrder(ctx, orderConfirmRequest); err != nil {
			fmt.Printf("failed to process confirm event: %v \n", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to start order confirm consumer: %w", err)
	}

	return nil
}

func (s *OrderService) getCachedOrder(ctx context.Context, cacheKey string) (*schema.OrderResponse, error) {

	order, err := s.redisClient.Get(ctx, cacheKey)
	if err != nil {
		return nil, errors.New("order not found or expired")
	}

	var orderResponse *schema.OrderResponse
	err = json.Unmarshal([]byte(order), &orderResponse)
	if err != nil {
		return orderResponse, errors.New("failed to unmarshal order: " + err.Error())
	}

	return orderResponse, nil
}

func (s *OrderService) updateCacheOrderStaus(ctx context.Context, cacheKey string, orderResponse *schema.OrderResponse) error {
	orderResponseJSON, err := json.Marshal(orderResponse)

	if err != nil {
		return errors.New("failed to marshal order response: " + err.Error())
	}

	err = s.redisClient.Set(ctx, cacheKey, orderResponseJSON, _cacheTime)
	if err != nil {
		return err
	}

	return nil
}

func sendOrderResponse(url string, orderResponse *schema.OrderResponse) error {
	orderProviderRequest := mapper.OrderProviderRequestFromOrderResponse(orderResponse, _callbackURL)
	orderProviderRequestJSON, err := json.Marshal(orderProviderRequest)

	if err != nil {
		return errors.New("failed to marshal order response: " + err.Error())
	}

	util.SendPostRequest(url, orderProviderRequestJSON)
	return nil
}

func sendFailedOrder(url string, orderUpdateInfo schema.OrderUpdateRequest) {
	payload, _ := json.Marshal(map[string]interface{}{
		"order_id": orderUpdateInfo.OrderID,
		"status":   orderUpdateInfo.Status,
	})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
}

func getCachKey(orderID string) string {
	return _cacheKey + orderID
}
