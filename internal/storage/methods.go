package storage

import (
	"context"
	"strings"

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

func (ts *TinkoffStorage) GetBondByTicker(ctx context.Context, ticker string) (*models.Bond, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for i := range ts.bonds {
		if strings.EqualFold(ts.bonds[i].Ticker, ticker) {
			return &ts.bonds[i], nil
		}
	}
	return nil, nil
}

func (ts *TinkoffStorage) GetBondByTickerArray(ctx context.Context, tickers []string) ([]models.Bond, error) {
	if err := ts.EnsureInitialized(ctx); err != nil {
		return nil, err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	tickerMap := make(map[string]bool)
	for _, t := range tickers {
		tickerMap[strings.ToUpper(t)] = true
	}

	var result []models.Bond
	for i := range ts.bonds {
		if tickerMap[strings.ToUpper(ts.bonds[i].Ticker)] {
			result = append(result, ts.bonds[i])
		}
	}
	return result, nil
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
	}

	return nil, nil
}
