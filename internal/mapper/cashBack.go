package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func CashBackFromModel(cashBack model.CashBack) schema.CashBackInterface {
	if cashBack.Type == model.CashBackTypePercentage {
		return &schema.CashBackPercentage{
			Type:  cashBack.Type,
			Code:  cashBack.Code,
			Value: cashBack.Value,
		}
	}
	return &schema.CashBackFixed{
		Type:  model.CashBackTypeFixed,
		Code:  cashBack.Code,
		Value: cashBack.Value,
	}
}
