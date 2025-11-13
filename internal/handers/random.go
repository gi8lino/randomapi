package handlers

import (
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
	elements data.Elements,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(elements) == 0 {
			http.Error(w, "no elements available", http.StatusInternalServerError)
			return
		}

		elem := elements[rnd.Intn(len(elements))]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(elem)
	}
}
