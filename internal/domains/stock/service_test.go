package stock

// import (
// 	"bytes"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"strings"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"

// 	"sun-stockanalysis-api/internal/models"
// 	repositorymock "sun-stockanalysis-api/internal/mocks/repository"
// )

// type StockServiceSuite struct {
// 	suite.Suite
// 	repo       *repositorymock.MockStockRepository
// 	service    StockService
// 	httpClient *mockHTTPClient
// }

// func (s *StockServiceSuite) SetupTest() {
// 	s.repo = repositorymock.NewMockStockRepository(s.T())
// 	s.httpClient = &mockHTTPClient{}
// 	s.service = NewStockService(s.repo, s.httpClient, "token")
// }

// func (s *StockServiceSuite) TestGetStock_ReturnsStock() {
// 	id := uuid.New()
// 	expected := &models.Stock{
// 		ID:     id,
// 		Symbol: "AAPL",
// 		Name:   "Apple Inc.",
// 	}

// 	s.repo.EXPECT().FindByID(id).Return(expected, nil)

// 	result, err := s.service.GetStock(id)

// 	s.NoError(err)
// 	s.Equal(expected, result)
// }

// func (s *StockServiceSuite) TestGetStock_ReturnsError() {
// 	id := uuid.New()
// 	wantErr := errors.New("not found")

// 	s.repo.EXPECT().FindByID(id).Return((*models.Stock)(nil), wantErr)

// 	result, err := s.service.GetStock(id)

// 	s.Nil(result)
// 	s.EqualError(err, wantErr.Error())
// }

// func (s *StockServiceSuite) TestCreateStock_Persists() {
// 	input := CreateStockInput{}
// 	input.Body.Symbol = "TSLA"

// 	s.httpClient.do = func(req *http.Request) (*http.Response, error) {
// 		s.Equal("token", req.Header.Get("X-Finnhub-Token"))
// 		if strings.Contains(req.URL.Path, "/stock/profile2") {
// 			return newHTTPResponse(http.StatusOK, `{"exchange":"NASDAQ","finnhubIndustry":"Automotive","currency":"USD","name":"Tesla, Inc.","symbol":"TSLA"}`), nil
// 		}
// 		if strings.Contains(req.URL.Path, "/search") {
// 			return newHTTPResponse(http.StatusOK, `{"count":1,"result":[{"symbol":"TSLA","type":"Stock"}]}`), nil
// 		}
// 		return newHTTPResponse(http.StatusNotFound, `{}`), nil
// 	}

// 	s.repo.EXPECT().EnsureMasterExchange("NASDAQ").Return(nil)
// 	s.repo.EXPECT().EnsureMasterSector("Automotive").Return(nil)
// 	s.repo.EXPECT().EnsureMasterAssetType("Stock").Return(nil)

// 	s.repo.EXPECT().Create(mock.MatchedBy(func(stock *models.Stock) bool {
// 		s.NotNil(stock)
// 		return stock.Symbol == "TSLA" &&
// 			stock.Name == "Tesla, Inc." &&
// 			stock.Sector == "Automotive" &&
// 			stock.Exchange == "NASDAQ" &&
// 			stock.AssetType == "Stock" &&
// 			stock.Currency == "USD"
// 	})).Return(nil)

// 	err := s.service.CreateStock(input)

// 	s.NoError(err)
// }

// func (s *StockServiceSuite) TestCreateStock_ReturnsError() {
// 	input := CreateStockInput{}
// 	input.Body.Symbol = "NVDA"

// 	s.httpClient.do = func(req *http.Request) (*http.Response, error) {
// 		if strings.Contains(req.URL.Path, "/stock/profile2") {
// 			return newHTTPResponse(http.StatusOK, `{"exchange":"NASDAQ","finnhubIndustry":"Technology","currency":"USD","name":"NVIDIA Corporation","symbol":"NVDA"}`), nil
// 		}
// 		if strings.Contains(req.URL.Path, "/search") {
// 			return newHTTPResponse(http.StatusOK, `{"count":1,"result":[{"symbol":"NVDA","type":"Stock"}]}`), nil
// 		}
// 		return newHTTPResponse(http.StatusNotFound, `{}`), nil
// 	}

// 	wantErr := errors.New("create failed")

// 	s.repo.EXPECT().EnsureMasterExchange("NASDAQ").Return(nil)
// 	s.repo.EXPECT().EnsureMasterSector("Technology").Return(nil)
// 	s.repo.EXPECT().EnsureMasterAssetType("Stock").Return(nil)
// 	s.repo.EXPECT().Create(mock.Anything).Return(wantErr)

// 	err := s.service.CreateStock(input)

// 	s.EqualError(err, wantErr.Error())
// }

// func TestStockServiceSuite(t *testing.T) {
// 	suite.Run(t, new(StockServiceSuite))
// }

// type mockHTTPClient struct {
// 	do func(req *http.Request) (*http.Response, error)
// }

// func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
// 	return m.do(req)
// }

// func newHTTPResponse(status int, body string) *http.Response {
// 	return &http.Response{
// 		StatusCode: status,
// 		Body:       io.NopCloser(bytes.NewBufferString(body)),
// 	}
// }
