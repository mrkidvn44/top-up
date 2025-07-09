package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
	pb "top-up-api/proto/order"
	providerpb "top-up-api/proto/provider"
)

func OrderResponseFromOrderRequest(orderRequest schema.OrderRequest, sku *model.Sku, orderID uint) *schema.OrderResponse {
	skuResponse := SkuResponseFromModel(*sku)

	return &schema.OrderResponse{
		OrderID:       orderID,
		UserID:        orderRequest.UserID,
		Sku:           *skuResponse,
		TotalPrice:    skuResponse.Price,
		Status:        model.PurchaseHistoryStatusPending,
		PhoneNumber:   orderRequest.PhoneNumber,
		CashBackValue: skuResponse.CashBackInterface.CalculateCashBack(sku.Price),
	}

}

func OrderProviderRequestFromOrderResponse(orderResponse *schema.OrderResponse, callbackUrl string) *schema.OrderProviderRequest {
	return &schema.OrderProviderRequest{
		OrderID:     orderResponse.OrderID,
		PhoneNumber: orderResponse.PhoneNumber,
		TotalPrice:  orderResponse.TotalPrice,
		Price:       orderResponse.Sku.Price,
		CallBackUrl: callbackUrl,
	}
}

func OrderConfirmRequestFromProto(order *pb.OrderConfirmRequest) *schema.OrderConfirmRequest {
	return &schema.OrderConfirmRequest{
		OrderID:       uint(order.OrderId),
		UserID:        uint(order.UserId),
		SkuID:         uint(order.SkuId),
		TotalPrice:    int(order.TotalPrice),
		Status:        model.PurchaseHistoryStatus(order.Status),
		PhoneNumber:   order.PhoneNumber,
		CashBackValue: int(order.CashBackValue),
	}
}

func OrderUpdateRequestFromProto(order *pb.OrderUpdateRequest) *schema.OrderUpdateRequest {
	return &schema.OrderUpdateRequest{
		OrderID:     uint(order.OrderId),
		Status:      model.PurchaseHistoryStatus(order.Status),
		PhoneNumber: order.PhoneNumber,
	}
}

func OrderProcessRequestFromOrder(order *schema.OrderResponse, callbackUrl string) *providerpb.OrderProcessRequest {
	req := OrderProviderRequestFromOrderResponse(order, callbackUrl)
	return &providerpb.OrderProcessRequest{
		OrderId:     uint64(req.OrderID),
		PhoneNumber: req.PhoneNumber,
		TotalPrice:  int64(req.TotalPrice),
		Price:       int64(req.Price),
		CallBackUrl: req.CallBackUrl,
	}
}
