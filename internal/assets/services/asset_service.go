package services

import (
	"context"
	"sync"

	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/storage"
	"invest-mate/pkg/services"
)

type AssetService interface {
	GetAssets(ctx context.Context, page, limit int) ([]domain.Asset, int64, error)
	GetBonds(ctx context.Context, page, limit int) ([]domain.Bond, int64, error)
	GetShares(ctx context.Context, page, limit int) ([]domain.Share, int64, error)
	GetEtfs(ctx context.Context, page, limit int) ([]domain.Etf, int64, error)
	GetCurrencies(ctx context.Context, page, limit int) ([]domain.Currency, int64, error)
}

type assetService struct {
	repo           repository.AssetRepository
	tinkoffStorage *storage.TinkoffStorage
	mu             sync.RWMutex
}

// Создание нового сервиса
func NewAssetService(repo repository.AssetRepository, tinkoffStorage *storage.TinkoffStorage) AssetService {
	return &assetService{
		repo:           repo,
		tinkoffStorage: tinkoffStorage,
	}
}

// Получение инструментов
func (s *assetService) GetAssets(ctx context.Context, page, limit int) ([]domain.Asset, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetAssets, page, limit)
}

// Получение облигаций
func (s *assetService) GetBonds(ctx context.Context, page, limit int) ([]domain.Bond, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetBonds, page, limit)
}

// Получение акций
func (s *assetService) GetShares(ctx context.Context, page, limit int) ([]domain.Share, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetShares, page, limit)
}

// Получение фондов
func (s *assetService) GetEtfs(ctx context.Context, page, limit int) ([]domain.Etf, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetEtfs, page, limit)
}

// Получение валют
func (s *assetService) GetCurrencies(ctx context.Context, page, limit int) ([]domain.Currency, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetCurrencies, page, limit)
}
