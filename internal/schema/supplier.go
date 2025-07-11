package schema

import "top-up-api/internal/model"

type SupplierResponse struct {
	Code   string               `json:"code"`
	Name   string               `json:"name"`
	Logo   string               `json:"logo"`
	Status model.SupplierStatus `json:"status"`
}

type SupplierInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
