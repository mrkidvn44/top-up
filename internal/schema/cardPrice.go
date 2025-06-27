package schema

import "top-up-api/internal/model"

type CardPriceResponse struct {
	Code  string `json:"code"`
	Value int    `json:"value"`
}

func CardPriceResponseFromModel(cardPrice model.CardPrice) *CardPriceResponse {
	return &CardPriceResponse{
		Code:  cardPrice.Code,
		Value: cardPrice.Value,
	}
}
