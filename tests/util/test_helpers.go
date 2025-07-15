package util

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	grpcClient "top-up-api/internal/grpc/client"
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	mockGrpc "top-up-api/tests/mock"
	mockRepo "top-up-api/tests/repository/mock"
)

// Factory functions
func CreateMockSku(id uint, supplierCode string, price int, cashbackType model.CashBackType, cashbackValue int, supplierName string) *model.Sku {
	return &model.Sku{
		Model:        gorm.Model{ID: id},
		SupplierCode: supplierCode,
		Price:        price,
		CashBack: model.CashBack{
			Code:  "CB" + fmt.Sprintf("%03d", id),
			Type:  cashbackType,
			Value: cashbackValue,
		},
		Supplier: model.Supplier{
			Code: supplierCode,
			Name: supplierName,
		},
	}
}

func CreateMockProvider(id uint, code, source, providerType string, weight int, suppliers []model.Supplier) model.Provider {
	return model.Provider{
		Model:     gorm.Model{ID: id},
		Code:      code,
		Source:    source,
		Type:      providerType,
		Weight:    weight,
		Suppliers: suppliers,
	}
}

func CreateMockSupplier(code, name string) model.Supplier {
	return model.Supplier{
		Code: code,
		Name: name,
	}
}

// Common mock setup functions
func SetupBasicMocks(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, sku *model.Sku, providers []model.Provider) {
	skuRepo.On("GetSkuByID", mock.Anything, sku.ID).Return(sku, nil)
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
}

func SetupErrorMocks(skuRepo *mockRepo.SkuRepositoryMock, providerRepo *mockRepo.ProviderRepositoryMock, skuID uint, skuError error) {
	skuRepo.On("GetSkuByID", mock.Anything, skuID).Return((*model.Sku)(nil), skuError)
	providers := []model.Provider{
		CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 100, []model.Supplier{
			CreateMockSupplier("VTL", "Viettel"),
		}),
	}
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
}

func SetupCacheErrorMocks(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, sku *model.Sku, providers []model.Provider, cacheError error) {
	skuRepo.On("GetSkuByID", mock.Anything, sku.ID).Return(sku, nil)
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(cacheError)
}

// ConfirmOrder test helpers
func CreateCachedOrderResponse(orderID, userID uint, totalPrice int, phoneNumber string, cashbackValue int, skuID uint, supplierCode, supplierName string, cashbackType model.CashBackType, cashbackTypeValue int) *schema.OrderResponse {
	response := &schema.OrderResponse{
		OrderID:              orderID,
		UserID:               userID,
		TotalPrice:           totalPrice,
		PhoneNumber:          phoneNumber,
		CashBackValue:        cashbackValue,
		Status:               model.PurchaseHistoryStatusPending,
		RandomProviderWeight: 50,
		Sku: schema.SkuResponse{
			ID:    skuID,
			Price: totalPrice,
			SupplierInfo: schema.SupplierInfo{
				Code: supplierCode,
				Name: supplierName,
			},
		},
	}
	response.Sku.CashBackInterface = CreateCashBackInterface(cashbackType, skuID, cashbackTypeValue)
	return response
}

func CreateCashBackInterface(cashbackType model.CashBackType, skuID uint, cashbackTypeValue int) schema.CashBackInterface {
	code := "CB" + fmt.Sprintf("%03d", skuID)
	switch cashbackType {
	case model.CashBackTypePercentage:
		return &schema.CashBackPercentage{
			Type:  model.CashBackTypePercentage,
			Code:  code,
			Value: cashbackTypeValue,
		}
	case model.CashBackTypeFixed:
		return &schema.CashBackFixed{
			Type:  model.CashBackTypeFixed,
			Code:  code,
			Value: cashbackTypeValue,
		}
	default:
		return nil
	}
}

