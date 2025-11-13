package app_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gi8lino/randomapi/internal/app"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("Success (minimal config, ephemeral port)", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
		defer cancel()

		tmp := t.TempDir()
		dataPath := filepath.Join(tmp, "data.json")

		// Minimal valid JSON array for LoadElements.
		require.NoError(t, os.WriteFile(dataPath, []byte(`[1, 2, 3]`), 0o600))

		args := []string{
			"--data-path=" + dataPath,
			"--listen-address=127.0.0.1:0", // avoid port conflicts
		}

		var out bytes.Buffer
		err := app.Run(ctx, "v1", args, &out)
		require.NoError(t, err)
	})

	t.Run("Help requested prints usage and returns nil", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		var out bytes.Buffer
		err := app.Run(ctx, "v1.2.3", []string{"--help"}, &out)
		require.NoError(t, err)

		// tinyflags prints usage text into the writer via err.Error().
		assert.Contains(t, out.String(), "Usage")
	})

	t.Run("Version requested prints version and returns nil", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		var out bytes.Buffer
		err := app.Run(ctx, "v9.8.7", []string{"--version"}, &out)
		require.NoError(t, err)

		// tinyflags version error string should include the version.
		assert.Contains(t, out.String(), "v9.8.7")
	})

	t.Run("Unknown flag surfaces parsing error wrapped as configuration error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		var out bytes.Buffer
		err := app.Run(ctx, "vX", []string{"--totally-unknown"}, &out)
		require.Error(t, err)

		msg := err.Error()
		assert.Contains(t, msg, "configuration error:")
		assert.Contains(t, msg, "unknown flag: --totally-unknown")
	})

	t.Run("Missing data file surfaces load error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()

		tmp := t.TempDir()
		// Intentionally point to a non-existent file
		dataPath := filepath.Join(tmp, "does-not-exist.json")

		args := []string{
			"--data-path=" + dataPath,
			"--listen-address=127.0.0.1:0",
		}

		var out bytes.Buffer
		err := app.Run(ctx, "v1", args, &out)
		require.Error(t, err)

		msg := err.Error()
		assert.Contains(t, msg, "load elements: read data:")
		assert.Contains(t, msg, dataPath)
	})
}
