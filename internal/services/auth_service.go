package services

import (
	"errors"
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/models"
	"fiber-api-boilerplate/internal/repository"
	"fiber-api-boilerplate/internal/utils"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

// NewAuthService creates new auth service
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// RegisterInput represents registration request
type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

// LoginInput represents login request
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TokenResponse represents auth token response
type TokenResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	User         models.UserResponse `json:"user"`
}

// Register creates new user account
func (s *AuthService) Register(input RegisterInput) (*TokenResponse, error) {
	// Check if email already exists
	if _, err := s.userRepo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTAccessExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTRefreshExpire)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(input LoginInput) (*TokenResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verify password
	if !utils.VerifyPassword(input.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTAccessExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTRefreshExpire)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user from database
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTAccessExpire)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTRefreshExpire)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user.ToResponse(),
	}, nil
}
