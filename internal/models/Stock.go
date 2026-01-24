package models

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Symbol      string    `gorm:"type:varchar(64);" json:"symbol"`
	Name        string    `gorm:"type:varchar(128);" json:"name"`
	Sector      string    `gorm:"type:varchar(64);" json:"sector"`
	Price       float64   `gorm:"not null;" json:"price"`
	Exchange    string    `gorm:"type:varchar(64);not null;" json:"exchange"`
	AssetType   string    `gorm:"type:varchar(64);not null;" json:"asset_type"`
	Currency    string    `gorm:"type:varchar(10);not null;" json:"currency"`
	IsArchive   bool      `gorm:"not null;default:false;" json:"is_archive"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
