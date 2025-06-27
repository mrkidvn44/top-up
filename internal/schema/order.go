package schema

import (
	"encoding/json"
	"top-up-api/internal/model"
)

type OrderConfirmRequest struct {
	OrderID            uint                        `json:"order_id"`
	UserID             uint                        `json:"user_id"`
	CardDetailResponse CardDetailResponse          `json:"card_detail"`
	TotalPrice         int                         `json:"total_price"`
	Status             model.PurchaseHistoryStatus `json:"status"`
	PhoneNumber        string                      `json:"phone_number"`
	CashBackValue      int                         `json:"cash_back_value"`
}

type OrderRequest struct {
	UserID       uint   `json:"user_id"`
	CardDetailID uint   `json:"card_detail_id"`
	PhoneNumber  string `json:"phone_number"`
}

type OrderResponse struct {
	OrderID            uint                        `json:"order_id"`
	UserID             uint                        `json:"user_id"`
	CardDetailResponse CardDetailResponse          `json:"card_detail"`
	TotalPrice         int                         `json:"total_price"`
	Status             model.PurchaseHistoryStatus `json:"status"`
	PhoneNumber        string                      `json:"phone_number"`
	CashBackValue      int                         `json:"cash_back_value"`
}

func PurchaseHistoryFromOrderConfirmRequest(orderConfirmRequest OrderConfirmRequest) *model.PurchaseHistory {
	return &model.PurchaseHistory{
		UserID:       orderConfirmRequest.UserID,
		CardDetailID: orderConfirmRequest.CardDetailResponse.ID,
		PhoneNumber:  orderConfirmRequest.PhoneNumber,
		TotalPrice:   orderConfirmRequest.TotalPrice,
		Status:       model.PurchaseHistoryStatusPending,
	}
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
		o.Status == orderConfirmRequest.Status &&
		o.PhoneNumber == orderConfirmRequest.PhoneNumber &&
		o.CashBackValue == orderConfirmRequest.CashBackValue
}

func (o *OrderConfirmRequest) UnmarshalJSON(data []byte) error {
	type alias OrderConfirmRequest
	if err := json.Unmarshal(data, (*alias)(o)); err != nil {
		return err
	}
	return nil
}
