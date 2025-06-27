package model

import (
	"gorm.io/gorm"
)

type PurchaseHistoryStatus string

const (
	PurchaseHistoryStatusPending PurchaseHistoryStatus = "pending"
	PurchaseHistoryStatusSuccess PurchaseHistoryStatus = "success"
	PurchaseHistoryStatusFailed  PurchaseHistoryStatus = "failed"
)

type PurchaseHistory struct {
	gorm.Model
	OrderID      uint                  `json:"order_id" gorm:"not null"`
	UserID       uint                  `json:"user_id" gorm:"not null"`
	CardDetailID uint                  `json:"card_detail_id" gorm:"not null"`
	TotalPrice   int                   `json:"total_price" gorm:"not null"`
	PhoneNumber  string                `json:"phone_number" gorm:"not null"`
	Status       PurchaseHistoryStatus `json:"status" gorm:"type:purchase_history_status; not null"`
	User         User                  `json:"user" gorm:"foreignKey:UserID;references:ID"`
	CardDetail   CardDetail            `json:"card_detail" gorm:"foreignKey:CardDetailID;references:ID"`
}

func (PurchaseHistory) TableName() string {
	return "purchase_history"
}
