package handler

import (
	"log/slog"
	"net/http"
)

// ReloadHandler triggers a config reload.
func ReloadHandler(reload func() error, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := reload(); err != nil {
			logger.Error("reload not successful", "error", err)
			http.Error(w, "reload not successful", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			logger.Error("write response", "error", err)
		}
	})
}
