package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"invest-mate/internal/api"
	"invest-mate/internal/models"
	"invest-mate/internal/repository"
	repositoryMappers "invest-mate/internal/repository/mappers"
	"invest-mate/pkg/logger"
)

func (ts *TinkoffStorage) Initialize(ctx context.Context) error {
	var initErr error

	ts.initOnce.Do(func() {
		logger.InfoLog("Initializing Tinkoff storage...")
		start := time.Now()

		if ts.repo != nil {
			logger.InfoLog("Attempting to load from database...")
			if repoBonds, err := ts.repo.GetBonds(ctx, 5000, 0); err == nil && len(repoBonds) > 0 {
				ts.mu.Lock()
				for _, b := range repoBonds {
					// TODO: Добавить маппер с repository на models
					ts.bonds = append(ts.bonds, models.Bond{
						Uid:      b.Uid,
						Ticker:   b.Ticker,
						Name:     b.Name,
						Currency: b.Currency,
					})
				}
				ts.initialized = true
				ts.mu.Unlock()

				logger.InfoLog("✅ Loaded %d bonds from database in %v",
					len(repoBonds), time.Since(start))
				return
			}
			logger.InfoLog("Database empty or unavailable, loading from API...")
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		var initErrors []error

		addError := func(err error) {
			mu.Lock()
			initErrors = append(initErrors, err)
			mu.Unlock()
		}

		wg.Add(4)

		// Bonds
		go func() {
			defer wg.Done()
			loaded, err := api.GetBonds(ctx)

			if err != nil {
				addError(err)
				return
			}

			ts.mu.Lock()
			ts.bonds = loaded
			ts.mu.Unlock()

			if ts.repo != nil && len(loaded) > 0 {
				dbBonds := repositoryMappers.BondToRepositoryMapper(loaded)

				if err := repository.SaveEntities(ctx, ts.repo.DB(), dbBonds); err != nil {
					logger.ErrorLog("Failed to save bonds: %v", err)
				}
			}
		}()

		// Shares
		go func() {
			defer wg.Done()
			loaded, err := api.GetShares(ctx)

			if err != nil {
				addError(err)
				return
			}

			ts.mu.Lock()
			ts.shares = loaded
			ts.mu.Unlock()

			if ts.repo != nil && len(loaded) > 0 {
				dbShares := repositoryMappers.ShareToRepositoryMapper(loaded)

				if err := repository.SaveEntities(ctx, ts.repo.DB(), dbShares); err != nil {
					logger.ErrorLog("Failed to save shares: %v", err)
				}
			}
		}()

		// ETFs
		go func() {
			defer wg.Done()
			loaded, err := api.GetEtfs(ctx)
			if err != nil {
				addError(err)
				return
			}
			ts.mu.Lock()
			ts.etfs = loaded
			ts.mu.Unlock()
		}()

		// Currencies
		go func() {
			defer wg.Done()
			loaded, err := api.GetCurrencies(ctx)
			if err != nil {
				addError(err)
				return
			}
			ts.mu.Lock()
			ts.currencies = loaded
			ts.mu.Unlock()
		}()

		wg.Wait()

		if len(initErrors) > 0 {
			initErr = fmt.Errorf("initialization failed with %d errors", len(initErrors))
			logger.ErrorLog("Storage initialization failed: %v", initErr)
			return
		}

		ts.mu.Lock()
		ts.initialized = true
		ts.mu.Unlock()

		logger.InfoLog("✅ Tinkoff storage initialized",
			"duration", time.Since(start),
			"bonds", len(ts.bonds),
			"shares", len(ts.shares),
			"etfs", len(ts.etfs),
			"currencies", len(ts.currencies),
			"source", "api",
		)
	})

	return initErr
}

func (ts *TinkoffStorage) EnsureInitialized(ctx context.Context) error {
	ts.mu.RLock()
	initialized := ts.initialized
	ts.mu.RUnlock()

	if !initialized {
		return ts.Initialize(ctx)
	}

	return nil
}
