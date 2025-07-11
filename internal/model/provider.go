package model

import "gorm.io/gorm"

type ProviderType string

const (
	ProviderTypeGRCP ProviderType = "grcp"
	ProviderTypeHTTP ProviderType = "http"
)

type Provider struct {
	gorm.Model
	Code      string     `json:"code" gorm:"unique;not null"`
	Source    string     `json:"source"`
	Type      string     `json:"type" orm:"type:provider_type; not null"`
	Weight    int        `json:"weight"`
	Suppliers []Supplier `json:"suppliers" gorm:"many2many:provider_suppliers;"`
}

func (Provider) TableName() string {
	return "provider"
}
