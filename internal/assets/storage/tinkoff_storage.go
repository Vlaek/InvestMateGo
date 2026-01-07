package storage

import (
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

	repo *repository.PostgresRepository
}

var (
	instance *TinkoffStorage
	once     sync.Once
)

func NewTinkoffStorage(repo *repository.PostgresRepository) *TinkoffStorage {
	return &TinkoffStorage{repo: repo}
}

func GetInstance(repo *repository.PostgresRepository) *TinkoffStorage {
	once.Do(func() {
		instance = NewTinkoffStorage(repo)
	})
	return instance
}

func GetInstanceWithoutRepo() *TinkoffStorage {
	once.Do(func() {
		instance = &TinkoffStorage{}
	})
	return instance
}

func (ts *TinkoffStorage) IsInitialized() bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.initialized
}

func (ts *TinkoffStorage) Clear() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.bonds = nil
	ts.shares = nil
	ts.etfs = nil
	ts.currencies = nil
	ts.initialized = false
}

func (ts *TinkoffStorage) GetStats() map[string]interface{} {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return map[string]interface{}{
		"initialized": ts.initialized,
		"bonds":       len(ts.bonds),
		"shares":      len(ts.shares),
		"etfs":        len(ts.etfs),
		"currencies":  len(ts.currencies),
		"has_db":      ts.repo != nil,
	}
}
