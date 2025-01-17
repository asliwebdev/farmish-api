package services

import (
	"errors"
	"farmish/internal/models"
	"farmish/internal/repository"
	"farmish/pkg/utils"
	"fmt"

	"github.com/google/uuid"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *UserService) SignUp(user *models.User) (string, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = hashedPassword
	user.ID = uuid.New()
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return "", err
	}

	token, err := utils.CreateToken(user.Email, user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) Login(credentials *models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.UserRepo.GetUserByEmail(credentials.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}

	if user == nil {
		return models.LoginResponse{}, ErrInvalidCredentials
	}

	if err = utils.ComparePasswords(user.Password, credentials.Password); err != nil {
		return models.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := utils.CreateToken(user.Email, user.ID)
	if err != nil {
		return models.LoginResponse{}, fmt.Errorf("failed to generate JWT token: %v", err)
	}

	resp := models.LoginResponse{
		Token: token,
		ID:    user.ID,
	}

	return resp, nil
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.UserRepo.GetAllUsers()
}

func (s *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.UserRepo.GetUserByID(userID)
}

func (s *UserService) UpdateUser(user *models.UpdateUser) error {
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}
	return s.UserRepo.UpdateUser(user)
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	return s.UserRepo.DeleteUser(userID)
}
