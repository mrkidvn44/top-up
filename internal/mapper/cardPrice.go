package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func CardPriceResponseFromModel(cardPrice model.CardPrice) *schema.CardPriceResponse {
	return &schema.CardPriceResponse{
		Code:  cardPrice.Code,
		Value: cardPrice.Value,
	}
}
