package model

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Sku{},
		&CashBack{},
		&Supplier{},
		&PurchaseHistory{},
	}
}
