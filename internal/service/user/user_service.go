package user

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"vk-backend/internal/domain"
	"vk-backend/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, name string, password string) (*domain.User, error)
	GetUserByName(ctx context.Context, name string) (*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, name string, password string) (*domain.User, error) {
	if name == "" {
		return nil, domain.ErrEmptyName
	}
	if password == "" {
		return nil, domain.ErrEmptyPassword

	}

	ok, err := s.repo.UserExists(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("user service can't check if user exists: %w", err)
	}
	if ok {
		return nil, domain.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("user service can't hash password: %w", err)
	}

	user, err := s.repo.AddUser(ctx, name, string(hash))
	if err != nil {
		return nil, fmt.Errorf("user service can't add user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	if name == "" {
		return nil, domain.ErrEmptyName
	}
	ok, err := s.repo.UserExists(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("user service can't check if user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrUserNotExists
	}

	user, err := s.repo.GetUserByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("user service can't get user by name: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrUserNotExists
	}
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user service can't get user by id: %w", err)
	}

	return user, nil
}
