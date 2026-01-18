package repository

import (
	"context"

	"sun-stockanalysis-api/internal/database"
	"gorm.io/gorm"
)

type HealthRepository interface {
	CheckDBStatus(ctx context.Context) error
}

type HealthRepositoryImpl struct {
	DB *gorm.DB
}

func (hr HealthRepositoryImpl) CheckDBStatus(ctx context.Context) error {

	return database.CheckPostgresHealth(hr.DB, ctx)
}

func NewHealthRepository(DB *gorm.DB) HealthRepository {
	return HealthRepositoryImpl{
		DB: DB,
	}
}
