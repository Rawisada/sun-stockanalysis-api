package main

import (
	"log"
	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/database"
	"sun-stockanalysis-api/internal/repository"
	"sun-stockanalysis-api/internal/domains/stock"
	"sun-stockanalysis-api/internal/server"
	"sun-stockanalysis-api/internal/models"
)

func main() {
	cfg := configurations.ConfigGetting()
	// DB
	db := database.NewPostgresDatabase(cfg.Database).ConnectionGetting()

	// (optional) migrate
	if err := db.AutoMigrate(&models.Stock{}); err != nil {
		log.Fatalf("migrate error: %v", err)
	}

	// DI wiring
	stockRepo := repository.NewStockRepository(db)
	stockService := stock.NewStockService(stockRepo)
	stockController := controllers.NewStockController(stockService)

	// Fiber server
	srv := server.NewServer(cfg, stockController)

	log.Printf("server starting on :%d", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
