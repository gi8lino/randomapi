package handlers_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexElement(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	t.Run("returns element by index with correct headers", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`{"msg":"first"}`),
			[]byte(`{"msg":"second"}`),
			[]byte(`{"msg":"third"}`),
		}

		req := httptest.NewRequest(http.MethodGet, "/index/1", nil)
		req.SetPathValue("nr", "1")
		w := httptest.NewRecorder()

		handler := handlers.IndexElement(elements, logger)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close() // nolint:errcheck

		require.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.Equal(t, `{"msg":"second"}`, w.Body.String())
	})

	t.Run("returns 400 for invalid index", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`{"msg":"first"}`),
		}

		req := httptest.NewRequest(http.MethodGet, "/index/nope", nil)
		req.SetPathValue("nr", "nope")
		w := httptest.NewRecorder()

		handler := handlers.IndexElement(elements, logger)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "invalid index\n", w.Body.String())
	})

	t.Run("returns 404 for out of range index", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`{"msg":"first"}`),
		}

		req := httptest.NewRequest(http.MethodGet, "/index/3", nil)
		req.SetPathValue("nr", "3")
		w := httptest.NewRecorder()

		handler := handlers.IndexElement(elements, logger)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "index out of range\n", w.Body.String())
	})

	t.Run("returns 500 when no elements available", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{}

		req := httptest.NewRequest(http.MethodGet, "/index/0", nil)
		w := httptest.NewRecorder()

		handler := handlers.IndexElement(elements, logger)
		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "no elements available\n", w.Body.String())
	})
}
