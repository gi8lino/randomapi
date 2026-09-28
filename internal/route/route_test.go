package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForRequest(t *testing.T) {
	t.Parallel()

	t.Run("without prefix context", func(t *testing.T) {
		t.Parallel()

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		assert.Equal(t, "/random", ForRequest(request, "/random"))
	})

	for _, prefix := range []string{"", "/api"} {
		prefix := prefix
		t.Run("prefix="+prefix, func(t *testing.T) {
			t.Parallel()

			handler := WithPrefix(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(ForRequest(r, "/index/1")))
			}), prefix)

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

			assert.Equal(t, prefix+"/index/1", response.Body.String())
		})
	}
}

func TestRedirect(t *testing.T) {
	t.Parallel()

	for _, prefix := range []string{"", "/api"} {
		prefix := prefix
		t.Run("prefix="+prefix, func(t *testing.T) {
			t.Parallel()

			handler := WithPrefix(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Redirect(w, r, "/random?format=json", http.StatusSeeOther)
			}), prefix)

			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/", nil))

			assert.Equal(t, http.StatusSeeOther, response.Code)
			assert.Equal(t, prefix+"/random?format=json", response.Header().Get("Location"))
		})
	}
}

func TestWithPrefixRewritesHandlerRedirect(t *testing.T) {
	t.Parallel()

	handler := WithPrefix(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/random", http.StatusPermanentRedirect)
	}), "/api")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusPermanentRedirect, response.Code)
	assert.Equal(t, "/api/random", response.Header().Get("Location"))
}
