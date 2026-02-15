package push_subscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

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
	Notify(event *models.AlertEvent, message string)
	StartSimulation(ctx context.Context)
}

type PushSubscriptionServiceImpl struct {
	subRepo            repository.PushSubscriptionRepository
	vapidPublicKey     string
	vapidPrivate       string
	subject            string
	triggerScore       int
	simulationEnabled  bool
	simulationInterval time.Duration
}

func NewPushSubscriptionService(
	subRepo repository.PushSubscriptionRepository,
	pushCfg *configurations.Push,
) (PushSubscriptionService, error) {
	if subRepo == nil {
		return nil, errors.New("push subscription repository is required")
	}

	service := &PushSubscriptionServiceImpl{
		subRepo:            subRepo,
		subject:            "mailto:admin@example.com",
		triggerScore:       4,
		simulationInterval: 5 * time.Minute,
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
		service.simulationEnabled = pushCfg.SimulationEnabled
		if pushCfg.SimulationInterval > 0 {
			service.simulationInterval = pushCfg.SimulationInterval
		}
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

func (s *PushSubscriptionServiceImpl) StartSimulation(ctx context.Context) {
	if !s.simulationEnabled {
		return
	}
	interval := s.simulationInterval
	if interval <= 0 {
		interval = 5 * time.Minute
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// synthetic event for quick end-to-end push verification
				simulated := &models.AlertEvent{
					Symbol:       "SIMULATION",
					TrendEMA20:   4,
					TrendTanhEMA: 1,
					ScoreEMA:     4,
				}
				payload, err := s.buildPopupPayload("Simulation Alert", simulated, "ทดสอบแจ้งเตือนทุก 5 นาที")
				if err != nil {
					continue
				}
				s.sendToSubscriptions(payload)
			}
		}
	}()
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
