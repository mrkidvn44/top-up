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

	"gorm.io/gorm"
)

const (
	_lockTimeOut      = 5 * time.Minute
	_cacheTime        = 30 * time.Minute
	_cacheKey         = "order_id"
	_paymentCreateURL = "http://localhost:8081/v1/api/order/create"
	_paymentUpdateURL = "http://localhost:8081/v1/api/order/update"
	_providerURL      = "http://localhost:8082/v1/api/order/"
	_callbackURL      = "http://localhost:8080/v1/api/order/update-status"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error)
	ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error
	UpdateOrderStatus(ctx context.Context, orderUpdateInfo schema.OrderUpdateRequest) error
}

type OrderService struct {
	skuRepo             repository.ISkuRepository
	purchaseHistoryRepo repository.IPurchaseHistoryRepository
	redisClient         redis.Interface
}

var _ IOrderService = (*OrderService)(nil)

func NewOrderService(
	skuRepo repository.ISkuRepository,
	purchaseHistoryRepo repository.IPurchaseHistoryRepository,
	redisClient redis.Interface,
) *OrderService {
	return &OrderService{
		skuRepo:             skuRepo,
		purchaseHistoryRepo: purchaseHistoryRepo,
		redisClient:         redisClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error) {
	sku, err := s.skuRepo.GetSkuByID(ctx, order.SkuID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("sku not found")
		}
		return nil, err
	}

	orderID := util.GenerateOrderID()
	orderResponse := mapper.OrderResponseFromOrderRequest(order, sku, orderID)

	orderResponseJSON, err := json.Marshal(orderResponse)
	if err != nil {
		return nil, errors.New("failed to marshal order response: " + err.Error())
	}

	cacheKey := getCachKey(strconv.Itoa(int(orderID)))
	err = s.redisClient.Set(ctx, cacheKey, orderResponseJSON, _cacheTime)
	if err != nil {
		return nil, err
	}

	go util.SendPostRequest(_paymentCreateURL, orderResponseJSON)

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
		go sendFailedOrder(_paymentUpdateURL, orderUpdateInfo)
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
