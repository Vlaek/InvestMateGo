package storage

import (
	"invest-mate/internal/models"
	"sync"
)

type TinkoffStorage struct {
	mu sync.RWMutex

	bonds      []models.Bond
	shares     []models.Share
	etfs       []models.Etf
	currencies []models.Currency

	initialized bool
	initOnce    sync.Once
}

var (
	instance *TinkoffStorage
	once     sync.Once
)

func GetInstance() *TinkoffStorage {
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
