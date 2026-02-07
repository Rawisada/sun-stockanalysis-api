package models

import (
	"github.com/google/uuid"
)

type RelationNews struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Symbol    string    `gorm:"type:varchar(64);" json:"symbol"`
	RelationSymbol string `gorm:"type:varchar(64);" json:"relation_symbol"`
	IsActive  bool      `gorm:"not null;default:true;" json:"is_active"`
	CreatedAt LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (RelationNews) TableName() string {
	return "relation_news"
}
