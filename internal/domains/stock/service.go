package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

type StockService interface {
	GetStock(id uuid.UUID) (*models.Stock, error)
	CreateStock(input CreateStockInput) error
	ListAll() ([]models.Stock, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type StockServiceImpl struct {
	repo         repository.StockRepository
	httpClient   HTTPClient
	finnhubToken string
}

func NewStockService(repo repository.StockRepository, httpClient HTTPClient, finnhubToken string) StockService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &StockServiceImpl{
		repo:         repo,
		httpClient:   httpClient,
		finnhubToken: finnhubToken,
	}
}

func (s *StockServiceImpl) GetStock(id uuid.UUID) (*models.Stock, error) {
	return s.repo.FindByID(id)
}

func (s *StockServiceImpl) CreateStock(input CreateStockInput) error {
	profile, err := s.fetchProfile(input.Body.Symbol)
	if err != nil {
		return err
	}

	assetType, err := s.fetchAssetType(input.Body.Symbol)
	if err != nil {
		return err
	}

	if profile.Exchange != "" {
		if err := s.repo.EnsureMasterExchange(profile.Exchange); err != nil {
			return err
		}
	}
	if profile.FinnhubIndustry != "" {
		if err := s.repo.EnsureMasterSector(profile.FinnhubIndustry); err != nil {
			return err
		}
	}
	if assetType != "" {
		if err := s.repo.EnsureMasterAssetType(assetType); err != nil {
			return err
		}
	}

	symbol := profile.Symbol
	if symbol == "" {
		symbol = input.Body.Symbol
	}

	return s.repo.Create(&models.Stock{
		Symbol:    symbol,
		Name:      profile.Name,
		Sector:    profile.FinnhubIndustry,
		Exchange:  profile.Exchange,
		AssetType: assetType,
		Currency:  profile.Currency,
	})
}

func (s *StockServiceImpl) ListAll() ([]models.Stock, error) {
	return s.repo.FindAll()
}

type finnhubProfileResponse struct {
	Exchange        string `json:"exchange"`
	FinnhubIndustry string `json:"finnhubIndustry"`
	Currency        string `json:"currency"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
}

type finnhubSearchResponse struct {
	Count  int `json:"count"`
	Result []struct {
		Symbol string `json:"symbol"`
		Type   string `json:"type"`
	} `json:"result"`
}

func (s *StockServiceImpl) fetchProfile(symbol string) (*finnhubProfileResponse, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/stock/profile2?symbol=%s", symbol)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
		return nil, fmt.Errorf("finnhub profile request failed: %s", strings.TrimSpace(string(body)))
	}

	var result finnhubProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *StockServiceImpl) fetchAssetType(symbol string) (string, error) {
	url := fmt.Sprintf("https://finnhub.io/api/v1/search?q=%s&exchange=US", symbol)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Finnhub-Token", s.finnhubToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("finnhub search request failed: %s", strings.TrimSpace(string(body)))
	}

	var result finnhubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, item := range result.Result {
		if strings.EqualFold(item.Symbol, symbol) {
			return item.Type, nil
		}
	}
	if len(result.Result) > 0 {
		return result.Result[0].Type, nil
	}
	return "", nil
}
