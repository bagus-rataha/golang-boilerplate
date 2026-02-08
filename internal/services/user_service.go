package services

import (
	"errors"
	"fiber-api-boilerplate/internal/dto"
	"fiber-api-boilerplate/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *UserService) UpdateProfile(userID uint, input dto.UpdateProfileInput) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = input.Name

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

func (s *UserService) ListUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.List()
	if err != nil {
		return nil, err
	}

	return dto.ToUserResponseList(users), nil
}
