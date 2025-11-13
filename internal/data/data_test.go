package data_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadElements(t *testing.T) {
	t.Parallel()

	t.Run("loads valid JSON array", func(t *testing.T) {
		t.Parallel()

		tmp := t.TempDir()
		path := filepath.Join(tmp, "data.json")

		content := `[
			{"msg": "hello"},
			{"msg": "world"},
			42,
			"string value"
		]`

		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

		elements, err := data.LoadElements(path)
		require.NoError(t, err)

		assert.Len(t, elements, 4)

		// And ensure it preserved raw JSON values
		assert.JSONEq(t, `{"msg":"hello"}`, string(elements[0]))
		assert.JSONEq(t, `{"msg":"world"}`, string(elements[1]))
		assert.Equal(t, "42", string(elements[2]))
		assert.Equal(t, `"string value"`, string(elements[3]))
	})

	t.Run("returns error on missing file", func(t *testing.T) {
		t.Parallel()

		_, err := data.LoadElements("/does/not/exist.json")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read data")
	})

	t.Run("returns error on invalid JSON", func(t *testing.T) {
		t.Parallel()

		tmp := t.TempDir()
		path := filepath.Join(tmp, "invalid.json")

		// not an array, not even valid JSON
		require.NoError(t, os.WriteFile(path, []byte(`{ this is not valid json ]`), 0o600))

		_, err := data.LoadElements(path)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal data")
	})

	t.Run("returns error when JSON is valid but not an array", func(t *testing.T) {
		t.Parallel()

		tmp := t.TempDir()
		path := filepath.Join(tmp, "not-array.json")

		require.NoError(t, os.WriteFile(path, []byte(`{"a":1,"b":2}`), 0o600))

		_, err := data.LoadElements(path)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal data")
	})
}
