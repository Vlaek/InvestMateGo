package storage

import (
	"context"

	"invest-mate/internal/models"
)

type InstrumentFinder interface {
	FindByTicker(ticker string) (models.Instrument, bool)
	FindByFigi(figi string) (models.Instrument, bool)
	FindByIsin(isin string) (models.Instrument, bool)
}

// ---------- Бонды ----------
func (ts *TinkoffStorage) GetBonds(ctx context.Context) ([]models.Bond, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.bonds, nil
}

// ---------- Акции ----------
func (ts *TinkoffStorage) GetShares(ctx context.Context) ([]models.Share, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.shares, nil
}

// ---------- Фонды ----------
func (ts *TinkoffStorage) GetEtfs(ctx context.Context) ([]models.Etf, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.etfs, nil
}

// ---------- Валюта ----------
func (ts *TinkoffStorage) GetCurrencies(ctx context.Context) ([]models.Currency, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.currencies, nil
}

// ---------- Общие методы ----------
func (ts *TinkoffStorage) GetInstrumentByFigiAndType(
	ctx context.Context,
	figi string,
	instrumentType models.InstrumentType,
) (models.Instrument, error) {

	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	switch instrumentType {
	case models.InstrumentTypeBond:
		for i := range ts.bonds {
			if ts.bonds[i].Figi == figi {
				return ts.bonds[i], nil
			}
		}
	case models.InstrumentTypeShare:
		for i := range ts.shares {
			if ts.shares[i].Figi == figi {
				return ts.shares[i], nil
			}
		}
	case models.InstrumentTypeETF:
		for i := range ts.etfs {
			if ts.etfs[i].Figi == figi {
				return ts.etfs[i], nil
			}
		}
	case models.InstrumentTypeCurrency:
		for i := range ts.currencies {
			if ts.currencies[i].Figi == figi {
				return ts.currencies[i], nil
			}
		}
	}

	return nil, nil
}
