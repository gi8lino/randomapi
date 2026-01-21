package handler_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomElement(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	t.Run("returns random element with correct headers", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`{"msg":"first"}`),
			[]byte(`{"msg":"second"}`),
			[]byte(`{"msg":"third"}`),
		}
		store := data.NewElementsStore(elements)

		req := httptest.NewRequest(http.MethodGet, "/random", nil)
		w := httptest.NewRecorder()
		h := handler.RandomElement(store, logger)
		h.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close() // nolint:errcheck

		require.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		body := w.Body.String()

		// Response must be exactly one of the provided elements.
		validBodies := map[string]bool{
			`{"msg":"first"}`:  true,
			`{"msg":"second"}`: true,
			`{"msg":"third"}`:  true,
		}
		assert.True(t, validBodies[body], "unexpected body: %q", body)
	})

	t.Run("returns 500 when no elements available", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{} // empty
		store := data.NewElementsStore(elements)

		req := httptest.NewRequest(http.MethodGet, "/random", nil)
		w := httptest.NewRecorder()

		h := handler.RandomElement(store, logger)
		h.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		// http.Error adds a trailing newline
		assert.Equal(t, "no elements available\n", w.Body.String())
	})
}
