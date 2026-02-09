package alert_events

import (
	"context"
	"time"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/realtime"
	"sun-stockanalysis-api/internal/repository"
)

type AlertEventService interface {
	BuildForSymbol(ctx context.Context, symbol string) error
}

type AlertEventServiceImpl struct {
	quoteRepo repository.StockQuoteRepository
	eventRepo repository.AlertEventRepository
	notifier  realtime.AlertEventNotifier
}

func NewAlertEventService(
	quoteRepo repository.StockQuoteRepository,
	eventRepo repository.AlertEventRepository,
	notifier realtime.AlertEventNotifier,
) AlertEventService {
	return &AlertEventServiceImpl{
		quoteRepo: quoteRepo,
		eventRepo: eventRepo,
		notifier:  notifier,
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
	if err != nil || len(quotes) < 5 {
		return err
	}

	trendEMA20 := 0
	for _, q := range quotes {
		trendEMA20 += signFloat(q.ChangeEMA20)
	}

	latest := quotes[0]
	trendTanhEMA := signFloat(latest.ChangeTanhEMA)

	scoreEMA, ok := scoreFromTrend(trendEMA20, trendTanhEMA)
	if !ok {
		return nil
	}

	scorePCrossEMA := 0
	if latest.PriceCurrent == latest.EMA100 {
		scorePCrossEMA = 1
	}

	event := &models.AlertEvent{
		Symbol:         symbol,
		TrendEMA20:     trendEMA20,
		TrendTanhEMA:   trendTanhEMA,
		ScoreEMA:       float64(scoreEMA),
		ScorePCrossEMA: float64(scorePCrossEMA),
	}
	if scoreEMA == 2 || scoreEMA == 1 || scoreEMA == -1 || scoreEMA == -2 {
		if err := s.eventRepo.Create(event); err != nil {
			return err
		}
	}
	if s.notifier != nil {
		message := messageForScore(scoreEMA)
		if message != "" {
			s.notifier.Notify(event, message)
		}
	}
	return nil
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

func messageForScore(scoreEMA int) string {
	switch scoreEMA {
	case 1:
		return "ควรซื้อ/จับตามอง"
	case 2:
		return "ต้องซื้อ"
	case -1:
		return "ควรขาย/จับตามอง"
	case -2:
		return "ต้องขาย"
	default:
		return ""
	}
}
