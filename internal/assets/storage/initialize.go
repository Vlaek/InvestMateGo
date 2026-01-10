package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"invest-mate/internal/assets/api"
	"invest-mate/internal/assets/mappers/bonds"
	"invest-mate/internal/assets/mappers/currencies"
	"invest-mate/internal/assets/mappers/etfs"
	"invest-mate/internal/assets/mappers/shares"
	"invest-mate/internal/assets/models/domain"
	"invest-mate/internal/assets/models/entity"
	"invest-mate/internal/assets/repository"
	"invest-mate/pkg/logger"
)

// Инициализация хранилища
func (ts *TinkoffStorage) Initialize(ctx context.Context) error {
	var initErr error

	ts.initOnce.Do(func() {
		logger.InfoLog("Initializing Tinkoff storage...")
		start := time.Now()

		if ts.repo != nil {
			logger.InfoLog("Attempting to load from database...")
			loadedFromDB := ts.loadFromDatabase(ctx)

			if loadedFromDB {
				logger.InfoLog("✅ Tinkoff storage initialized from database in %v", time.Since(start))
				return
			}

			logger.InfoLog("Database empty or unavailable, loading from API...")
		}

		loadedFromAPI := ts.loadFromAPI(ctx)

		if !loadedFromAPI {
			initErr = fmt.Errorf("failed to initialize Tinkoff storage from API")
			return
		}

		ts.mu.Lock()
		ts.initialized = true
		ts.mu.Unlock()

		logger.InfoLog("✅ Tinkoff storage initialized: duration=%v, bonds=%d, shares=%d, etfs=%d, currencies=%d, source=%s",
			time.Since(start),
			len(ts.bonds),
			len(ts.shares),
			len(ts.etfs),
			len(ts.currencies),
			"api",
		)
	})

	return initErr
}

// Загрузка данных из БД
func (ts *TinkoffStorage) loadFromDatabase(ctx context.Context) bool {
	repoBonds, bondsErr := ts.repo.GetBonds(ctx, 5000, 0)
	repoShares, sharesErr := ts.repo.GetShares(ctx, 5000, 0)
	repoEtfs, etfsErr := ts.repo.GetEtfs(ctx, 5000, 0)
	repoCurrencies, currenciesErr := ts.repo.GetCurrencies(ctx, 5000, 0)

	hasBonds := bondsErr == nil && len(repoBonds) > 0
	hasShares := sharesErr == nil && len(repoShares) > 0
	hasEtfs := etfsErr == nil && len(repoEtfs) > 0
	hasCurrencies := currenciesErr == nil && len(repoCurrencies) > 0

	if !hasBonds && !hasShares && !hasEtfs && !hasCurrencies {
		logger.InfoLog("No data found in database")
		return false
	}

	ts.mu.Lock()

	defer ts.mu.Unlock()

	if hasBonds {
		ts.bonds = make([]domain.Bond, 0, len(repoBonds))

		for _, entity := range repoBonds {
			ts.bonds = append(ts.bonds, bonds.FromEntityToDomain(entity))
		}
	} else if bondsErr != nil {
		logger.ErrorLog("Failed to load bonds from database: %v", bondsErr)
	}

	if hasShares {
		ts.shares = make([]domain.Share, 0, len(repoShares))

		for _, entity := range repoShares {
			ts.shares = append(ts.shares, shares.FromEntityToDomain(entity))
		}
	} else if sharesErr != nil {
		logger.ErrorLog("Failed to load shares from database: %v", sharesErr)
	}

	if hasEtfs {
		ts.etfs = make([]domain.Etf, 0, len(repoEtfs))

		for _, entity := range repoEtfs {
			ts.etfs = append(ts.etfs, etfs.FromEntityToDomain(entity))
		}
	} else if etfsErr != nil {
		logger.ErrorLog("Failed to load ETFs from database: %v", etfsErr)
	}

	if hasCurrencies {
		ts.currencies = make([]domain.Currency, 0, len(repoCurrencies))

		for _, entity := range repoCurrencies {
			ts.currencies = append(ts.currencies, currencies.FromEntityToDomain(entity))
		}
	} else if currenciesErr != nil {
		logger.ErrorLog("Failed to load currencies from database: %v", currenciesErr)
	}

	logger.InfoLog("Loaded bonds: %d, shares: %d, etfs: %d, currencies: %d",
		len(ts.bonds), len(ts.shares), len(ts.etfs), len(ts.currencies))
	ts.initialized = true

	return true
}

