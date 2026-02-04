package services

import (
	"errors"
	"fiber-api-boilerplate/internal/models"
	"fiber-api-boilerplate/internal/repository"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates new user service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// UpdateProfileInput represents update profile request
type UpdateProfileInput struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

// GetProfile returns user profile
func (s *UserService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(userID uint, input UpdateProfileInput) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	user.Name = input.Name

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ListUsers returns all users
func (s *UserService) ListUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.List()
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}
