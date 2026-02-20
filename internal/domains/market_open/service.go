package market_open

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
	"sun-stockanalysis-api/pkg/logger"

	"github.com/google/uuid"
)

const (
	defaultMarketStatusURL = "https://finnhub.io/api/v1/stock/market-status?exchange=US"
	defaultPollSeconds     = 60
	defaultStopHour        = 4
	defaultStopMinute      = 30
	defaultSchedulerHour   = 00
	defaultSchedulerMinute = 16
)

var (
	marketStatusURL = getEnvString("MARKET_STATUS_URL", defaultMarketStatusURL)
	pollInterval    = time.Duration(getEnvInt("MARKET_POLL_SECONDS", defaultPollSeconds)) * time.Second
	stopHour        = getEnvInt("MARKET_STOP_HOUR", defaultStopHour)
	stopMinute      = getEnvInt("MARKET_STOP_MINUTE", defaultStopMinute)
	schedulerHour   = getEnvInt("MARKET_SCHEDULER_HOUR", defaultSchedulerHour)
	schedulerMinute = getEnvInt("MARKET_SCHEDULER_MINUTE", defaultSchedulerMinute)
)

type MarketOpenService interface {
	Start(ctx context.Context)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type StockQuoteService interface {
	Start(ctx context.Context)
	RunOnce(ctx context.Context)
	Stop()
}

type StockDailyService interface {
	BuildForWindow(ctx context.Context, start, end time.Time) error
}

type MarketOpenNotifier interface {
	NotifyMarketOpen(message string)
	NotifyMarketClose(message string)
}

type MarketOpenServiceImpl struct {
	repo         repository.MarketOpenRepository
	httpClient   HTTPClient
	finnhubToken string
	quoteService StockQuoteService
	dailyService StockDailyService
	notifier     MarketOpenNotifier
	log          *logger.Logger
}

func NewMarketOpenService(
	repo repository.MarketOpenRepository,
	httpClient HTTPClient,
	finnhubToken string,
	quoteService StockQuoteService,
	dailyService StockDailyService,
	notifier MarketOpenNotifier,
	log *logger.Logger,
) MarketOpenService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &MarketOpenServiceImpl{
		repo:         repo,
		httpClient:   httpClient,
		finnhubToken: finnhubToken,
		quoteService: quoteService,
		dailyService: dailyService,
		notifier:     notifier,
		log:          log,
	}
}

func (s *MarketOpenServiceImpl) Start(ctx context.Context) {
	go s.runScheduler(ctx)
}

func (s *MarketOpenServiceImpl) runScheduler(ctx context.Context) {
	// TEMP: run immediately on startup.
	// s.runDailyPolling(ctx)

	for {
		wait := nextRunDuration(schedulerHour, schedulerMinute, time.Local)
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}

		s.runDailyPolling(ctx)
	}
}

