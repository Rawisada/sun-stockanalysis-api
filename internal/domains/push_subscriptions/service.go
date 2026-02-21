package push_subscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Delete(ctx context.Context, userID, deviceID string) error
	Notify(event *models.AlertEvent, message string)
	NotifyCompanyNewsReady(message string)
	NotifyMarketOpen(message string)
	NotifyMarketClose(message string)
	StartSimulation(ctx context.Context, interval time.Duration, message string)
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
		subject:      "admin@example.com",
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
	service.subject = normalizeWebPushSubject(service.subject)

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

	endpoint := strings.TrimSpace(input.Endpoint)
	p256dhKey := strings.TrimSpace(input.P256DHKey)
	authKey := strings.TrimSpace(input.AuthKey)
	deviceID := strings.TrimSpace(input.DeviceID)
	userAgent := strings.TrimSpace(input.UserAgent)
	log.Printf(
		"push subscription upsert request user_id=%s device_id=%s endpoint=%s p256dh=%s auth=%s user_agent=%q",
		userUUID.String(),
		deviceID,
		endpoint,
		maskKey(p256dhKey),
		maskKey(authKey),
		userAgent,
	)

	return s.subRepo.Upsert(&models.PushSubscription{
		UserID:    userUUID,
		DeviceID:  deviceID,
		Endpoint:  endpoint,
		P256DHKey: p256dhKey,
		AuthKey:   authKey,
		UserAgent: userAgent,
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

	s.sendToSubscriptions("Stock Alert", payload)
}

func (s *PushSubscriptionServiceImpl) NotifyCompanyNewsReady(message string) {
	if strings.TrimSpace(message) == "" {
		message = "ข่าววันนี้มาเเล้ว"
	}
	payload, err := s.buildPopupPayload("Company News", nil, message)
	if err != nil {
		return
	}
	s.sendToSubscriptions("Company News", payload)
}

func (s *PushSubscriptionServiceImpl) NotifyMarketOpen(message string) {
	if strings.TrimSpace(message) == "" {
		message = "ตลาดเปิดแล้ว"
	}
	payload, err := s.buildPopupPayload("Market Open", nil, message)
	if err != nil {
		return
	}
	s.sendToSubscriptions("Market Open", payload)
}

func (s *PushSubscriptionServiceImpl) NotifyMarketClose(message string) {
	if strings.TrimSpace(message) == "" {
		message = "ตลาดปิดแล้ว"
	}
	payload, err := s.buildPopupPayload("Market Close", nil, message)
	if err != nil {
		return
	}
	s.sendToSubscriptions("Market Close", payload)
}

func (s *PushSubscriptionServiceImpl) StartSimulation(ctx context.Context, interval time.Duration, message string) {
	if interval <= 0 {
		interval = time.Minute
	}
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "Simulation push notification"
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.Printf("push simulation started interval=%s", interval)
		for {
			select {
			case <-ctx.Done():
				log.Printf("push simulation stopped: %v", ctx.Err())
				return
			case <-ticker.C:
				timestampedMessage := fmt.Sprintf("%s (%s)", msg, time.Now().Format(time.RFC3339))
				payload, err := s.buildPopupPayload("Simulation", nil, timestampedMessage)
				if err != nil {
					log.Printf("push simulation payload build failed err=%v", err)
					continue
				}
				s.sendToSubscriptions("Simulation", payload)
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

type pushSendResult struct {
	total   int
	success int
	failed  int
	removed int
	forbidden int
	err     error
}

func (s *PushSubscriptionServiceImpl) sendToSubscriptions(title string, payload []byte) pushSendResult {
	result := pushSendResult{}
	subscriptions, err := s.subRepo.ListActive()
	if err != nil {
		result.err = err
		log.Printf("push notify result title=%s err=%v", title, result.err)
		return result
	}
	if len(subscriptions) == 0 {
		log.Printf("push notify result title=%s total=0 message=no_active_subscriptions", title)
		return result
	}

	result.total = len(subscriptions)
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
		statusCode := 0
		responseBody := ""
		if resp != nil {
			statusCode = resp.StatusCode
			body, readErr := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if readErr != nil {
				log.Printf("push notify response read failed endpoint=%s status=%d err=%v", sub.Endpoint, statusCode, readErr)
			} else {
				responseBody = strings.TrimSpace(string(body))
			}
			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
				_ = s.subRepo.DeleteByEndpoint(sub.Endpoint)
				result.removed++
			}
		}
		if statusCode == http.StatusForbidden {
			result.forbidden++
		}
		if sendErr != nil {
			if statusCode > 0 {
				log.Printf("push notify failed endpoint=%s status=%d reason=%q err=%v", sub.Endpoint, statusCode, responseBody, sendErr)
			} else {
				log.Printf("push notify failed endpoint=%s err=%v", sub.Endpoint, sendErr)
			}
			result.failed++
			continue
		}
		if statusCode > 0 && (statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices) {
			log.Printf("push notify non-2xx endpoint=%s status=%d reason=%q", sub.Endpoint, statusCode, responseBody)
			result.failed++
			continue
		}
		result.success++
	}
	log.Printf(
		"push notify result title=%s total=%d success=%d failed=%d removed=%d forbidden=%d",
		title,
		result.total,
		result.success,
		result.failed,
		result.removed,
		result.forbidden,
	)
	return result
}

func maskKey(v string) string {
	if len(v) <= 10 {
		return v
	}
	return v[:6] + "..." + v[len(v)-4:]
}

func normalizeWebPushSubject(subject string) string {
	s := strings.TrimSpace(subject)
	s = strings.TrimPrefix(s, "mailto:")
	s = strings.TrimSpace(s)
	if s == "" {
		return "admin@example.com"
	}
	return s
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


