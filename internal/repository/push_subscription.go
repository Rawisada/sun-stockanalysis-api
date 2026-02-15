package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sun-stockanalysis-api/internal/models"
)

type PushSubscriptionRepository interface {
	Upsert(subscription *models.PushSubscription) error
	ListActive() ([]models.PushSubscription, error)
	DeleteByEndpoint(endpoint string) error
	DeleteByUserAndDevice(userID uuid.UUID, deviceID string) error
}

type PushSubscriptionRepositoryImpl struct {
	db *gorm.DB
}

func NewPushSubscriptionRepository(db *gorm.DB) PushSubscriptionRepository {
	return &PushSubscriptionRepositoryImpl{db: db}
}

func (r *PushSubscriptionRepositoryImpl) Upsert(subscription *models.PushSubscription) error {
	if subscription == nil {
		return errors.New("push subscription is nil")
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "device_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"endpoint":   subscription.Endpoint,
			"p256dh_key": subscription.P256DHKey,
			"auth_key":   subscription.AuthKey,
			"user_agent": subscription.UserAgent,
			"is_active":  true,
		}),
	}).Create(subscription).Error
}

func (r *PushSubscriptionRepositoryImpl) ListActive() ([]models.PushSubscription, error) {
	var subscriptions []models.PushSubscription
	if err := r.db.
		Where("is_active = ?", true).
		Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *PushSubscriptionRepositoryImpl) DeleteByEndpoint(endpoint string) error {
	if endpoint == "" {
		return nil
	}
	return r.db.
		Where("endpoint = ?", endpoint).
		Delete(&models.PushSubscription{}).Error
}

func (r *PushSubscriptionRepositoryImpl) DeleteByUserAndDevice(userID uuid.UUID, deviceID string) error {
	if userID == uuid.Nil || deviceID == "" {
		return nil
	}
	return r.db.
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Delete(&models.PushSubscription{}).Error
}
