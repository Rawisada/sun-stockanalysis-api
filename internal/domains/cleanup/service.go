package cleanup

import (
	"context"
	"time"

	"sun-stockanalysis-api/internal/repository"
)

type CleanupService interface {
	Start(ctx context.Context)
}

type CleanupServiceImpl struct {
	quoteRepo              repository.StockQuoteRepository
	companyRepo            repository.CompanyNewsRepository
	alertEventRepo         repository.AlertEventRepository
	marketOpenRepo         repository.MarketOpenRepository
	refreshTokenRepo       repository.RefreshTokenRepository
	retainDays             int
	alertRetainDays        int
	marketOpenRetainDays   int
	refreshTokenRetainDays int
}

func NewCleanupService(
	quoteRepo repository.StockQuoteRepository,
	companyRepo repository.CompanyNewsRepository,
	alertEventRepo repository.AlertEventRepository,
	marketOpenRepo repository.MarketOpenRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	retainDays int,
	alertRetainDays int,
	marketOpenRetainDays int,
	refreshTokenRetainDays int,
) CleanupService {
	if retainDays <= 0 {
		retainDays = 15
	}
	if alertRetainDays <= 0 {
		alertRetainDays = 7
	}
	if marketOpenRetainDays <= 0 {
		marketOpenRetainDays = 7
	}
	if refreshTokenRetainDays <= 0 {
		refreshTokenRetainDays = 30
	}
	return &CleanupServiceImpl{
		quoteRepo:              quoteRepo,
		companyRepo:            companyRepo,
		alertEventRepo:         alertEventRepo,
		marketOpenRepo:         marketOpenRepo,
		refreshTokenRepo:       refreshTokenRepo,
		retainDays:             retainDays,
		alertRetainDays:        alertRetainDays,
		marketOpenRetainDays:   marketOpenRetainDays,
		refreshTokenRetainDays: refreshTokenRetainDays,
	}
}

func (s *CleanupServiceImpl) Start(ctx context.Context) {
	go s.runScheduler(ctx)
}

func (s *CleanupServiceImpl) runScheduler(ctx context.Context) {
	for {
		wait := nextRunDuration(0, 5, thailandLocation())
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
		s.runOnce(ctx)
	}
}

func (s *CleanupServiceImpl) runOnce(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	loc := thailandLocation()
	now := time.Now().In(loc)
	cutoffDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).
		AddDate(0, 0, -s.retainDays)
	alertEventCutoffDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).
		AddDate(0, 0, -s.alertRetainDays)
	marketOpenCutoffDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).
		AddDate(0, 0, -s.marketOpenRetainDays)
	refreshTokenCutoffDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).
		AddDate(0, 0, -s.refreshTokenRetainDays)

	if s.quoteRepo != nil {
		_ = s.quoteRepo.DeleteBefore(cutoffDate)
	}
	if s.companyRepo != nil {
		_ = s.companyRepo.DeleteBefore(cutoffDate)
	}
	if s.alertEventRepo != nil {
		_ = s.alertEventRepo.DeleteBefore(alertEventCutoffDate)
	}
	if s.marketOpenRepo != nil {
		_ = s.marketOpenRepo.DeleteBefore(marketOpenCutoffDate)
	}
	if s.refreshTokenRepo != nil {
		_ = s.refreshTokenRepo.DeleteBefore(refreshTokenCutoffDate)
	}
}

func nextRunDuration(hour, minute int, loc *time.Location) time.Duration {
	now := time.Now().In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return time.Until(next)
}

func thailandLocation() *time.Location {
	return time.FixedZone("Asia/Bangkok", 7*60*60)
}
