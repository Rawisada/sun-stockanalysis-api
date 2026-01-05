package stock

type StockRepository interface {
	FindByID(id uint) (*Stock, error)
	Create(stock *Stock) error
}
