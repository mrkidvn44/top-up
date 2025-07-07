package schema

import (
	"encoding/json"
	"top-up-api/internal/model"
)

type OrderConfirmRequest struct {
	OrderID       uint                        `json:"order_id"`
	UserID        uint                        `json:"user_id"`
	SkuID         uint                        `json:"sku_id"`
	TotalPrice    int                         `json:"total_price"`
	Status        model.PurchaseHistoryStatus `json:"status" validate:"purchasehistorystatus"`
	PhoneNumber   string                      `json:"phone_number"`
	CashBackValue int                         `json:"cash_back_value"`
}

type OrderRequest struct {
	UserID      uint   `json:"user_id"`
	SkuID       uint   `json:"sku_id"`
	PhoneNumber string `json:"phone_number"`
}

type OrderResponse struct {
	OrderID       uint                        `json:"order_id"`
	UserID        uint                        `json:"user_id"`
	Sku           SkuResponse                 `json:"sku"`
	TotalPrice    int                         `json:"total_price"`
	Status        model.PurchaseHistoryStatus `json:"status"`
	PhoneNumber   string                      `json:"phone_number"`
	CashBackValue int                         `json:"cash_back_value"`
}

type OrderProviderRequest struct {
	OrderID     uint   `json:"order_id"`
	PhoneNumber string `json:"phone_number"`
	TotalPrice  int    `json:"total_price"`
	Price       int    `json:"price"`
	CallBackUrl string `json:"callback_url"`
}

type OrderUpdateRequest struct {
	OrderID     uint                        `json:"order_id"`
	Status      model.PurchaseHistoryStatus `json:"status"`
	PhoneNumber string                      `json:"phone_number"`
}

func (o *OrderResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OrderResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

func (o *OrderResponse) CompareWithOrderConfirmRequest(orderConfirmRequest OrderConfirmRequest) bool {
	return o.OrderID == orderConfirmRequest.OrderID &&
		o.UserID == orderConfirmRequest.UserID &&
		o.TotalPrice == orderConfirmRequest.TotalPrice &&
		o.PhoneNumber == orderConfirmRequest.PhoneNumber &&
		o.CashBackValue == orderConfirmRequest.CashBackValue
}
