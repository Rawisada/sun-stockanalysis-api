package company_news

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
	"sun-stockanalysis-api/pkg/logger"
)

const (
	companyNewsURL = "https://finnhub.io/api/v1/company-news"
)

type CompanyNewsService interface {
	Start(ctx context.Context)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CompanyNewsServiceImpl struct {
	relationRepo repository.RelationNewsRepository
	companyRepo  repository.CompanyNewsRepository
	httpClient   HTTPClient
	finnhubToken string
	log          *logger.Logger
}

func NewCompanyNewsService(
	relationRepo repository.RelationNewsRepository,
	companyRepo repository.CompanyNewsRepository,
	httpClient HTTPClient,
	finnhubToken string,
	log *logger.Logger,
) CompanyNewsService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &CompanyNewsServiceImpl{
		relationRepo: relationRepo,
		companyRepo:  companyRepo,
		httpClient:   httpClient,
		finnhubToken: finnhubToken,
		log:          log,
	}
}

func (s *CompanyNewsServiceImpl) Start(ctx context.Context) {
	go s.runScheduler(ctx)
}

func (s *CompanyNewsServiceImpl) runScheduler(ctx context.Context) {
	for {
		wait := nextRunDuration(20, 0, thailandLocation())
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

func (s *CompanyNewsServiceImpl) runOnce(ctx context.Context) {
	symbols, err := s.relationRepo.ListDistinctRelationSymbols()
	if err != nil || len(symbols) == 0 {
		return
	}

	today := time.Now().In(thailandLocation()).Format("2006-01-02")
	if s.log != nil {
		s.log.Infof("company_news fetch started: symbols=%d date=%s", len(symbols), today)
	}
	totalSaved := 0
	for _, symbol := range symbols {
		select {
		case <-ctx.Done():
			if s.log != nil {
				s.log.Infof("company_news fetch canceled")
			}
			return
		default:
		}
		symbol = strings.TrimSpace(symbol)
		if symbol == "" {
			continue
		}
		news, err := s.fetchCompanyNews(symbol, today, today)
		if err != nil || len(news) == 0 {
			continue
		}
		items := make([]models.CompanyNews, 0, len(news))
		for _, n := range news {
			items = append(items, models.CompanyNews{
				Symbol:   n.Related,
				Headline: n.Headline,
				Source:   n.Source,
				Summary:  n.Summary,
				Url:      n.URL,
			})
		}
		if err := s.companyRepo.CreateMany(items); err == nil {
			totalSaved += len(items)
		}
	}
	if s.log != nil {
		s.log.Infof("company_news fetch completed: saved=%d date=%s", totalSaved, today)
	}
}

type finnhubCompanyNewsResponse struct {
	Headline string `json:"headline"`
	Source   string `json:"source"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
	Related  string `json:"related"`
}

func (s *CompanyNewsServiceImpl) fetchCompanyNews(symbol, fromDate, toDate string) ([]finnhubCompanyNewsResponse, error) {
	reqURL, err := url.Parse(companyNewsURL)
	if err != nil {
		return nil, err
	}
	q := reqURL.Query()
	q.Set("symbol", symbol)
	q.Set("from", fromDate)
	q.Set("to", toDate)
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
		return nil, fmt.Errorf("finnhub company-news request failed: %s", strings.TrimSpace(string(body)))
	}

	var result []finnhubCompanyNewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
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
