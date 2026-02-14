package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type MarketOpenRepository interface {
	FindByTradeDate(tradeDate time.Time) (*models.MarketOpen, error)
	Create(record *models.MarketOpen) error
	UpdateCloseAt(id uuid.UUID, closeAt time.Time, isTradingDay bool) error
	DeleteBefore(t time.Time) error
}

type MarketOpenRepositoryImpl struct {
	db *gorm.DB
}

func NewMarketOpenRepository(db *gorm.DB) MarketOpenRepository {
	return &MarketOpenRepositoryImpl{db: db}
}

func (r *MarketOpenRepositoryImpl) FindByTradeDate(tradeDate time.Time) (*models.MarketOpen, error) {
	var record models.MarketOpen
	if err := r.db.Where("trade_date = ?", tradeDate).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *MarketOpenRepositoryImpl) Create(record *models.MarketOpen) error {
	if record == nil {
		return errors.New("market_open record is nil")
	}
	return r.db.Create(record).Error
}

func (r *MarketOpenRepositoryImpl) UpdateCloseAt(id uuid.UUID, closeAt time.Time, isTradingDay bool) error {
	return r.db.Model(&models.MarketOpen{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"close_at":       closeAt,
			"is_trading_day": isTradingDay,
		}).Error
}

func (r *MarketOpenRepositoryImpl) DeleteBefore(t time.Time) error {
	return r.db.
		Where("created_at < ?", t).
		Delete(&models.MarketOpen{}).Error
}
