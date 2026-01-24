package controllers

import (
	"context"
	"net/http"

	common "sun-stockanalysis-api/internal/common"
	"sun-stockanalysis-api/internal/domains/auth"
)

type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type LoginResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginResponse struct {
	Status int `status:"default"`
	Body   common.DataResponse[LoginResponseBody]
}

type RegisterResponseBody struct {
	UserID string `json:"user_id"`
}

type RegisterResponse struct {
	Status int `status:"default"`
	Body   common.DataResponse[RegisterResponseBody]
}

type RefreshResponse struct {
	Status int `status:"default"`
	Body   common.DataResponse[LoginResponseBody]
}

func (c *AuthController) Login(ctx context.Context, input *auth.LoginInput) (*LoginResponse, error) {
	_ = ctx

	if input.Body.Email == "" || input.Body.Password == "" {
		return nil, common.NewBadRequest("email and password required")
	}

	result, err := c.authService.Login(*input)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			return nil, common.NewUnauthorized("invalid email or password")
		}
		return nil, common.NewInternalError(err.Error())
	}

	return &LoginResponse{
		Status: http.StatusOK,
		Body: common.SuccessResponse(LoginResponseBody{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    result.ExpiresIn,
		}),
	}, nil
}

func (c *AuthController) Register(ctx context.Context, input *auth.RegisterInput) (*RegisterResponse, error) {
	_ = ctx

	if input.Body.Email == "" || input.Body.Password == "" {
		return nil, common.NewBadRequest("email and password required")
	}

	result, err := c.authService.Register(*input)
	if err != nil {
		if err == auth.ErrEmailAlreadyExists {
			return nil, common.NewBadRequest("email already exists")
		}
		return nil, common.NewInternalError(err.Error())
	}

	return &RegisterResponse{
		Status: http.StatusCreated,
		Body: common.SuccessResponse(RegisterResponseBody{
			UserID: result.UserID,
		}),
	}, nil
}

func (c *AuthController) Refresh(ctx context.Context, input *auth.RefreshInput) (*RefreshResponse, error) {
	_ = ctx

	if input.Body.RefreshToken == "" {
		return nil, common.NewBadRequest("refresh_token required")
	}

	result, err := c.authService.Refresh(*input)
	if err != nil {
		if err == auth.ErrInvalidRefreshToken {
			return nil, common.NewUnauthorized("invalid refresh token")
		}
		return nil, common.NewInternalError(err.Error())
	}

	return &RefreshResponse{
		Status: http.StatusOK,
		Body: common.SuccessResponse(LoginResponseBody{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    result.ExpiresIn,
		}),
	}, nil
}
