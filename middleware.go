package noodle

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

// Handler provides context-aware http.Handler with error return value for
// enhanced chaining
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// Middleware behaves like standard closure middleware pattern, only with
// context-aware handler type
type Middleware func(Handler) Handler

// Chain composes middlewares into a single context-aware handler
type Chain []Middleware

// Recover is basic middleware that catches panics and converts them into
// errors
func Recover(next Handler) Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("panic: %v", e)
			}
		}()
		err = next(c, w, r)
		return
	}
}

// Logger is a middleware that logs requests, along with
// request URI, handler return value and request timing
func Logger(next Handler) Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		start := time.Now()
		err = next(c, w, r)
		end := time.Now()
		log.Printf("%s %s from %s (%s) error = %v", r.Method, r.URL.String(), r.RemoteAddr, end.Sub(start), err)
		return
	}
}

// New creates new middleware Chain and initalizes it with its parameters
func New(mws ...Middleware) Chain {
	return mws
}

// Default creates new middleware Chain with Recover middleware on top
func Default(mws ...Middleware) Chain {
	return New(Logger, Recover).Use(mws...)
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

// ServeHTTP creates empty context and applies Handler to it, satisfying
// http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	_ = h(ctx, w, r)
}

// Adapt converts generic "dumb" middleware to context-aware, so that context
// is maintained throgout calling chain and error value is propagated correctly
func Adapt(mw func(http.Handler) http.Handler) Middleware {
	return func(next Handler) Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			var err error
			wrappedNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err = next(c, w, r)
			})
			mw(wrappedNext).ServeHTTP(w, r)
			return err
		}
	}
}

// AdaptNegroni converts negroni.HandlerFunc to context-aware Middleware, similar to
// Adapt().
func AdaptNegroni(mw func(http.ResponseWriter, *http.Request, http.HandlerFunc)) Middleware {
	return func(next Handler) Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			var err error
			wrappedNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err = next(c, w, r)
			})
			mw(w, r, wrappedNext)
			return err
		}
	}
}
