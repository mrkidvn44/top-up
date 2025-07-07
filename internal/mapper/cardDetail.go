package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func CardDetailResponseFromModel(cardDetail model.CardDetail) *schema.CardDetailResponse {
	return &schema.CardDetailResponse{
		ID:                cardDetail.ID,
		CardPriceResponse: *CardPriceResponseFromModel(cardDetail.CardPrice),
		CashBackInterface: CashBackFromModel(cardDetail.CashBack),
		ProviderInfo: schema.ProviderInfo{
			Code: cardDetail.Provider.Code,
			Name: cardDetail.Provider.Name,
		},
	}
}

func CardDetailsGroupByProviderFromModel(cardDetails []model.CardDetail) *[]schema.CardDetailsGroupByProvider {
	if len(cardDetails) == 0 {
		return nil
	}

	groupedDetails := make(map[string]schema.CardDetailsGroupByProvider)
	for _, cardDetail := range cardDetails {
		providerCode := cardDetail.Provider.Code
		if _, exists := groupedDetails[providerCode]; !exists {
			groupedDetails[providerCode] = schema.CardDetailsGroupByProvider{
				ProviderCode:    providerCode,
				ProviderName:    cardDetail.Provider.Name,
				ProviderLogoUrl: cardDetail.Provider.LogoUrl,
				CardDetails:     []schema.CardDetailMiniatureResponse{},
			}
		}
		entry := groupedDetails[providerCode]
		entry.CardDetails = append(entry.CardDetails, schema.CardDetailMiniatureResponse{
			ID:                cardDetail.ID,
			CardPriceResponse: *CardPriceResponseFromModel(cardDetail.CardPrice),
			CashBackInterface: CashBackFromModel(cardDetail.CashBack),
		})
		groupedDetails[providerCode] = entry
	}

	var result []schema.CardDetailsGroupByProvider
	for _, group := range groupedDetails {
		result = append(result, group)
	}
	return &result
}
