package schema

import (
	"encoding/json"
)

type Order struct {
	OrderID       uint    `json:"order_id"`
	UserID        uint    `json:"user_id"`
	CardDetailID  uint    `json:"card_detail_id"`
	TotalPrice    int    `json:"total_price"`
	Status        string `json:"status"`
	PhoneNumber   string `json:"phone_number"`
	CashBackValue int    `json:"cash_back_value"`
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var rawCardDetail struct {
		OrderID    uint `json:"order_id"`
		UserID     uint `json:"user_id"`
		CardDetail struct {
			ID uint `json:"id"`
		} `json:"card_detail"`
		TotalPrice    int    `json:"total_price"`
		Status        string `json:"status"`
		PhoneNumber   string `json:"phone_number"`
		CashBackValue int    `json:"cash_back_value"`
	}
	if err := json.Unmarshal(data, &rawCardDetail); err != nil {
		return err
	}

	o.OrderID = rawCardDetail.OrderID
	o.UserID = rawCardDetail.UserID
	o.CardDetailID = rawCardDetail.CardDetail.ID
	o.TotalPrice = rawCardDetail.TotalPrice
	o.Status = rawCardDetail.Status
	o.PhoneNumber = rawCardDetail.PhoneNumber
	o.CashBackValue = rawCardDetail.CashBackValue

	return nil
}
