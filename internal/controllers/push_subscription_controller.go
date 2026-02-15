package controllers

import (
	"context"
	"net/http"

	"sun-stockanalysis-api/internal/authctx"
	"sun-stockanalysis-api/internal/domains/push_subscriptions"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
)

type PushSubscriptionController struct {
	service push_subscriptions.PushSubscriptionService
}

func NewPushSubscriptionController(service push_subscriptions.PushSubscriptionService) *PushSubscriptionController {
	return &PushSubscriptionController{service: service}
}

type PushSubscriptionUpsertInput struct {
	Body struct {
		DeviceID     string `json:"device_id"`
		UserAgent    string `json:"user_agent"`
		Subscription struct {
			Endpoint string `json:"endpoint"`
			Keys     struct {
				P256DH string `json:"p256dh"`
				Auth   string `json:"auth"`
			} `json:"keys"`
		} `json:"subscription"`
	}
}

type PushSubscriptionUpsertResponseBody struct {
	Saved bool `json:"saved"`
}

type PushSubscriptionUpsertResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[PushSubscriptionUpsertResponseBody]
}

type PushSubscriptionDeleteInput struct {
	DeviceID string `query:"device_id" doc:"Device identifier to unsubscribe"`
}

type PushSubscriptionDeleteResponseBody struct {
	Deleted bool `json:"deleted"`
}

type PushSubscriptionDeleteResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[PushSubscriptionDeleteResponseBody]
}

type VAPIDPublicKeyResponseBody struct {
	PublicKey string `json:"public_key"`
}

type VAPIDPublicKeyResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[VAPIDPublicKeyResponseBody]
}

func (c *PushSubscriptionController) Upsert(ctx context.Context, input *PushSubscriptionUpsertInput) (*PushSubscriptionUpsertResponse, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok {
		return nil, apierror.NewUnauthorized("invalid token context")
	}
	if input == nil {
		return nil, apierror.NewBadRequest("request body required")
	}

	err := c.service.Save(ctx, userID, push_subscriptions.SaveSubscriptionInput{
		DeviceID:  input.Body.DeviceID,
		Endpoint:  input.Body.Subscription.Endpoint,
		P256DHKey: input.Body.Subscription.Keys.P256DH,
		AuthKey:   input.Body.Subscription.Keys.Auth,
		UserAgent: input.Body.UserAgent,
	})
	if err != nil {
		return nil, apierror.NewBadRequest(err.Error())
	}

	return &PushSubscriptionUpsertResponse{
		Status: http.StatusOK,
		Body: response.Success(PushSubscriptionUpsertResponseBody{
			Saved: true,
		}),
	}, nil
}

func (c *PushSubscriptionController) GetVAPIDPublicKey(ctx context.Context, _ *EmptyRequest) (*VAPIDPublicKeyResponse, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok || userID == "" {
		return nil, apierror.NewUnauthorized("invalid token context")
	}
	publicKey, err := c.service.GetPublicKey(ctx)
	if err != nil {
		return nil, apierror.NewInternalError(err.Error())
	}

	return &VAPIDPublicKeyResponse{
		Status: http.StatusOK,
		Body: response.Success(VAPIDPublicKeyResponseBody{
			PublicKey: publicKey,
		}),
	}, nil
}

func (c *PushSubscriptionController) Delete(ctx context.Context, input *PushSubscriptionDeleteInput) (*PushSubscriptionDeleteResponse, error) {
	userID, ok := authctx.UserIDFromContext(ctx)
	if !ok {
		return nil, apierror.NewUnauthorized("invalid token context")
	}
	if input == nil || input.DeviceID == "" {
		return nil, apierror.NewBadRequest("device_id required")
	}

	if err := c.service.Delete(ctx, userID, input.DeviceID); err != nil {
		return nil, apierror.NewBadRequest(err.Error())
	}

	return &PushSubscriptionDeleteResponse{
		Status: http.StatusOK,
		Body: response.Success(PushSubscriptionDeleteResponseBody{
			Deleted: true,
		}),
	}, nil
}
