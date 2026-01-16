package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"invest-mate/internal/users/models"
	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/repository"
	"invest-mate/pkg/logger"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *domain.RegisterRequest) (*domain.UserResponse, error)
	LoginUser(ctx context.Context, req *domain.LoginRequest) (*domain.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error)
	UpdateUser(ctx context.Context, id string, updates *domain.User) (*domain.UserResponse, error)
	GetListUsers(ctx context.Context, page, limit int) ([]*domain.UserResponse, int64, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// Создание нового сервиса
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Регистрация нового пользователя
func (s *userService) RegisterUser(ctx context.Context, req *domain.RegisterRequest) (*domain.UserResponse, error) {
	if err := ValidateRegisterUserRequest(req); err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Username:  req.Username,
		Role:      models.Default,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.HashPassword(req.Password); err != nil {
		logger.ErrorLog("Failed to hash password: %v", err)
		return nil, errors.New("failed to process password")
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	logger.InfoLog("User registered: %s (%s)", user.Email, user.ID)

	return user.ToResponse(), nil
}

// Авторизация пользователя
func (s *userService) LoginUser(ctx context.Context, req *domain.LoginRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByField(ctx, "email", req.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	logger.InfoLog("User logged in: %s", user.Email)

	return user.ToResponse(), nil
}

// Получение пользователя по идентификатору
func (s *userService) DeleteUser(ctx context.Context, id string) (bool, error) {
	return s.userRepo.Delete(ctx, id)
}

// Получение пользователя по идентификатору
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByField(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// Получение пользователя по почте
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByField(ctx, "email", email)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// Обновление данных пользователя
func (s *userService) UpdateUser(ctx context.Context, id string, updates *domain.User) (*domain.UserResponse, error) {
	if err := ValidateUpdateUserRequest(updates); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByField(ctx, "id", id)
	if err != nil {
		return nil, err
	}

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

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// Получение списка пользователей
func (s *userService) GetListUsers(ctx context.Context, page, limit int) ([]*domain.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	offset := 0
	if limit > 0 {
		offset = (page - 1) * limit
	}

	users, err := s.userRepo.GetList(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, int64(len(users)), nil
}
