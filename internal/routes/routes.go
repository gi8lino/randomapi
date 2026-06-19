package routes

import (
	"log/slog"
	"net/http"

	"github.com/containeroo/httpprefix"
	"github.com/gi8lino/randomapi/internal/data"
	"github.com/gi8lino/randomapi/internal/handlers"
)

// NewRouter creates and wires the HTTP mux with handlers and middleware;
// mounts under routePrefix if provided.
func NewRouter(
	logger *slog.Logger,
	routePrefix string,
	elements data.Elements,
) http.Handler {
	root := http.NewServeMux()

	root.Handle("GET /healthz", handlers.Healthz())
	root.Handle("POST /healthz", handlers.Healthz())

	root.Handle("GET /random", handlers.RandomElement(elements, logger))
	root.Handle("GET /index/{nr}", handlers.IndexElement(elements, logger))

	return httpprefix.MountUnderPrefix(root, routePrefix)
}
