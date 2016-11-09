package noodle

import (
	"context"
	"net/http"
)

// Handler provides context-aware http.Handler with error return value for
// enhanced chaining
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// ServeHTTP applies Handler to request's own context, satisfying
// http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = h(r.Context(), w, r)
}

// Middleware behaves like standard closure middleware pattern, only with
// context-aware handler type
type Middleware func(Handler) Handler

// Chain composes middlewares into a single context-aware handler
type Chain []Middleware

// New creates new middleware Chain and initalizes it with its parameters
func New(mws ...Middleware) Chain {
	return mws
}

// Use appends its parameters to middleware chain. Returns new separate
// middleware chain
func (c Chain) Use(mws ...Middleware) Chain {
	res := make([]Middleware, len(c)+len(mws))
	copy(res[:len(c)], c)
	copy(res[len(c):], mws)
	return res
}

// Then finalizes middleware Chain converting it to context-aware Handler
func (c Chain) Then(final Handler) Handler {
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}
