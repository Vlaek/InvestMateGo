package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"invest-mate/internal/api"
	"invest-mate/internal/mappers"
	"invest-mate/internal/models/domain"
	"invest-mate/internal/models/entity"
	"invest-mate/internal/repository"
	"invest-mate/pkg/logger"
)

// TODO: Отрефакторить, сделать унифицированное решение. Вынести мапперы в mappers

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
		for _, b := range repoBonds {
			if domainBond := ts.mapBondFromEntity(b); domainBond != nil {
				ts.bonds = append(ts.bonds, *domainBond)
			}
		}

		logger.InfoLog("Loaded %d bonds from database", len(repoBonds))
	} else if bondsErr != nil {
		logger.ErrorLog("Failed to load bonds from database: %v", bondsErr)
	}

	if hasShares {
		for _, s := range repoShares {
			if domainShare := ts.mapShareFromEntity(s); domainShare != nil {
				ts.shares = append(ts.shares, *domainShare)
			}
		}

		logger.InfoLog("Loaded %d shares from database", len(repoShares))
	} else if sharesErr != nil {
		logger.ErrorLog("Failed to load shares from database: %v", sharesErr)
	}

	if hasEtfs {
		for _, e := range repoEtfs {
			if domainEtf := ts.mapEtfFromEntity(e); domainEtf != nil {
				ts.etfs = append(ts.etfs, *domainEtf)
			}
		}

		logger.InfoLog("Loaded %d ETFs from database", len(repoEtfs))
	} else if etfsErr != nil {
		logger.ErrorLog("Failed to load ETFs from database: %v", etfsErr)
	}

	if hasCurrencies {
		for _, c := range repoCurrencies {
			if domainCurrency := ts.mapCurrencyFromEntity(c); domainCurrency != nil {
				ts.currencies = append(ts.currencies, *domainCurrency)
			}
		}

		logger.InfoLog("Loaded %d currencies from database", len(repoCurrencies))
	} else if currenciesErr != nil {
		logger.ErrorLog("Failed to load currencies from database: %v", currenciesErr)
	}

	ts.initialized = true

	return true
}

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
			if dbBonds := mappers.BondToRepositoryMapper(loaded); dbBonds != nil {
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
			if dbShares := mappers.ShareToRepositoryMapper(loaded); dbShares != nil {
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

		if len(loaded) > 0 {
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

		if len(loaded) > 0 {
			incrementSuccess()
		}
	}()

	wg.Wait()

	if len(initErrors) > 0 {
		logger.ErrorLog("API initialization errors: %v", initErrors)
	}

	return successCount > 0
}

func (ts *TinkoffStorage) mapBondFromEntity(b entity.Bond) *domain.Bond {
	return &domain.Bond{
		Uid:               b.Uid,
		Ticker:            b.Ticker,
		Name:              b.Name,
		Currency:          b.Currency,
		Figi:              b.Figi,
		Isin:              b.Isin,
		MinPriceIncrement: b.MinPriceIncrement,
		Lot:               b.Lot,
		Nominal:           b.Nominal,
		CountryOfRisk:     b.CountryOfRisk,
	}
}

func (ts *TinkoffStorage) mapShareFromEntity(s entity.Share) *domain.Share {
	return &domain.Share{
		Uid:               s.Uid,
		Ticker:            s.Ticker,
		Name:              s.Name,
		Currency:          s.Currency,
		Figi:              s.Figi,
		Isin:              s.Isin,
		MinPriceIncrement: s.MinPriceIncrement,
		Lot:               s.Lot,
		Sector:            s.Sector,
	}
}

func (ts *TinkoffStorage) mapEtfFromEntity(e entity.Share) *domain.Etf {
	return &domain.Etf{
		Uid:               e.Uid,
		Ticker:            e.Ticker,
		Name:              e.Name,
		Currency:          e.Currency,
		Figi:              e.Figi,
		Isin:              e.Isin,
		MinPriceIncrement: e.MinPriceIncrement,
		Lot:               e.Lot,
		Sector:            e.Sector,
	}
}

func (ts *TinkoffStorage) mapCurrencyFromEntity(c entity.Share) *domain.Currency {
	return &domain.Currency{
		Uid:               c.Uid,
		Ticker:            c.Ticker,
		Name:              c.Name,
		Currency:          c.Currency,
		Figi:              c.Figi,
		Isin:              c.Isin,
		MinPriceIncrement: c.MinPriceIncrement,
		Lot:               c.Lot,
	}
}

func (ts *TinkoffStorage) saveBondsToDB(ctx context.Context, dbBonds []entity.Bond) error {
	if dbBonds == nil || len(dbBonds) == 0 {
		return nil
	}

	tx := ts.repo.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorLog("Panic during bond save: %v", r)
		}
	}()

	if err := repository.SaveEntities(ctx, tx, dbBonds); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save bonds: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit bond transaction: %w", err)
	}

	logger.InfoLog("Saved %d bonds to database", len(dbBonds))

	return nil
}

func (ts *TinkoffStorage) saveSharesToDB(ctx context.Context, dbShares []entity.Share) error {
	if dbShares == nil || len(dbShares) == 0 {
		return nil
	}

	tx := ts.repo.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.ErrorLog("Panic during share save: %v", r)
		}
	}()

	if err := repository.SaveEntities(ctx, tx, dbShares); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save shares: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit share transaction: %w", err)
	}

	logger.InfoLog("Saved %d shares to database", len(dbShares))

	return nil
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
