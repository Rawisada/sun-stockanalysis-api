package models

import (
	"github.com/google/uuid"
)

type MarketOpen struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TradeDate    LocalDate `gorm:"type:date;not null;index:idx_market_date,unique" json:"trade_date"`
	IsTradingDay bool      `gorm:"not null;default:true" json:"is_trading_day"`
	OpenAt       LocalTime `gorm:"type:timestamptz" json:"open_at"`
	CloseAt      LocalTime `gorm:"type:timestamptz" json:"close_at"`
	CreatedAt    LocalTime `gorm:"autoCreateTime" json:"created_at"`
}

func (MarketOpen) TableName() string {
	return "market_open"
}
