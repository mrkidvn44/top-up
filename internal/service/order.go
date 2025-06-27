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
	"top-up-api/pkg/redis"
	"top-up-api/pkg/util"
)

const (
	CacheTime        = 30 * time.Minute
	CacheKey         = "order_id"
	CashierCreateURL = "http://localhost:8081/v1/api/order/create"
	CashierUpdateURL = "http://localhost:8081/v1/api/order/update"
	ProviderURL      = "http://localhost:8082/v1/api/order/"
	CallbackURL      = "http://localhost:8080/v1/api/order/update-status"
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

	go util.SendPostRequest(CashierCreateURL, orderResponseJSON)

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
			if orderConfirmRequest.Status == model.PurchaseHistoryStatusConfirm {
				go SendProviderRequest(&orderResponse)
			}
			return nil
		}
		return errors.New("order already confirmed or failed")
	}

	purchaseHistory := mapper.PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest)
	err = s.purchaseHistoryRepo.CreatePurchaseHistory(ctx, purchaseHistory)
	if err != nil {
		return err
	}

	if orderConfirmRequest.Status == model.PurchaseHistoryStatusConfirm {
		go SendProviderRequest(&orderResponse)
	}

	return nil
}

func SendProviderRequest(orderResponse *schema.OrderResponse) error {
	orderProviderRequest := mapper.OrderProviderRequestFromOrderResponse(orderResponse, CallbackURL)
	orderProviderRequestJSON, err := json.Marshal(orderProviderRequest)
	fmt.Print("Order Provider Request: ", string(orderProviderRequestJSON), "\n")
	if err != nil {
		return errors.New("failed to marshal order response: " + err.Error())
	}
	util.SendPostRequest(ProviderURL, orderProviderRequestJSON)
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
			req, err := http.NewRequest(http.MethodPatch, CashierUpdateURL, bytes.NewBuffer(payload))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				return
			}
		}()
	}

	return nil
}
