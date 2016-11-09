package middleware

import (
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/store"
	"net/http"
)

// LocalStore is a middleware that injects common data store into
// request context
func LocalStore(next noodle.Handler) noodle.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.Set(r, storeKey, store.New()))
	}
}

// GetStore extracts common store from context
func GetStore(r *http.Request) *store.Store {
	res, _ := noodle.Get(r, storeKey).(*store.Store)
	return res
}
