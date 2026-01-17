package services

import (
	"invest-mate/internal/portfolios/repository"
)

type PortfoliosService interface {
}

type portfoliosService struct {
	portfoliosRepo repository.PortfoliosRepository
}

// Создание нового сервиса
func NewPortfoliosService(portfoliosRepo repository.PortfoliosRepository) PortfoliosService {
	return &portfoliosService{portfoliosRepo: portfoliosRepo}
}
