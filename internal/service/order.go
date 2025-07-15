package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	pb "top-up-api/internal/grpc/client"
	"top-up-api/internal/mapper"
	"top-up-api/internal/model"
	"top-up-api/internal/repository"
	"top-up-api/internal/schema"
	"top-up-api/pkg/redis"
	"top-up-api/pkg/util"

	"gorm.io/gorm"
)

const (
	_lockTimeOut              = 5 * time.Minute
	_orderCacheTime           = 30 * time.Minute
	_idempotencyCacheTime     = 24 * time.Hour
	_orderRequestKeyPrefix    = "order_id"
	_providerRequestKeyPrefix = "order_req_id"
	_paymentCreateURL         = "http://localhost:8081/v1/api/order/create"
	_paymentUpdateURL         = "http://localhost:8081/v1/api/order/update"
	_callbackURL              = "http://localhost:8080/v1/api/order/update-status"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error)
	ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error
	UpdateOrderStatus(ctx context.Context, orderUpdateInfo schema.OrderUpdateRequest) error
}

type orderService struct {
	skuRepo             repository.SkuRepository
	purchaseHistoryRepo repository.PurchaseHistoryRepository
	redisClient         redis.Interface
	providerClients     map[string]providerServiceList
}

type providerServiceList struct {
	totalWeight     int
	providerClients []providerClient
}

type providerClient interface {
	getCumulativeWeight() int
	sendRequest(order *schema.OrderResponse) error
}

var _ OrderService = (*orderService)(nil)
var _ providerClient = (*httpProviderClient)(nil)
var _ providerClient = (*grpcProviderClient)(nil)

func NewOrderService(
	skuRepo repository.SkuRepository,
	purchaseHistoryRepo repository.PurchaseHistoryRepository,
	redisClient redis.Interface,
	grpcClients pb.GRPCServiceClient,
	providerRepo repository.ProviderRepository,
) *orderService {

	return &orderService{
		skuRepo:             skuRepo,
		purchaseHistoryRepo: purchaseHistoryRepo,
		redisClient:         redisClient,
		providerClients:     getProviderClientsListMapping(providerRepo, grpcClients),
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order schema.OrderRequest) (*schema.OrderResponse, error) {
	sku, err := s.skuRepo.GetSkuByID(ctx, order.SkuID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("sku not found")
		}
		return nil, err
	}

	orderID := util.GenerateOrderID()
	orderResponse := mapper.OrderResponseFromOrderRequest(order, sku, orderID)
	SupplierCode := orderResponse.Sku.SupplierInfo.Code
	orderResponse.RandomProviderWeight = getRandomWeight(s.providerClients[SupplierCode].totalWeight)

	orderResponseJSON, err := json.Marshal(orderResponse)
	if err != nil {
		return nil, errors.New("failed to marshal order response: " + err.Error())
	}

	cacheKey := getCachKey(_orderRequestKeyPrefix, strconv.Itoa(int(orderID)))
	err = s.redisClient.Set(ctx, cacheKey, orderResponseJSON, _orderCacheTime)
	if err != nil {
		return nil, err
	}

	go util.SendPostRequest(_paymentCreateURL, orderResponseJSON)

	return orderResponse, nil
}

func (s *orderService) ConfirmOrder(ctx context.Context, orderConfirmRequest schema.OrderConfirmRequest) error {
	orderID := strconv.Itoa(int(orderConfirmRequest.OrderID))
	err := s.redisClient.TryAcquireLock(ctx, orderID, _lockTimeOut)
	if err != nil {
		return err
	}
	defer s.redisClient.ReleaseLock(ctx, orderID)

	cacheKey := getCachKey(_orderRequestKeyPrefix, orderID)

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
		go s.sendRequestToProvider(orderResponse)
	}

	return nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderUpdateInfo schema.OrderUpdateRequest) error {
	orderID := strconv.Itoa(int(orderUpdateInfo.OrderID))

	idempotencyKey := getCachKey(_providerRequestKeyPrefix, orderID)
	cachedResponse, err := s.redisClient.Get(ctx, idempotencyKey)
	if err == nil && cachedResponse != "" {
		return getIdempotencyResponseValue(cachedResponse)
	}

	err = s.redisClient.TryAcquireLock(ctx, orderID, _lockTimeOut)
	if err != nil {
		return err
	}
	defer s.redisClient.ReleaseLock(ctx, orderID)

	orderCacheKey := getCachKey(_orderRequestKeyPrefix, orderID)
	orderResponse, err := s.getCachedOrder(ctx, orderCacheKey)
	if err != nil {
		s.cacheIdempotencyResponse(ctx, idempotencyKey, false, err.Error())
		return err
	}

	if orderResponse.Status != model.PurchaseHistoryStatusConfirm {
		err := errors.New("order is not confirmed or failed")
		s.cacheIdempotencyResponse(ctx, idempotencyKey, false, err.Error())
		return err
	}

	err = s.purchaseHistoryRepo.UpdatePurchaseHistoryStatusByOrderID(ctx, orderUpdateInfo.OrderID, orderUpdateInfo.Status)
	if err != nil {
		s.cacheIdempotencyResponse(ctx, idempotencyKey, false, err.Error())
		return err
	}

	orderResponse.Status = orderUpdateInfo.Status
	err = s.updateCacheOrderStaus(ctx, orderCacheKey, orderResponse)
	if err != nil {
		s.cacheIdempotencyResponse(ctx, idempotencyKey, false, err.Error())
		return err
	}

	if orderUpdateInfo.Status == model.PurchaseHistoryStatusFailed {
		go sendFailedOrder(_paymentUpdateURL, orderUpdateInfo)
	}

	s.cacheIdempotencyResponse(ctx, idempotencyKey, true, "")

	return nil
}

