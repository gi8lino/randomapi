package flag_test

import (
	"strings"
	"testing"

	"github.com/gi8lino/randomapi/internal/flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	t.Run("minimal defaults", func(t *testing.T) {
		t.Parallel()

		var out strings.Builder
		cfg, err := flag.ParseArgs("v1.2.3", nil, &out)
		require.NoError(t, err)

		// Defaults from flag.Config / ParseArgs
		assert.Equal(t, "", cfg.RoutePrefix)
		assert.Equal(t, "text", string(cfg.LogFormat))
		assert.Equal(t, ":8080", cfg.ListenAddr)
		assert.Equal(t, "/app/data.json", cfg.DataPath)
	})

	t.Run("route prefix default empty", func(t *testing.T) {
		t.Parallel()

		var out strings.Builder
		cfg, err := flag.ParseArgs("v", []string{}, &out)
		require.NoError(t, err)

		assert.Empty(t, cfg.RoutePrefix)
	})

	t.Run("route-prefix normalized value", func(t *testing.T) {
		t.Parallel()

		// e.g. "random-api/" → "/random-api"
		args := []string{"--route-prefix=random-api/"}

		var out strings.Builder
		cfg, err := flag.ParseArgs("v", args, &out)
		require.NoError(t, err)

		assert.Equal(t, "/random-api", cfg.RoutePrefix)
	})

	t.Run("listen address override", func(t *testing.T) {
		t.Parallel()

		args := []string{"--listen-address=127.0.0.1:9090"}
		var out strings.Builder
		cfg, err := flag.ParseArgs("1.0.0", args, &out)
		require.NoError(t, err)

		assert.Equal(t, "127.0.0.1:9090", cfg.ListenAddr)
	})

	t.Run("log format json", func(t *testing.T) {
		t.Parallel()

		args := []string{"--log-format=json"}
		var out strings.Builder
		cfg, err := flag.ParseArgs("dev", args, &out)
		require.NoError(t, err)

		assert.Equal(t, "json", string(cfg.LogFormat))
	})

	t.Run("log format short flag", func(t *testing.T) {
		t.Parallel()

		args := []string{"-l", "json"}
		var out strings.Builder
		cfg, err := flag.ParseArgs("dev", args, &out)
		require.NoError(t, err)

		assert.Equal(t, "json", string(cfg.LogFormat))
	})

	t.Run("invalid log format", func(t *testing.T) {
		t.Parallel()

		args := []string{"--log-format=xml"} // not in {text,json}
		var out strings.Builder
		_, err := flag.ParseArgs("dev", args, &out)
		require.Error(t, err)
	})

	t.Run("data path flag value", func(t *testing.T) {
		t.Parallel()

		args := []string{"--data-path=/path/to/data.json"}
		var out strings.Builder
		cfg, err := flag.ParseArgs("dev", args, &out)
		require.NoError(t, err)

		assert.Equal(t, "/path/to/data.json", cfg.DataPath)
	})

	t.Run("invalid listen address", func(t *testing.T) {
		t.Parallel()

		args := []string{"--listen-address=not-an-addr"}
		var out strings.Builder
		_, err := flag.ParseArgs("dev", args, &out)
		require.Error(t, err)
	})
}
