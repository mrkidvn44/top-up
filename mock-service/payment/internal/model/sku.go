package model

import "gorm.io/gorm"

type Sku struct {
	gorm.Model
	ProviderCode string   `json:"provider_code" gorm:"not null"`
	CashBackCode string   `json:"cash_back_code"`
	Price        int      `json:"price" gorm:"not null"`
	CashBack     CashBack `json:"cash_back" gorm:"foreignKey:CashBackCode;references:Code;default:null"`
	Provider     Provider `json:"provider" gorm:"foreignKey:ProviderCode;references:Code"`
}

func (Sku) TableName() string {
	return "card_detail"
}
