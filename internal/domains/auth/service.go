package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"sun-stockanalysis-api/internal/configurations"
	"sun-stockanalysis-api/internal/models"
	"sun-stockanalysis-api/internal/repository"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

const refreshTokenTTL = 7 * 24 * time.Hour

type AuthService interface {
	Login(input LoginInput) (*LoginResult, error)
	Register(input RegisterInput) (*RegisterResult, error)
	Refresh(input RefreshInput) (*LoginResult, error)
}

type AuthServiceImpl struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	stateConfig      *configurations.State
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	stateConfig *configurations.State,
) AuthService {
	return &AuthServiceImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		stateConfig:      stateConfig,
	}
}

func (s *AuthServiceImpl) Login(input LoginInput) (*LoginResult, error) {
	if s.stateConfig == nil || s.stateConfig.Secret == "" {
		return nil, errors.New("auth secret not configured")
	}

	user, err := s.userRepo.FindByEmail(input.Body.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Body.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, expiresAt, err := s.createAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenHash, refreshExpiresAt, err := newRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Create(&models.RefreshTokens{
		UserID:    user.ID.String(),
		TokenHash: refreshTokenHash,
		ExpiresAt: refreshExpiresAt.Format(time.RFC3339),
		RevokedAt: 0,
	}); err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID, time.Now())

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(expiresAt).Seconds()),
	}, nil
}

func (s *AuthServiceImpl) Register(input RegisterInput) (*RegisterResult, error) {
	if input.Body.Email == "" || input.Body.Password == "" {
		return nil, errors.New("email and password required")
	}

	exists, err := s.userRepo.ExistsByEmail(input.Body.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     input.Body.Email,
		Password:  string(hashed),
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
		Role:      "USER",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterResult{UserID: user.ID.String()}, nil
}

func (s *AuthServiceImpl) Refresh(input RefreshInput) (*LoginResult, error) {
	if s.stateConfig == nil || s.stateConfig.Secret == "" {
		return nil, errors.New("auth secret not configured")
	}
	if input.Body.RefreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	hashBytes := sha256.Sum256([]byte(input.Body.RefreshToken))
	hash := hex.EncodeToString(hashBytes[:])

	stored, err := s.refreshTokenRepo.FindByHash(hash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if stored.RevokedAt > 0 {
		return nil, ErrInvalidRefreshToken
	}

	expiresAt, err := time.Parse(time.RFC3339, stored.ExpiresAt)
	if err != nil || time.Now().After(expiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(stored.UserID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, accessExpiresAt, err := s.createAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, newRefreshTokenHash, newRefreshExpiresAt, err := newRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.RevokeByHash(hash, float64(time.Now().Unix())); err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Create(&models.RefreshTokens{
		UserID:    user.ID.String(),
		TokenHash: newRefreshTokenHash,
		ExpiresAt: newRefreshExpiresAt.Format(time.RFC3339),
		RevokedAt: 0,
	}); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(time.Until(accessExpiresAt).Seconds()),
	}, nil
}

func (s *AuthServiceImpl) createAccessToken(user *models.User) (string, time.Time, error) {
	expiry := normalizeDuration(s.stateConfig.ExpiredsAt)
	if expiry <= 0 {
		expiry = 15 * time.Minute
	}
	expiresAt := time.Now().Add(expiry)

	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		Issuer:    s.stateConfig.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   claims.Subject,
		"iss":   claims.Issuer,
		"iat":   claims.IssuedAt.Unix(),
		"exp":   claims.ExpiresAt.Unix(),
		"email": user.Email,
		"role":  user.Role,
	})

	signed, err := token.SignedString([]byte(s.stateConfig.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func normalizeDuration(d time.Duration) time.Duration {
	if d > 0 && d < time.Second {
		return d * time.Second
	}
	return d
}

func newRefreshToken() (string, string, time.Time, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", time.Time{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	hash := sha256.Sum256([]byte(token))
	return token, hex.EncodeToString(hash[:]), time.Now().Add(refreshTokenTTL), nil
}
