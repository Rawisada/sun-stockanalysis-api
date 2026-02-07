package controllers

import (
	"context"
	"net/http"

	"sun-stockanalysis-api/internal/domains/auth"
	"sun-stockanalysis-api/pkg/apierror"
	"sun-stockanalysis-api/pkg/response"
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
	Body   response.ApiResponse[LoginResponseBody]
}

type RegisterResponseBody struct {
	UserID string `json:"user_id"`
}

type RegisterResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[RegisterResponseBody]
}

type RefreshResponse struct {
	Status int `status:"default"`
	Body   response.ApiResponse[LoginResponseBody]
}

func (c *AuthController) Login(ctx context.Context, input *auth.LoginInput) (*LoginResponse, error) {
	_ = ctx

	if input.Body.Email == "" || input.Body.Password == "" {
		return nil, apierror.NewBadRequest("email and password required")
	}

	result, err := c.authService.Login(*input)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			return nil, apierror.NewUnauthorized("invalid email or password")
		}
		return nil, apierror.NewInternalError(err.Error())
	}

	return &LoginResponse{
		Status: http.StatusOK,
		Body: response.Success(LoginResponseBody{
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
		return nil, apierror.NewBadRequest("email and password required")
	}

	result, err := c.authService.Register(*input)
	if err != nil {
		if err == auth.ErrEmailAlreadyExists {
			return nil, apierror.NewBadRequest("email already exists")
		}
		return nil, apierror.NewInternalError(err.Error())
	}

	return &RegisterResponse{
		Status: http.StatusCreated,
		Body: response.Success(RegisterResponseBody{
			UserID: result.UserID,
		}),
	}, nil
}

func (c *AuthController) Refresh(ctx context.Context, input *auth.RefreshInput) (*RefreshResponse, error) {
	_ = ctx

	if input.Body.RefreshToken == "" {
		return nil, apierror.NewBadRequest("refresh_token required")
	}

	result, err := c.authService.Refresh(*input)
	if err != nil {
		if err == auth.ErrInvalidRefreshToken {
			return nil, apierror.NewUnauthorized("invalid refresh token")
		}
		return nil, apierror.NewInternalError(err.Error())
	}

	return &RefreshResponse{
		Status: http.StatusOK,
		Body: response.Success(LoginResponseBody{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    result.ExpiresIn,
		}),
	}, nil
}
