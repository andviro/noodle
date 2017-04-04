package noodle

import (
	"context"
	"net/http"
)

type wrContextKey struct{}

var wrKey wrContextKey

type wrTuple struct {
	w http.ResponseWriter
	r *http.Request
}

// Middleware behaves like standard closure middleware pattern, only with
// context-aware handler type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain composes middlewares into a single context-aware handler
type Chain []Middleware

// New creates new middleware Chain and initalizes it with its parameters
func New(mws ...Middleware) Chain {
	return mws
}

// WithValue replaces request's context with a copy in which the value associated with key is val
func WithValue(r *http.Request, key, value interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

// Value conveniently calls the Value method on the request's context
func Value(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

// Wrap saves http.ResponseWriter and *http.Request into context
// taken from the original http.Request
func Wrap(w http.ResponseWriter, r *http.Request) context.Context {
	return context.WithValue(r.Context(), wrKey, wrTuple{w, r})
}

// Unwrap extracts http.ResponseWriter and *http.Request from context
func Unwrap(ctx context.Context) (w http.ResponseWriter, r *http.Request) {
	wr := ctx.Value(wrKey).(wrTuple)
	return wr.w, wr.r
}

// Use appends its parameters to middleware chain. Returns new separate
// middleware chain
func (c Chain) Use(mws ...Middleware) Chain {
	res := make([]Middleware, len(c)+len(mws))
	copy(res[:len(c)], c)
	copy(res[len(c):], mws)
	return res
}

// Then finalizes middleware Chain converting it to context-aware http.HandlerFunc
func (c Chain) Then(final http.HandlerFunc) http.HandlerFunc {
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}
