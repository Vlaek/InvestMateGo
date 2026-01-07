package mappers

import (
	"invest-mate/internal/users/models/domain"
)

func RegisterRequestToDomain(req *domain.RegisterRequest) *domain.User {
	return &domain.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		IsAdmin:   false,
	}
}
