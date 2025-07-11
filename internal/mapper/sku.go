package mapper

import (
	"top-up-api/internal/model"
	"top-up-api/internal/schema"
)

func SkuResponseFromModel(sku model.Sku) *schema.SkuResponse {
	return &schema.SkuResponse{
		ID:                sku.ID,
		Price:             sku.Price,
		CashBackInterface: CashBackFromModel(sku.CashBack),
		SupplierInfo: schema.SupplierInfo{
			Code: sku.Supplier.Code,
			Name: sku.Supplier.Name,
		},
	}
}

func SkusGroupBySupplierFromModel(skus []model.Sku) *[]schema.SkusGroupBySupplier {
	if len(skus) == 0 {
		return nil
	}

	groupedDetails := make(map[string]schema.SkusGroupBySupplier)
	for _, sku := range skus {
		supplierCode := sku.Supplier.Code
		if _, exists := groupedDetails[supplierCode]; !exists {
			groupedDetails[supplierCode] = schema.SkusGroupBySupplier{
				SupplierCode:    supplierCode,
				SupplierName:    sku.Supplier.Name,
				SupplierLogoUrl: sku.Supplier.LogoUrl,
				Skus:            []schema.SkuMiniatureResponse{},
			}
		}
		entry := groupedDetails[supplierCode]
		entry.Skus = append(entry.Skus, schema.SkuMiniatureResponse{
			ID:                sku.ID,
			Price:             sku.Price,
			CashBackInterface: CashBackFromModel(sku.CashBack),
		})
		groupedDetails[supplierCode] = entry
	}

	var result []schema.SkusGroupBySupplier
	for _, group := range groupedDetails {
		result = append(result, group)
	}
	return &result
}
