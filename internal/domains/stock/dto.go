package stock

type CreateStockInput struct {
	Body struct {
		Symbol    string  `json:"symbol"`
		Name      string  `json:"name"`
		Sector    string  `json:"sector"`
		Price     float64 `json:"price"`
		Exchange  string  `json:"exchange"`
		AssetType string  `json:"assetType"`
		Currency  string  `json:"currency"`
	}
}
