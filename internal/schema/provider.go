package schema

import "top-up-api/internal/model"

type ProviderResponse struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Status string `json:"status"`
}

func ProviderResponseFromModel(provider *model.Provider) *ProviderResponse {
	return &ProviderResponse{
		Code:   provider.Code,
		Name:   provider.Name,
		Logo:   provider.LogoUrl,
		Status: provider.Status,
	}
}

