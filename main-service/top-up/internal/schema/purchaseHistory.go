package schema

type PurchaseHistoryResponse struct {
	OrderID       uint        `json:"order_id"`
	UserID        uint        `json:"user_id"`
	SkuID         uint        `json:"sku_id"`
	TotalPrice    int         `json:"total_price"`
	PhoneNumber   string      `json:"phone_number"`
	Status        string      `json:"status"`
	CashBackValue int         `json:"cash_back_value"`
	Sku           SkuResponse `json:"sku"`
}
