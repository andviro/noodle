package wok

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/andviro/noodle.v2"
	"net/http"
)

// Wok is a simple wrapper for httprouter with route groups and native support for http.HandlerFunc
type Wok struct {
	prefix string
	parent *Wok
	chain  noodle.Chain
	*httprouter.Router
}

// RouteClosure is a convenience type that allows setting route handlers by calling the function
type RouteClosure func(http.HandlerFunc)

// Handle allows route closure to attach a http.Handler
func (rc RouteClosure) Handle(h http.Handler) {
	rc(h.ServeHTTP)
}

type key int

const (
	paramKey key = iota
)

// convert turns http.Handler into httprouter.Handle
func (wok *Wok) convert(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), paramKey, p)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

// New creates new Wok initialized with middlewares.
// The resulting middleware chain will be called for all routes in Wok
func New(mws ...noodle.Middleware) *Wok {
	return &Wok{
		Router: httprouter.New(),
		chain:  noodle.New(mws...),
	}
}

// Handle allows to attach some noodle Middlewares and a Handle to a route
func (wok *Wok) Handle(method, path string, mws ...noodle.Middleware) RouteClosure {
	chain := noodle.New(mws...)
	for router := wok; router != nil; router = router.parent {
		chain = router.chain.Use(chain...)
		path = URLJoin(router.prefix, path)
	}
	return func(h http.HandlerFunc) {
		h = chain.Then(h)
		wok.Router.Handle(method, path, wok.convert(h))
	}
}

// Group starts new route group with common prefix.
// Middleware passed to Group will be used for all routes in it.
func (wok *Wok) Group(prefix string, mws ...noodle.Middleware) *Wok {
	return &Wok{
		prefix: prefix,
		parent: wok,
		Router: wok.Router,
		chain:  noodle.New(mws...),
	}
}

// Var returns route variable for context or empty string
func Var(r *http.Request, name string) string {
	return noodle.Value(r, paramKey).(httprouter.Params).ByName(name)
}
