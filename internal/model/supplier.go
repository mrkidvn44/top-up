package model

import (
	"gorm.io/gorm"
)

type SupplierStatus string

const (
	SupplierStatusActive   SupplierStatus = "active"
	SupplierStatusInactive SupplierStatus = "inactive"
)

type Supplier struct {
	gorm.Model
	Code    string         `json:"code" gorm:"not null;unique"`
	Name    string         `json:"name" gorm:"not null;unique"`
	LogoUrl string         `json:"logo_url" gorm:"not null"`
	Status  SupplierStatus `json:"status" gorm:"type:supplier_status; not null"`
	Providers []Provider `json:"providers" gorm:"many2many:provider_suppliers;"`
}

func (Supplier) TableName() string {
	return "supplier"
}
