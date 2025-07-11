package model

func GetModels() []interface{} {
	return []interface{}{
		&Sku{},
		&CashBack{},
		&Provider{},
		&Supplier{},
		&PurchaseHistory{},
	}
}