// Загрузка данных из API
func (ts *TinkoffStorage) loadFromAPI(ctx context.Context) bool {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var initErrors []error
	successCount := 0

	addError := func(err error) {
		mu.Lock()
		initErrors = append(initErrors, err)
		mu.Unlock()
	}

	incrementSuccess := func() {
		mu.Lock()
		successCount++
		mu.Unlock()
	}

	wg.Add(4)

	go func() {
		defer wg.Done()
		logger.InfoLog("Loading bonds from API...")
		loaded, err := api.GetBonds(ctx)

		if err != nil {
			addError(fmt.Errorf("failed to load bonds: %w", err))
			return
		}

		if loaded == nil {
			addError(fmt.Errorf("bonds loaded is nil"))
			return
		}

		ts.mu.Lock()
		ts.bonds = loaded
		ts.mu.Unlock()

		if ts.repo != nil && len(loaded) > 0 {
			if dbBonds := bonds.FromDomainToEntitySlice(loaded); dbBonds != nil {
				if err := ts.saveBondsToDB(ctx, dbBonds); err == nil {
					incrementSuccess()
				}
			}
		} else {
			incrementSuccess()
		}
	}()

	go func() {
		defer wg.Done()
		logger.InfoLog("Loading shares from API...")
		loaded, err := api.GetShares(ctx)

		if err != nil {
			addError(fmt.Errorf("failed to load shares: %w", err))
			return
		}

		if loaded == nil {
			addError(fmt.Errorf("shares loaded is nil"))
			return
		}

		ts.mu.Lock()
		ts.shares = loaded
		ts.mu.Unlock()

		if ts.repo != nil && len(loaded) > 0 {
			if dbShares := shares.FromDomainToEntitySlice(loaded); dbShares != nil {
				if err := ts.saveSharesToDB(ctx, dbShares); err == nil {
					incrementSuccess()
				}
			}
		} else {
			incrementSuccess()
		}
	}()

	go func() {
		defer wg.Done()
		logger.InfoLog("Loading ETFs from API...")
		loaded, err := api.GetEtfs(ctx)

		if err != nil {
			addError(fmt.Errorf("failed to load ETFs: %w", err))
			return
		}

		if loaded == nil {
			addError(fmt.Errorf("ETFs loaded is nil"))
			return
		}

		ts.mu.Lock()
		ts.etfs = loaded
		ts.mu.Unlock()

		if ts.repo != nil && len(loaded) > 0 {
			if dbEtfs := etfs.FromDomainToEntitySlice(loaded); dbEtfs != nil {
				if err := ts.saveEtfsToDB(ctx, dbEtfs); err == nil {
					incrementSuccess()
				}
			}
		} else {
			incrementSuccess()
		}
	}()

	go func() {
		defer wg.Done()
		logger.InfoLog("Loading currencies from API...")
		loaded, err := api.GetCurrencies(ctx)

		if err != nil {
			addError(fmt.Errorf("failed to load currencies: %w", err))
			return
		}

		if loaded == nil {
			addError(fmt.Errorf("currencies loaded is nil"))
			return
		}

		ts.mu.Lock()
		ts.currencies = loaded
		ts.mu.Unlock()

		if ts.repo != nil && len(loaded) > 0 {
			if dbCurrencies := currencies.FromDomainToEntitySlice(loaded); dbCurrencies != nil {
				if err := ts.saveCurrenciesToDB(ctx, dbCurrencies); err == nil {
					incrementSuccess()
				}
			}
		} else {
			incrementSuccess()
		}
	}()

	wg.Wait()

	if len(initErrors) > 0 {
		logger.ErrorLog("API initialization errors: %v", initErrors)
	}

	return successCount > 0
}

// Сохранение облигаций в базу данных
func (ts *TinkoffStorage) saveBondsToDB(ctx context.Context, dbBonds []entity.Bond) error {
	return repository.SaveToDB(ctx, ts.repo.GetDB(), dbBonds, "bonds")
}

// Сохранение акций в базу данных
func (ts *TinkoffStorage) saveSharesToDB(ctx context.Context, dbShares []entity.Share) error {
	return repository.SaveToDB(ctx, ts.repo.GetDB(), dbShares, "shares")
}

// Сохранение фондов в базу данных
func (ts *TinkoffStorage) saveEtfsToDB(ctx context.Context, dbEtfs []entity.Etf) error {
	return repository.SaveToDB(ctx, ts.repo.GetDB(), dbEtfs, "etfs")
}

// Сохранение валют в базу данных
func (ts *TinkoffStorage) saveCurrenciesToDB(ctx context.Context, dbCurrencies []entity.Currency) error {
	return repository.SaveToDB(ctx, ts.repo.GetDB(), dbCurrencies, "currencies")
}

// Проверка инициализации хранилища и инициализация, если не инициализировано.
func (ts *TinkoffStorage) EnsureInitialized(ctx context.Context) error {
	ts.mu.RLock()
	initialized := ts.initialized
	ts.mu.RUnlock()

	if !initialized {
		return ts.Initialize(ctx)
	}

	return nil
}
