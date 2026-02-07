package controllers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"sun-stockanalysis-api/internal/domains/auth"
	authmock "sun-stockanalysis-api/internal/mocks/domains/auth"
	"sun-stockanalysis-api/pkg/apierror"
)

type AuthControllerSuite struct {
	suite.Suite
	authService *authmock.MockAuthService
	controller  *AuthController
}

func (s *AuthControllerSuite) SetupTest() {
	s.authService = authmock.NewMockAuthService(s.T())
	s.controller = NewAuthController(s.authService)
}

func (s *AuthControllerSuite) TestLogin_ValidationError() {
	input := &auth.LoginInput{}
	input.Body.Email = ""
	input.Body.Password = ""

	resp, err := s.controller.Login(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeBadRequest, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestLogin_InvalidCredentials() {
	input := &auth.LoginInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "wrong"

	s.authService.EXPECT().Login(*input).Return((*auth.LoginResult)(nil), auth.ErrInvalidCredentials)

	resp, err := s.controller.Login(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeUnauthorized, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestLogin_InternalError() {
	input := &auth.LoginInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	s.authService.EXPECT().Login(*input).Return((*auth.LoginResult)(nil), errors.New("db down"))

	resp, err := s.controller.Login(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeInternalError, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestLogin_Success() {
	input := &auth.LoginInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	s.authService.EXPECT().Login(*input).Return(&auth.LoginResult{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresIn:    3600,
	}, nil)

	resp, err := s.controller.Login(context.Background(), input)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.Status)
	s.Equal("access", resp.Body.Data.AccessToken)
	s.Equal("refresh", resp.Body.Data.RefreshToken)
	s.Equal("Bearer", resp.Body.Data.TokenType)
	s.Equal(int64(3600), resp.Body.Data.ExpiresIn)
}

func (s *AuthControllerSuite) TestRegister_ValidationError() {
	input := &auth.RegisterInput{}
	input.Body.Email = ""
	input.Body.Password = ""

	resp, err := s.controller.Register(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeBadRequest, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestRegister_EmailExists() {
	input := &auth.RegisterInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	s.authService.EXPECT().Register(*input).Return((*auth.RegisterResult)(nil), auth.ErrEmailAlreadyExists)

	resp, err := s.controller.Register(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeBadRequest, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestRegister_InternalError() {
	input := &auth.RegisterInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	s.authService.EXPECT().Register(*input).Return((*auth.RegisterResult)(nil), errors.New("db down"))

	resp, err := s.controller.Register(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeInternalError, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestRegister_Success() {
	input := &auth.RegisterInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"
	input.Body.FirstName = "Jane"
	input.Body.LastName = "Doe"

	s.authService.EXPECT().Register(*input).Return(&auth.RegisterResult{
		UserID: "user-id",
	}, nil)

	resp, err := s.controller.Register(context.Background(), input)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusCreated, resp.Status)
	s.Equal("user-id", resp.Body.Data.UserID)
}

func (s *AuthControllerSuite) TestRefresh_ValidationError() {
	input := &auth.RefreshInput{}
	input.Body.RefreshToken = ""

	resp, err := s.controller.Refresh(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeBadRequest, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestRefresh_InvalidToken() {
	input := &auth.RefreshInput{}
	input.Body.RefreshToken = "bad-token"

	s.authService.EXPECT().Refresh(*input).Return((*auth.LoginResult)(nil), auth.ErrInvalidRefreshToken)

	resp, err := s.controller.Refresh(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeUnauthorized, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestRefresh_InternalError() {
	input := &auth.RefreshInput{}
	input.Body.RefreshToken = "token"

	s.authService.EXPECT().Refresh(*input).Return((*auth.LoginResult)(nil), errors.New("db down"))

	resp, err := s.controller.Refresh(context.Background(), input)

	s.Nil(resp)
	s.Error(err)
	s.Equal(apierror.ErrCodeInternalError, err.(*apierror.APIError).Code)
}

func (s *AuthControllerSuite) TestRefresh_Success() {
	input := &auth.RefreshInput{}
	input.Body.RefreshToken = "token"

	s.authService.EXPECT().Refresh(*input).Return(&auth.LoginResult{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresIn:    7200,
	}, nil)

	resp, err := s.controller.Refresh(context.Background(), input)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.Status)
	s.Equal("access", resp.Body.Data.AccessToken)
	s.Equal("refresh", resp.Body.Data.RefreshToken)
	s.Equal("Bearer", resp.Body.Data.TokenType)
	s.Equal(int64(7200), resp.Body.Data.ExpiresIn)
}

func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}
