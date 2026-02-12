package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/database"
	"sun-stockanalysis-api/internal/domains/alert_events"
	"sun-stockanalysis-api/internal/domains/auth"
	"sun-stockanalysis-api/internal/domains/cleanup"
	"sun-stockanalysis-api/internal/domains/company_news"
	"sun-stockanalysis-api/internal/domains/market_open"
	"sun-stockanalysis-api/internal/domains/relation_news"
	"sun-stockanalysis-api/internal/domains/stock"
	"sun-stockanalysis-api/internal/domains/stock_daily"
	"sun-stockanalysis-api/internal/domains/stock_quotes"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/realtime"
	"sun-stockanalysis-api/internal/repository"
	"sun-stockanalysis-api/internal/server"
	"sun-stockanalysis-api/pkg/logger"
)

func main() {
	if loc, err := time.LoadLocation("Asia/Bangkok"); err == nil {
		time.Local = loc
	}
	cfg := configurations.ConfigGetting()
	logg := logger.NewLogger(getEnvString("LOG_LEVEL", "info"))
	defer logg.Sync()
	// DB
	db := database.NewPostgresDatabase(cfg.Database).ConnectionGetting()

	// (optional) migrate
	if err := db.AutoMigrate(
		&models.Stock{},
		&models.StockQuote{},
		&models.User{},
		&models.RefreshTokens{},
		&models.MasterAssetType{},
		&models.MasterExchange{},
		&models.MasterSector{},
		&models.MarketOpen{},
		&models.StockDaily{},
		&models.RelationNews{},
		&models.CompanyNews{},
		&models.AlertEvent{},
	); err != nil {
		logg.Fatalf("migrate error: %v", err)
	}

	// DI wiring
	stockRepo := repository.NewStockRepository(db)
	stockService := stock.NewStockService(stockRepo, nil, cfg.Finnhub.Token)
	stockController := controllers.NewStockController(stockService)
	stockQuoteRepo := repository.NewStockQuoteRepository(db)
	alertEventRepo := repository.NewAlertEventRepository(db)
	alertHub := realtime.NewAlertHub()
	stockQuoteHub := realtime.NewStockQuoteHub()
	alertEventService := alert_events.NewAlertEventService(stockQuoteRepo, alertEventRepo, alertHub)
	stockQuoteService := stock_quotes.NewStockQuoteService(stockRepo, stockQuoteRepo, alertEventService, stockQuoteHub, nil, cfg.Finnhub.Token)
	stockQuoteController := controllers.NewStockQuoteController(stockQuoteService)
	stockDailyRepo := repository.NewStockDailyRepository(db)
	stockDailyService := stock_daily.NewStockDailyService(stockRepo, stockQuoteRepo, stockDailyRepo)
	stockDailyController := controllers.NewStockDailyController(stockDailyService)
	relationNewsRepo := repository.NewRelationNewsRepository(db)
	relationNewsService := relation_news.NewRelationNewsService(relationNewsRepo)
	companyNewsRepo := repository.NewCompanyNewsRepository(db)
	companyNewsService := company_news.NewCompanyNewsService(relationNewsRepo, companyNewsRepo, nil, cfg.Finnhub.Token, logg)
	companyNewsController := controllers.NewCompanyNewsController(companyNewsService)
	healthRepo := repository.NewHealthRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, cfg.State)
	authController := controllers.NewAuthController(authService)
	relationNewsController := controllers.NewRelationNewsController(relationNewsService)
	marketOpenRepo := repository.NewMarketOpenRepository(db)
	marketOpenService := market_open.NewMarketOpenService(marketOpenRepo, nil, cfg.Finnhub.Token, stockQuoteService, stockDailyService, logg)
	cleanupService := cleanup.NewCleanupService(stockQuoteRepo, companyNewsRepo, 15)
	marketOpenService.Start(context.Background())
	companyNewsService.Start(context.Background())
	cleanupService.Start(context.Background())

	healthController := controllers.NewHealthController(healthRepo, "1.0.0")
	appControllers := controllers.NewControllers(healthController, stockController, stockQuoteController, stockDailyController, companyNewsController, authController, relationNewsController)

	// Fiber server
	srv := server.NewServer(cfg, appControllers, alertHub, stockQuoteHub, logg)

	go func() {
		logg.Infof("server starting on :%d", cfg.Server.Port)
		if err := srv.Start(); err != nil {
			logg.Fatalf("server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	logg.Infof("shutdown signal received: %v", sig)
	if err := srv.Stop(); err != nil {
		logg.Errorf("server shutdown error: %v", err)
	}
	logg.Info("server stopped gracefully")
}

func getEnvString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