func (s *orderService) getCachedOrder(ctx context.Context, cacheKey string) (*schema.OrderResponse, error) {

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

func (s *orderService) updateCacheOrderStaus(ctx context.Context, cacheKey string, orderResponse *schema.OrderResponse) error {
	orderResponseJSON, err := json.Marshal(orderResponse)

	if err != nil {
		return errors.New("failed to marshal order response: " + err.Error())
	}

	err = s.redisClient.Set(ctx, cacheKey, orderResponseJSON, _orderCacheTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *orderService) cacheIdempotencyResponse(ctx context.Context, key string, success bool, errorMessage string) {
	response := schema.OrderIdempotencyResponse{
		Success:      success,
		ErrorMessage: errorMessage,
		Timestamp:    time.Now().Unix(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return
	}

	s.redisClient.Set(ctx, key, responseJSON, _idempotencyCacheTime)
}

func (s *orderService) sendRequestToProvider(orderResponse *schema.OrderResponse) error {
	supplierCode := orderResponse.Sku.SupplierInfo.Code

	for _, client := range s.providerClients[supplierCode].providerClients {
		if orderResponse.RandomProviderWeight <= client.getCumulativeWeight() {
			err := client.sendRequest(orderResponse)
			return err
		}

	}
	return errors.New("can't find suitable provider")
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

func getCachKey(prefix string, orderID string) string {
	return prefix + orderID
}

func getRandomWeight(totalWeight int) int {
	if totalWeight <= 0 {
		return 0
	}
	return rand.IntN(totalWeight)
}

func getIdempotencyResponseValue(cachedResponse string) error {
	var idempotencyResponse schema.OrderIdempotencyResponse
	if err := json.Unmarshal([]byte(cachedResponse), &idempotencyResponse); err != nil {
		return errors.New("failed to unmarshal idempotency response: " + err.Error())
	}
	if idempotencyResponse.Success {
		return nil
	}
	return errors.New(idempotencyResponse.ErrorMessage)
}

type httpProviderClient struct {
	url              string
	callbacks        string
	cumulativeWeight int
}

func (h *httpProviderClient) sendRequest(order *schema.OrderResponse) error {
	orderProviderRequest := mapper.OrderProviderRequestFromOrderResponse(order, h.callbacks)
	orderProviderRequestJSON, err := json.Marshal(orderProviderRequest)
	if err != nil {
		return err
	}
	util.SendPostRequest(h.url, orderProviderRequestJSON)
	return nil
}

func (h *httpProviderClient) getCumulativeWeight() int {
	return h.cumulativeWeight
}

type grpcProviderClient struct {
	client           pb.ProviderGRPCClient
	callbacks        string
	cumulativeWeight int
}

func (g *grpcProviderClient) sendRequest(order *schema.OrderResponse) error {
	ctx := context.Background()
	defer ctx.Done()
	req := mapper.OrderProcessRequestFromOrder(order, g.callbacks)
	return g.client.ProcessOrder(ctx, req)
}

func (g *grpcProviderClient) getCumulativeWeight() int {
	return g.cumulativeWeight
}

func getProviderClientsListMapping(providerRepo repository.ProviderRepository, grpcClients pb.GRPCServiceClient) map[string]providerServiceList {
	ctx := context.Background()
	providers, err := providerRepo.GetProvidersWithSuppliers(ctx)
	if err != nil {
		panic("failed to load providers: " + err.Error())
	}

	grpcClients.BuildProviderGRPCClients(providers)
	supplierClients := make(map[string]providerServiceList)

	for _, provider := range providers {
		for _, supplier := range provider.Suppliers {
			entry := supplierClients[supplier.Code]
			cumulativeWeight := entry.totalWeight + provider.Weight

			client := createProviderClient(provider, grpcClients, cumulativeWeight)
			entry.providerClients = append(entry.providerClients, client)
			entry.totalWeight = cumulativeWeight
			supplierClients[supplier.Code] = entry
		}
	}

	return supplierClients
}

func createProviderClient(provider model.Provider, grpcClients pb.GRPCServiceClient, cumulativeWeight int) providerClient {
	switch provider.Type {
	case "http":
		return &httpProviderClient{
			url: provider.Source, callbacks: _callbackURL, cumulativeWeight: cumulativeWeight,
		}
	case "grpc":
		return &grpcProviderClient{
			client: grpcClients.ProviderGRPCClients[provider.Code], callbacks: _callbackURL, cumulativeWeight: cumulativeWeight,
		}
	default:
		panic("unsupported provider type: " + provider.Type)
	}
}
