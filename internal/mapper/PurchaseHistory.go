package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func PurchaseHistoryResponseFromModel(purchaseHistory *model.PurchaseHistory) *schema.PurchaseHistoryResponse {
	return &schema.PurchaseHistoryResponse{
		OrderID:       purchaseHistory.OrderID,
		UserID:        purchaseHistory.UserID,
		SkuID:         purchaseHistory.SkuID,
		TotalPrice:    purchaseHistory.TotalPrice,
		PhoneNumber:   purchaseHistory.PhoneNumber,
		Status:        string(purchaseHistory.Status),
		CashBackValue: purchaseHistory.CashBackValue,
		Sku:           *SkuResponseFromModel(purchaseHistory.Sku),
	}
}

func PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest schema.OrderConfirmRequest) *model.PurchaseHistory {
	return &model.PurchaseHistory{
		UserID:        orderConfirmRequest.UserID,
		OrderID:       orderConfirmRequest.OrderID,
		SkuID:         orderConfirmRequest.SkuID,
		PhoneNumber:   orderConfirmRequest.PhoneNumber,
		TotalPrice:    orderConfirmRequest.TotalPrice,
		Status:        orderConfirmRequest.Status,
		CashBackValue: orderConfirmRequest.CashBackValue,
	}
}