func (s *MarketOpenServiceImpl) runDailyPolling(ctx context.Context) {
	quoteStarted := false
	postHandled := false

	for {
		// Hard stop for the current market window (e.g. 04:00 next day).
		if shouldStopForDay(time.Now()) {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		status, err := s.fetchMarketStatus()
		if err != nil {
			return
		}
		if status == nil {
			return
		}

		session := strings.ToLower(strings.TrimSpace(status.session()))
		isOpen := status.isOpen()
		s.logStatus(status)
		switch {
		case session == "pre-market":
			sleepContext(ctx, pollInterval)
			continue
		case session == "regular" && isOpen:
			_ = s.ensureOpenRecord(status)
			if s.quoteService != nil && !quoteStarted {
				s.quoteService.Start(ctx)
				quoteStarted = true
				if s.notifier != nil {
					s.notifier.NotifyMarketOpen("The market is open. Prices are being updated.")
				}
			}
			sleepContext(ctx, pollInterval)
			continue
		case session == "post-market" || (session == "regular" && !isOpen):
			_ = s.updateCloseRecord(status)
			if s.quoteService != nil && !postHandled {
				s.quoteService.RunOnce(ctx)
				s.quoteService.Stop()
				postHandled = true
				if s.notifier != nil {
					s.notifier.NotifyMarketClose("ตลาดปิดแล้ว")
				}
				if s.dailyService != nil {
					start, end := metricsWindow(time.Now())
					_ = s.dailyService.BuildForWindow(ctx, start, end)
				}
			}
			sleepContext(ctx, pollInterval)
			continue
		default:
			sleepContext(ctx, pollInterval)
			continue
		}
	}
}

func (s *MarketOpenServiceImpl) logStatus(status *finnhubMarketStatusResponse) {
	correlationID := uuid.NewString()
	session := ""
	exchange := ""
	isOpen := false
	timestamp := int64(0)
	timezone := ""
	if status != nil {
		session = status.session()
		exchange = status.Exchange
		isOpen = status.IsOpen
		timestamp = status.T
		timezone = status.Timezone
	}
	message := fmt.Sprintf("market status: session=%s exchange=%s isOpen=%t t=%d timezone=%s",
		session,
		exchange,
		isOpen,
		timestamp,
		timezone,
	)
	if s.log != nil {
		s.log.With("correlation_id", correlationID).Infof(message)
		return
	}
	_ = correlationID
}

type finnhubMarketStatusResponse struct {
	Exchange string  `json:"exchange"`
	Holiday  *string `json:"holiday"`
	IsOpen   bool    `json:"isOpen"`
	Session  *string `json:"session"`
	T        int64   `json:"t"`
	Timezone string  `json:"timezone"`
}

func (r *finnhubMarketStatusResponse) session() string {
	if r == nil || r.Session == nil {
		return ""
	}
	return *r.Session
}

func (r *finnhubMarketStatusResponse) isOpen() bool {
	if r == nil {
		return false
	}
	return r.IsOpen
}

func (s *MarketOpenServiceImpl) fetchMarketStatus() (*finnhubMarketStatusResponse, error) {
	req, err := http.NewRequest(http.MethodGet, marketStatusURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Finnhub-Token", s.finnhubToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("finnhub market-status request failed: %s", strings.TrimSpace(string(body)))
	}

	var result *finnhubMarketStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MarketOpenServiceImpl) ensureOpenRecord(status *finnhubMarketStatusResponse) error {
	tradeDate, openAt := tradeDateAndTime(status)
	tradeDate = tradeDateForMarketWindow(openAt)
	_, err := s.repo.FindByTradeDate(tradeDate)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	record := &models.MarketOpen{
		TradeDate:    models.NewLocalDate(tradeDate),
		IsTradingDay: true,
		OpenAt:       models.NewLocalTime(openAt),
	}
	return s.repo.Create(record)
}

func (s *MarketOpenServiceImpl) updateCloseRecord(status *finnhubMarketStatusResponse) error {
	tradeDate, closeAt := tradeDateAndTime(status)
	tradeDate = tradeDateForMarketWindow(closeAt)
	record, err := s.repo.FindByTradeDate(tradeDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return s.repo.UpdateCloseAt(record.ID, closeAt, false)
}

func tradeDateAndTime(status *finnhubMarketStatusResponse) (time.Time, time.Time) {
	if status == nil {
		now := time.Now()
		date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return date, now
	}

	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	timestamp := status.T
	if timestamp <= 0 {
		now := time.Now().In(loc)
		date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		return date, now
	}

	at := time.Unix(timestamp, 0).In(loc)
	date := time.Date(at.Year(), at.Month(), at.Day(), 0, 0, 0, 0, loc)
	return date, at
}

func tradeDateForMarketWindow(at time.Time) time.Time {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	t := at.In(loc)
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	if t.Hour() < 4 {
		return date.AddDate(0, 0, -1)
	}
	return date
}

func shouldStopForDay(now time.Time) bool {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	current := now.In(loc)
	stopAt := time.Date(current.Year(), current.Month(), current.Day(), stopHour, stopMinute, 0, 0, loc)
	return !current.Before(stopAt)
}

func metricsWindow(now time.Time) (time.Time, time.Time) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	current := now.In(loc)
	start := time.Date(current.Year(), current.Month(), current.Day(), 20, 0, 0, 0, loc).Add(-24 * time.Hour)
	end := time.Date(current.Year(), current.Month(), current.Day(), 4, 30, 0, 0, loc)
	return start, end
}

func nextRunDuration(hour, minute int, loc *time.Location) time.Duration {
	now := time.Now().In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return time.Until(next)
}

func sleepContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func getEnvString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}
