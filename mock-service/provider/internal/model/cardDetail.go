package model

import "gorm.io/gorm"

type CardDetail struct {
	gorm.Model
	ProviderCode  string    `json:"provider_code" gorm:"not null"`
	CashBackCode  string    `json:"cash_back_code"`
	CardPriceCode string    `json:"card_price_code" gorm:"not null"`
	CashBack      CashBack  `json:"cash_back" gorm:"foreignKey:CashBackCode;references:Code;default:null"`
	CardPrice     CardPrice `json:"card_price" gorm:"foreignKey:CardPriceCode;references:Code;"`
	Provider      Provider  `json:"provider" gorm:"foreignKey:ProviderCode;references:Code"`
}

func (CardDetail) TableName() string {
	return "card_detail"
}
