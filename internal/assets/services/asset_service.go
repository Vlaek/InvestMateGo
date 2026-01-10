package services

import (
	"context"
	"sync"

	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/storage"
)

type AssetService interface {
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

// Получение облигаций
func (s *assetService) GetBonds(ctx context.Context, page, limit int) ([]domain.Bond, int64, error) {
	return getWithPagination(ctx, s.tinkoffStorage.GetBonds, page, limit)
}

// Получение акций
func (s *assetService) GetShares(ctx context.Context, page, limit int) ([]domain.Share, int64, error) {
	return getWithPagination(ctx, s.tinkoffStorage.GetShares, page, limit)
}

// Получение фондов
func (s *assetService) GetEtfs(ctx context.Context, page, limit int) ([]domain.Etf, int64, error) {
	return getWithPagination(ctx, s.tinkoffStorage.GetEtfs, page, limit)
}

// Получение валют
func (s *assetService) GetCurrencies(ctx context.Context, page, limit int) ([]domain.Currency, int64, error) {
	return getWithPagination(ctx, s.tinkoffStorage.GetCurrencies, page, limit)
}

// Функция для получения данных с пагинацией
func getWithPagination[T any](
	ctx context.Context,
	getFunc func(context.Context) ([]T, error),
	page, limit int,
) ([]T, int64, error) {
	items, err := getFunc(ctx)

	if err != nil {
		return nil, 0, err
	}

	result, total := paginate(items, page, limit)

	return result, total, nil
}

// Функция пагинации
func paginate[T any](items []T, page, limit int) ([]T, int64) {
	if len(items) == 0 {
		return []T{}, 0
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit

	if start >= len(items) {
		return []T{}, int64(len(items))
	}

	end := start + limit

	if end > len(items) {
		end = len(items)
	}

	return items[start:end], int64(len(items))
}
