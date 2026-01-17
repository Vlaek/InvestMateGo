package services

import (
	"context"
	"errors"

	"invest-mate/internal/users/models/domain"

	"golang.org/x/crypto/bcrypt"
)

// Валидация запроса регистрации
func ValidateRegisterUserRequest(req *domain.RegisterRequest) error {
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

// Валидация запроса регистрации
func ValidateUpdateUserRequest(req *domain.User) error {
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
	if req.Role == "" {
		return errors.New("role is required")
	}
	if req.Role.IsValid() {
		return errors.New("invalid role")
	}

	return nil
}

// Валидация пароля
func (s *userService) VerifyPassword(ctx context.Context, userID, password string) error {
	user, err := s.userRepo.FindByField(ctx, "id", userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}
