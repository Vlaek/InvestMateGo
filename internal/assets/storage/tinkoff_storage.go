package storage

import (
	"context"
	"sync"

	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/repository"
)

type TinkoffStorage struct {
	mu sync.RWMutex

	bonds      []domain.Bond
	shares     []domain.Share
	etfs       []domain.Etf
	currencies []domain.Currency

	initialized bool
	initOnce    sync.Once

	repo repository.AssetRepository
}

var (
	instance *TinkoffStorage
	once     sync.Once
)

func NewTinkoffStorage(repo repository.AssetRepository) *TinkoffStorage {
	return &TinkoffStorage{repo: repo}
}

// Получение облигаций из хранилища
func (ts *TinkoffStorage) GetBonds(ctx context.Context) ([]domain.Bond, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.bonds, nil
}

// Получение акций из хранилища
func (ts *TinkoffStorage) GetShares(ctx context.Context) ([]domain.Share, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.shares, nil
}

// Получение фондов из хранилища
func (ts *TinkoffStorage) GetEtfs(ctx context.Context) ([]domain.Etf, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.etfs, nil
}

// Получение валют из хранилища
func (ts *TinkoffStorage) GetCurrencies(ctx context.Context) ([]domain.Currency, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.currencies, nil
}
