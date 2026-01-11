package stock

import "github.com/google/uuid"

type StockRepository interface {
	FindByID(id uuid.UUID) (*Stock, error)
	Create(stock *Stock) error
}
