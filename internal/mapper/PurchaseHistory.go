package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func PurchaseHistoryResponseFromModel(purchaseHistory *model.PurchaseHistory) *schema.PurchaseHistoryResponse {
	return &schema.PurchaseHistoryResponse{
		OrderID:       purchaseHistory.OrderID,
		UserID:        purchaseHistory.UserID,
		CardDetailID:  purchaseHistory.CardDetailID,
		TotalPrice:    purchaseHistory.TotalPrice,
		PhoneNumber:   purchaseHistory.PhoneNumber,
		Status:        string(purchaseHistory.Status),
		CashBackValue: purchaseHistory.CashBackValue,
		CardDetail:    *CardDetailResponseFromModel(purchaseHistory.CardDetail),
	}
}
