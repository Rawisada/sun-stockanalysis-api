package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email       string    `gorm:"type:varchar(64);uniqueIndex;" json:"email"`
	Password    string    `gorm:"type:varchar(128);" json:"password"`
	FirstName   string    `gorm:"type:varchar(64);" json:"first_name"`
	LastName    string    `gorm:"type:varchar(64);" json:"last_name"`
	LastLoginAt LocalTime `gorm:"autoUpdateTime" json:"last_login_at"`
	Role        string    `gorm:"not null;" json:"role"`
	IsActive    bool      `gorm:"not null;default:true;" json:"is_active"`
	CreatedAt   LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}
