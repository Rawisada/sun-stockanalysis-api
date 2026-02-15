package routes

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"sun-stockanalysis-api/internal/controllers"
)

func RegisterPushSubscriptionRoutes(api huma.API, controllers *controllers.Controllers, middleware func(huma.Context, func(huma.Context))) {
	protected := huma.NewGroup(api, "")
	protected.UseMiddleware(middleware)

	huma.Register(protected, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/push/vapid-public-key",
		Summary: "Get VAPID public key",
		Tags:    v1Tags(),
	}, controllers.PushSubscriptionController.GetVAPIDPublicKey)

	huma.Register(protected, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/push/subscriptions",
		Summary: "Create or update push subscription",
		Tags:    v1Tags(),
	}, controllers.PushSubscriptionController.Upsert)
}
