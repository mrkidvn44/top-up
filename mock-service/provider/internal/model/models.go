package model

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&CardDetail{},
		&CardPrice{},
		&CashBack{},
		&Provider{},
		&PurchaseHistory{},
	}
}
