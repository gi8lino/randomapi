package data

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// NewReloadFunc builds a reload function with logging and serialization.
func NewReloadFunc(
	path string,
	logger *slog.Logger,
) func() error {
	return NewReloadFuncWithStore(path, logger, nil)
}

// NewReloadFuncWithStore builds a reload function that updates the provided store.
func NewReloadFuncWithStore(
	path string,
	logger *slog.Logger,
	store *ElementsStore,
) func() error {
	var reloadMu sync.Mutex
	return func() error {
		reloadMu.Lock()
		defer reloadMu.Unlock()

		elements, err := LoadElements(path)
		if err != nil {
			return err
		}
		if len(elements) == 0 {
			return fmt.Errorf("no elements available")
		}
		logger.Info("reloading data", "count", len(elements))
		if store != nil {
			store.Set(elements)
		}

		return nil
	}
}

// WatchReload listens for reload signals and invokes the reload function.
func WatchReload(ctx context.Context, reloadCh <-chan os.Signal, logger *slog.Logger, reloadFn func() error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-reloadCh:
			logger.Info("reload requested")
			if err := reloadFn(); err != nil {
				logger.Error("reload failed", "err", err)
			} else {
				logger.Info("reload completed")
			}
		}
	}
}
