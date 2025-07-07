package model

import (
	"gorm.io/gorm"
)

type CashBackType string

const (
	CashBackTypePercentage CashBackType = "percentage"
	CashBackTypeFixed      CashBackType = "fixed"
)

type CashBack struct {
	gorm.Model
	Code  string       `json:"code" gorm:"not null;unique"`
	Type  CashBackType `json:"type" gorm:"type:cash_back_type; not null"`
	Value int          `json:"value" gorm:"not null"`
}

func (CashBack) TableName() string {
	return "cash_back"
}
