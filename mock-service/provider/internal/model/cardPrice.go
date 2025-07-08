package model

import (
	"gorm.io/gorm"
)

type CardPrice struct {
	gorm.Model
	Code  string `json:"code" gorm:"not null;unique"`
	Value int    `json:"value" gorm:"not null"`
}

func (CardPrice) TableName() string {
	return "card_price"
}
