package repository

import (
	"time"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshTokens) error
	FindByHash(hash string) (*models.RefreshTokens, error)
	RevokeByHash(hash string, revokedAt float64) error
	DeleteBefore(t time.Time) error
}

type RefreshTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{db: db}
}

func (r *RefreshTokenRepositoryImpl) Create(token *models.RefreshTokens) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepositoryImpl) FindByHash(hash string) (*models.RefreshTokens, error) {
	var t models.RefreshTokens
	if err := r.db.First(&t, "token_hash = ?", hash).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepositoryImpl) RevokeByHash(hash string, revokedAt float64) error {
	return r.db.Model(&models.RefreshTokens{}).
		Where("token_hash = ?", hash).
		Update("revoked_at", revokedAt).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteBefore(t time.Time) error {
	return r.db.
		Where("created_at < ?", t).
		Delete(&models.RefreshTokens{}).Error
}
