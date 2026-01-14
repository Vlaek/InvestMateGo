package services

import (
	"context"
	"sync"

	"invest-mate/internal/assets/mappers/bonds"
	"invest-mate/internal/assets/mappers/currencies"
	"invest-mate/internal/assets/mappers/etfs"
	"invest-mate/internal/assets/mappers/shares"
	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/assets/repository"
	"invest-mate/internal/assets/storage"
	"invest-mate/pkg/services"
)

type AssetService interface {
	GetAssets(ctx context.Context, page, limit int) ([]domain.Asset, int64, error)
	GetAssetByField(ctx context.Context, fieldName string, fieldValue string) (*entity.AssetInstrument, error)

	GetBonds(ctx context.Context, page, limit int) ([]domain.Bond, int64, error)
	GetBondByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Bond, error)

	GetShares(ctx context.Context, page, limit int) ([]domain.Share, int64, error)
	GetShareByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Share, error)

	GetEtfs(ctx context.Context, page, limit int) ([]domain.Etf, int64, error)
	GetEtfByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Etf, error)

	GetCurrencies(ctx context.Context, page, limit int) ([]domain.Currency, int64, error)
	GetCurrencyByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Currency, error)
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

// Получение инструмента по идентификатору
func (s *assetService) GetAssetByField(ctx context.Context, fieldName string, fieldValue string) (*entity.AssetInstrument, error) {
	entity, err := s.repo.GetAssetByField(ctx, fieldName, fieldValue)

	if err == nil && entity != nil {
		return &entity, nil
	}

	return nil, err
}

// Получение облигаций
func (s *assetService) GetBonds(ctx context.Context, page, limit int) ([]domain.Bond, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetBonds, page, limit)
}

// Получение облигации по идентификатору
func (s *assetService) GetBondByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Bond, error) {
	entity, err := s.repo.GetBondByField(ctx, fieldName, fieldValue)

	if err == nil && entity != nil {
		domain := bonds.FromEntityToDomain(*entity)

		return &domain, nil
	}

	return nil, err
}

// Получение акций
func (s *assetService) GetShares(ctx context.Context, page, limit int) ([]domain.Share, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetShares, page, limit)
}

// Получение акции по идентификатору
func (s *assetService) GetShareByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Share, error) {
	entity, err := s.repo.GetShareByField(ctx, fieldName, fieldValue)

	if err == nil && entity != nil {
		domain := shares.FromEntityToDomain(*entity)

		return &domain, nil
	}

	return nil, err
}

// Получение фондов
func (s *assetService) GetEtfs(ctx context.Context, page, limit int) ([]domain.Etf, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetEtfs, page, limit)
}

// Получение фонда по идентификатору
func (s *assetService) GetEtfByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Etf, error) {
	entity, err := s.repo.GetEtfByField(ctx, fieldName, fieldValue)

	if err == nil {
		domain := etfs.FromEntityToDomain(*entity)

		return &domain, nil
	}

	return nil, err
}

// Получение валют
func (s *assetService) GetCurrencies(ctx context.Context, page, limit int) ([]domain.Currency, int64, error) {
	return services.GetWithPagination(ctx, s.tinkoffStorage.GetCurrencies, page, limit)
}

// Получение валюты по идентификатору
func (s *assetService) GetCurrencyByField(ctx context.Context, fieldName string, fieldValue string) (*domain.Currency, error) {
	entity, err := s.repo.GetCurrencyByField(ctx, fieldName, fieldValue)

	if err == nil {
		domain := currencies.FromEntityToDomain(*entity)

		return &domain, nil
	}

	return nil, err
}