func SetupConfirmOrderMocks(redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, orderID string, cachedOrder *schema.OrderResponse, purchaseHistoryError error) {
	providers := []model.Provider{
		CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 100, []model.Supplier{
			CreateMockSupplier("VTL", "Viettel"),
		}),
	}
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("TryAcquireLock", mock.Anything, orderID, mock.AnythingOfType("time.Duration")).Return(nil)
	redis.On("ReleaseLock", mock.Anything, orderID).Return(nil)
	cachedOrderBytes, _ := json.Marshal(cachedOrder)
	redis.On("Get", mock.Anything, "order_id"+orderID).Return(string(cachedOrderBytes), nil)
	redis.On("Set", mock.Anything, "order_id"+orderID, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
	if purchaseHistoryError != nil {
		purchaseRepo.On("CreatePurchaseHistory", mock.Anything, mock.AnythingOfType("*model.PurchaseHistory")).Return(purchaseHistoryError)
	} else {
		purchaseRepo.On("CreatePurchaseHistory", mock.Anything, mock.AnythingOfType("*model.PurchaseHistory")).Return(nil)
	}
}

// Helper for single provider setup used in multiple tests
var SingleProvider = func(code, name string) []model.Provider {
	return []model.Provider{
		CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 100, []model.Supplier{
			CreateMockSupplier(code, name),
		}),
	}
}

// Helper for gRPC provider setup
var SingleGrpcProvider = func(providerCode, code, name string) []model.Provider {
	return []model.Provider{
		CreateMockProvider(1, providerCode, "grpc://provider1.com:9090", "grpc", 100, []model.Supplier{
			CreateMockSupplier(code, name),
		}),
	}
}

// Helper for mixed providers setup
var MixedProviders = func(code, name string) []model.Provider {
	return []model.Provider{
		CreateMockProvider(1, "HTTP_PROVIDER1", "http://provider1.com", "http", 30, []model.Supplier{
			CreateMockSupplier(code, name),
		}),
		CreateMockProvider(2, "GRPC_PROVIDER1", "grpc://provider2.com:9090", "grpc", 70, []model.Supplier{
			CreateMockSupplier(code, name),
		}),
	}
}

