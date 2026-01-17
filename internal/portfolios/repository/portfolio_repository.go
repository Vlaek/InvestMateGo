package repository

import (
	"gorm.io/gorm"
)

type PortfoliosRepository interface {
}

type portfoliosRepository struct {
	db *gorm.DB
}

// Создание нового репозитория
func NewPortfoliosRepository(db *gorm.DB) PortfoliosRepository {
	return &portfoliosRepository{db: db}
}
