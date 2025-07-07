package schema

import (
	"encoding/json"
	"errors"
	"top-up-api/internal/model"
)

type CardDetailResponse struct {
	ID                uint `json:"id"`
	CardPriceResponse `json:"card_price"`
	CashBackInterface `json:"cash_back"`
	ProviderInfo      `json:"provider"`
}

type CardDetailMiniatureResponse struct {
	ID                uint `json:"id"`
	CardPriceResponse `json:"card_price"`
	CashBackInterface `json:"provider"`
}

type CardDetailsGroupByProvider struct {
	ProviderCode    string                        `json:"provider_code"`
	ProviderName    string                        `json:"provider_name"`
	ProviderLogoUrl string                        `json:"provider_logo_url"`
	CardDetails     []CardDetailMiniatureResponse `json:"card_details"`
}

func (c *CardDetailResponse) UnmarshalJSON(data []byte) error {
	var rawCardDetail struct {
		ID        uint              `json:"id"`
		CardPrice CardPriceResponse `json:"card_price"`
		CashBack  json.RawMessage   `json:"cash_back"`
		Provider  ProviderInfo      `json:"provider"`
	}
	if err := json.Unmarshal(data, &rawCardDetail); err != nil {
		return err
	}

	c.ID = rawCardDetail.ID
	c.CardPriceResponse = rawCardDetail.CardPrice
	c.ProviderInfo = rawCardDetail.Provider
	var typeDetector struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(rawCardDetail.CashBack, &typeDetector); err != nil {
		return err
	}

	switch typeDetector.Type {
	case string(model.CashBackTypePercentage):
		var cb CashBackPercentage
		if err := json.Unmarshal(rawCardDetail.CashBack, &cb); err != nil {
			return err
		}
		c.CashBackInterface = &cb
	case string(model.CashBackTypeFixed):
		var cb CashBackFixed
		if err := json.Unmarshal(rawCardDetail.CashBack, &cb); err != nil {
			return err
		}
		c.CashBackInterface = &cb
	default:
		return errors.New("unknown cashback type")
	}

	return nil
}
