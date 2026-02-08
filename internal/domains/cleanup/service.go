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
	quoteRepo   repository.StockQuoteRepository
	companyRepo repository.CompanyNewsRepository
	retainDays  int
}

func NewCleanupService(
	quoteRepo repository.StockQuoteRepository,
	companyRepo repository.CompanyNewsRepository,
	retainDays int,
) CleanupService {
	if retainDays <= 0 {
		retainDays = 15
	}
	return &CleanupServiceImpl{
		quoteRepo:   quoteRepo,
		companyRepo: companyRepo,
		retainDays:  retainDays,
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

	if s.quoteRepo != nil {
		_ = s.quoteRepo.DeleteBefore(cutoffDate)
	}
	if s.companyRepo != nil {
		_ = s.companyRepo.DeleteBefore(cutoffDate)
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
