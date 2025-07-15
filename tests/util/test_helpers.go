package util

import (
	"encoding/json"
	"fmt"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

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
