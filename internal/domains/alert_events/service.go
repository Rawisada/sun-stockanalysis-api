package alert_events

import (
	"context"
	"time"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

type AlertEventService interface {
	BuildForSymbol(ctx context.Context, symbol string) error
}

type AlertEventServiceImpl struct {
	quoteRepo repository.StockQuoteRepository
	eventRepo repository.AlertEventRepository
}

func NewAlertEventService(
	quoteRepo repository.StockQuoteRepository,
	eventRepo repository.AlertEventRepository,
) AlertEventService {
	return &AlertEventServiceImpl{
		quoteRepo: quoteRepo,
		eventRepo: eventRepo,
	}
}

func (s *AlertEventServiceImpl) BuildForSymbol(ctx context.Context, symbol string) error {
	if symbol == "" {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	start, end := dayBoundsBangkok(time.Now())
	quotes, err := s.quoteRepo.FindLatestBySymbolBetween(symbol, start, end, 5)
	if err != nil || len(quotes) == 0 {
		return err
	}

	trendEMA20 := 0
	for _, q := range quotes {
		trendEMA20 += signFloat(q.ChangeEMA20)
	}

	latest := quotes[0]
	trendTanhEMA := signFloat(latest.ChangeTanhEMA)

	scoreema, ok := scoreFromTrend(trendEMA20, trendTanhEMA)
	if !ok {
		return nil
	}

	if latest.PriceCurrency == latest.EMA100 {
		scorepcrossema += 1
	}

	event := &models.AlertEvent{
		Symbol:       symbol,
		TrendEMA20:   trendEMA20,
		TrendTanhEMA: trendTanhEMA,
		ScoreEMA:        float64(scoreema),
		ScorePCrossEMA: float64(scorepcrossema),
	}
	return s.eventRepo.Create(event)
}

func dayBoundsBangkok(t time.Time) (time.Time, time.Time) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	local := t.In(loc)
	start := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return start, end
}

func signFloat(v float64) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	default:
		return 0
	}
}

func scoreFromTrend(trendEMA20, trendTanhEMA int) (int, bool) {
	switch trendEMA20 {
	case -2:
		switch trendTanhEMA {
		case -1:
			return -2, true
		case 0, 1:
			return -1, true
		}
	case 2:
		switch trendTanhEMA {
		case 1:
			return 2, true
		case 0, -1:
			return 1, true
		}
	}
	return 0, false
}
