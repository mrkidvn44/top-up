package model

import "gorm.io/gorm"

type Sku struct {
	gorm.Model
	SupplierCode string   `json:"supplier_code" gorm:"not null"`
	CashBackCode string   `json:"cash_back_code"`
	Price        int      `json:"price" gorm:"not null"`
	CashBack     CashBack `json:"cash_back" gorm:"foreignKey:CashBackCode;references:Code;default:null"`
	Supplier     Supplier `json:"supplier" gorm:"foreignKey:SupplierCode;references:Code"`
}

func (Sku) TableName() string {
	return "sku"
}
