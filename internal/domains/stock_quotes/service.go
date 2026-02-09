package stock_quotes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"sun-stockanalysis-api/internal/domains/alert_events"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

const (
	defaultQuoteURL        = "https://finnhub.io/api/v1/quote"
	defaultQuotePollSec    = 60
	defaultQuoteTimeoutSec = 10
	defaultEMAPeriod20     = 20
	defaultEMAPeriod100    = 100
)

var (
	quoteURL     = getEnvString("QUOTE_URL", defaultQuoteURL)
	quotePoll    = time.Duration(getEnvInt("QUOTE_POLL_SECONDS", defaultQuotePollSec)) * time.Second
	quoteTimeout = time.Duration(getEnvInt("QUOTE_TIMEOUT_SECONDS", defaultQuoteTimeoutSec)) * time.Second
	emaPeriod20  = getEnvInt("EMA_PERIOD20", defaultEMAPeriod20)
	emaPeriod100 = getEnvInt("EMA_PERIOD100", defaultEMAPeriod100)
)

type StockQuoteService interface {
	Start(ctx context.Context)
	RunOnce(ctx context.Context)
	Stop()
	List(ctx context.Context, symbol string) ([]models.StockQuote, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type StockQuoteServiceImpl struct {
	stockRepo     repository.StockRepository
	quoteRepo     repository.StockQuoteRepository
	alertService  alert_events.AlertEventService
	httpClient    HTTPClient
	finnhubToken  string
	pollInterval  time.Duration
	requestTimout time.Duration
	mu            sync.Mutex
	cancel        context.CancelFunc
}

func NewStockQuoteService(
	stockRepo repository.StockRepository,
	quoteRepo repository.StockQuoteRepository,
	alertService alert_events.AlertEventService,
	httpClient HTTPClient,
	finnhubToken string,
) StockQuoteService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: quoteTimeout}
	}
	return &StockQuoteServiceImpl{
		stockRepo:     stockRepo,
		quoteRepo:     quoteRepo,
		alertService:  alertService,
		httpClient:    httpClient,
		finnhubToken:  finnhubToken,
		pollInterval:  quotePoll,
		requestTimout: quoteTimeout,
	}
}

func (s *StockQuoteServiceImpl) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		return
	}
	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	go s.run(runCtx)
}

func (s *StockQuoteServiceImpl) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel == nil {
		return
	}
	s.cancel()
	s.cancel = nil
}

func (s *StockQuoteServiceImpl) RunOnce(ctx context.Context) {
	s.fetchAndStoreAll(ctx)
}

func (s *StockQuoteServiceImpl) List(ctx context.Context, symbol string) ([]models.StockQuote, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if strings.TrimSpace(symbol) == "" {
		return s.quoteRepo.FindAll()
	}
	return s.quoteRepo.FindBySymbol(strings.TrimSpace(symbol))
}

func (s *StockQuoteServiceImpl) run(ctx context.Context) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		s.fetchAndStoreAll(ctx)
	}
}

type finnhubQuoteResponse struct {
	C  float64 `json:"c"`
	D  float64 `json:"d"`
	DP float64 `json:"dp"`
	H  float64 `json:"h"`
	L  float64 `json:"l"`
	O  float64 `json:"o"`
	PC float64 `json:"pc"`
	T  int64   `json:"t"`
}

func (s *StockQuoteServiceImpl) fetchQuote(symbol string) (*finnhubQuoteResponse, error) {
	reqURL, err := url.Parse(quoteURL)
	if err != nil {
		return nil, err
	}
	q := reqURL.Query()
	q.Set("symbol", symbol)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
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
		return nil, fmt.Errorf("finnhub quote request failed: %s", strings.TrimSpace(string(body)))
	}

	var result *finnhubQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *StockQuoteServiceImpl) fetchAndStoreAll(ctx context.Context) {
	symbols, err := s.stockRepo.ListSymbols()
	if err != nil || len(symbols) == 0 {
		return
	}

	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			return
		default:
		}
		symbol = strings.TrimSpace(symbol)
		if symbol == "" {
			continue
		}
		quote, err := s.fetchQuote(symbol)
		if err != nil || quote == nil {
			continue
		}
		prev, err := s.quoteRepo.FindLatestBySymbol(symbol)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		ema20 := s.calculateEMA(quote.C, emaPeriod20, prev)
		ema100 := s.calculateEMA(quote.C, emaPeriod100, prev)
		tanhEMA := math.Tanh(ema20-ema100) / 5.0
		changeEMA20 := 0.0
		changeTanhEMA := 0.0
		emaTrend := 0
		if ema20 > ema100 {
			emaTrend = 1
		} else if ema20 < ema100 {
			emaTrend = -1
		} else if ema20 == ema100 {
			emaTrend = 0
		}
		if prev != nil {
			changeEMA20 = ema20 - prev.EMA20
			changeTanhEMA = tanhEMA - prev.TanhEMA
		}
		createdAt := time.Now().In(time.FixedZone("Asia/Bangkok", 7*60*60)).Truncate(time.Minute)
		changePrice := quote.D
		changePercent := quote.DP
		if err := s.quoteRepo.Create(&models.StockQuote{
			Symbol:        symbol,
			PriceCurrent:  quote.C,
			ChangePrice:   &changePrice,
			ChangePercent: &changePercent,
			EMA20:         ema20,
			EMA100:        ema100,
			TanhEMA:       tanhEMA,
			ChangeEMA20:   changeEMA20,
			ChangeTanhEMA: changeTanhEMA,
			EMATrend:      emaTrend,
			CreatedAt:     models.NewLocalTime(createdAt.Truncate(time.Minute)),
		}); err != nil {
			continue
		}
		if s.alertService != nil {
			_ = s.alertService.BuildForSymbol(ctx, symbol)
		}
	}
}

func (s *StockQuoteServiceImpl) calculateEMA(current float64, period int, prev *models.StockQuote) float64 {
	if period <= 1 {
		return current
	}
	emaPrev := current
	if prev == nil {
		return current
	}
	switch period {
	case emaPeriod20:
		emaPrev = prev.EMA20
	case emaPeriod100:
		emaPrev = prev.EMA100
	}

	alpha := 2.0 / float64(period-1)
	return alpha*current + (1-alpha)*emaPrev
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
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
