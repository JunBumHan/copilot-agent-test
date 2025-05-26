package usecase

import (
	"errors"

	"github.com/JunBumHan/copilot-agent-test/internal/domain"
)

// Common errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	GetByID(id uint) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id uint) error
}

// UserService implements the UserUseCase interface
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new UserService with the given repository
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by their ID
func (s *UserService) GetByID(id uint) (*domain.User, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}
	
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

// GetAll retrieves all users
func (s *UserService) GetAll() ([]*domain.User, error) {
	return s.userRepo.FindAll()
}

// Create creates a new user
func (s *UserService) Create(user *domain.User) error {
	if user.Name == "" || user.Email == "" {
		return ErrInvalidInput
	}
	
	return s.userRepo.Create(user)
}

// Update updates an existing user
func (s *UserService) Update(user *domain.User) error {
	if user.ID == 0 || user.Name == "" || user.Email == "" {
		return ErrInvalidInput
	}
	
	_, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return ErrUserNotFound
	}
	
	return s.userRepo.Update(user)
}

// Delete removes a user by ID
func (s *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidInput
	}
	
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	
	return s.userRepo.Delete(id)
}