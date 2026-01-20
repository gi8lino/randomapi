package handler

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/gi8lino/randomapi/internal/data"
)

// rnd is a local pseudo-random number generator seeded once at startup.
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomElement returns a handler that responds with a single random JSON element
// from the provided in-memory list.
func RandomElement(
	store *data.ElementsStore,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		elements := store.Get()
		if len(elements) == 0 {
			logger.Error("no elements available")
			http.Error(w, "no elements available", http.StatusInternalServerError)
			return
		}

		elem := elements[rnd.Intn(len(elements))]
		logger.Debug("random element", "element", string(elem))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(elem); err != nil {
			logger.Error("write response", "error", err)
		}
	}
}
