package schema

import (
	"encoding/json"
)

type Order struct {
	OrderID       uint   `json:"order_id"`
	UserID        uint   `json:"user_id"`
	SkuID         uint   `json:"sku_id"`
	TotalPrice    int    `json:"total_price"`
	Status        string `json:"status"`
	PhoneNumber   string `json:"phone_number"`
	CashBackValue int    `json:"cash_back_value"`
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var rawSku struct {
		OrderID uint `json:"order_id"`
		UserID  uint `json:"user_id"`
		Sku     struct {
			ID uint `json:"id"`
		} `json:"Sku"`
		TotalPrice    int    `json:"total_price"`
		Status        string `json:"status"`
		PhoneNumber   string `json:"phone_number"`
		CashBackValue int    `json:"cash_back_value"`
	}
	if err := json.Unmarshal(data, &rawSku); err != nil {
		return err
	}

	o.OrderID = rawSku.OrderID
	o.UserID = rawSku.UserID
	o.SkuID = rawSku.Sku.ID
	o.TotalPrice = rawSku.TotalPrice
	o.Status = rawSku.Status
	o.PhoneNumber = rawSku.PhoneNumber
	o.CashBackValue = rawSku.CashBackValue

	return nil
}
