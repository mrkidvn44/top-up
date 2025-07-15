package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	grpcClient "top-up-api/internal/grpc/client"
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	"top-up-api/internal/service"
	mockGrpc "top-up-api/tests/mock"
	mockRepo "top-up-api/tests/repository/mock"
	"top-up-api/tests/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	orderReqPercentage = schema.OrderRequest{
		UserID:      1,
		SkuID:       1,
		PhoneNumber: "081234567890",
	}
	orderReqFixed = schema.OrderRequest{
		UserID:      2,
		SkuID:       2,
		PhoneNumber: "082345678901",
	}
	orderReqZeroWeight = schema.OrderRequest{
		UserID:      3,
		SkuID:       3,
		PhoneNumber: "083456789012",
	}
	orderReqNotFound = schema.OrderRequest{
		UserID:      1,
		SkuID:       999,
		PhoneNumber: "081234567890",
	}
	orderReqDBError = schema.OrderRequest{
		UserID:      1,
		SkuID:       1,
		PhoneNumber: "081234567890",
	}
	orderReqTimeout = schema.OrderRequest{
		UserID:      1,
		SkuID:       1,
		PhoneNumber: "081234567890",
	}
	orderReqRedisSetError = schema.OrderRequest{
		UserID:      1,
		SkuID:       1,
		PhoneNumber: "081234567890",
	}
	orderReqRedisTimeout = schema.OrderRequest{
		UserID:      4,
		SkuID:       4,
		PhoneNumber: "084567890123",
	}
	orderReqLarge = schema.OrderRequest{
		UserID:      999999,
		SkuID:       100,
		PhoneNumber: "+6285567890123",
	}
	orderReqEmptyPhone = schema.OrderRequest{
		UserID:      5,
		SkuID:       5,
		PhoneNumber: "",
	}

	confirmReqVTLConfirmStatus = schema.OrderConfirmRequest{
		OrderID:       1001,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqMBFFailedStatus = schema.OrderConfirmRequest{
		OrderID:       1002,
		UserID:        2,
		SkuID:         2,
		TotalPrice:    20000,
		Status:        model.PurchaseHistoryStatusFailed,
		PhoneNumber:   "082345678901",
		CashBackValue: 1000,
	}
	confirmReqVTLFailedLock = schema.OrderConfirmRequest{
		OrderID:       1003,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLNotFound = schema.OrderConfirmRequest{
		OrderID:       1004,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLCorrupt = schema.OrderConfirmRequest{
		OrderID:       1005,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLUserMismatchMain = schema.OrderConfirmRequest{
		OrderID:       1006,
		UserID:        2,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLPriceMismatchMain = schema.OrderConfirmRequest{
		OrderID:       1007,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    15000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLPendingMain = schema.OrderConfirmRequest{
		OrderID:       1008,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusPending,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLAlreadyConfirmedMain = schema.OrderConfirmRequest{
		OrderID:       1009,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}
	confirmReqVTLDBErrorMain = schema.OrderConfirmRequest{
		OrderID:       1010,
		UserID:        1,
		SkuID:         1,
		TotalPrice:    10000,
		Status:        model.PurchaseHistoryStatusConfirm,
		PhoneNumber:   "081234567890",
		CashBackValue: 500,
	}

	updateReqSuccess = schema.OrderUpdateRequest{
		OrderID:     1001,
		Status:      model.PurchaseHistoryStatusSuccess,
		PhoneNumber: "081234567890",
	}
	updateReqFailed = schema.OrderUpdateRequest{
		OrderID:     1001,
		Status:      model.PurchaseHistoryStatusFailed,
		PhoneNumber: "081234567890",
	}
)

type CreateOrderTestCase struct {
	Name          string
	OrderRequest  schema.OrderRequest
	SetupMocks    func(*mockRepo.SkuRepositoryMock, *mockGrpc.RedisMock, *mockRepo.ProviderRepositoryMock)
	ExpectedError string
	Assert        func(*testing.T, *schema.OrderResponse)
}

type ConfirmOrderTestCase struct {
	Name                string
	OrderConfirmRequest schema.OrderConfirmRequest
	SetupMocks          func(*mockRepo.SkuRepositoryMock, *mockRepo.PurchaseHistoryRepositoryMock, *mockGrpc.RedisMock, *mockRepo.ProviderRepositoryMock)
	ExpectedError       string
	GRPCSetup           *util.GRPCClientSetup // Optional gRPC setup configuration
}

type UpdateOrderStatusTestCase struct {
	Name               string
	OrderUpdateRequest schema.OrderUpdateRequest
	SetupMocks         func(*mockRepo.SkuRepositoryMock, *mockRepo.PurchaseHistoryRepositoryMock, *mockGrpc.RedisMock, *mockRepo.ProviderRepositoryMock)
	ExpectedError      string
}

func runTableDrivenTests[T any](t *testing.T, cases []T, run func(*testing.T, T)) {
	for _, tc := range cases {
		t.Run(getTestCaseName(tc), func(t *testing.T) {
			run(t, tc)
		})
	}
}

// Helper to extract the Name field from test case structs
func getTestCaseName(tc any) string {
	switch v := tc.(type) {
	case CreateOrderTestCase:
		return v.Name
	case ConfirmOrderTestCase:
		return v.Name
	case UpdateOrderStatusTestCase:
		return v.Name
	default:
		return "unnamed"
	}
}

func TestOrderService_CreateOrder(t *testing.T) {
	cases := []CreateOrderTestCase{
		{
			Name:         "successful order creation with percentage cashback",
			OrderRequest: orderReqPercentage,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(1, "VTL", 10000, model.CashBackTypePercentage, 5, "Viettel")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 100, []model.Supplier{
						util.CreateMockSupplier("VTL", "Viettel"),
					}),
				}
				util.SetupBasicMocks(skuRepo, redis, providerRepo, mockSku, providers)
			},
			ExpectedError: "",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, orderReqPercentage.UserID, result.UserID)
				assert.Equal(t, orderReqPercentage.SkuID, result.Sku.ID)
				assert.Equal(t, orderReqPercentage.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, 10000, result.TotalPrice)
				assert.Equal(t, model.PurchaseHistoryStatusPending, result.Status)
				assert.Greater(t, result.OrderID, uint(0))
				assert.GreaterOrEqual(t, result.RandomProviderWeight, 0)
				assert.LessOrEqual(t, result.RandomProviderWeight, 100)
				assert.Equal(t, 500, result.CashBackValue) // 5% of 10000 = 500
				assert.Equal(t, "VTL", result.Sku.SupplierInfo.Code)
				assert.Equal(t, "Viettel", result.Sku.SupplierInfo.Name)
			},
		},
		{
			Name:         "successful order creation with fixed cashback",
			OrderRequest: orderReqFixed,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(2, "MBF", 20000, model.CashBackTypeFixed, 1000, "MBF")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 50, []model.Supplier{
						util.CreateMockSupplier("MBF", "Mobifone"),
					}),
					util.CreateMockProvider(2, "PROVIDER2", "grpc://provider2.com", "grpc", 30, []model.Supplier{
						util.CreateMockSupplier("MBF", "Mobifone"),
					}),
				}
				util.SetupBasicMocks(skuRepo, redis, providerRepo, mockSku, providers)
			},
			ExpectedError: "",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, orderReqFixed.UserID, result.UserID)
				assert.Equal(t, orderReqFixed.SkuID, result.Sku.ID)
				assert.Equal(t, orderReqFixed.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, 20000, result.TotalPrice)
				assert.Equal(t, model.PurchaseHistoryStatusPending, result.Status)
				assert.Greater(t, result.OrderID, uint(0))
				assert.GreaterOrEqual(t, result.RandomProviderWeight, 0)
				assert.LessOrEqual(t, result.RandomProviderWeight, 80) // Total weight 50+30=80
				assert.Equal(t, 1000, result.CashBackValue)            // Fixed cashback
				assert.Equal(t, "MBF", result.Sku.SupplierInfo.Code)
			},
		},
		{
			Name:         "successful order creation with zero weight provider",
			OrderRequest: orderReqZeroWeight,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(3, "VTL", 50000, model.CashBackTypePercentage, 0, "Viettel")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 0, []model.Supplier{
						util.CreateMockSupplier("VTL", "Viettel"),
					}),
				}
				util.SetupBasicMocks(skuRepo, redis, providerRepo, mockSku, providers)
			},
			ExpectedError: "",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, orderReqZeroWeight.UserID, result.UserID)
				assert.Equal(t, orderReqZeroWeight.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, 50000, result.TotalPrice)
				assert.Equal(t, model.PurchaseHistoryStatusPending, result.Status)
				assert.Equal(t, 0, result.CashBackValue)        // No cashback
				assert.Equal(t, 0, result.RandomProviderWeight) // Zero total weight
			},
		},
		{
			Name:         "sku not found - gorm.ErrRecordNotFound",
			OrderRequest: orderReqNotFound,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				util.SetupErrorMocks(skuRepo, providerRepo, 999, gorm.ErrRecordNotFound)
			},
			ExpectedError: "sku not found",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.Nil(t, result)
			},
		},
		{
			Name:         "database connection error from sku repository",
			OrderRequest: orderReqDBError,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				util.SetupErrorMocks(skuRepo, providerRepo, 1, errors.New("database connection failed"))
			},
			ExpectedError: "database connection failed",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.Nil(t, result)
			},
		},
		{
			Name:         "timeout error from sku repository",
			OrderRequest: orderReqTimeout,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				util.SetupErrorMocks(skuRepo, providerRepo, 1, errors.New("context deadline exceeded"))
			},
			ExpectedError: "context deadline exceeded",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.Nil(t, result)
			},
		},
		{
			Name:         "redis cache set error",
			OrderRequest: orderReqRedisSetError,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(1, "VTL", 10000, model.CashBackTypePercentage, 1, "Viettel")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 100, []model.Supplier{
						util.CreateMockSupplier("VTL", "Viettel"),
					}),
				}
				util.SetupCacheErrorMocks(skuRepo, redis, providerRepo, mockSku, providers, errors.New("redis connection failed"))
			},
			ExpectedError: "redis connection failed",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.Nil(t, result)
			},
		},
		{
			Name:         "redis timeout error",
			OrderRequest: orderReqRedisTimeout,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(4, "XL", 15000, model.CashBackTypeFixed, 500, "XL Axiata")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 75, []model.Supplier{
						util.CreateMockSupplier("XL", "XL Axiata"),
					}),
				}
				util.SetupCacheErrorMocks(skuRepo, redis, providerRepo, mockSku, providers, errors.New("redis timeout"))
			},
			ExpectedError: "redis timeout",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.Nil(t, result)
			},
		},
		{
			Name:         "successful order with large phone number and high price",
			OrderRequest: orderReqLarge,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(100, "THREE", 500000, model.CashBackTypePercentage, 10, "3 (Tri)")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 1000, []model.Supplier{
						util.CreateMockSupplier("THREE", "3 (Tri)"),
					}),
				}
				util.SetupBasicMocks(skuRepo, redis, providerRepo, mockSku, providers)
			},
			ExpectedError: "",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, orderReqLarge.UserID, result.UserID)
				assert.Equal(t, orderReqLarge.SkuID, result.Sku.ID)
				assert.Equal(t, orderReqLarge.PhoneNumber, result.PhoneNumber)
				assert.Equal(t, 500000, result.TotalPrice)
				assert.Equal(t, 50000, result.CashBackValue) // 10% of 500000
				assert.GreaterOrEqual(t, result.RandomProviderWeight, 0)
				assert.LessOrEqual(t, result.RandomProviderWeight, 1000)
			},
		},
		{
			Name:         "empty phone number - edge case",
			OrderRequest: orderReqEmptyPhone,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				mockSku := util.CreateMockSku(5, "SMARTFREN", 25000, model.CashBackTypeFixed, 2500, "Smartfren")
				providers := []model.Provider{
					util.CreateMockProvider(1, "PROVIDER1", "http://provider1.com", "http", 200, []model.Supplier{
						util.CreateMockSupplier("SMARTFREN", "Smartfren"),
					}),
				}
				util.SetupBasicMocks(skuRepo, redis, providerRepo, mockSku, providers)
			},
			ExpectedError: "",
			Assert: func(t *testing.T, result *schema.OrderResponse) {
				assert.NotNil(t, result)
				assert.Equal(t, orderReqEmptyPhone.UserID, result.UserID)
				assert.Equal(t, "", result.PhoneNumber) // Empty phone number should be preserved
				assert.Equal(t, 25000, result.TotalPrice)
				assert.Equal(t, 2500, result.CashBackValue)
			},
		},
	}
	runTableDrivenTests(t, cases, func(t *testing.T, tc CreateOrderTestCase) {
		skuRepo := new(mockRepo.SkuRepositoryMock)
		purchaseRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
		redis := new(mockGrpc.RedisMock)
		grpcClients := &grpcClient.GRPCServiceClient{
			ProviderGRPCClients: make(map[string]grpcClient.ProviderGRPCClient),
		}
		providerRepo := new(mockRepo.ProviderRepositoryMock)

		tc.SetupMocks(skuRepo, redis, providerRepo)

		orderService := service.NewOrderService(skuRepo, purchaseRepo, redis, *grpcClients, providerRepo)
		result, err := orderService.CreateOrder(context.Background(), tc.OrderRequest)

		if tc.ExpectedError != "" {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.ExpectedError)
		} else {
			assert.NoError(t, err)
		}

		tc.Assert(t, result)

		skuRepo.AssertExpectations(t)
		redis.AssertExpectations(t)
		providerRepo.AssertExpectations(t)
	})
}

