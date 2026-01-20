package data

import "sync/atomic"

// ElementsStore holds the current elements slice for lock-free reads.
type ElementsStore struct {
	v atomic.Value
}

// NewElementsStore initializes a store with the provided elements.
func NewElementsStore(elements Elements) *ElementsStore {
	store := &ElementsStore{}
	store.v.Store(elements)
	return store
}

// Get returns the latest elements snapshot.
func (s *ElementsStore) Get() Elements {
	if s == nil {
		return Elements{}
	}
	val := s.v.Load()
	if val == nil {
		return Elements{}
	}
	return val.(Elements)
}

// Set replaces the elements snapshot.
func (s *ElementsStore) Set(elements Elements) {
	if s == nil {
		return
	}
	s.v.Store(elements)
}
