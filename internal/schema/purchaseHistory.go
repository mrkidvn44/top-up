package schema

type PurchaseHistoryResponse struct {
	OrderID       uint               `json:"order_id"`
	UserID        uint               `json:"user_id"`
	CardDetailID  uint               `json:"card_detail_id"`
	TotalPrice    int                `json:"total_price"`
	PhoneNumber   string             `json:"phone_number"`
	Status        string             `json:"status"`
	CashBackValue int                `json:"cash_back_value"`
	CardDetail    CardDetailResponse `json:"card_detail"`
}
