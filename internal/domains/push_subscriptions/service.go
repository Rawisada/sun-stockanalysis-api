package push_subscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

type SaveSubscriptionInput struct {
	DeviceID  string
	Endpoint  string
	P256DHKey string
	AuthKey   string
	UserAgent string
}

type PushSubscriptionService interface {
	GetPublicKey(ctx context.Context) (string, error)
	Save(ctx context.Context, userID string, input SaveSubscriptionInput) error
	Delete(ctx context.Context, userID, deviceID string) error
	Notify(event *models.AlertEvent, message string)
	NotifyCompanyNewsReady(message string)
	NotifyMarketOpen(message string)
	NotifyMarketClose(message string)
}

type PushSubscriptionServiceImpl struct {
	subRepo        repository.PushSubscriptionRepository
	vapidPublicKey string
	vapidPrivate   string
	subject        string
	triggerScore   int
}

func NewPushSubscriptionService(
	subRepo repository.PushSubscriptionRepository,
	pushCfg *configurations.Push,
) (PushSubscriptionService, error) {
	if subRepo == nil {
		return nil, errors.New("push subscription repository is required")
	}

	service := &PushSubscriptionServiceImpl{
		subRepo:      subRepo,
		subject:      "mailto:admin@example.com",
		triggerScore: 4,
	}

	if pushCfg != nil {
		if strings.TrimSpace(pushCfg.Subject) != "" {
			service.subject = strings.TrimSpace(pushCfg.Subject)
		}
		if pushCfg.TriggerScore > 0 {
			service.triggerScore = pushCfg.TriggerScore
		}
		service.vapidPublicKey = strings.TrimSpace(pushCfg.VAPIDPublicKey)
		service.vapidPrivate = strings.TrimSpace(pushCfg.VAPIDPrivateKey)
	}

	if err := service.ensureVAPIDKeys(); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *PushSubscriptionServiceImpl) GetPublicKey(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	if s.vapidPublicKey == "" {
		return "", errors.New("vapid public key is not configured")
	}
	return s.vapidPublicKey, nil
}

func (s *PushSubscriptionServiceImpl) Save(ctx context.Context, userID string, input SaveSubscriptionInput) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if strings.TrimSpace(userID) == "" {
		return errors.New("user id is required")
	}
	userUUID, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return errors.New("invalid user id")
	}
	if strings.TrimSpace(input.DeviceID) == "" {
		return errors.New("device_id is required")
	}
	if strings.TrimSpace(input.Endpoint) == "" {
		return errors.New("subscription endpoint is required")
	}
	if strings.TrimSpace(input.P256DHKey) == "" || strings.TrimSpace(input.AuthKey) == "" {
		return errors.New("subscription keys are required")
	}

	return s.subRepo.Upsert(&models.PushSubscription{
		UserID:    userUUID,
		DeviceID:  strings.TrimSpace(input.DeviceID),
		Endpoint:  strings.TrimSpace(input.Endpoint),
		P256DHKey: strings.TrimSpace(input.P256DHKey),
		AuthKey:   strings.TrimSpace(input.AuthKey),
		UserAgent: strings.TrimSpace(input.UserAgent),
		IsActive:  true,
	})
}

func (s *PushSubscriptionServiceImpl) Delete(ctx context.Context, userID, deviceID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if strings.TrimSpace(userID) == "" {
		return errors.New("user id is required")
	}
	if strings.TrimSpace(deviceID) == "" {
		return errors.New("device_id is required")
	}

	userUUID, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return errors.New("invalid user id")
	}
	return s.subRepo.DeleteByUserAndDevice(userUUID, strings.TrimSpace(deviceID))
}

func (s *PushSubscriptionServiceImpl) Notify(event *models.AlertEvent, message string) {
	if event == nil {
		return
	}
	score := int(event.ScoreEMA)
	if score != s.triggerScore && score != -s.triggerScore {
		return
	}

	payload, err := s.buildPopupPayload("Stock Alert", event, message)
	if err != nil {
		return
	}

	s.sendToSubscriptions(payload)
}

func (s *PushSubscriptionServiceImpl) NotifyCompanyNewsReady(message string) {
	if strings.TrimSpace(message) == "" {
		message = "ข่าววันนี้มาเเล้ว"
	}
	payload, err := s.buildPopupPayload("Company News", nil, message)
	if err != nil {
		return
	}
	s.sendToSubscriptions(payload)
}

func (s *PushSubscriptionServiceImpl) NotifyMarketOpen(message string) {
	if strings.TrimSpace(message) == "" {
		message = "ตลาดเปิดแล้ว"
	}
	payload, err := s.buildPopupPayload("Market Open", nil, message)
	if err != nil {
		return
	}
	s.sendToSubscriptions(payload)
}

func (s *PushSubscriptionServiceImpl) NotifyMarketClose(message string) {
	if strings.TrimSpace(message) == "" {
		message = "ตลาดปิดแล้ว"
	}
	payload, err := s.buildPopupPayload("Market Close", nil, message)
	if err != nil {
		return
	}
	s.sendToSubscriptions(payload)
}

func (s *PushSubscriptionServiceImpl) buildPopupPayload(title string, event *models.AlertEvent, message string) ([]byte, error) {
	return json.Marshal(struct {
		Type    string             `json:"type"`
		Title   string             `json:"title"`
		Event   *models.AlertEvent `json:"event"`
		Message string             `json:"message"`
	}{
		Type:    "popup",
		Title:   title,
		Event:   event,
		Message: message,
	})
}

func (s *PushSubscriptionServiceImpl) sendToSubscriptions(payload []byte) {
	subscriptions, err := s.subRepo.ListActive()
	if err != nil || len(subscriptions) == 0 {
		return
	}

	for _, sub := range subscriptions {
		resp, sendErr := webpush.SendNotification(payload, &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				Auth:   sub.AuthKey,
				P256dh: sub.P256DHKey,
			},
		}, &webpush.Options{
			Subscriber:      s.subject,
			VAPIDPublicKey:  s.vapidPublicKey,
			VAPIDPrivateKey: s.vapidPrivate,
			TTL:             30,
		})
		if resp != nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
				_ = s.subRepo.DeleteByEndpoint(sub.Endpoint)
			}
		}
		if sendErr != nil {
			log.Printf("push notify failed endpoint=%s err=%v", sub.Endpoint, sendErr)
		}
	}
}

func (s *PushSubscriptionServiceImpl) ensureVAPIDKeys() error {
	if s.vapidPublicKey == "" {
		return errors.New("PUSH_VAPIDPUBLICKEY is required")
	}
	if s.vapidPrivate == "" {
		return errors.New("PUSH_VAPIDPRIVATEKEY is required")
	}
	return nil
}
