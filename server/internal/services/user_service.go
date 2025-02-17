package services

import (
	"github.com/gauravst/real-time-chat/internal/models"
	"github.com/gauravst/real-time-chat/internal/repositories"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	// if user.Name == "" {
	// 	return errors.New("user name cannot be empty")
	// }
	//
	// err := s.userRepo.CreateUser(user)
	// if err != nil {
	// 	return fmt.Errorf("failed to create user: %w", err)
	// }

	return nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(id int) (*models.User, error) {
	// user, err := s.userRepo.GetUserByID(id)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get user: %w", err)
	// }
	//
	return nil, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(user *models.User) error {
	// if user.Name == "" {
	// 	return errors.New("user name cannot be empty")
	// }
	//
	// err := s.userRepo.UpdateUser(user)
	// if err != nil {
	// 	return fmt.Errorf("failed to update user: %w", err)
	// }
	//
	return nil
}

// DeleteUser deletes a user by their ID
func (s *userService) DeleteUser(id int) error {
	// err := s.userRepo.DeleteUser(id)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete user: %w", err)
	// }
	//
	return nil
}
