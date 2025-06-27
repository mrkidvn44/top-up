package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest schema.OrderConfirmRequest) *model.PurchaseHistory {
	return &model.PurchaseHistory{
		UserID:        orderConfirmRequest.UserID,
		OrderID:       orderConfirmRequest.OrderID,
		CardDetailID:  orderConfirmRequest.CardDetail.ID,
		PhoneNumber:   orderConfirmRequest.PhoneNumber,
		TotalPrice:    orderConfirmRequest.TotalPrice,
		Status:        model.PurchaseHistoryStatusPending,
		CashBackValue: orderConfirmRequest.CashBackValue,
	}
}

func OrderResponseFromOrderRequest(orderRequest schema.OrderRequest, cardDetail *model.CardDetail, orderID uint) *schema.OrderResponse {
	cardDetailResponse := CardDetailResponseFromModel(*cardDetail)

	return &schema.OrderResponse{
		OrderID:       orderID,
		UserID:        orderRequest.UserID,
		CardDetail:    *cardDetailResponse,
		TotalPrice:    cardDetailResponse.CardPriceResponse.Value,
		Status:        model.PurchaseHistoryStatusPending,
		PhoneNumber:   orderRequest.PhoneNumber,
		CashBackValue: cardDetailResponse.CashBackInterface.CalculateCashBack(cardDetail.CardPrice.Value),
	}

}
