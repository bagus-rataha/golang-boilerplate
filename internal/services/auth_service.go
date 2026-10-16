package services

import (
	"errors"
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/dto"
	"fiber-api-boilerplate/internal/models"
	"fiber-api-boilerplate/internal/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepository is the data-access contract AuthService depends on.
// Defined consumer-side so it can be satisfied by the real repository
// (structural typing) or by a mock in tests.
type userRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	List() ([]models.User, error)
}

// refreshTokenRepository is the refresh-token data-access contract used by AuthService.
type refreshTokenRepository interface {
	Create(rt *models.RefreshToken) error
	DeleteByToken(token string) (int64, error)
	DeleteAllByUserID(userID uuid.UUID) error
}

// ErrInvalidRefreshToken covers every refresh token that cannot be exchanged:
// malformed, expired, signed with the wrong secret, already used or revoked,
// or belonging to a user that no longer exists.
var ErrInvalidRefreshToken = errors.New("invalid or revoked refresh token")

type AuthService struct {
	userRepo         userRepository
	refreshTokenRepo refreshTokenRepository
	config           *config.Config
}

func NewAuthService(
	userRepo userRepository,
	refreshTokenRepo refreshTokenRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           cfg,
	}
}

// generateAndStoreTokens issues access + refresh tokens and persists the
// refresh token row so it can be rotated and revoked.
func (s *AuthService) generateAndStoreTokens(user *models.User) (*dto.TokenResponse, error) {
	accessToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTAccessSecret, s.config.JWTAccessExpire,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTRefreshSecret, s.config.JWTRefreshExpire,
	)
	if err != nil {
		return nil, err
	}

	rt := &models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(s.config.JWTRefreshExpire),
	}
	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.ToUserResponse(user),
	}, nil
}

func (s *AuthService) Register(input dto.RegisterInput) (*dto.TokenResponse, error) {
	if _, err := s.userRepo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.generateAndStoreTokens(user)
}

func (s *AuthService) Login(input dto.LoginInput) (*dto.TokenResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !utils.VerifyPassword(input.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAndStoreTokens(user)
}

// RefreshToken exchanges a refresh token for a new pair. Deleting the stored
// row is the consumption step: only the caller whose delete removed it may
// continue, so concurrent requests with the same token cannot all succeed.
func (s *AuthService) RefreshToken(refreshToken string) (*dto.TokenResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, s.config.JWTRefreshSecret)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	deleted, err := s.refreshTokenRepo.DeleteByToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("consume refresh token: %w", err)
	}
	if deleted == 0 {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("load user for refresh: %w", err)
	}

	return s.generateAndStoreTokens(user)
}

// Logout revokes a single refresh token. Idempotent — it does not error if
// the token was already absent from the store.
func (s *AuthService) Logout(refreshToken string) error {
	_, err := s.refreshTokenRepo.DeleteByToken(refreshToken)
	return err
}

// LogoutAll revokes every refresh token belonging to the user.
func (s *AuthService) LogoutAll(userID uuid.UUID) error {
	return s.refreshTokenRepo.DeleteAllByUserID(userID)
}
