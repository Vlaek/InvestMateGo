package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"invest-mate/internal/users/models"
)

type User struct {
	ID           string          `json:"id"`
	Email        string          `json:"email" validate:"required,email"`
	Username     string          `json:"username" validate:"required,min=3,max=50"`
	PasswordHash string          `json:"-" validate:"required"`
	Role         models.UserRole `json:"role"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	TokenType    string        `json:"tokenType"`
	ExpiresIn    int64         `json:"expiresIn"`
}

type UserResponse struct {
	ID        string          `json:"id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	Role      models.UserRole `json:"role"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type DeleteRequest struct {
	ID string `json:"id" validate:"required"`
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
