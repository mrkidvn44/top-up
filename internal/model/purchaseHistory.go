package model

import (
	"gorm.io/gorm"
)

type PurchaseHistoryStatus string

const (
	PurchaseHistoryStatusPending PurchaseHistoryStatus = "pending"
	PurchaseHistoryStatusConfirm PurchaseHistoryStatus = "confirm"
	PurchaseHistoryStatusSuccess PurchaseHistoryStatus = "success"
	PurchaseHistoryStatusFailed  PurchaseHistoryStatus = "failed"
)

type PurchaseHistory struct {
	gorm.Model
	OrderID       uint                  `json:"order_id" gorm:"not null"`
	UserID        uint                  `json:"user_id" gorm:"not null"`
	SkuID         uint                  `json:"sku_id" gorm:"not null"`
	TotalPrice    int                   `json:"total_price" gorm:"not null"`
	PhoneNumber   string                `json:"phone_number" gorm:"not null"`
	CashBackValue int                   `json:"cash_back_value" gorm:"default:0"`
	Status        PurchaseHistoryStatus `json:"status" gorm:"type:purchase_history_status; not null"`
	Sku           Sku                   `json:"sku" gorm:"foreignKey:SkuID;references:ID"`
}

func (PurchaseHistory) TableName() string {
	return "purchase_history"
}
