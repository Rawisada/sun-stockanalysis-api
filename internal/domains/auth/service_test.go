package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/models"
	repositorymock "sun-stockanalysis-api/internal/mocks/repository"
)

type AuthServiceSuite struct {
	suite.Suite
	userRepo    *repositorymock.MockUserRepository
	refreshRepo *repositorymock.MockRefreshTokenRepository
	state       *configurations.State
	service     AuthService
}

func (s *AuthServiceSuite) SetupTest() {
	s.userRepo = repositorymock.NewMockUserRepository(s.T())
	s.refreshRepo = repositorymock.NewMockRefreshTokenRepository(s.T())
	s.state = &configurations.State{
		Secret:     "test-secret",
		ExpiredsAt: 15 * time.Minute,
		Issuer:     "test-issuer",
	}
	s.service = NewAuthService(s.userRepo, s.refreshRepo, s.state)
}

func (s *AuthServiceSuite) TestRegister_RequiresEmailAndPassword() {
	input := RegisterInput{}
	input.Body.Email = ""
	input.Body.Password = ""

	result, err := s.service.Register(input)

	s.Nil(result)
	s.Error(err)
}

func (s *AuthServiceSuite) TestRegister_EmailAlreadyExists() {
	input := RegisterInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	s.userRepo.EXPECT().ExistsByEmail(input.Body.Email).Return(true, nil)

	result, err := s.service.Register(input)

	s.Nil(result)
	s.ErrorIs(err, ErrEmailAlreadyExists)
}

func (s *AuthServiceSuite) TestRegister_Success() {
	input := RegisterInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"
	input.Body.FirstName = "Jane"
	input.Body.LastName = "Doe"

	s.userRepo.EXPECT().ExistsByEmail(input.Body.Email).Return(false, nil)
	s.userRepo.EXPECT().Create(mock.MatchedBy(func(user *models.User) bool {
		s.NotNil(user)
		return user.Email == input.Body.Email &&
			user.FirstName == input.Body.FirstName &&
			user.LastName == input.Body.LastName &&
			user.Role == "USER" &&
			user.Password != "" &&
			user.Password != input.Body.Password
	})).Return(nil)

	result, err := s.service.Register(input)

	s.NoError(err)
	s.NotNil(result)
	s.NotEmpty(result.UserID)
}

func (s *AuthServiceSuite) TestLogin_InvalidSecret() {
	service := NewAuthService(s.userRepo, s.refreshRepo, &configurations.State{})

	input := LoginInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	result, err := service.Login(input)

	s.Nil(result)
	s.Error(err)
}

func (s *AuthServiceSuite) TestLogin_InvalidCredentials() {
	input := LoginInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "wrong"

	s.userRepo.EXPECT().FindByEmail(input.Body.Email).Return((*models.User)(nil), errors.New("not found"))

	result, err := s.service.Login(input)

	s.Nil(result)
	s.ErrorIs(err, ErrInvalidCredentials)
}

func (s *AuthServiceSuite) TestLogin_Success() {
	input := LoginInput{}
	input.Body.Email = "user@example.com"
	input.Body.Password = "secret"

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
	s.Require().NoError(err)

	userID := uuid.New()
	user := &models.User{
		ID:       userID,
		Email:    input.Body.Email,
		Password: string(hashed),
		Role:     "USER",
	}

	s.userRepo.EXPECT().FindByEmail(input.Body.Email).Return(user, nil)
	s.refreshRepo.EXPECT().Create(mock.Anything).Return(nil)
	s.userRepo.EXPECT().UpdateLastLogin(userID, mock.Anything).Return(nil)

	result, err := s.service.Login(input)

	s.NoError(err)
	s.NotNil(result)
	s.NotEmpty(result.AccessToken)
	s.NotEmpty(result.RefreshToken)
	s.Greater(result.ExpiresIn, int64(0))
}

func (s *AuthServiceSuite) TestRefresh_InvalidToken() {
	input := RefreshInput{}
	input.Body.RefreshToken = ""

	result, err := s.service.Refresh(input)

	s.Nil(result)
	s.ErrorIs(err, ErrInvalidRefreshToken)
}

func (s *AuthServiceSuite) TestRefresh_Success() {
	input := RefreshInput{}
	input.Body.RefreshToken = "refresh-token"

	hashBytes := sha256.Sum256([]byte(input.Body.RefreshToken))
	hash := hex.EncodeToString(hashBytes[:])

	userID := uuid.New()
	user := &models.User{
		ID:    userID,
		Email: "user@example.com",
		Role:  "USER",
	}

	s.refreshRepo.EXPECT().FindByHash(hash).Return(&models.RefreshTokens{
		UserID:    userID.String(),
		TokenHash: hash,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		RevokedAt: 0,
	}, nil)
	s.userRepo.EXPECT().FindByID(userID).Return(user, nil)
	s.refreshRepo.EXPECT().RevokeByHash(hash, mock.Anything).Return(nil)
	s.refreshRepo.EXPECT().Create(mock.Anything).Return(nil)

	result, err := s.service.Refresh(input)

	s.NoError(err)
	s.NotNil(result)
	s.NotEmpty(result.AccessToken)
	s.NotEmpty(result.RefreshToken)
	s.Greater(result.ExpiresIn, int64(0))
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}
