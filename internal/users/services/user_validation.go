package services

import (
	"errors"

	"invest-mate/internal/users/models/domain"
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
