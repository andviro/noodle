package store

import (
	"errors"
	"golang.org/x/net/context"
	"sync"
)

const (
	globalContextKey = iota
	localContextKey
)

var NotFoundError = errors.New("Key not found in store")

// Atomic is a convenience type for function that's called in Read or Update
// store transaction
type Atomic func(map[string]interface{}) error

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
func (s *Store) Must(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	data, ok := s.data[key]
	if !ok {
		panic(NotFoundError)
	}
	return data
}

// Set saves value value to the store
func (s *Store) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

// Read runs supplied function in read-only atomic context. Note that function
// is granted arbitrary access to underlying map, and storing values in it will
// be possible but thread-unsafe.
func (s *Store) Read(f Atomic) error {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return f(s.data)
}

// Update runs supplied function in read-write atomic context. Use this method
// to read and modify underlying map in thread-safe way.
func (s *Store) Update(f Atomic) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return f(s.data)
}

// WithGlobal is a convenience function that creates a new Store and binds it
// to Context. There's no real difference from WithLocal function below.
func WithGlobal(ctx context.Context) context.Context {
	return context.WithValue(ctx, globalContextKey, New())
}

// WithLocal is a convenience function that creates a new Store and binds it
// to Context.
func WithLocal(ctx context.Context) context.Context {
	return context.WithValue(ctx, localContextKey, New())
}

// Global extracts Store from Context that was bound to it earlier with
// WithGlobal function call
func Global(ctx context.Context) (*Store, bool) {
	res, ok := ctx.Value(globalContextKey).(*Store)
	return res, ok
}

// Local extracts Store from Context that was bound to it earlier with
// WithLocal function call
func Local(ctx context.Context) (*Store, bool) {
	res, ok := ctx.Value(localContextKey).(*Store)
	return res, ok
}
