package models

import "github.com/google/uuid"

type MasterAssetType struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(120);not null;unique" json:"name"`
	IsActive bool      `gorm:"not null;default:true" json:"is_active"`
}

func (MasterAssetType) TableName() string {
	return "master_asset_type"
}

type MasterExchange struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(120);not null;unique" json:"name"`
	IsActive bool      `gorm:"not null;default:true" json:"is_active"`
}

func (MasterExchange) TableName() string {
	return "master_exchange"
}

type MasterSector struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(120);not null;unique" json:"name"`
	IsActive bool      `gorm:"not null;default:true" json:"is_active"`
}

func (MasterSector) TableName() string {
	return "master_sector"
}
