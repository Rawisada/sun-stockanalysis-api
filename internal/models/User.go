package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email       string    `gorm:"type:varchar(64);uniqueIndex;" json:"email"`
	Password    string    `gorm:"type:varchar(128);" json:"password"`
	FirstName   string    `gorm:"type:varchar(64);" json:"first_name"`
	LastName    string    `gorm:"type:varchar(64);" json:"last_name"`
	LastLoginAt time.Time `gorm:"autoUpdateTime" json:"last_login_at"`
	Role        string    `gorm:"not null;" json:"role"`
	IsArchive   bool      `gorm:"not null;default:false;" json:"is_archive"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
