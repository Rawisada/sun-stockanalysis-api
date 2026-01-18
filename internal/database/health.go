package database

import (
	"context"

	"gorm.io/gorm"
)

func CheckPostgresHealth(db *gorm.DB, ctx context.Context) error {
	return db.WithContext(ctx).Exec("SELECT 1").Error
}