// Enhanced confirm order mock setup with provider type flexibility
func SetupConfirmOrderMocksWithProviders(redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, orderID string, cachedOrder *schema.OrderResponse, providers []model.Provider, purchaseHistoryError error) {
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("TryAcquireLock", mock.Anything, orderID, mock.AnythingOfType("time.Duration")).Return(nil)
	redis.On("ReleaseLock", mock.Anything, orderID).Return(nil)
	cachedOrderBytes, _ := json.Marshal(cachedOrder)
	redis.On("Get", mock.Anything, "order_id"+orderID).Return(string(cachedOrderBytes), nil)
	redis.On("Set", mock.Anything, "order_id"+orderID, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
	if purchaseHistoryError != nil {
		purchaseRepo.On("CreatePurchaseHistory", mock.Anything, mock.AnythingOfType("*model.PurchaseHistory")).Return(purchaseHistoryError)
	} else {
		purchaseRepo.On("CreatePurchaseHistory", mock.Anything, mock.AnythingOfType("*model.PurchaseHistory")).Return(nil)
	}
}

// Helper to setup mocks with specific purchase history matcher
func SetupConfirmOrderMocksWithMatcher(redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, orderID string, cachedOrder *schema.OrderResponse, providers []model.Provider, purchaseHistoryMatcher interface{}, purchaseHistoryError error) {
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("TryAcquireLock", mock.Anything, orderID, mock.AnythingOfType("time.Duration")).Return(nil)
	redis.On("ReleaseLock", mock.Anything, orderID).Return(nil)
	cachedOrderBytes, _ := json.Marshal(cachedOrder)
	redis.On("Get", mock.Anything, "order_id"+orderID).Return(string(cachedOrderBytes), nil)
	redis.On("Set", mock.Anything, "order_id"+orderID, mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
	if purchaseHistoryError != nil {
		purchaseRepo.On("CreatePurchaseHistory", mock.Anything, purchaseHistoryMatcher).Return(purchaseHistoryError)
	} else {
		purchaseRepo.On("CreatePurchaseHistory", mock.Anything, purchaseHistoryMatcher).Return(nil)
	}
}

// Simplified helper for basic confirm order test cases
func SetupBasicConfirmOrderTest(redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, orderID string, confirmReq schema.OrderConfirmRequest, supplierCode, supplierName string, cashbackType model.CashBackType, cashbackValue int, providers []model.Provider) {
	cachedOrder := CreateCachedOrderResponse(confirmReq.OrderID, confirmReq.UserID, confirmReq.TotalPrice, confirmReq.PhoneNumber, confirmReq.CashBackValue, confirmReq.SkuID, supplierCode, supplierName, cashbackType, cashbackValue)

	// Add random provider weight for mixed provider tests
	if len(providers) > 1 && confirmReq.OrderID == 1001 {
		cachedOrder.RandomProviderWeight = 50 // This should route to gRPC provider (weight range 30-100)
	}

	purchaseHistoryMatcher := mock.MatchedBy(func(ph *model.PurchaseHistory) bool {
		return ph.OrderID == confirmReq.OrderID && ph.UserID == confirmReq.UserID && ph.Status == confirmReq.Status
	})

	SetupConfirmOrderMocksWithMatcher(redis, providerRepo, purchaseRepo, orderID, cachedOrder, providers, purchaseHistoryMatcher, nil)
}

// Helper for error test cases that only need basic provider and lock setup
func SetupErrorConfirmOrderTest(redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, orderID string, lockError error, cacheError error, cachedData string) {
	providers := SingleProvider("VTL", "Viettel")
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)

	if lockError != nil {
		redis.On("TryAcquireLock", mock.Anything, orderID, mock.AnythingOfType("time.Duration")).Return(lockError)
		return
	}

	redis.On("TryAcquireLock", mock.Anything, orderID, mock.AnythingOfType("time.Duration")).Return(nil)
	redis.On("ReleaseLock", mock.Anything, orderID).Return(nil)

	if cacheError != nil {
		redis.On("Get", mock.Anything, "order_id"+orderID).Return("", cacheError)
	} else {
		redis.On("Get", mock.Anything, "order_id"+orderID).Return(cachedData, nil)
	}
}

// Helper for order mismatch test cases
func SetupOrderMismatchTest(redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock, confirmReq schema.OrderConfirmRequest, orderID string, cachedOrder *schema.OrderResponse) {
	providers := SingleProvider("VTL", "Viettel")
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("TryAcquireLock", mock.Anything, orderID, mock.AnythingOfType("time.Duration")).Return(nil)
	redis.On("ReleaseLock", mock.Anything, orderID).Return(nil)
	cachedOrderJSON, _ := json.Marshal(cachedOrder)
	redis.On("Get", mock.Anything, "order_id"+orderID).Return(string(cachedOrderJSON), nil)
}

type GRPCClientSetup struct {
	ProviderCode string
	ShouldError  bool
	ErrorMessage string
}

func SetupGRPCMockClient(grpcClients *grpcClient.GRPCServiceClient, setup *GRPCClientSetup) {
	mockGrpcClient := new(mockGrpc.ProviderGRPCClientMock)
	if setup.ShouldError {
		mockGrpcClient.On("ProcessOrder", mock.Anything, mock.AnythingOfType("*provider.OrderProcessRequest")).Return(errors.New(setup.ErrorMessage))
	} else {
		mockGrpcClient.On("ProcessOrder", mock.Anything, mock.AnythingOfType("*provider.OrderProcessRequest")).Return(nil)
	}
	grpcClients.ProviderGRPCClients[setup.ProviderCode] = mockGrpcClient
}
