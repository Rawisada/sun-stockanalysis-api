package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushSubscription struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uidx_push_subscription_user_device,priority:1;index" json:"user_id"`
	DeviceID  string    `gorm:"type:varchar(128);not null;uniqueIndex:uidx_push_subscription_user_device,priority:2" json:"device_id"`
	Endpoint  string    `gorm:"type:text;not null;index" json:"endpoint"`
	P256DHKey string    `gorm:"column:p256dh_key;type:text;not null" json:"p256dh_key"`
	AuthKey   string    `gorm:"column:auth_key;type:text;not null" json:"auth_key"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	IsActive  bool      `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt LocalTime `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt LocalTime `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PushSubscription) TableName() string {
	return "push_subscriptions"
}

func (s *PushSubscription) BeforeCreate(_ *gorm.DB) error {
	now := NewLocalTime(time.Now())
	if time.Time(s.CreatedAt).IsZero() {
		s.CreatedAt = now
	}
	s.UpdatedAt = now
	return nil
}

func (s *PushSubscription) BeforeUpdate(_ *gorm.DB) error {
	s.UpdatedAt = NewLocalTime(time.Now())
	return nil
}
