package stock_daily

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

const (
	defaultEMAPeriod20  = 20
	defaultEMAPeriod100 = 100
)

var (
	emaPeriod20  = getEnvInt("STOCK_DAILY_EMA_PERIOD20", defaultEMAPeriod20)
	emaPeriod100 = getEnvInt("STOCK_DAILY_EMA_PERIOD100", defaultEMAPeriod100)
)

type StockDailyService interface {
	BuildForWindow(ctx context.Context, start, end time.Time) error
	ListBySymbol(ctx context.Context, symbol string) ([]models.StockDaily, error)
}

type StockDailyServiceImpl struct {
	stockRepo  repository.StockRepository
	quoteRepo  repository.StockQuoteRepository
	metricRepo repository.StockDailyRepository
}

func NewStockDailyService(
	stockRepo repository.StockRepository,
	quoteRepo repository.StockQuoteRepository,
	metricRepo repository.StockDailyRepository,
) StockDailyService {
	return &StockDailyServiceImpl{
		stockRepo:  stockRepo,
		quoteRepo:  quoteRepo,
		metricRepo: metricRepo,
	}
}

func (s *StockDailyServiceImpl) BuildForWindow(ctx context.Context, start, end time.Time) error {
	symbols, err := s.stockRepo.ListSymbols()
	if err != nil || len(symbols) == 0 {
		return err
	}

	tradeDate := tradeDateFromEnd(end)

	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		quotes, err := s.quoteRepo.FindBySymbolBetween(symbol, start, end)
		if err != nil || len(quotes) == 0 {
			continue
		}

		first := quotes[0]
		last := quotes[len(quotes)-1]

		avg, high, low := summarizePrices(quotes)

		ema20 := s.calculateEMA(symbol, last.PriceCurrent, emaPeriod20)
		ema100 := s.calculateEMA(symbol, last.PriceCurrent, emaPeriod100)
		emaTrend := 0
		if ema20 > ema100 {
			emaTrend = 1
		} else if ema20 < ema100 {
			emaTrend = -1
		} else if ema20 == ema100 {
			emaTrend = 0
		}

		metric := &models.StockDaily{
			Symbol:         symbol,
			PriceAverage:   avg,
			PriceHigh:      high,
			PriceLow:       low,
			PriceOpen:      first.PriceCurrent,
			PricePrevClose: last.PriceCurrent,
			ChangePrice:    last.ChangePrice,
			ChangePercent:  last.ChangePercent,
			DeltaPrice:     high - low,
			EMA20:          ema20,
			EMA100:         ema100,
			TradeDate:      tradeDate,
			EMATrend:       emaTrend,
		}

		if err := s.metricRepo.Create(metric); err != nil {
			return err
		}
	}
	return nil
}

func (s *StockDailyServiceImpl) ListBySymbol(ctx context.Context, symbol string) ([]models.StockDaily, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return nil, errors.New("symbol is empty")
	}
	return s.metricRepo.FindPreviousBySymbol(symbol)
}

func (s *StockDailyServiceImpl) calculateEMA(symbol string, current float64, period int) float64 {
	if period <= 1 {
		return current
	}
	prev, err := s.metricRepo.FindLatestBySymbol(symbol)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return current
		}
		return current
	}
	if prev == nil {
		return current
	}
	emaPrev := current
	switch period {
	case emaPeriod20:
		emaPrev = prev.EMA20
	case emaPeriod100:
		emaPrev = prev.EMA100
	}

	alpha := 2.0 / float64(period-1)
	return alpha*current + (1-alpha)*emaPrev
}

func summarizePrices(quotes []models.StockQuote) (avg float64, high float64, low float64) {
	if len(quotes) == 0 {
		return 0, 0, 0
	}
	sum := 0.0
	high = quotes[0].PriceCurrent
	low = quotes[0].PriceCurrent
	for _, q := range quotes {
		price := q.PriceCurrent
		sum += price
		if price > high {
			high = price
		}
		if price < low {
			low = price
		}
	}
	avg = sum / float64(len(quotes))
	return avg, high, low
}

func tradeDateFromEnd(end time.Time) models.LocalDate {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	endLocal := end.In(loc)
	yesterday := endLocal.AddDate(0, 0, -1)
	date := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	return models.NewLocalDate(date)
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
