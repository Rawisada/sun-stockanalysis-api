package controllers

type EmptyRequest struct{}

type EmptyResponse struct {
	Status int `status:"default"`
}

type Response[T any] struct {
	Body T
}

type Controllers struct {
	HealthController       *HealthController
	StockController        *StockController
	StockQuoteController   *StockQuoteController
	StockDailyController   *StockDailyController
	AuthController         *AuthController
	RelationNewsController *RelationNewsController
}

func NewControllers(
	healthController *HealthController,
	stockController *StockController,
	stockQuoteController *StockQuoteController,
	stockDailyController *StockDailyController,
	authController *AuthController,
	relationNewsController *RelationNewsController,
) *Controllers {
	return &Controllers{
		HealthController:       healthController,
		StockController:        stockController,
		StockQuoteController:   stockQuoteController,
		StockDailyController:   stockDailyController,
		AuthController:         authController,
		RelationNewsController: relationNewsController,
	}
}
