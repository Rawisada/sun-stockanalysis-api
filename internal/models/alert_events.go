package models

import "github.com/google/uuid"

type AlertEvent struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Symbol       string    `gorm:"type:varchar(64);not null;index" json:"symbol"`
	TrendEMA20   int       `gorm:"column:trend_ema_20;not null" json:"trend_ema_20"`
	TrendTanhEMA int       `gorm:"column:trend_tanh_ema;not null" json:"trend_tanh_ema"`
	TrendCurrency int       `gorm:"column:trend_currency;not null" json:"trend_currency"`
	ScoreEMA        float64   `gorm:"not null" json:"score_ema"`
	ScorePCrossEMA        float64   `gorm:"not null" json:"score_p_cross_ema"`
	CreatedAt    LocalTime `gorm:"autoCreateTime" json:"created_at"`
}

func (AlertEvent) TableName() string {
	return "alert_events"
}
