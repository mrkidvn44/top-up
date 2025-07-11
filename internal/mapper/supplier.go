package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func SupplierResponseFromModel(supplier *model.Supplier) *schema.SupplierResponse {
	return &schema.SupplierResponse{
		Code:   supplier.Code,
		Name:   supplier.Name,
		Logo:   supplier.LogoUrl,
		Status: supplier.Status,
	}
}
