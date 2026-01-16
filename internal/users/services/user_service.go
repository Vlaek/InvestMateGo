package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"invest-mate/internal/shared/models"
	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/repository"
	"invest-mate/pkg/logger"
)

type UserService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.UserResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error)
	UpdateUser(ctx context.Context, id string, updates *domain.User) (*domain.UserResponse, error)
	ListUsers(ctx context.Context, page, limit int) ([]*domain.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.UserResponse, error) {
	// Валидация
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Создаем доменную модель
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Username:  req.Username,
		Role:      models.Default,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Хешируем пароль
	if err := user.HashPassword(req.Password); err != nil {
		logger.ErrorLog("Failed to hash password: %v", err)
		return nil, errors.New("failed to process password")
	}

	// Сохраняем пользователя
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	logger.InfoLog("User registered: %s (%s)", user.Email, user.ID)

	return user.ToResponse(), nil
}

func (s *userService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.UserResponse, error) {
	// Находим пользователя по email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Проверяем пароль
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	logger.InfoLog("User logged in: %s", user.Email)

	return user.ToResponse(), nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, updates *domain.User) (*domain.UserResponse, error) {
	// Получаем существующего пользователя
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля (кроме пароля)
	if updates.Email != "" {
		user.Email = updates.Email
	}
	if updates.Username != "" {
		user.Username = updates.Username
	}
	if updates.Role != "" {
		user.Role = updates.Role
	}

	user.UpdatedAt = time.Now()

	// Сохраняем изменения
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) ([]*domain.UserResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// Валидация запроса регистрации
func validateRegisterRequest(req *domain.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(req.Username) > 50 {
		return errors.New("username must be at most 50 characters")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}
