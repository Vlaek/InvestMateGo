package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"invest-mate/internal/api"
	"invest-mate/pkg/logger"
)

func (ts *TinkoffStorage) Initialize(ctx context.Context) error {
	var initErr error

	ts.initOnce.Do(func() {
		logger.InfoLog("Initializing Tinkoff storage...")
		start := time.Now()

		var wg sync.WaitGroup
		var mu sync.Mutex
		var initErrors []error

		addError := func(err error) {
			mu.Lock()
			initErrors = append(initErrors, err)
			mu.Unlock()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			bonds, err := api.GetBonds(ctx)
			if err != nil {
				logger.ErrorLog("Failed to get bonds", "error", err)
				addError(err)
				return
			}
			ts.mu.Lock()
			ts.bonds = bonds
			ts.mu.Unlock()
		}()

		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			shares, err := api.GetShares(ctx)
			if err != nil {
				logger.ErrorLog("Failed to get shares", "error", err)
				addError(err)
				return
			}
			ts.mu.Lock()
			ts.shares = shares
			ts.mu.Unlock()
		}()

		wg.Wait()

		if len(initErrors) > 0 {
			initErr = fmt.Errorf("initialization failed with %d errors", len(initErrors))
			return
		}

		ts.mu.Lock()
		ts.initialized = true
		ts.mu.Unlock()

		logger.InfoLog("Tinkoff storage initialized",
			"duration", time.Since(start),
			"bonds", len(ts.bonds),
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
