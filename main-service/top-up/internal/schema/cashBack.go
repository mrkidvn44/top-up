package schema

import "top-up-api/internal/model"

type CashBackInterface interface {
	CalculateCashBack(value int) int
}

type CashBackPercentage struct {
	Type  model.CashBackType `json:"type"`
	Code  string             `json:"code"`
	Value int                `json:"value"`
}

type CashBackFixed struct {
	Type  model.CashBackType `json:"type"`
	Code  string             `json:"code"`
	Value int                `json:"value"`
}

func (c *CashBackPercentage) CalculateCashBack(value int) int {
	return int(float64(value) * float64(c.Value) / 100)
}

func (c *CashBackFixed) CalculateCashBack(value int) int {
	return c.Value
}
