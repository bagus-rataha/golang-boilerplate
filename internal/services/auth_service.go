package services

import (
	"errors"
	"fiber-api-boilerplate/internal/config"
	"fiber-api-boilerplate/internal/dto"
	"fiber-api-boilerplate/internal/models"
	"fiber-api-boilerplate/internal/repository"
	"fiber-api-boilerplate/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
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

	accessToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTSecret, s.config.JWTAccessExpire,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTSecret, s.config.JWTRefreshExpire,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.ToUserResponse(user),
	}, nil
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

	accessToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTSecret, s.config.JWTAccessExpire,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTSecret, s.config.JWTRefreshExpire,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.ToUserResponse(user),
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*dto.TokenResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTSecret, s.config.JWTAccessExpire,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateToken(
		user.ID, user.Email, user.Role,
		s.config.JWTSecret, s.config.JWTRefreshExpire,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         dto.ToUserResponse(user),
	}, nil
}
