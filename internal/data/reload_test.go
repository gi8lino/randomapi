package data

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReloadFuncWithStore(t *testing.T) {
	t.Parallel()

	t.Run("with existing store", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "data.json")
		require.NoError(t, os.WriteFile(path, []byte(`["first","second"]`), 0o600))

		logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
		store := NewElementsStore(Elements{[]byte(`"old"`)})

		reload := NewReloadFuncWithStore(path, logger, store)
		require.NoError(t, reload())

		assert.Equal(t, Elements{[]byte(`"first"`), []byte(`"second"`)}, store.Get())
	})

	t.Run("Empty elements error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "data.json")
		require.NoError(t, os.WriteFile(path, []byte(`[]`), 0o600))

		logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
		store := NewElementsStore(Elements{[]byte(`"old"`)})

		reload := NewReloadFuncWithStore(path, logger, store)
		err := reload()
		require.Error(t, err)
		assert.Equal(t, "no elements available", err.Error())
		assert.Equal(t, Elements{[]byte(`"old"`)}, store.Get())
	})

	t.Run("Read error", func(t *testing.T) {
		t.Parallel()

		logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
		store := NewElementsStore(Elements{[]byte(`"old"`)})

		reload := NewReloadFuncWithStore("/does/not/exist.json", logger, store)
		err := reload()
		require.Error(t, err)
		assert.Equal(t, Elements{[]byte(`"old"`)}, store.Get())
	})
}

func TestWatchReload(t *testing.T) {
	t.Parallel()
	t.Run("Calls reload function", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
		reloadCh := make(chan os.Signal, 1)
		called := make(chan struct{}, 1)

		go WatchReload(ctx, reloadCh, logger, func() error {
			called <- struct{}{}
			return nil
		})

		reloadCh <- syscall.SIGHUP

		select {
		case <-called:
			// ok
		case <-time.After(1 * time.Second):
			t.Fatal("reload function was not called")
		}
	})
}
