package controllers

type EmptyRequest struct{}

type EmptyResponse struct {
	Status int `status:"default"`
}

type Response[T any] struct {
	Body T
}

type Controllers struct {
	HealthController           *HealthController
	StockController            *StockController
	StockQuoteController       *StockQuoteController
	StockDailyController       *StockDailyController
	CompanyNewsController      *CompanyNewsController
	AuthController             *AuthController
	RelationNewsController     *RelationNewsController
	PushSubscriptionController *PushSubscriptionController
}

func NewControllers(
	healthController *HealthController,
	stockController *StockController,
	stockQuoteController *StockQuoteController,
	stockDailyController *StockDailyController,
	companyNewsController *CompanyNewsController,
	authController *AuthController,
	relationNewsController *RelationNewsController,
	pushSubscriptionController *PushSubscriptionController,
) *Controllers {
	return &Controllers{
		HealthController:           healthController,
		StockController:            stockController,
		StockQuoteController:       stockQuoteController,
		StockDailyController:       stockDailyController,
		CompanyNewsController:      companyNewsController,
		AuthController:             authController,
		RelationNewsController:     relationNewsController,
		PushSubscriptionController: pushSubscriptionController,
	}
}
