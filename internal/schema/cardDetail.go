package schema

import (
	"encoding/json"
	"errors"
	"top-up-api/internal/model"
)

type CardDetailResponse struct {
	ID                uint              `json:"id"`
	CardPriceResponse CardPriceResponse `json:"card_price"`
	CashBackResponse  CashBackInterface `json:"cash_back"`
}

func (c *CardDetailResponse) UnmarshalJSON(data []byte) error {
	var rawCardDetail struct {
		ID                uint              `json:"id"`
		CardPriceResponse CardPriceResponse `json:"card_price"`
		CashBackResponse  json.RawMessage   `json:"cash_back"`
	}
	if err := json.Unmarshal(data, &rawCardDetail); err != nil {
		return err
	}

	c.ID = rawCardDetail.ID
	c.CardPriceResponse = rawCardDetail.CardPriceResponse

	var typeDetector struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(rawCardDetail.CashBackResponse, &typeDetector); err != nil {
		return err
	}

	switch typeDetector.Type {
	case string(model.CashBackTypePercentage):
		var cb CashBackPercentage
		if err := json.Unmarshal(rawCardDetail.CashBackResponse, &cb); err != nil {
			return err
		}
		c.CashBackResponse = &cb
	case string(model.CashBackTypeFixed):
		var cb CashBackFixed
		if err := json.Unmarshal(rawCardDetail.CashBackResponse, &cb); err != nil {
			return err
		}
		c.CashBackResponse = &cb
	default:
		return errors.New("unknown cashback type")
	}

	return nil
}

func CardDetailResponseFromModel(cardDetail model.CardDetail) *CardDetailResponse {
	return &CardDetailResponse{
		ID:                cardDetail.ID,
		CardPriceResponse: *CardPriceResponseFromModel(cardDetail.CardPrice),
		CashBackResponse:  BindCashBack(cardDetail.CashBack),
	}
}
