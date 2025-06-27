package schema

type ProviderResponse struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Status string `json:"status"`
}

type ProviderInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
