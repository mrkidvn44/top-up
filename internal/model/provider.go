package model

import (
	"gorm.io/gorm"
)

type ProviderStatus string

const (
	ProviderStatusActive   ProviderStatus = "active"
	ProviderStatusInactive ProviderStatus = "inactive"
)

type Provider struct {
	gorm.Model
	Code    string `json:"code" gorm:"not null;unique"`
	Name    string `json:"name" gorm:"not null;unique"`
	LogoUrl string `json:"logo_url" gorm:"not null"`
	Status  string `json:"status" gorm:"not null"`
}

func (Provider) TableName() string {
	return "provider"
}
