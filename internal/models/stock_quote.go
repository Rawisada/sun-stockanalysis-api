package models

import (
	"github.com/google/uuid"
)

type StockQuote struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Symbol        string    `gorm:"type:varchar(64);not null;index" json:"symbol"`
	PriceCurrency float64   `gorm:"not null" json:"price_currency"`
	ChangePrice   float64   `gorm:"not null" json:"change_price"`
	ChangePercent float64   `gorm:"not null" json:"change_percent"`
	EMA20         float64   `gorm:"column:ema_20;not null" json:"ema_20"`
	EMA100        float64   `gorm:"column:ema_100;not null" json:"ema_100"`
	TanhEMA       float64   `gorm:"column:tanh_ema;not null" json:"tanh_ema"`
	ChangeEMA20   float64   `gorm:"column:change_ema_20;not null" json:"change_ema_20"`
	ChangeTanhEMA float64   `gorm:"column:change_tanh_ema;not null" json:"change_tanh_ema"`
	EMATrend      int       `gorm:"column:ema_trend;not null" json:"ema_trend"`
	CreatedAt     LocalTime `gorm:"autoCreateTime" json:"created_at"`
}

func (StockQuote) TableName() string {
	return "stock_quotes"
}
