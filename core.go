package noodle

import (
	"golang.org/x/net/context"
	"net/http"
)

// Handler provides context-aware http.Handler with error return value for
// enhanced chaining
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

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
func (c Chain) Use(mws ...Middleware) (res Chain) {
	res = make([]Middleware, len(c)+len(mws))
	copy(res[:len(c)], c)
	copy(res[len(c):], mws)
	return
}

// Then finalizes middleware Chain converting it to context-aware Handler
func (c Chain) Then(final Handler) Handler {
	for i := len(c) - 1; i >= 0; i-- {
		final = c[i](final)
	}
	return final
}

// origin is the root context for all requests. By default it contains
// reference to global thread-safe Store. Origin can be extended or overwritten to
// provide common application-wide initial context.
var origin = context.TODO()

// ServeHTTP creates empty context and applies Handler to it, satisfying
// http.Handler interface

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = h(origin, w, r)
}