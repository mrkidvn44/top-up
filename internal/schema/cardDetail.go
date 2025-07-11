package schema

import (
	"encoding/json"
	"errors"
	"top-up-api/internal/model"
)

type SkuResponse struct {
	ID                uint `json:"id"`
	Price             int  `json:"price"`
	CashBackInterface `json:"cash_back"`
	SupplierInfo      `json:"supplier"`
}

type SkuMiniatureResponse struct {
	ID                uint `json:"id"`
	Price             int  `json:"price"`
	CashBackInterface `json:"cashback"`
}

type SkusGroupBySupplier struct {
	SupplierCode    string                 `json:"supplier_code"`
	SupplierName    string                 `json:"supplier_name"`
	SupplierLogoUrl string                 `json:"supplier_logo_url"`
	Skus            []SkuMiniatureResponse `json:"sku"`
}

func (c *SkuResponse) UnmarshalJSON(data []byte) error {
	var rawSku struct {
		ID       uint            `json:"id"`
		Price    int             `json:"price"`
		CashBack json.RawMessage `json:"cash_back"`
		Supplier SupplierInfo    `json:"supplier"`
	}
	if err := json.Unmarshal(data, &rawSku); err != nil {
		return err
	}

	c.ID = rawSku.ID
	c.Price = rawSku.Price
	c.SupplierInfo = rawSku.Supplier
	var typeDetector struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(rawSku.CashBack, &typeDetector); err != nil {
		return err
	}

	switch typeDetector.Type {
	case string(model.CashBackTypePercentage):
		var cb CashBackPercentage
		if err := json.Unmarshal(rawSku.CashBack, &cb); err != nil {
			return err
		}
		c.CashBackInterface = &cb
	case string(model.CashBackTypeFixed):
		var cb CashBackFixed
		if err := json.Unmarshal(rawSku.CashBack, &cb); err != nil {
			return err
		}
		c.CashBackInterface = &cb
	default:
		return errors.New("unknown cashback type")
	}

	return nil
}
