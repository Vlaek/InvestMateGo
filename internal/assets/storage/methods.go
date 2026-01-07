package storage

import (
	"context"

	"invest-mate/internal/assets/models/domain"
)

type InstrumentFinder interface {
	FindByTicker(ticker string) (domain.Instrument, bool)
	FindByFigi(figi string) (domain.Instrument, bool)
	FindByIsin(isin string) (domain.Instrument, bool)
}

func (ts *TinkoffStorage) GetBonds(ctx context.Context) ([]domain.Bond, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.bonds, nil
}

func (ts *TinkoffStorage) GetShares(ctx context.Context) ([]domain.Share, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.shares, nil
}

func (ts *TinkoffStorage) GetEtfs(ctx context.Context) ([]domain.Etf, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.etfs, nil
}

func (ts *TinkoffStorage) GetCurrencies(ctx context.Context) ([]domain.Currency, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.currencies, nil
}
