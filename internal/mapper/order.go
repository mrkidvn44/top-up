package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	pb "top-up-api/proto/order"
)

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

func OrderProviderRequestFromOrderResponse(orderResponse *schema.OrderResponse, callbackUrl string) *schema.OrderProviderRequest {
	return &schema.OrderProviderRequest{
		OrderID:     orderResponse.OrderID,
		PhoneNumber: orderResponse.PhoneNumber,
		TotalPrice:  orderResponse.TotalPrice,
		CardPrice:   orderResponse.CardDetail.CardPriceResponse.Value,
		CallBackUrl: callbackUrl,
	}
}

func OrderConfirmRequestFromProto(order *pb.OrderConfirmRequest) *schema.OrderConfirmRequest {
	return &schema.OrderConfirmRequest{
		OrderID:       uint(order.OrderId),
		UserID:        uint(order.UserId),
		CardDetailID:  uint(order.CardDetailId),
		TotalPrice:    int(order.TotalPrice),
		Status:        model.PurchaseHistoryStatus(order.Status),
		PhoneNumber:   order.PhoneNumber,
		CashBackValue: int(order.CashBackValue),
	}
}
