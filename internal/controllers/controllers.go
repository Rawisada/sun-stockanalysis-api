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
	AuthController         *AuthController
	RelationNewsController *RelationNewsController
}

func NewControllers(
	healthController *HealthController,
	stockController *StockController,
	authController *AuthController,
	relationNewsController *RelationNewsController,
) *Controllers {
	return &Controllers{
		HealthController:       healthController,
		StockController:        stockController,
		AuthController:         authController,
		RelationNewsController: relationNewsController,
	}
}
