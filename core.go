package noodle

import (
	"context"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
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

func Set(r *http.Request, key, value interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

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

// Then finalizes middleware Chain converting it to context-aware Handler
func (c Chain) Then(final Handler) Handler {
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}