func TestOrderService_ConfirmOrder(t *testing.T) {
	tests := []ConfirmOrderTestCase{
		{
			Name:                "successful order confirmation with status confirm - HTTP provider",
			OrderConfirmRequest: confirmReqVTLConfirmStatus,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.SingleProvider("VTL", "Viettel")
				util.SetupBasicConfirmOrderTest(redis, providerRepo, purchaseRepo, "1001", confirmReqVTLConfirmStatus, "VTL", "Viettel", model.CashBackTypePercentage, 5, providers)
			},
			ExpectedError: "",
		},
		{
			Name:                "successful order confirmation with status confirm - gRPC provider",
			OrderConfirmRequest: confirmReqVTLConfirmStatus,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.SingleGrpcProvider("GRPC_PROVIDER1", "VTL", "Viettel")
				util.SetupBasicConfirmOrderTest(redis, providerRepo, purchaseRepo, "1001", confirmReqVTLConfirmStatus, "VTL", "Viettel", model.CashBackTypePercentage, 5, providers)
			},
			GRPCSetup: &util.GRPCClientSetup{
				ProviderCode: "GRPC_PROVIDER1",
				ShouldError:  false,
			},
			ExpectedError: "",
		},
		{
			Name:                "successful order confirmation with status failed - gRPC provider",
			OrderConfirmRequest: confirmReqMBFFailedStatus,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.SingleGrpcProvider("GRPC_PROVIDER2", "MBF", "Mobifone")
				util.SetupBasicConfirmOrderTest(redis, providerRepo, purchaseRepo, "1002", confirmReqMBFFailedStatus, "MBF", "Mobifone", model.CashBackTypeFixed, 1000, providers)
			},
			GRPCSetup: &util.GRPCClientSetup{
				ProviderCode: "GRPC_PROVIDER2",
				ShouldError:  false,
			},
			ExpectedError: "",
		},
		{
			Name:                "successful order confirmation with mixed providers - should use gRPC",
			OrderConfirmRequest: confirmReqVTLConfirmStatus,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.MixedProviders("VTL", "Viettel")
				util.SetupBasicConfirmOrderTest(redis, providerRepo, purchaseRepo, "1001", confirmReqVTLConfirmStatus, "VTL", "Viettel", model.CashBackTypePercentage, 5, providers)
			},
			GRPCSetup: &util.GRPCClientSetup{
				ProviderCode: "GRPC_PROVIDER1",
				ShouldError:  false,
			},
			ExpectedError: "",
		},
		{
			Name:                "successful order confirmation with status failed",
			OrderConfirmRequest: confirmReqMBFFailedStatus,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.SingleProvider("MBF", "Mobifone")
				util.SetupBasicConfirmOrderTest(redis, providerRepo, purchaseRepo, "1002", confirmReqMBFFailedStatus, "MBF", "Mobifone", model.CashBackTypeFixed, 1000, providers)
			},
			ExpectedError: "",
		},
		{
			Name:                "gRPC provider error during order processing",
			OrderConfirmRequest: confirmReqVTLConfirmStatus,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.SingleGrpcProvider("GRPC_PROVIDER_ERROR", "VTL", "Viettel")
				util.SetupBasicConfirmOrderTest(redis, providerRepo, purchaseRepo, "1001", confirmReqVTLConfirmStatus, "VTL", "Viettel", model.CashBackTypePercentage, 5, providers)
			},
			GRPCSetup: &util.GRPCClientSetup{
				ProviderCode: "GRPC_PROVIDER_ERROR",
				ShouldError:  true,
				ErrorMessage: "gRPC connection failed",
			},
			ExpectedError: "", // gRPC errors don't fail the confirmation, they're handled asynchronously
		},
		{
			Name:                "lock acquisition failed",
			OrderConfirmRequest: confirmReqVTLFailedLock,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				util.SetupErrorConfirmOrderTest(redis, providerRepo, "1003", errors.New("lock acquisition failed"), nil, "")
			},
			ExpectedError: "lock acquisition failed",
		},
		{
			Name:                "order not found in cache",
			OrderConfirmRequest: confirmReqVTLNotFound,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				util.SetupErrorConfirmOrderTest(redis, providerRepo, "1004", nil, errors.New("key not found"), "")
			},
			ExpectedError: "order not found or expired",
		},
		{
			Name:                "corrupted cached order data",
			OrderConfirmRequest: confirmReqVTLCorrupt,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				util.SetupErrorConfirmOrderTest(redis, providerRepo, "1005", nil, nil, "invalid json data")
			},
			ExpectedError: "failed to unmarshal order",
		},
		{
			Name:                "order mismatch - different user ID",
			OrderConfirmRequest: confirmReqVTLUserMismatchMain,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				cachedOrder := util.CreateCachedOrderResponse(confirmReqVTLUserMismatchMain.OrderID, 1, 10000, confirmReqVTLUserMismatchMain.PhoneNumber, 500, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				util.SetupOrderMismatchTest(redis, providerRepo, confirmReqVTLUserMismatchMain, "1006", cachedOrder)
			},
			ExpectedError: "order mismatch",
		},
		{
			Name:                "order mismatch - different total price",
			OrderConfirmRequest: confirmReqVTLPriceMismatchMain,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				cachedOrder := util.CreateCachedOrderResponse(confirmReqVTLPriceMismatchMain.OrderID, 1, 10000, confirmReqVTLPriceMismatchMain.PhoneNumber, 500, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				util.SetupOrderMismatchTest(redis, providerRepo, confirmReqVTLPriceMismatchMain, "1007", cachedOrder)
			},
			ExpectedError: "order mismatch",
		},
		{
			Name:                "order status pending - invalid status transition",
			OrderConfirmRequest: confirmReqVTLPendingMain,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				cachedOrder := util.CreateCachedOrderResponse(confirmReqVTLPendingMain.OrderID, 1, 10000, confirmReqVTLPendingMain.PhoneNumber, 500, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				util.SetupOrderMismatchTest(redis, providerRepo, confirmReqVTLPendingMain, "1008", cachedOrder)
			},
			ExpectedError: "order is pending",
		},
		{
			Name:                "order already confirmed",
			OrderConfirmRequest: confirmReqVTLAlreadyConfirmedMain,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				cachedOrder := util.CreateCachedOrderResponse(confirmReqVTLAlreadyConfirmedMain.OrderID, 1, 10000, confirmReqVTLAlreadyConfirmedMain.PhoneNumber, 500, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				cachedOrder.Status = model.PurchaseHistoryStatusConfirm // Already confirmed
				util.SetupOrderMismatchTest(redis, providerRepo, confirmReqVTLAlreadyConfirmedMain, "1009", cachedOrder)
			},
			ExpectedError: "order already confirmed or failed",
		},
		{
			Name:                "database error during purchase history creation",
			OrderConfirmRequest: confirmReqVTLDBErrorMain,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providers := util.SingleProvider("VTL", "Viettel")
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
				redis.On("TryAcquireLock", mock.Anything, "1010", mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("ReleaseLock", mock.Anything, "1010").Return(nil)

				cachedOrder := util.CreateCachedOrderResponse(confirmReqVTLDBErrorMain.OrderID, 1, 10000, confirmReqVTLDBErrorMain.PhoneNumber, 500, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				cachedOrderJSON, _ := json.Marshal(cachedOrder)
				redis.On("Get", mock.Anything, "order_id1010").Return(string(cachedOrderJSON), nil)

				purchaseRepo.On("CreatePurchaseHistory", mock.Anything, mock.MatchedBy(func(ph *model.PurchaseHistory) bool {
					return ph.OrderID == confirmReqVTLDBErrorMain.OrderID && ph.UserID == confirmReqVTLDBErrorMain.UserID && ph.Status == model.PurchaseHistoryStatusConfirm
				})).Return(errors.New("database connection failed"))
			},
			ExpectedError: "database connection failed",
		},
	}
	runTableDrivenTests(t, tests, func(t *testing.T, tc ConfirmOrderTestCase) {
		skuRepo := new(mockRepo.SkuRepositoryMock)
		purchaseRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
		redis := new(mockGrpc.RedisMock)

		grpcClients := &grpcClient.GRPCServiceClient{
			ProviderGRPCClients: make(map[string]grpcClient.ProviderGRPCClient),
		}
		providerRepo := new(mockRepo.ProviderRepositoryMock)

		tc.SetupMocks(skuRepo, purchaseRepo, redis, providerRepo)

		// Setup gRPC client mocks if needed
		if tc.GRPCSetup != nil {
			util.SetupGRPCMockClient(grpcClients, tc.GRPCSetup)
		}

		orderService := service.NewOrderService(skuRepo, purchaseRepo, redis, *grpcClients, providerRepo)
		err := orderService.ConfirmOrder(context.Background(), tc.OrderConfirmRequest)

		if tc.ExpectedError != "" {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.ExpectedError)
		} else {
			assert.NoError(t, err)
		}

		redis.AssertExpectations(t)
		if tc.ExpectedError == "" {
			purchaseRepo.AssertExpectations(t)
			providerRepo.AssertExpectations(t)
		}
	})
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	tests := []UpdateOrderStatusTestCase{
		{
			Name:               "successful order status update to success",
			OrderUpdateRequest: updateReqSuccess,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(util.SingleProvider("VTL", "Viettel"), nil)
				redis.On("Get", mock.Anything, "order_req_id1001").Return("", errors.New("not found"))
				redis.On("TryAcquireLock", mock.Anything, "1001", mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("ReleaseLock", mock.Anything, "1001").Return(nil)
				cachedOrder := util.CreateCachedOrderResponse(updateReqSuccess.OrderID, 1, 10000, updateReqSuccess.PhoneNumber, 0, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				cachedOrder.Status = model.PurchaseHistoryStatusConfirm
				cachedOrderJSON, _ := json.Marshal(cachedOrder)
				redis.On("Get", mock.Anything, "order_id1001").Return(string(cachedOrderJSON), nil)
				purchaseRepo.On("UpdatePurchaseHistoryStatusByOrderID", mock.Anything, uint(1001), model.PurchaseHistoryStatusSuccess).Return(nil)
				redis.On("Set", mock.Anything, "order_id1001", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("Set", mock.Anything, "order_req_id1001", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
			},
			ExpectedError: "",
		},
		{
			Name:               "successful order status update to failed",
			OrderUpdateRequest: updateReqFailed,
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(util.SingleProvider("VTL", "Viettel"), nil)
				redis.On("Get", mock.Anything, "order_req_id1001").Return("", errors.New("not found"))
				redis.On("TryAcquireLock", mock.Anything, "1001", mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("ReleaseLock", mock.Anything, "1001").Return(nil)
				cachedOrder := util.CreateCachedOrderResponse(updateReqFailed.OrderID, 1, 10000, updateReqFailed.PhoneNumber, 0, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				cachedOrder.Status = model.PurchaseHistoryStatusConfirm
				cachedOrderJSON, _ := json.Marshal(cachedOrder)
				redis.On("Get", mock.Anything, "order_id1001").Return(string(cachedOrderJSON), nil)
				purchaseRepo.On("UpdatePurchaseHistoryStatusByOrderID", mock.Anything, uint(1001), model.PurchaseHistoryStatusFailed).Return(nil)
				redis.On("Set", mock.Anything, "order_id1001", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("Set", mock.Anything, "order_req_id1001", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
			},
			ExpectedError: "",
		},
		{
			Name: "idempotency - request already processed successfully",
			OrderUpdateRequest: schema.OrderUpdateRequest{
				OrderID:     1001,
				Status:      model.PurchaseHistoryStatusSuccess,
				PhoneNumber: "081234567890",
			},
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(util.SingleProvider("VTL", "Viettel"), nil)
				idempotencyResponse := schema.OrderIdempotencyResponse{
					Success:      true,
					ErrorMessage: "",
					Timestamp:    time.Now().Unix(),
				}
				responseJSON, _ := json.Marshal(idempotencyResponse)
				redis.On("Get", mock.Anything, "order_req_id1001").Return(string(responseJSON), nil)
			},
			ExpectedError: "",
		},
		{
			Name: "idempotency - request already processed with error",
			OrderUpdateRequest: schema.OrderUpdateRequest{
				OrderID:     1001,
				Status:      model.PurchaseHistoryStatusSuccess,
				PhoneNumber: "081234567890",
			},
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(util.SingleProvider("VTL", "Viettel"), nil)
				idempotencyResponse := schema.OrderIdempotencyResponse{
					Success:      false,
					ErrorMessage: "previous error",
					Timestamp:    time.Now().Unix(),
				}
				responseJSON, _ := json.Marshal(idempotencyResponse)
				redis.On("Get", mock.Anything, "order_req_id1001").Return(string(responseJSON), nil)
			},
			ExpectedError: "previous error",
		},
		{
			Name: "order not confirmed",
			OrderUpdateRequest: schema.OrderUpdateRequest{
				OrderID:     1001,
				Status:      model.PurchaseHistoryStatusSuccess,
				PhoneNumber: "081234567890",
			},
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(util.SingleProvider("VTL", "Viettel"), nil)
				redis.On("Get", mock.Anything, "order_req_id1001").Return("", errors.New("not found"))
				redis.On("TryAcquireLock", mock.Anything, "1001", mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("ReleaseLock", mock.Anything, "1001").Return(nil)
				cachedOrder := util.CreateCachedOrderResponse(1001, 1, 10000, "081234567890", 0, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				cachedOrder.Status = model.PurchaseHistoryStatusPending // Not confirmed
				cachedOrderJSON, _ := json.Marshal(cachedOrder)
				redis.On("Get", mock.Anything, "order_id1001").Return(string(cachedOrderJSON), nil)
				redis.On("Set", mock.Anything, "order_req_id1001", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
			},
			ExpectedError: "order is not confirmed or failed",
		},
		{
			Name: "database error during status update",
			OrderUpdateRequest: schema.OrderUpdateRequest{
				OrderID:     1001,
				Status:      model.PurchaseHistoryStatusSuccess,
				PhoneNumber: "081234567890",
			},
			SetupMocks: func(skuRepo *mockRepo.SkuRepositoryMock, purchaseRepo *mockRepo.PurchaseHistoryRepositoryMock, redis *mockGrpc.RedisMock, providerRepo *mockRepo.ProviderRepositoryMock) {
				providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(util.SingleProvider("VTL", "Viettel"), nil)
				redis.On("Get", mock.Anything, "order_req_id1001").Return("", errors.New("not found"))
				redis.On("TryAcquireLock", mock.Anything, "1001", mock.AnythingOfType("time.Duration")).Return(nil)
				redis.On("ReleaseLock", mock.Anything, "1001").Return(nil)
				cachedOrder := util.CreateCachedOrderResponse(1001, 1, 10000, "081234567890", 0, 1, "VTL", "Viettel", model.CashBackTypePercentage, 5)
				cachedOrder.Status = model.PurchaseHistoryStatusConfirm
				cachedOrderJSON, _ := json.Marshal(cachedOrder)
				redis.On("Get", mock.Anything, "order_id1001").Return(string(cachedOrderJSON), nil)
				purchaseRepo.On("UpdatePurchaseHistoryStatusByOrderID", mock.Anything, uint(1001), model.PurchaseHistoryStatusSuccess).Return(errors.New("database error"))
				redis.On("Set", mock.Anything, "order_req_id1001", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)
			},
			ExpectedError: "database error",
		},
	}
	runTableDrivenTests(t, tests, func(t *testing.T, tc UpdateOrderStatusTestCase) {
		skuRepo := new(mockRepo.SkuRepositoryMock)
		purchaseRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
		redis := new(mockGrpc.RedisMock)
		grpcClients := &grpcClient.GRPCServiceClient{
			ProviderGRPCClients: make(map[string]grpcClient.ProviderGRPCClient),
		}
		providerRepo := new(mockRepo.ProviderRepositoryMock)

		tc.SetupMocks(skuRepo, purchaseRepo, redis, providerRepo)

		orderService := service.NewOrderService(skuRepo, purchaseRepo, redis, *grpcClients, providerRepo)
		err := orderService.UpdateOrderStatus(context.Background(), tc.OrderUpdateRequest)

		if tc.ExpectedError != "" {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.ExpectedError)
		} else {
			assert.NoError(t, err)
		}

		redis.AssertExpectations(t)
		providerRepo.AssertExpectations(t)
	})
}
func BenchmarkOrderService_CreateOrder(b *testing.B) {
	// Setup
	skuRepo := new(mockRepo.SkuRepositoryMock)
	purchaseRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
	redis := new(mockGrpc.RedisMock)
	grpcClients := &grpcClient.GRPCServiceClient{
		ProviderGRPCClients: make(map[string]grpcClient.ProviderGRPCClient),
	}
	providerRepo := new(mockRepo.ProviderRepositoryMock)

	mockSku := &model.Sku{
		Model:        gorm.Model{ID: 1},
		SupplierCode: "VTL",
		Price:        10000,
		CashBack: model.CashBack{
			Code:  "CB001",
			Type:  model.CashBackTypePercentage,
			Value: 5,
		},
		Supplier: model.Supplier{
			Code: "VTL",
			Name: "Viettel",
		},
	}

	providers := []model.Provider{
		{
			Model:  gorm.Model{ID: 1},
			Code:   "PROVIDER1",
			Source: "http://provider1.com",
			Type:   "http",
			Weight: 100,
			Suppliers: []model.Supplier{
				{Code: "VTL", Name: "Viettel"},
			},
		},
	}

	skuRepo.On("GetSkuByID", mock.Anything, uint(1)).Return(mockSku, nil)
	providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)
	redis.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).Return(nil)

	orderService := service.NewOrderService(skuRepo, purchaseRepo, redis, *grpcClients, providerRepo)

	orderRequest := schema.OrderRequest{
		UserID:      1,
		SkuID:       1,
		PhoneNumber: "081234567890",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := orderService.CreateOrder(context.Background(), orderRequest)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type InitializationTestCase struct {
	Name           string
	SetupProviders func() []model.Provider
	ExpectPanic    bool
	PanicMessage   string
}

func runInitializationTestCases(t *testing.T, cases []InitializationTestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			skuRepo := new(mockRepo.SkuRepositoryMock)
			purchaseRepo := new(mockRepo.PurchaseHistoryRepositoryMock)
			redis := new(mockGrpc.RedisMock)
			grpcClients := &grpcClient.GRPCServiceClient{
				ProviderGRPCClients: make(map[string]grpcClient.ProviderGRPCClient),
			}
			providerRepo := new(mockRepo.ProviderRepositoryMock)

			providers := tc.SetupProviders()
			providerRepo.On("GetProvidersWithSuppliers", mock.Anything).Return(providers, nil)

			if tc.ExpectPanic {
				assert.PanicsWithValue(t, tc.PanicMessage, func() {
					service.NewOrderService(skuRepo, purchaseRepo, redis, *grpcClients, providerRepo)
				})
			} else {
				assert.NotPanics(t, func() {
					service.NewOrderService(skuRepo, purchaseRepo, redis, *grpcClients, providerRepo)
				})
			}

			providerRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_Initialization(t *testing.T) {
	httpProvider := func(id uint, code, url string, weight int, supplierCode, supplierName string) model.Provider {
		return util.CreateMockProvider(id, code, url, "http", weight, []model.Supplier{
			util.CreateMockSupplier(supplierCode, supplierName),
		})
	}
	grpcProvider := func(id uint, code, url string, weight int, supplierCode, supplierName string) model.Provider {
		return util.CreateMockProvider(id, code, url, "grpc", weight, []model.Supplier{
			util.CreateMockSupplier(supplierCode, supplierName),
		})
	}
	unknownProvider := func(id uint, code, url, typ string, weight int, supplierCode, supplierName string) model.Provider {
		return util.CreateMockProvider(id, code, url, typ, weight, []model.Supplier{
			util.CreateMockSupplier(supplierCode, supplierName),
		})
	}

	testCases := []InitializationTestCase{
		{
			Name: "successful initialization with http and grpc providers",
			SetupProviders: func() []model.Provider {
				return []model.Provider{
					httpProvider(1, "HTTP_PROVIDER", "http://test.com", 100, "VTL", "Viettel"),
					grpcProvider(2, "GRPC_PROVIDER", "grpc://test.com:9090", 150, "MBF", "Mobifone"),
				}
			},
			ExpectPanic: false,
		},
		{
			Name: "unsupported provider type",
			SetupProviders: func() []model.Provider {
				return []model.Provider{
					unknownProvider(1, "UNKNOWN_PROVIDER", "unknown://test.com", "unknown", 100, "VTL", "Viettel"),
				}
			},
			ExpectPanic:  true,
			PanicMessage: "unsupported provider type: unknown",
		},
	}

	runInitializationTestCases(t, testCases)
}
