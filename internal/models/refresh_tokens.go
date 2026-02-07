package models

import (
	"github.com/google/uuid"
)

type RefreshTokens struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(64);" json:"user_id"`
	TokenHash string    `gorm:"type:varchar(128);" json:"token_hash"`
	ExpiresAt string    `gorm:"type:varchar(64);" json:"expires_at"`
	RevokedAt float64   `gorm:"not null;" json:"revoked_at"`
	CreatedAt LocalTime `gorm:"autoCreateTime" json:"created_at"`
}
