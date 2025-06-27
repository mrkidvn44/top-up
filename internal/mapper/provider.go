package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func ProviderResponseFromModel(provider *model.Provider) *schema.ProviderResponse {
	return &schema.ProviderResponse{
		Code:   provider.Code,
		Name:   provider.Name,
		Logo:   provider.LogoUrl,
		Status: provider.Status,
	}
}
