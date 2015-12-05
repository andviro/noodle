package middleware

import (
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/store"
	"golang.org/x/net/context"
	"net/http"
)

const (
	storeKey int = iota
)

// LocalStore is a middleware that injects common data store into
// request context
func LocalStore(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		return next(context.WithValue(c, storeKey, store.New()), w, r)
	}
}

// GetStore extracts common store from context
func GetStore(c context.Context) (*store.Store, bool) {
	res, ok := c.Value(storeKey).(*store.Store)
	return res, ok
}
