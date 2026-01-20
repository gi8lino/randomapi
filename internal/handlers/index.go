package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gi8lino/randomapi/internal/data"
)

// IndexElement returns a handler that responds with the JSON element at the
// provided index in the in-memory list.
func IndexElement(
	elements data.Elements,
	logger *slog.Logger,
) http.HandlerFunc {
	if len(elements) == 0 {
		// return a handler that returns an error. Faster than check for len on each request.
		return func(w http.ResponseWriter, r *http.Request) {
			logger.Error("no elements available")
			http.Error(w, "no elements available", http.StatusInternalServerError)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		rawIndex := r.PathValue("nr")
		idx, err := strconv.Atoi(rawIndex)
		if err != nil || idx < 0 {
			logger.Warn("invalid index", "index", rawIndex, "error", err)
			http.Error(w, "invalid index", http.StatusBadRequest)
			return
		}
		if idx >= len(elements) {
			logger.Warn("index out of range", "index", idx, "max", len(elements)-1)
			http.Error(w, "index out of range", http.StatusBadRequest)
			return
		}

		elem := elements[idx]
		logger.Debug("index element", "index", idx, "element", string(elem))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(elem); err != nil {
			logger.Error("write response", "error", err)
		}
	}
}
