package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElementsStore(t *testing.T) {
	t.Parallel()

	t.Run("Successfull Set and get", func(t *testing.T) {
		t.Parallel()

		initial := Elements{
			[]byte(`{"msg":"first"}`),
		}
		store := NewElementsStore(initial)
		assert.Equal(t, initial, store.Get())

		next := Elements{
			[]byte(`{"msg":"second"}`),
			[]byte(`{"msg":"third"}`),
		}
		store.Set(next)
		assert.Equal(t, next, store.Get())
	})

	t.Run("Nil store", func(t *testing.T) {
		t.Parallel()

		var store *ElementsStore
		assert.Empty(t, store.Get())
	})

	t.Run("Set to a nil store", func(t *testing.T) {
		t.Parallel()

		var store *ElementsStore
		store.Set(Elements{[]byte(`"value"`)})
		assert.Nil(t, store)
	})

	t.Run("Get uninitialized store", func(t *testing.T) {
		t.Parallel()

		var store ElementsStore
		assert.Empty(t, store.Get())
	})
}
