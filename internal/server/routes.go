package server

import (
	"log/slog"
	"net/http"

	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/handler"
)

// NewRouter creates and wires the HTTP mux with handlers and middleware;
// mounts under routePrefix if provided.
func NewRouter(
	logger *slog.Logger,
	routePrefix string,
	store *data.ElementsStore,
	reloadFn func() error,
) http.Handler {
	mux := http.NewServeMux()

	// Health checks.
	mux.Handle("GET /healthz", handler.Healthz())
	mux.Handle("POST /healthz", handler.Healthz())

	// Random element endpoint.
	mux.Handle("GET /random", handler.RandomElement(store, logger))
	// Element by index endpoint.
	mux.Handle("GET /index/{nr}", handler.IndexElement(store, logger))
	// Data reload endpoint.
	mux.Handle("POST /-/reload", handler.ReloadHandler(reloadFn, logger))

	// Mount the whole app under the prefix if provided.
	var handler http.Handler = mux
	if routePrefix != "" {
		handler = mountUnderPrefix(mux, routePrefix)
	}

	return handler
}
