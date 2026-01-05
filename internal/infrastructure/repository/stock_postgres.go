package repository

import (
    stock "project-go-basic/internal/domain/stock";
	"gorm.io/gorm";
)

type stockRepositoryImpl struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) stock.StockRepository {
	return &stockRepositoryImpl{
		db: db,
	}
}

func (r *stockRepositoryImpl) FindByID(id uint) (*stock.Stock, error) {
	var s stock.Stock
	err := r.db.First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *stockRepositoryImpl) Create(s *stock.Stock) error {
	return r.db.Create(s).Error
}