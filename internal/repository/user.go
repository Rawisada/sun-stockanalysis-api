package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	UpdateLastLogin(id uuid.UUID, when time.Time) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, "email = ? AND is_active = true", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositoryImpl) FindByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, "id = ? AND is_active = true", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepositoryImpl) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).
		Where("email = ? AND is_active = true", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) UpdateLastLogin(id uuid.UUID, when time.Time) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login_at", when).Error
}
