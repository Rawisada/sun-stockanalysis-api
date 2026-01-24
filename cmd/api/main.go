package main

import (
	"log"
	"strconv"
	"strings"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/controllers"
	"sun-stockanalysis-api/internal/database"
	"sun-stockanalysis-api/internal/domains/auth"
	"sun-stockanalysis-api/internal/domains/stock"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
	"sun-stockanalysis-api/internal/server"
)

func main() {
	cfg := configurations.ConfigGetting()
	// DB
	db := database.NewPostgresDatabase(cfg.Database).ConnectionGetting()

	// (optional) migrate
	if err := db.AutoMigrate(&models.Stock{}, &models.User{}, &models.RefreshTokens{}); err != nil {
		log.Fatalf("migrate error: %v", err)
	}

	// DI wiring
	stockRepo := repository.NewStockRepository(db)
	stockService := stock.NewStockService(stockRepo)
	stockController := controllers.NewStockController(stockService)
	healthRepo := repository.NewHealthRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, cfg.State)
	authController := controllers.NewAuthController(authService)

	serverConfig := server.ServerConfig{
		Title:            server.ServerTitle("sun-stockanalysis-api"),
		Version:          server.ServerVersion("1.0.0"),
		Port:             toPortString(cfg.Server.Port),
		MaxPayloadSizeKB: parseBodyLimitKB(cfg.Server.BodyLimit),
		TimeoutSeconds:   int(cfg.Server.TimeOut.Seconds()),
		AuthSecret:       cfg.State.Secret,
		AuthIssuer:       cfg.State.Issuer,
	}

	healthController := controllers.NewHealthController(healthRepo, string(serverConfig.Version))
	appControllers := controllers.NewControllers(healthController, stockController, authController)

	// Fiber server
	srv := server.NewServer(serverConfig, appControllers)

	log.Printf("server starting on :%d", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

func parseBodyLimitKB(v string) int {
	value := strings.TrimSpace(strings.ToUpper(v))
	if strings.HasSuffix(value, "M") {
		n := strings.TrimSuffix(value, "M")
		return atoiOrZero(n) * 1024
	}
	if strings.HasSuffix(value, "K") {
		n := strings.TrimSuffix(value, "K")
		return atoiOrZero(n)
	}
	return atoiOrZero(value) * 1024
}

func atoiOrZero(v string) int {
	var n int
	for _, r := range v {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}

func toPortString(port int) string {
	if port <= 0 {
		return "8080"
	}
	return strconv.Itoa(port)
}
