package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName       string            `json:"first_name" gorm:"not null"`
	LastName        string            `json:"last_name" gorm:"not null"`
	PhoneNumber     string            `json:"phone_number" gorm:"not null;unique"`
	Balance         int               `json:"balance" gorm:"not null"`
	CashBack        int               `json:"cash_back" gorm:"not null"`
	Password        string            `json:"password" gorm:"not null"`
	PurchaseHistory []PurchaseHistory `json:"purchase_history" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "user"
}
