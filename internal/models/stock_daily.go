package models

import (
	"github.com/google/uuid"
)

type StockDaily struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Symbol         string    `gorm:"type:varchar(64);not null;index" json:"symbol"`
	PriceAverage   float64   `gorm:"not null" json:"price_average"`
	PriceHigh      float64   `gorm:"not null" json:"price_high"`
	PriceLow       float64   `gorm:"not null" json:"price_low"`
	PriceOpen      float64   `gorm:"not null" json:"price_open"`
	PricePrevClose float64   `gorm:"not null" json:"price_prev_close"`
	ChangePrice    float64   `gorm:"not null" json:"change_price"`
	ChangePercent  float64   `gorm:"not null" json:"change_percent"`
	DeltaPrice     float64   `gorm:"column:dalta_price;not null" json:"dalta_price"`
	EMA20          float64   `gorm:"column:ema_20;not null" json:"ema_20"`
	EMA100         float64   `gorm:"column:ema_100;not null" json:"ema_100"`
	EMATrend       int       `gorm:"column:ema_trend;not null" json:"ema_trend"`
	TradeDate      LocalDate `gorm:"column:trend_date;not null" json:"trend_date"`
	CreatedAt      LocalTime `gorm:"autoCreateTime" json:"created_at"`
}

func (StockDaily) TableName() string {
	return "stock_daily"
}
