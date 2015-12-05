package middleware

import (
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/store"
	"golang.org/x/net/context"
	"net/http"
)

// LocalStore is a middleware that injects common data store into
// request context
func LocalStore(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		return next(context.WithValue(c, storeKey, store.New()), w, r)
	}
}

// GetStore returns common store extracted from context or nil if no store was found
func GetStore(c context.Context) *store.Store {
	res, _ := c.Value(storeKey).(*store.Store)
	return res
}
