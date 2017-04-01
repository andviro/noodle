package middleware

import (
	"gopkg.in/andviro/noodle.v2"
	"gopkg.in/andviro/noodle.v2/store"
	"net/http"
)

// LocalStore is a middleware that injects common data store into
// request context
func LocalStore(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.WithValue(r, storeKey, store.New()))
	}
}

// GetStore extracts common store from context
func GetStore(r *http.Request) *store.Store {
	res, _ := noodle.Value(r, storeKey).(*store.Store)
	return res
}
