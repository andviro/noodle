package noodle

import (
	"context"
	"net/http"
)

// Middleware behaves like standard closure middleware pattern, only with
// context-aware handler type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain composes middlewares into a single context-aware handler
type Chain []Middleware

// New creates new middleware Chain and initalizes it with its parameters
func New(mws ...Middleware) Chain {
	return mws
}

// Set replaces request's context with a copy in which the value associated with key is val
func Set(r *http.Request, key, value interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

// Get calls Value() method on request's context and returns the result
func Get(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
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
