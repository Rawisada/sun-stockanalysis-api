package stock

type CreateStockInput struct {
	Body struct {
		Symbol string `json:"symbol"`
	}
}
