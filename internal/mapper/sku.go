package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func SkuResponseFromModel(sku model.Sku) *schema.SkuResponse {
	return &schema.SkuResponse{
		ID:                sku.ID,
		Price:             sku.Price,
		CashBackInterface: CashBackFromModel(sku.CashBack),
		ProviderInfo: schema.ProviderInfo{
			Code: sku.Provider.Code,
			Name: sku.Provider.Name,
		},
	}
}

func SkusGroupByProviderFromModel(skus []model.Sku) *[]schema.SkusGroupByProvider {
	if len(skus) == 0 {
		return nil
	}

	groupedDetails := make(map[string]schema.SkusGroupByProvider)
	for _, sku := range skus {
		providerCode := sku.Provider.Code
		if _, exists := groupedDetails[providerCode]; !exists {
			groupedDetails[providerCode] = schema.SkusGroupByProvider{
				ProviderCode:    providerCode,
				ProviderName:    sku.Provider.Name,
				ProviderLogoUrl: sku.Provider.LogoUrl,
				Skus:            []schema.SkuMiniatureResponse{},
			}
		}
		entry := groupedDetails[providerCode]
		entry.Skus = append(entry.Skus, schema.SkuMiniatureResponse{
			ID:                sku.ID,
			Price:             sku.Price,
			CashBackInterface: CashBackFromModel(sku.CashBack),
		})
		groupedDetails[providerCode] = entry
	}

	var result []schema.SkusGroupByProvider
	for _, group := range groupedDetails {
		result = append(result, group)
	}
	return &result
}
