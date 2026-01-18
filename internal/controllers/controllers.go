package controllers

type EmptyRequest struct{}

type EmptyResponse struct {
	Status int `status:"default"`
}

type Response[T any] struct {
	Body T
}

type Controllers struct {
	HealthController *HealthController
	StockController *StockController
}

func NewControllers(
	healthController *HealthController,
	stockController *StockController,
) *Controllers {
	return &Controllers{
		HealthController: healthController,
		StockController: stockController,
	}
}
