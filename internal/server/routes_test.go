package server_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))

	t.Run("GET /healthz", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{} // not used by health handler
		router := server.NewRouter(logger, "", elements)

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "ok", rec.Body.String())
	})

	t.Run("POST /healthz", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{}
		router := server.NewRouter(logger, "", elements)

		req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close() // nolint:errcheck

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "ok", rec.Body.String())
	})

	t.Run("GET /random returns one of provided elements", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`{"msg":"first"}`),
			[]byte(`{"msg":"second"}`),
		}

		router := server.NewRouter(logger, "", elements)

		req := httptest.NewRequest(http.MethodGet, "/random", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close() // nolint:errcheck

		require.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		body := strings.TrimSpace(rec.Body.String())

		validBodies := map[string]bool{
			`{"msg":"first"}`:  true,
			`{"msg":"second"}`: true,
		}
		assert.True(t, validBodies[body], "unexpected body: %q", body)
	})

	t.Run("GET /index/{nr} returns element at index", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`{"msg":"first"}`),
			[]byte(`{"msg":"second"}`),
		}

		router := server.NewRouter(logger, "", elements)

		req := httptest.NewRequest(http.MethodGet, "/index/1", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close() // nolint:errcheck

		require.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.Equal(t, `{"msg":"second"}`, strings.TrimSpace(rec.Body.String()))
	})

	t.Run("route prefix applied", func(t *testing.T) {
		t.Parallel()

		elements := data.Elements{
			[]byte(`"value"`),
		}

		router := server.NewRouter(logger, "/api", elements)

		t.Run("health under prefix", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/api/healthz", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close() // nolint:errcheck

			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "ok", rec.Body.String())
		})

		t.Run("random under prefix", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/api/random", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close() // nolint:errcheck

			require.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			assert.Equal(t, `"value"`, strings.TrimSpace(rec.Body.String()))
		})

		t.Run("index under prefix", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/api/index/0", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close() // nolint:errcheck

			require.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			assert.Equal(t, `"value"`, strings.TrimSpace(rec.Body.String()))
		})
	})
}
