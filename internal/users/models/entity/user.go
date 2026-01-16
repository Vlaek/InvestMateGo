package entity

import (
	"invest-mate/internal/shared/models"
	"time"
)

type User struct {
	ID           string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string          `gorm:"uniqueIndex;not null;size:255"`
	Username     string          `gorm:"size:50"`
	PasswordHash string          `gorm:"not null;size:255"`
	Role         models.UserRole `gorm:"not null;size:10"`
	CreatedAt    time.Time       `gorm:"autoCreateTime"`
	UpdatedAt    time.Time       `gorm:"autoUpdateTime"`
}
