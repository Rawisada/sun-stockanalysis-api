package models

import (
	"github.com/google/uuid"
)

type CompanyNews struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Symbol  string      `gorm:"not null;varchar(64);" json:"symbol"`
	Headline string `gorm:"type:varchar(200);" json:"headline"`
	Source  string      `gorm:"not null;varchar(200);" json:"source"`
	Summary  string      `gorm:"not null;dvarchar(500);" json:"summary"`
	Url  string      `gorm:"not null;varchar(200);" json:"url"`
	CreatedAt LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CompanyNews) TableName() string {
	return "company_news"
}
