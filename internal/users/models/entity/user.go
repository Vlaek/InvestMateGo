package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string         `gorm:"uniqueIndex;not null;size:255"`
	Username     string         `gorm:"uniqueIndex;not null;size:50"`
	PasswordHash string         `gorm:"not null;size:255"`
	FirstName    string         `gorm:"size:100"`
	LastName     string         `gorm:"size:100"`
	IsActive     bool           `gorm:"default:true"`
	IsAdmin      bool           `gorm:"default:false"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
