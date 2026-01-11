package main

import (
	"log"
	"sun-stockanalysis-api/app/http"
	"sun-stockanalysis-api/config"
	"sun-stockanalysis-api/internal/infrastructure/database"
	"sun-stockanalysis-api/internal/repository"
	"sun-stockanalysis-api/internal/usecase/stock"
	"sun-stockanalysis-api/app/http/handler"
	stockDomain "sun-stockanalysis-api/internal/domain/stock"
)

func main() {
	cfg := config.ConfigGetting()

	// DB
	db := database.NewPostgresDatabase(cfg.Database).ConnectionGetting()

	// (optional) migrate
	if err := db.AutoMigrate(&stockDomain.Stock{}); err != nil {
		log.Fatalf("migrate error: %v", err)
	}

	// DI wiring
	stockRepo := repository.NewStockRepository(db)
	stockService := stock.NewStockService(stockRepo)
	stockHandler := handler.NewStockHandlerImpl(stockService)

	// Fiber server
	srv := http.NewServer(cfg, stockHandler)

	log.Printf("server starting on :%d", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
