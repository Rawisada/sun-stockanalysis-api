package stock

type CreateStockInput struct {
	Body struct {
		Symbol    string  `json:"symbol"`
		Name      string  `json:"name"`
		Sector    string  `json:"sector"`
		Exchange  string  `json:"exchange"`
		AssetType string  `json:"assetType"`
		Currency  string  `json:"currency"`
	}
}
