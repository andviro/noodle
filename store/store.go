package store

import (
	"sync"
)

type KeyError struct {
	Key string
}

func (ke KeyError) Error() string {
	return "Key `" + ke.Key + "` not found in store"
}

// Store provides simple thread-safe storage for application. It supposed to be
// injected into context using WithStore function
type Store struct {
	data map[string]interface{}
	lock sync.RWMutex
}

// New creates new empty store
func New() *Store {
	return &Store{data: make(map[string]interface{})}
}

// Get reads value from the store, returns value and boolean flag
func (s *Store) Get(key string) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	data, ok := s.data[key]
	return data, ok
}

// Must reads value from the store and panics if there's no such key
func (s *Store) MustGet(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	data, ok := s.data[key]
	if !ok {
		panic(KeyError{key})
	}
	return data
}

// Set saves value to the store
func (s *Store) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

// Delete removes value from store
func (s *Store) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
}

// View executes a function within read-only atomic transaction. Note that function
// is granted arbitrary access to the underlying map, and writing operations
// are possible but thread-unsafe.
func (s *Store) View(f func(map[string]interface{}) error) error {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return f(s.data)
}

// Update executes a function in read-write atomic transaction. Use this method
// to modify underlying map in thread-safe way.
func (s *Store) Update(f func(map[string]interface{}) error) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return f(s.data)
}
